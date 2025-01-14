package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// runPyCode executes a single Python source code file with the specified input.
// It runs the Python script using the provided file and captures its output.
// 
// Parameters:
//   - codeFile: A file pointer to the Python source code to be executed
//   - input: A string to be used as standard input for the Python script
//
// Returns:
//   - The standard output of the executed Python script as a string
//   - An error if the script execution fails, including execution error details
//
// The function uses the system's Python interpreter to run the script and captures
// both standard output and standard error streams. If an error occurs during
// execution, it returns a detailed error message including the error and stderr content.
func runPyCode(codeFile *os.File, input string) (string, error) {
	cmd := exec.Command("python", codeFile.Name())

	var stdout, stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("运行时错误: %v, %s", err, stderr.String())
	}

	return stdout.String(), nil
}
