package cli

import (
	"fmt"
	glvarsapi "github.com/erminson/gitlab-vars/internal/client"
	"github.com/spf13/cobra"
	"os"
)

var client *glvarsapi.VarsAPI

var rootCmd = &cobra.Command{
	Use:   "glvars",
	Short: "glvars CI/CD variables (Gitlab API)",
	Long:  `glvars CLI application for working with CI/CD variables (Gitlab API)`,
}

func Execute() {
	client = Client()
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
