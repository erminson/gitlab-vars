package cli

import (
	"encoding/json"
	"fmt"
	"github.com/erminson/gitlab-vars/internal/types"
	"github.com/erminson/gitlab-vars/internal/usecase"
	"github.com/spf13/cobra"
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
				fmt.Printf(" #{err}")
				os.Exit(1)
			}

			projectID, _ := cmd.Flags().GetInt64("project")
			uc := usecase.NewUseCase(projectID, client)

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
	addCmd.PersistentFlags().Int64VarP(&ProjectId, "project", "p", 0, "verbose output")
	rootCmd.AddCommand(addCmd)
}
