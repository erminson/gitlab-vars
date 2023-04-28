package cli

import (
	"fmt"
	"github.com/erminson/gitlab-vars/internal/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show list of variables",
	Long:  `Show list of variables`,
	Run: func(cmd *cobra.Command, args []string) {
		projectId := viper.GetInt64("project-id")
		uc := usecase.NewUseCase(projectId, client)

		vars, err := uc.ListVariables()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(vars)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
