package smartling

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Errors []struct {
		Key     string
		Message string
	}
}

func (err ValidationError) Error() string {
	messages := []string{}

	for _, err := range err.Errors {
		message := "\n- "

		if err.Key != "" {
			message += fmt.Sprintf("%s: ", err.Key)
		}

		message += err.Message

		messages = append(messages, message)
	}

	return "Smartling replies with validation error" +
		strings.Join(messages, "")
}
