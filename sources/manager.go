package sources

import (
	"context"
	"sort"

	"github.com/pkg/errors"
	gogpt "github.com/sashabaranov/go-openai"
	"github.com/troylelandshields/hardconversations/internal/tokens"
	"github.com/troylelandshields/hardconversations/logger"
)

type Manager struct {
	ai            *gogpt.Client
	textProviders []source[TextEmbeddingProvider]
}

func New(openAIClient *gogpt.Client) *Manager {
	return &Manager{
		ai:            openAIClient,
		textProviders: []source[TextEmbeddingProvider]{},
	}
}

func NewFromParent(m *Manager) *Manager {
	copiedProviders := make([]source[TextEmbeddingProvider], len(m.textProviders))
	copy(copiedProviders, m.textProviders)

	return &Manager{
		ai:            m.ai,
		textProviders: copiedProviders,
	}
}

type source[T any] struct {
	provider    T
	weight      float64
	maxTokens   int
	allowErrors bool
}

// add other source options, like description
type SourceOption[T any] func(*source[T])

// WithWeight sets the weight of a source. The weight is used to give more or less priority to a source. The default weight is 1.0.
// Sources with a higher weight are used first, and if text embeddings are being used, the weight is multiplied by the cosine similarity
// between the prompt and the source to make it more or less likely that the source will be used.
// TODO: I don't like this using a generic.
func WithWeight(w float64) SourceOption[TextEmbeddingProvider] {
	return func(s *source[TextEmbeddingProvider]) {
		s.weight = w
	}
}

// WithMaxTokens sets the max amount of tokens this source can contribute to the contextual info. The default is 0, which means there is no limit
func WithMaxTokens(m int) SourceOption[TextEmbeddingProvider] {
	return func(s *source[TextEmbeddingProvider]) {
		s.maxTokens = m
	}
}

// WithAllowErrors means that if this source errors, it will be ignored and the next source will be used instead
func WithAllowErrors() SourceOption[TextEmbeddingProvider] {
	return func(s *source[TextEmbeddingProvider]) {
		s.allowErrors = true
	}
}

// GetSourceText pulls text from the sources in order of weight until we run out of tokens. If sortByRelevance is true, then we will consider cosine similarity between prompt and text.
// TODO: this is a bit of a mess, clean it up; also this might be the wrong place to be creating text embeddings since it could lead to a lot of repeated work
// TODO: support chunking text into smaller pieces
// TODO: I'm slapping userID as an optional param in here so I can pass it to OpenAI but I don't like it, figure out a better way
func (t *Manager) GetSourceText(ctx context.Context, sortByRelevance bool, minCosineSimilarityThreshold float64, allowedTokens int, prompt string, userID string) ([]string, error) {
	if !sortByRelevance {
		return t.getSourceTextSimple(ctx, allowedTokens, prompt)
	}

	return t.getSourceTextRelevant(ctx, minCosineSimilarityThreshold, allowedTokens, prompt, userID)
}

// getSourceTextSimple just pulls text from the sources in order of weight until we run out of tokens
func (t *Manager) getSourceTextSimple(ctx context.Context, allowedTokens int, prompt string) ([]string, error) {
	var contextualInfos []string

	logger.Debugf("Pulling contextual info from source %d text providers...", len(t.textProviders))
	for _, source := range t.textProviders {
		var sourceUsedTokens int

		// use either the max tokens set by the source, or the allowed tokens left over, which ever is smaller
		sourceMaxTokens := allowedTokens
		if source.maxTokens != 0 && source.maxTokens < allowedTokens {
			sourceMaxTokens = source.maxTokens
		}

		sourceTextEmbeddings, err := source.provider.Sources(ctx, prompt)
		if err != nil {
			if !source.allowErrors {
				return nil, err
			}
			logger.Debugf("Source %T errored: %v", source.provider, err)
		}

		for _, sourceTextEmbedding := range sourceTextEmbeddings {
			tokenCnt, err := tokens.Count(sourceTextEmbedding.Text)
			if err != nil {
				if !source.allowErrors {
					return nil, err
				}
				logger.Debugf("Source %T errored: %v", source.provider, err)
			}

			// this source has used up all of its tokens, move on to the next source
			if tokenCnt+sourceUsedTokens > sourceMaxTokens {
				logger.Debugf("Could not pull all contextual info for source, stopping at %d tokens (max of %d)", sourceUsedTokens, sourceMaxTokens)
				break
			}

			sourceUsedTokens += tokenCnt
			contextualInfos = append(contextualInfos, sourceTextEmbedding.Text)
		}
		allowedTokens -= sourceUsedTokens
		if allowedTokens <= 0 {
			break
		}
	}

	return contextualInfos, nil
}

