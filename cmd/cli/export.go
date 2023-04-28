package cli

import (
	"fmt"
	"github.com/erminson/gitlab-vars/internal/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// TODO: Rename to Import
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export list of variables (json)",
	Long:  `Export list of variables (json)`,
	Run: func(cmd *cobra.Command, args []string) {
		projectId := viper.GetInt64("project-id")
		uc := usecase.NewUseCase(projectId, client)

		err := uc.ImportVariablesFromFile(Filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
