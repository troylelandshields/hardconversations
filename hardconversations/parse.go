package hardconversations

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func Parse(text string, v interface{}) error {
	switch target := v.(type) {
	case *bool:
		text = strings.Trim(text, ".,!?")
		b, err := strconv.ParseBool(text)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse bool: [%s]", text))
		}
		*target = b
	case *int:
		text = strings.Trim(text, ".,!?")
		i, err := strconv.Atoi(text)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse int: [%s]", text))
		}
		*target = i
	case *string:
		*target = text
	default:
		err := json.Unmarshal([]byte(text), target)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to unmarshal JSON: [%s]", text))
		}
	}
	return nil
}

func ParseInstruction(v interface{}) (string, error) {
	switch target := v.(type) {
	case *bool, bool:
		return `Answer this with exactly "true" or "false" only, no punctuation: `, nil
	case *int, int:
		return `Answer this with an integer only, no punctuation or explanation: `, nil
	case *string, string:
		return "", nil
	default:
		prompt := "Provide the answer as a JSON object that looks like the following: \n"
		d, err := json.Marshal(target)
		if err != nil {
			return "", err
		}
		prompt += string(d)

		var fieldExplanations string
		targetValue := reflect.ValueOf(target)
		for i := 0; i < targetValue.NumField(); i++ {
			customInstruction := targetValue.Type().Field(i).Tag.Get("hardc-instruction")
			if customInstruction == "" {
				continue
			}

			if fieldExplanations != "" {
				fieldExplanations += "; "
			}

			fieldExplanations += targetValue.Type().Field(i).Name + " " + customInstruction
		}

		if fieldExplanations != "" {
			prompt += "\n where " + fieldExplanations
		}

		return prompt, nil
	}
}
