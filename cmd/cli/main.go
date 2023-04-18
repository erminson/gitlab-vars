package main

import (
	"fmt"
	glvarsapi "github.com/erminson/gitlab-vars/internal/client"
	usecase "github.com/erminson/gitlab-vars/internal/usecase"
	"os"
)

func main() {
	client, err := glvarsapi.NewVars(os.Getenv("PRIVATE_TOKEN"))
	if err != nil {
		fmt.Println(err)
		return
	}

	//params := types.Params{
	//	ProjectId: 45202690,
	//}
	//vars, err := client.GetVariables(params)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(vars)

	var projectId int64 = 45202690
	uc := usecase.NewUseCase(projectId, client)

	//str, err := uc.PrintVariablesToFile()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(str)

	err = uc.ForceLoadVariablesFromFile("/Users/erminson/vars.json")
	if err != nil {
		fmt.Println(err)
		return
	}
}
