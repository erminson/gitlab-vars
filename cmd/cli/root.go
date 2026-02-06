package cli

import (
	"fmt"
	 glvarsapi "github.com/erminson/gitlab-vars/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var version = "0.1.0"

var client *glvarsapi.VarsAPI

var (
	ProjectId int64
	Filename  string
	cfgFile   string
	VarKey 	  string
)

var rootCmd = &cobra.Command{
	Use:     "glvars",
	Short:   "glvars CLI tool for import and export project-level Gitlab CI/CD Variables",
	Long:    `glvars CLI tool for import and export project-level Gitlab CI/CD Variables`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".glvars")
	}

	viper.SetEnvPrefix("glvars")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
    fmt.Println("Error loading config file:", err)
	}

}
func initClient() {
	client = Client()
}

func init() {
	cobra.OnInitialize(initConfig, initClient)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().Int64P("project", "p", 0, "project id")
	err := viper.BindPFlag("project-id", rootCmd.PersistentFlags().Lookup("project"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringVarP(&VarKey, "key", "k", "", "env variable key")

	rootCmd.PersistentFlags().StringVarP(&Filename, "filename", "f", "", "path to file with variables")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.glvars.yaml)")
}

func Client() *glvarsapi.VarsAPI {
	host := viper.GetString("host")
	token := viper.GetString("private-token")

	var client *glvarsapi.VarsAPI
	var err error

	if host != "" {
		client, err = glvarsapi.NewVarsWithHost(token, host)
	} else {
		client, err = glvarsapi.NewVars(token)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil
	}

	return client
}
