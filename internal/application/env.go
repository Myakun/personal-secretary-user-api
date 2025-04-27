package application

import (
	"fmt"
)

const (
	Dev uint8 = iota
	Local
	Prod
	Stage
	Test
)

type env struct {
	value uint8
}

func (e *env) GetValue() uint8 {
	return e.value
}

//goland:noinspection GoExportedFuncWithUnexportedType
func EnvFromString(envStr string) (*env, error) {
	switch envStr {
	case "dev":
		return &env{Dev}, nil
	case "local":
		return &env{Local}, nil
	case "prod":
		return &env{Prod}, nil
	case "stage":
		return &env{Stage}, nil
	case "test":
		return &env{Test}, nil
	}

	return nil, fmt.Errorf("unknown env: %s", envStr)
}
