package cli

import (
	"fmt"
	glvarsapi "github.com/erminson/gitlab-vars/internal/client"
	"github.com/spf13/cobra"
	"os"
)

var version = "0.0.1"

var client *glvarsapi.VarsAPI

var ProjectId int64
var Filename string

var rootCmd = &cobra.Command{
	Use:     "glvars",
	Short:   "glvars CI/CD variables (Gitlab API)",
	Long:    `glvars CLI application for working with CI/CD variables (Gitlab API)`,
	Version: version,
}

func Execute() {
	client = Client()
	rootCmd.PersistentFlags().Int64VarP(&ProjectId, "project", "p", 0, "Project Id")
	rootCmd.PersistentFlags().StringVarP(&Filename, "filename", "f", "", "Path to file with variables")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Client() *glvarsapi.VarsAPI {
	client, err := glvarsapi.NewVars(os.Getenv("PRIVATE_TOKEN"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil
	}

	return client
}
