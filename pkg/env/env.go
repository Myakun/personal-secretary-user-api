package env

import (
	"fmt"
)

const (
	Dev uint8 = iota + 1
	Local
	Prod
	Stage
	Test
)

type Env struct {
	value uint8
}

func (e *Env) IsDev() bool {
	return e.value == Dev
}

func (e *Env) IsLocal() bool {
	return e.value == Local
}

func (e *Env) IsProd() bool {
	return e.value == Prod
}

func (e *Env) IsStage() bool {
	return e.value == Stage
}

func (e *Env) IsTest() bool {
	return e.value == Test
}

func FromString(envStr string) (*Env, error) {
	switch envStr {
	case "dev":
		return &Env{Dev}, nil
	case "local":
		return &Env{Local}, nil
	case "prod":
		return &Env{Prod}, nil
	case "stage":
		return &Env{Stage}, nil
	case "test":
		return &Env{Test}, nil
	}

	return nil, fmt.Errorf("unknown Env: %s", envStr)
}
