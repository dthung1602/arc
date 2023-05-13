package e2e

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/dthung1602/arc/pkg/app"
)

func executeAndCompare(t *testing.T, testFileName string) {
	inputBytes, err := os.ReadFile("test/data/" + testFileName + "_commands.txt")
	if err != nil {
		panic(err)
	}
	expectedBytes, _ := os.ReadFile("test/data/" + testFileName + "_expected_outputs.txt")
	if err != nil {
		panic(err)
	}
	input := string(inputBytes)
	expected := string(expectedBytes)

	cmd := exec.Command("redis-cli", "-p", "6378")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	_, err = io.WriteString(stdin, input)
	if err != nil {
		panic(err)
	}

	stdin.Close()
	builder := new(strings.Builder)
	_, err = io.Copy(builder, stdout)
	if err != nil {
		panic(err)
	}

	if expected != builder.String() {
		t.Error("redis-cli doesn't return what expected")
	}
}

func TestFunctionalities(t *testing.T) {
	_, err := exec.LookPath("redis-cli")
	if err != nil {
		t.Error("redis-cli is required to run e2e tests. make sure it's available in PATH")
	}

	arc, _ := app.NewApp()
	go arc.Serve()
	defer arc.Stop()

	executeAndCompare(t, "getset")
}
