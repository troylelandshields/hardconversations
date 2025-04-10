package sources

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/drewlanenga/govector"
	"github.com/pkg/errors"
	gogpt "github.com/sashabaranov/go-openai"
	"github.com/troylelandshields/hardconversations/internal/tokens"
)

const (
	maxEmbeddingTokenCount = 2048
)

type TextEmbedding struct {
	// Identifier that uniquely identifies this text embedding (e.g., a pageID or documentID). If this embedding is chunked, this Identifier will have the chunk idx appended to it
	Identifier string
	// Text that can be used as a data source for the AI
	Text string
	// Embedding of the text; if it will be used and it is not provided then API requests will be made to get it
	Embedding []float32
	// Optional metadata that can be used to identify the source; sources that get used will be returned in the metadata response, so this field can be used to pass more information about the source
	Metadata interface{}

	// Optional weight to use for this specific source text; defaults to the weight of the source provider
	Weight float64

	chunk       int // if the text is chunked, this is the chunk number
	totalChunks int
	tokenCount  int
}

func (t *TextEmbedding) TokenCount() int {
	return t.tokenCount
}

func (t *TextEmbedding) Chunk() int {
	return t.chunk
}

func (t *TextEmbedding) TotalChunks() int {
	return t.totalChunks
}

type TextEmbeddingProvider interface {
	Sources(ctx context.Context, prompt string) ([]TextEmbedding, error)
}

// TODO: figure out how to make this a good idea and make it public if so
func (t *Manager) AddSourceTextEmbeddingProvider(provider TextEmbeddingProvider, opt ...SourceOption[TextEmbeddingProvider]) {
	source := source[TextEmbeddingProvider]{provider: provider, weight: 1.0}
	for _, o := range opt {
		o(&source)
	}

	t.textProviders = append(t.textProviders, source)

	// sort by weight
	sort.Slice(t.textProviders, func(i, j int) bool {
		return t.textProviders[i].weight > t.textProviders[j].weight
	})
}

func (t *Manager) CreateTextEmbeddingsFromStrings(ctx context.Context, text []string, userID string) ([]TextEmbedding, error) {
	var embeddings []TextEmbedding
	for _, te := range text {
		embeddings = append(embeddings, TextEmbedding{Text: te})
	}

	return t.CreateTextEmbeddings(ctx, embeddings, userID)
}

func (t *Manager) CreateTextEmbeddings(ctx context.Context, textEmbeddings []TextEmbedding, userID string) ([]TextEmbedding, error) {
	return t.prepareForQuerying(ctx, textEmbeddings, userID, false)
}

// TODO: handle userID another way
func (t *Manager) prepareForQuerying(ctx context.Context, textEmbeddings []TextEmbedding, userID string, skipEmbeddings bool) ([]TextEmbedding, error) {
	var inputs []string
	var results []TextEmbedding
	var err error
	for _, te := range textEmbeddings {
		if len(te.Embedding) > 0 {
			te.tokenCount, err = tokens.Count(te.Text)
			if err != nil {
				return nil, errors.Wrap(err, "error counting tokens")
			}
			results = append(results, te)
			continue
		}

		identifier := te.Identifier
		text := te.Text

		chunks, err := tokens.Chunk(text, maxEmbeddingTokenCount)
		if err != nil {
			return nil, errors.Wrap(err, "error chunking text")
		}

		for i, chunk := range chunks {
			tokenCnt, err := tokens.Count(chunk)
			if err != nil {
				return nil, errors.Wrap(err, "error counting tokens")
			}
			inputs = append(inputs, textEmbeddingPrep(chunk))
			results = append(results, TextEmbedding{
				Identifier:  identifier + "--" + strconv.Itoa(i),
				Text:        chunk,
				Weight:      te.Weight,
				Metadata:    te.Metadata,
				chunk:       i,
				totalChunks: len(chunks),
				tokenCount:  tokenCnt,
			})
		}
	}

	// no inputs were missing embeddings, so no need to make a request
	if len(inputs) == 0 || skipEmbeddings {
		return results, nil
	}

	resp, err := t.ai.CreateEmbeddings(ctx, gogpt.EmbeddingRequest{
		Input: inputs,
		Model: gogpt.AdaEmbeddingV2,
		User:  userID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating embeddings")
	}

	for i, embedding := range resp.Data {
		if len(results[i].Embedding) > 0 {
			continue
		}
		results[i].Embedding = embedding.Embedding
	}

	return results, nil
}

func (t *Manager) cosineSimilarity(a, b []float32) (float64, error) {
	aVec, err := govector.AsVector(a)
	if err != nil {
		return 0, err
	}
	bVec, err := govector.AsVector(b)
	if err != nil {
		return 0, err
	}
	return govector.Cosine(aVec, bVec)
}

func textEmbeddingPrep(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\\n", " ")
	s = strings.ReplaceAll(s, "  ", " ")
	s = strings.ReplaceAll(s, "  ", " ")
	return s
}
