package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"ray-d-song.com/packer/hooks"
	"ray-d-song.com/packer/server"
	"ray-d-song.com/packer/utils"
)

var rootCmd = &cobra.Command{
	Use:   "pack",
	Short: "Packer is a simplified implementation of a package manager.",
	Long: `Packer is a simplified implementation of a package manager.
It includes a command-line program for managing project libraries, as well as a server program for managing repositories.`,
	Run: func(cmd *cobra.Command, args []string) {

		hooks.EveryStart()

		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		if args[0] == "server" {
			server.SetupServer()
			return
		}

		if !utils.CheckConfigExists() {
			fmt.Println("Configuration file does not exist")
			return
		}

		readConfig()
		if args[0] == "sync" {
			deps := viper.GetStringSlice("dependencies")
			diff := utils.DiffDeps(deps)

			registry := viper.GetString("registry")
			var wg sync.WaitGroup
			for _, dep := range diff {
				wg.Add(1)
				go func(dep string) {
					utils.Download(registry, dep)
				}(dep)
			}

			wg.Wait()
			hooks.AfterSync()
			return
		}

		if args[0] == "run" {
			if len(args) < 2 {
				fmt.Println("No script specified")
				return
			}
			script := args[1]
			scriptMap := viper.GetStringMapString("scripts")

			if value, exists := scriptMap[script]; exists {
				utils.Execute(value)
			} else {
				fmt.Println("Script does not exist")
			}
			return
		}

		if args[0] == "publish" {

		}

		cmd.Help()
		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readConfig() {
	viper.SetConfigName("pack")
	configPath, err := utils.GetExecPath()
	if err != nil {
		panic(err)
	}
	viper.AddConfigPath(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
