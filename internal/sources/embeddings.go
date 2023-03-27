package sources

import (
	"context"
	"sort"
	"strings"

	"github.com/drewlanenga/govector"
	"github.com/pkg/errors"
	gogpt "github.com/sashabaranov/go-openai"
)

type TextEmbedding struct {
	Text       string
	Embedding  []float32
	TokenCount int
}

type TextEmbeddingProvider interface {
	Sources(ctx context.Context, prompt string) ([]TextEmbedding, error)
}

// TODO: figure out how to make this a good idea and make it public if so
func (t *Manager) addSourceTextEmbeddingProvider(provider TextEmbeddingProvider, opt ...SourceOption[TextEmbeddingProvider]) {
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

// TODO: handle userID another way
func (t *Manager) getTextEmbeddings(ctx context.Context, textEmbeddings []TextEmbedding, userID string) ([]TextEmbedding, error) {
	var inputs []string
	for _, te := range textEmbeddings {
		if len(te.Embedding) > 0 {
			continue
		}

		text := te.Text

		text = strings.ReplaceAll(text, "\n", " ")
		text = strings.ReplaceAll(text, "\\n", " ")
		text = strings.ReplaceAll(text, "  ", " ")
		text = strings.ReplaceAll(text, "  ", " ")

		inputs = append(inputs, text)
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
		textEmbeddings[i].Embedding = embedding.Embedding
	}

	return textEmbeddings, nil
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
