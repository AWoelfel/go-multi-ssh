package utils

import (
	"strings"
)

func FromErrors(errs ...error) error {

	if len(errs) == 0 {
		return nil
	}

	if len(errs) == 1 {
		return errs[0]
	}

	var result Errors

	for i := 0; i < len(errs); i++ {

		switch errs[i].(type) {
		case Errors:
			if result == nil {
				result = errs[i].(Errors)
			} else {
				result = append(result, errs[i].(Errors)...)
			}
		case nil:
			continue
		default:
			if errs[i] != nil {
				result = append(result, errs[i])
			}
		}

	}

	if len(result) == 0 {
		return nil
	}

	if len(result) == 1 {
		return result[0]
	}

	return result
}

type Errors []error

func (e Errors) Error() string {

	var msgs []string

	for i := 0; i < len(e); i++ {
		eMsg := e[i].Error()
		if len(eMsg) == 0 {
			continue
		}
		msgs = append(msgs, eMsg)
	}

	return strings.Join(msgs, "\n")
}
