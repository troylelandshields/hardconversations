package tokens

import (
	"github.com/pkg/errors"
	tokenizer "github.com/samber/go-gpt-3-encoder"
)

var encoder *tokenizer.Encoder

func init() {
	var err error
	encoder, err = tokenizer.NewEncoder()
	if err != nil {
		panic(err)
	}
}

func MustCount(t string) int {
	count, err := Count(t)
	if err != nil {
		panic(err)
	}
	return count
}

func Count(t string) (int, error) {
	encoded, err := encoder.Encode(t)
	if err != nil {
		return 0, errors.Wrap(err, "error encoding text")
	}

	return len(encoded), nil
}

func Chunk(t string, maxTokenSize int) ([]string, error) {
	encoded, err := encoder.Encode(t)
	if err != nil {
		return nil, errors.Wrap(err, "error encoding text")
	}

	if len(encoded) <= maxTokenSize {
		return []string{t}, nil
	}

	var chunks []string
	var chunkSize int
	var chunk string
	for _, token := range encoded {
		if chunkSize+1 > maxTokenSize {
			chunks = append(chunks, chunk)
			chunkSize = 0
			chunk = ""
		}

		str := encoder.Decode([]int{token})
		chunkSize += 1
		chunk += str
	}

	if chunk != "" {
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}
