package cli

import (
	"fmt"
	"github.com/erminson/gitlab-vars/internal/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "",
	Long:  "",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		fmt.Println(key)
		var envScope string
		if len(args) > 1 {
			envScope = args[1]
		}
		projectId := viper.GetInt64("project-id")
		uc := usecase.NewUseCase(projectId, client)
		err := uc.DeleteVariable(key, envScope)
		if err != nil {
			fmt.Println(fmt.Errorf("error: %v", err))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
