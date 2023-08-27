package main

import (
	k8sConfig "github.com/AWoelfel/go-multi-ssh/cmd/k8s/config"
	k8sConnectionTest "github.com/AWoelfel/go-multi-ssh/cmd/k8s/connectionTest"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/k8sClient"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/single"
	"github.com/AWoelfel/go-multi-ssh/config"
	commonConnectionTest "github.com/AWoelfel/go-multi-ssh/connectionTest"
	"github.com/AWoelfel/go-multi-ssh/output"
	"github.com/spf13/cobra"
)

func MainCommand() *cobra.Command {

	cmdConfig := k8sConfig.DefaultConfig

	cmdOutputLabelMode := config.NoOutputLabel
	connectionTest := false

	var rootCmd = &cobra.Command{
		Use:  "kubectl mexec",
		Long: `Runs the given command on all selected Pods/Containers. The list of targets consists of all selected Containers in non-completed Pods.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := output.DefaultContextWithWriters(cmd.Context(), cmdOutputLabelMode)
			ctx = cmdConfig.WithContext(ctx)

			var err error
			ctx, err = k8sClient.AttachClient(ctx)

			if err != nil {
				return err
			}
			cmd.SetContext(ctx)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if connectionTest {
				return k8sConnectionTest.Execute(cmd.Context())
			}

			return single.Execute(cmd.Context(), args)
		},
	}

	commonConnectionTest.Flag(rootCmd.Flags(), &connectionTest)
	config.OutputLabelFlag(rootCmd.Flags(), &cmdOutputLabelMode)
	rootCmd.Flags().StringArrayVarP(&cmdConfig.SearchLabels, "selector", "s", cmdConfig.SearchLabels, "limits the search to pods with the given label selector i.e. \"app.kubernetes.io/component=database\"")

	return rootCmd
}

func main() {
	MainCommand().Execute()
}
