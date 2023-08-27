package main

import (
	dockerConfig "github.com/AWoelfel/go-multi-ssh/cmd/docker/config"
	dockerConnectionTest "github.com/AWoelfel/go-multi-ssh/cmd/docker/connectionTest"
	"github.com/AWoelfel/go-multi-ssh/cmd/docker/dockerClient"
	"github.com/AWoelfel/go-multi-ssh/cmd/docker/single"
	"github.com/AWoelfel/go-multi-ssh/config"
	"github.com/AWoelfel/go-multi-ssh/output"
	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
	"os"
)

import (
	commonConnectionTest "github.com/AWoelfel/go-multi-ssh/connectionTest"
)

func pluginMain() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {

		cmdConfig := dockerConfig.DefaultConfig

		cmdOutputLabelMode := config.NoOutputLabel
		connectionTest := false

		var rootCmd = &cobra.Command{
			Short: "docker mexec",
			Long:  `Runs the given command on all selected Containers. The list of targets consists of all selected running containers.`,
			Use:   "mexec",
			PreRunE: func(cmd *cobra.Command, args []string) error {

				ctx := output.ContextWithWriters(cmd.Context(), cmdOutputLabelMode, dockerCli.Out(), dockerCli.Err())

				ctx = cmdConfig.WithContext(ctx)
				ctx = dockerClient.WithClient(ctx, dockerCli)

				cmd.SetContext(ctx)
				return nil
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				if connectionTest {
					return dockerConnectionTest.Execute(cmd.Context())
				}

				return single.Execute(cmd.Context(), args)
			},
		}

		commonConnectionTest.Flag(rootCmd.Flags(), &connectionTest)
		config.OutputLabelFlag(rootCmd.Flags(), &cmdOutputLabelMode)
		rootCmd.Flags().StringArrayVarP(&cmdConfig.SearchLabels, "selector", "s", cmdConfig.SearchLabels, "limits the search to pods with the given label selector i.e. \"com.example.foo=bar\"")

		return rootCmd

	},
		manager.Metadata{
			SchemaVersion: "0.1.0",
			Vendor:        "AWoelfel",
			Version:       "0.0.1",
		})
}

func main() {
	if plugin.RunningStandalone() {
		os.Args = append([]string{"docker"}, os.Args[1:]...)
	}
	pluginMain()
}
