package cli

import (
	"fmt"
	"github.com/erminson/gitlab-vars/internal/usecase"
	"github.com/spf13/cobra"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show list of variables",
	Long:  `Show list of variables`,
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := cmd.Flags().GetInt64("project")
		uc := usecase.NewUseCase(projectID, client)

		vars, err := uc.PrintVariables()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(vars)
	},
}

var ProjectId int64

func init() {
	listCmd.PersistentFlags().Int64VarP(&ProjectId, "project", "p", 0, "verbose output")
	rootCmd.AddCommand(listCmd)
}
