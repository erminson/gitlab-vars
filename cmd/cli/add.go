package cli

import (
	"encoding/json"
	"fmt"
	"github.com/erminson/gitlab-vars/internal/types"
	"github.com/erminson/gitlab-vars/internal/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		jsonStr := string(data)
		jsonStr = strings.ReplaceAll(jsonStr, "\n", "")
		jsonStr = strings.ReplaceAll(jsonStr, " ", "")
		fmt.Println(jsonStr)

		if len(jsonStr) > 0 {
			var newVar types.Variable
			err = json.Unmarshal([]byte(jsonStr), &newVar)
			if err != nil {
				fmt.Println(fmt.Errorf("not correct json. error: %v", err))
				os.Exit(1)
			}

			projectId := viper.GetInt64("project-id")
			uc := usecase.NewUseCase(projectId, client)

			addedVar, err := uc.AddVariable(newVar)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(addedVar)
		}
	},
}

func init() {
	//rootCmd.AddCommand(addCmd)
}
