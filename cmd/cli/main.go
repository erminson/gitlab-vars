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

	var projectId int64 = 45202690
	uc := usecase.NewUseCase(projectId, client)
	
	err = uc.ForceLoadVariablesFromFile("/Users/erminson/vars.json")
	if err != nil {
		fmt.Println(err)
		return
	}
}
