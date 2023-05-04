package cli

import (
	"fmt"
	glvarsapi "github.com/erminson/gitlab-vars/internal/client"
	"github.com/erminson/gitlab-vars/internal/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// nolint
var rewriteCmd = &cobra.Command{
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

		projectId := viper.GetInt64("project-id")
		uc := usecase.NewUseCase(projectId, client)

		err = uc.ReWriteVariablesFromFile(Filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	//rootCmd.AddCommand(rewriteCmd)
}
