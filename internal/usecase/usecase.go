package usecase

import (
	"encoding/json"
	"fmt"
	glvarsapi "github.com/erminson/gitlab-vars/internal/client"
	"github.com/erminson/gitlab-vars/internal/types"
	"io"
	"os"
	"sync"
)

type UseCase struct {
	projectId int64
	client    *glvarsapi.VarsAPI
}

func NewUseCase(projectId int64, client *glvarsapi.VarsAPI) *UseCase {
	return &UseCase{
		projectId: projectId,
		client:    client,
	}
}

func (u *UseCase) SaveVariablesToFile(path string) error {
	params := types.Params{ProjectId: u.projectId}
	vars, err := u.client.GetVariables(params)

	cwd, _ := os.Getwd()
	fmt.Println(cwd)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(vars, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCase) ListVariables() (string, error) {
	params := types.Params{ProjectId: u.projectId}
	vars, err := u.client.GetVariables(params)

	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(vars, "", " ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (u *UseCase) ReWriteVariablesFromFile(filename string) error {
	newVars, err := loadVariablesFromFile(filename)
	if err != nil {
		return err
	}

	fmt.Println(newVars)

	params := types.Params{ProjectId: u.projectId}
	oldVars, err := u.client.GetVariables(params)
	if err != nil {
		return fmt.Errorf("getting error. %v. error: %v", params.String(), err)
	}

	wg := sync.WaitGroup{}
	for _, v := range oldVars {
		wg.Add(1)
		go func(v types.Variable) {
			defer wg.Done()
			params := types.Params{
				ProjectId: u.projectId,
				Key:       v.Key,
			}
			filter := types.Filter{types.FilterEnvScope: v.EnvironmentScope}
			err = u.client.DeleteVariable(params, filter)
			if err != nil {
				fmt.Println(fmt.Errorf("deleting error. %s, %s, error: %v", params.String(), filter.String(), err))
			}
		}(v)
	}
	wg.Wait()

	wg = sync.WaitGroup{}
	for _, v := range newVars {
		wg.Add(1)
		go func(v types.Variable) {
			defer wg.Done()
			params := types.Params{
				ProjectId: u.projectId,
			}

			_, err = u.client.CreateVariable(params, v)
			if err != nil {
				fmt.Println(fmt.Errorf("creating error. %s, %s, error: %v", params.String(), v.String(), err))
			}
		}(v)
	}
	wg.Wait()

	return nil
}

func (u *UseCase) ImportVariablesFromFile(filename string) error {
	newVars, err := loadVariablesFromFile(filename)
	if err != nil {
		return err
	}

	err = validateVariables(newVars)
	if err != nil {
		return err
	}

	params := types.Params{ProjectId: u.projectId}
	exportedVars, err := u.client.GetVariables(params)
	if err != nil {
		return err
	}

	var updateVars []types.Variable
	var createVars []types.Variable
	for _, v := range newVars {
		if containsVariableInSlice(v, exportedVars) {
			updateVars = append(updateVars, v)
		} else {
			createVars = append(createVars, v)
		}
	}

	wg := sync.WaitGroup{}
	for _, v := range updateVars {
		wg.Add(1)
		go func(v types.Variable) {
			defer wg.Done()

			params = types.Params{
				ProjectId: u.projectId,
				Key:       v.Key,
			}
			filter := types.Filter{
				types.FilterEnvScope: v.EnvironmentScope,
			}
			_, err := u.client.UpdateVariable(params, v, filter)
			if err != nil {
				fmt.Println(fmt.Errorf("updating error. %s, %s, %s, error: %v",
					params.String(), v.String(), filter.String(), err))
			}

		}(v)
	}
	wg.Wait()

	wg = sync.WaitGroup{}
	for _, v := range createVars {
		wg.Add(1)
		go func(v types.Variable) {
			defer wg.Done()
			params = types.Params{
				ProjectId: u.projectId,
			}
			_, err := u.client.CreateVariable(params, v)
			if err != nil {
				fmt.Println(fmt.Errorf("creating error. %s, %s, error: %v", params.String(), v.String(), err))
			}
		}(v)
	}
	wg.Wait()

	return nil
}

func (u *UseCase) AddVariable(newVar types.Variable) (types.Variable, error) {
	err := newVar.Validate()
	if err != nil {
		return types.Variable{}, err
	}

	params := types.Params{ProjectId: u.projectId}
	return u.client.CreateVariable(params, newVar)
}

func (u *UseCase) DeleteVariable(key, envScope string) error {
	params := types.Params{
		ProjectId: u.projectId,
		Key:       key,
	}

	filter := types.Filter{
		types.FilterEnvScope: envScope,
	}
	err := u.client.DeleteVariable(params, filter)
	if err != nil {
		return err
	}

	return nil
}

func loadVariablesFromFile(filename string) ([]types.Variable, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var variables []types.Variable
	err = json.Unmarshal(data, &variables)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

func containsVariableInSlice(in types.Variable, vars []types.Variable) bool {
	for _, v := range vars {
		if v.Key == in.Key && v.EnvironmentScope == in.EnvironmentScope {
			return true
		}
	}

	return false
}

func validateVariables(vars []types.Variable) error {
	for _, v := range vars {
		err := v.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
