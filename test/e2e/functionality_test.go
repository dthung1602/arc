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
	inputBytes, _ := os.ReadFile("test/data/" + testFileName + "_commands.txt")
	expectedBytes, _ := os.ReadFile("test/data/" + testFileName + "_expected_outputs.txt")
	input := string(inputBytes)
	expected := string(expectedBytes)

	cmd := exec.Command("redis-cli", "-p", "6378")

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	defer stdout.Close()

	cmd.Start()
	io.WriteString(stdin, input)
	stdin.Close()
	builder := new(strings.Builder)
	io.Copy(builder, stdout)

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