type contextualInfo struct {
	Text                     string
	TokenCnt                 int
	WeightedCosineSimilarity float64
	tokensLeft               *int
}

// getSourceTextRelevant gets all the sources, filter and sort by cosine similarity, then pull the top ones until we run out of tokens
func (t *Manager) getSourceTextRelevant(ctx context.Context, minCosineSimilarityThreshold float64, allowedTokens int, prompt string, userID string) ([]string, error) {
	var contextualInfos []contextualInfo

	promptEmbeddings, err := t.getTextEmbeddings(ctx, []TextEmbedding{{Text: prompt}}, userID)
	if err != nil {
		return nil, err
	}
	promptEmbedding := promptEmbeddings[0]

	logger.Debugf("Pulling contextual info from source %d text providers...", len(t.textProviders))
	for _, source := range t.textProviders {
		var allSourceInfo []TextEmbedding
		sourceMax := source.maxTokens
		if sourceMax == 0 {
			sourceMax = allowedTokens
		}
		sourceTextEmbeddings, err := source.provider.Sources(ctx, prompt)
		if err != nil {
			if !source.allowErrors {
				return nil, err
			}
			logger.Debugf("Source %T errored: %v", source.provider, err)
			continue
		}

		var missingEmbeddings bool
		for _, sourceTextEmbedding := range sourceTextEmbeddings {
			tokenCnt, err := tokens.Count(sourceTextEmbedding.Text)
			if err != nil {
				if !source.allowErrors {
					return nil, err
				}
				logger.Debugf("Source %T errored: %v", source.provider, err)
				continue
			}
			sourceTextEmbedding.TokenCount = tokenCnt

			// get text embeddings if we need them
			// if any of the sources don't have embeddings, then we will get them all
			if len(sourceTextEmbedding.Embedding) == 0 {
				missingEmbeddings = true

				// if the text is too long, skip it until we can support chunking
				if tokenCnt > 8192 {
					logger.Debugf("Source %T text too long to create embeddings and I haven't added support for chunking yet, skipping", source.provider)
					continue
				}
			}

			allSourceInfo = append(allSourceInfo, sourceTextEmbedding)
		}

		if missingEmbeddings {
			allSourceInfo, err = t.getTextEmbeddings(ctx, allSourceInfo, userID)
			if err != nil {
				if !source.allowErrors {
					return nil, err
				}
				logger.Debugf("Source %T errored: %v", source.provider, err)
				continue
			}
		}

		for _, sourceTextEmbedding := range allSourceInfo {
			// get cosine similarity
			cosineSimilarity, err := t.cosineSimilarity(sourceTextEmbedding.Embedding, promptEmbedding.Embedding)
			if err != nil {
				if !source.allowErrors {
					return nil, errors.Wrap(err, "failed to get cosine similarity")
				}
				logger.Debugf("Source %T errored: %v", source.provider, err)
				continue
			}

			weightedCosineSimilarity := float64(cosineSimilarity) * source.weight
			if weightedCosineSimilarity < minCosineSimilarityThreshold {
				logger.Debugf("Cosine similarity of %f is below threshold of %f, skipping", cosineSimilarity, minCosineSimilarityThreshold)
				continue
			}

			contextualInfos = append(contextualInfos, contextualInfo{
				Text:                     sourceTextEmbedding.Text,
				TokenCnt:                 sourceTextEmbedding.TokenCount,
				WeightedCosineSimilarity: weightedCosineSimilarity,
				tokensLeft:               &sourceMax,
			})
		}
	}

	// sort contextualInfos by weighted cosine similarity descending
	sort.Slice(contextualInfos, func(i, j int) bool {
		return contextualInfos[i].WeightedCosineSimilarity > contextualInfos[j].WeightedCosineSimilarity
	})

	// add contextualInfos to contextualText until we run out of tokens
	var contextualText []string
	for _, ci := range contextualInfos {
		// if not enough tokens left, skip it
		if ci.TokenCnt > *ci.tokensLeft || ci.TokenCnt > allowedTokens {
			continue
		}
		// reduce source's tokens left and allowed tokens
		*ci.tokensLeft -= ci.TokenCnt
		allowedTokens -= ci.TokenCnt
		contextualText = append(contextualText, ci.Text)
		if allowedTokens <= 0 {
			break
		}
	}

	return contextualText, nil
}
