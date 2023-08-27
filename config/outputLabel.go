package config

import (
	"errors"
	"github.com/spf13/pflag"
)

type OutputLabel string

const (
	NoOutputLabel     OutputLabel = "none"
	BlockOutputLabel  OutputLabel = "block"
	InlineOutputLabel OutputLabel = "inline"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *OutputLabel) String() string {
	return string(*e)
}

// Set must have pointer receiver bcs it should not change the value of a copy
func (e *OutputLabel) Set(v string) error {
	switch v {
	case "none", "block", "inline":
		*e = OutputLabel(v)
		return nil
	default:
		return errors.New(`must be one of "none", "block", or "inline"`)
	}
}

func (e *OutputLabel) Type() string {
	return "OutputLabel"
}

func OutputLabelFlag(flags *pflag.FlagSet, value *OutputLabel) {
	flags.VarP(value, "label", "l", `label mode for client outputs. Does not work in "connection test" mode. allowed values are: "none", "block", "inline"`)
}
