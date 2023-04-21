package cli

import (
	"fmt"
	glvarsapi "github.com/erminson/gitlab-vars/internal/client"
	"github.com/erminson/gitlab-vars/internal/usecase"
	"github.com/spf13/cobra"
	"os"
)

var createCmd = &cobra.Command{
	Use:     "rewrite",
	Short:   "Rewrite all variables",
	Long:    `Rewrite all variables`,
	Aliases: []string{"rw"},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := glvarsapi.NewVars(os.Getenv("PRIVATE_TOKEN"))
		if err != nil {
			fmt.Println(err)
			return
		}

		projectID, _ := cmd.Flags().GetInt64("project")
		filename, _ := cmd.Flags().GetString("filename")
		uc := usecase.NewUseCase(projectID, client)

		err = uc.ForceLoadVariablesFromFile(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var Filename string

func init() {
	createCmd.PersistentFlags().Int64VarP(&ProjectId, "project", "p", 0, "Project Id")
	createCmd.PersistentFlags().StringVarP(&Filename, "filename", "f", "", "Path to file with variables")
	rootCmd.AddCommand(createCmd)
}
