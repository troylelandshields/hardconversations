package sources

import "context"

type TextProvider interface {
	Sources(ctx context.Context, prompt string) ([]string, error)
}

func (t *Manager) AddSourceText(text string, opt ...SourceOption[TextEmbeddingProvider]) {
	t.addSourceTextEmbeddingProvider(hardCodedTextProvider{Text: text}, opt...)
}

func (t *Manager) AddSourceTextProvider(provider TextProvider, opt ...SourceOption[TextEmbeddingProvider]) {
	t.addSourceTextEmbeddingProvider(textProviderAdapter{provider}, opt...)
}

// adapters to convert a TextProvider to a TextEmbeddingProvider--the chat client will fill in empty embeddings if necessary
type hardCodedTextProvider TextEmbedding

func (p hardCodedTextProvider) Sources(ctx context.Context, prompt string) ([]TextEmbedding, error) {
	return []TextEmbedding{TextEmbedding(p)}, nil
}

type textProviderAdapter struct {
	TextProvider
}

func (a textProviderAdapter) Sources(ctx context.Context, prompt string) ([]TextEmbedding, error) {
	var result []TextEmbedding
	TextEmbeddings, err := a.TextProvider.Sources(ctx, prompt)
	if err != nil {
		return nil, err
	}

	for _, text := range TextEmbeddings {
		result = append(result, TextEmbedding{Text: text})
	}

	return result, nil
}
