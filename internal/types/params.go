package types

import (
	"errors"
	"fmt"
)

var (
	ErrParamsInvalidProjectId = errors.New("invalid params project id")
	ErrParamsInvalidKey       = errors.New("invalid params key")
)

type Params struct {
	ProjectId int64
	Key       string
}

func (p *Params) String() string {
	if p.Key == "" {
		return fmt.Sprintf("Params { ProjectId:%v }", p.ProjectId)
	}

	return fmt.Sprintf("Params { ProjectId:%v, Key:%v }", p.ProjectId, p.Key)
}

func (p *Params) Validate() error {
	if p.ProjectId <= 0 {
		return ErrParamsInvalidProjectId
	}

	if p.Key == "" {
		return ErrParamsInvalidKey
	}

	return nil
}

func (p *Params) ValidateProjectId() error {
	if p.ProjectId <= 0 {
		return ErrParamsInvalidProjectId
	}

	return nil
}
