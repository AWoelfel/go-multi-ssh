package main

import (
	sshConfig "github.com/AWoelfel/go-multi-ssh/cmd/ssh/config"
	sshConnectionTest "github.com/AWoelfel/go-multi-ssh/cmd/ssh/connectionTest"
	"github.com/AWoelfel/go-multi-ssh/cmd/ssh/single"
	"github.com/AWoelfel/go-multi-ssh/config"
	commonConnectionTest "github.com/AWoelfel/go-multi-ssh/connectionTest"
	"github.com/AWoelfel/go-multi-ssh/output"
	"github.com/spf13/cobra"
)

func MainCommand() *cobra.Command {

	cmdConfig := sshConfig.DefaultConfig

	cmdOutputLabelMode := config.NoOutputLabel
	connectionTest := false

	var rootCmd = &cobra.Command{
		Use: "mssh",
		PreRun: func(cmd *cobra.Command, args []string) {
			ctx := output.DefaultContextWithWriters(cmd.Context(), cmdOutputLabelMode)
			cmd.SetContext(cmdConfig.WithContext(ctx))
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if connectionTest {
				return sshConnectionTest.Execute(cmd.Context())
			}

			return single.Execute(cmd.Context(), args)
		},
	}

	commonConnectionTest.Flag(rootCmd.Flags(), &connectionTest)
	config.OutputLabelFlag(rootCmd.Flags(), &cmdOutputLabelMode)
	rootCmd.Flags().StringArrayVarP(&cmdConfig.IncludeTags, "tag", "t", cmdConfig.IncludeTags, "only access target hosts with the given tags (defaults to include all hosts)")
	rootCmd.Flags().StringArrayVarP(&cmdConfig.ExcludeTags, "!tag", "T", cmdConfig.ExcludeTags, "ignore target hosts with the given tags (defaults to ignore no hosts)")
	rootCmd.Flags().StringVarP(&cmdConfig.IndexFile, "index", "i", cmdConfig.IndexFile, "index file containing the target hosts")

	return rootCmd
}

func main() {
	MainCommand().Execute()
}
