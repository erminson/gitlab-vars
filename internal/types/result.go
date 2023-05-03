package types

type Result struct {
	Desc         string
	AcceptedVars []Variable
	FailedVars   map[Variable]error
}

type VariableError struct {
	Var   Variable
	Error error
}
