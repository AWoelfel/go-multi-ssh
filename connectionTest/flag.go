package connectionTest

import "github.com/spf13/pflag"

func Flag(flags *pflag.FlagSet, value *bool) {
	flags.BoolVarP(value, "connectionTest", "c", *value, "performs a connection test on all hosts in the index file fitting the selected tags")
}
