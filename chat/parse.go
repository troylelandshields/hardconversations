package chat

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	numberRegex = regexp.MustCompile(`-*[0-9\.]+`)
)

func Parse(text string, v interface{}) error {
	if strings.HasPrefix(text, "Error:") {
		return errors.Errorf("Unable to process request with error: %s", strings.TrimPrefix(text, "Error: "))
	}

	t := reflect.TypeOf(v)

	if t.Kind() != reflect.Pointer {
		return errors.New("v must be a pointer")
	}

	t = t.Elem()

	// just JSON unmarshal if it's a struct or slice of structs
	if t.Kind() == reflect.Struct || (t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Struct) {
		startChar := "{"
		endChar := "}"
		if t.Kind() == reflect.Slice {
			startChar = "["
			endChar = "]"
		}

		startIdx := strings.Index(text, startChar)
		if startIdx != -1 {
			text = text[startIdx:]
		}

		endIdx := strings.LastIndex(text, endChar)
		if endIdx != -1 {
			text = text[:endIdx+1]
		}

		if text == "" {
			text = startChar + endChar
		}

		err := json.Unmarshal([]byte(text), v)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to unmarshal JSON: %s", text))
		}
		return nil
	}

	switch t.Kind() {
	case reflect.Slice:
		splitText := strings.Split(text, ",")
		results := reflect.MakeSlice(t, len(splitText), len(splitText))
		for i, elem := range splitText {
			elem = strings.TrimSpace(elem)
			err := Parse(elem, results.Index(i).Addr().Interface())
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to parse slice element: [%s]", elem))
			}
		}
		reflect.ValueOf(v).Elem().Set(results)
	case reflect.Bool:
		text = strings.Trim(text, ".,!?")
		b, err := strconv.ParseBool(text)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse bool: [%s]", text))
		}
		reflect.ValueOf(v).Elem().SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		text = numberRegex.FindString(text)
		i, err := strconv.Atoi(text)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse int: [%s]", text))
		}
		reflect.ValueOf(v).Elem().SetInt(int64(i))
	case reflect.Float32, reflect.Float64:
		text = numberRegex.FindString(text)
		f, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse float: [%s]", text))
		}

		reflect.ValueOf(v).Elem().SetFloat(f)
	case reflect.String:
		reflect.ValueOf(v).Elem().SetString(text)

	}
	return nil
}

func ParseInstruction(v interface{}) (string, error) {
	t := reflect.TypeOf(v)

	var pluralityInstruction string
	if t.Kind() == reflect.Slice {
		pluralityInstruction = " (separate multiple answers with commas)"

		// TODO: I don't love this
		if t.Elem().Kind() == reflect.Struct {
			pluralityInstruction = " (provide answer as a JSON array)"
		}

		t = t.Elem()
	}

	answerTypeInstruction, err := answerTypeInstruction(t)
	if err != nil {
		return "", err
	}

	return answerTypeInstruction + pluralityInstruction + ": ", nil
}

func answerTypeInstruction(t reflect.Type) (string, error) {
	switch t.Kind() {
	case reflect.Pointer:
		return answerTypeInstruction(t.Elem())
	case reflect.Slice:
	case reflect.Bool:
		return `Answer this with exactly "true" or "false" only, no punctuation`, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return `Answer this with an integer only, no punctuation or explanation`, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return `Answer this with a positive integer only, no punctuation or explanation`, nil
	case reflect.Float32, reflect.Float64:
	case reflect.String:
		return "", nil
	case reflect.Struct:
		prompt := "Provide the answer as a JSON object or array that looks like the following\n"
		d, err := json.Marshal(reflect.Zero(t).Interface())
		if err != nil {
			return "", err
		}
		prompt += string(d)

		var fieldExplanations string
		// targetValue := t.Elem()
		for i := 0; i < t.NumField(); i++ {
			customInstruction := t.Field(i).Tag.Get("hardc-instruction")
			if customInstruction == "" {
				continue
			}

			if fieldExplanations != "" {
				fieldExplanations += "; "
			}

			fieldExplanations += t.Field(i).Name + " " + customInstruction
		}

		if fieldExplanations != "" {
			prompt += "\n where " + fieldExplanations
		}

		return prompt, nil
	}

	return "", errors.Errorf("unsupported type: %s", t.Kind().String())
}

type PromptInput interface {
	PromptInput() string
}

func ConvertInput(v interface{}) (string, error) {
	switch t := v.(type) {
	case PromptInput:
		return t.PromptInput(), nil
	case bool:
		return strconv.FormatBool(t), nil
	case int:
		return strconv.Itoa(t), nil
	case uint:
		strconv.FormatUint(uint64(t), 10)
	case string:
		return t, nil
	case float32, float64:
		return strconv.FormatFloat(t.(float64), 'f', -1, 64), nil
	case []string:
		return strings.Join(t, ", "), nil
	case []int:
		var s []string
		for _, i := range t {
			s = append(s, strconv.Itoa(i))
		}
		return strings.Join(s, ", "), nil
	case []float64:
		var s []string
		for _, f := range t {
			s = append(s, strconv.FormatFloat(f, 'f', -1, 64))
		}
		return strings.Join(s, ", "), nil
	case []float32:
		var s []string
		for _, f := range t {
			s = append(s, strconv.FormatFloat(float64(f), 'f', -1, 64))
		}
		return strings.Join(s, ", "), nil
	case []bool:
		var s []string
		for _, b := range t {
			s = append(s, strconv.FormatBool(b))
		}
		return strings.Join(s, ", "), nil
	}

	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
