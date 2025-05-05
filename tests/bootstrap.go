package tests

import (
	"os"
	"path/filepath"
	"personal-secretary-user-ap/internal/application"
	"runtime"
	"testing"
)

func TestMain(m *testing.M) {
	println("=== INIT BOOTSTRAP TEST MAIN ===")

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get current file info")
	}

	envFile := filepath.Join(filepath.Dir(currentFile), "../.env.test")

	_, err := application.GetInstance(&envFile)
	if nil != err {
		panic(err)
	}

	code := m.Run()

	println("=== TEARDOWN BOOTSTRAP TEST MAIN ===")

	os.Exit(code)
}
