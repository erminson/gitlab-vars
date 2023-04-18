package main

import (
	"fmt"
	glvarsapi "github.com/erminson/gitlab-vars/internal/client"
	"github.com/erminson/gitlab-vars/internal/types"
	"os"
)

func main() {
	client, err := glvarsapi.NewVars(os.Getenv("PRIVATE_TOKEN"))
	if err != nil {
		fmt.Println(err)
		return
	}

	params := types.Params{
		ProjectId: 45202690,
	}
	vars, err := client.GetVariables(params)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(vars)
}
