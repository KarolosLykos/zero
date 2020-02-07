package util

// @TODO split up and move into /pkg directory

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
	"github.com/google/uuid"
)

func CreateDirIfDoesNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		return err
	}
	return nil
}

func CleanGoIdentifier(identifier string) string {
	return strings.ReplaceAll(identifier, "-", "")
}

// @TODO how can we make these type of helpers extensible?
var FuncMap = template.FuncMap{
	"Title":             strings.Title,
	"ToLower":           strings.ToLower,
	"CleanGoIdentifier": CleanGoIdentifier,
	"GenerateUUID":      uuid.New,
}

func GetCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Getting working directory failed: %v\n", err)
		panic(err)
	}

	return dir
}

func ExecuteCommand(cmd *exec.Cmd, pathPrefix string, envars []string) {
	dir := GetCwd()

	cmd.Dir = path.Join(dir, pathPrefix)

	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	var errStdout, errStderr error

	if envars != nil {
		cmd.Env = envars
	}

	err := cmd.Start()
	if err != nil {
		log.Fatalf("Starting command failed: %v\n", err)
	}

	go func() {
		_, errStdout = io.Copy(os.Stdout, stdoutPipe)
	}()
	go func() {
		_, errStderr = io.Copy(os.Stderr, stderrPipe)
	}()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Executing command failed: %v\n", err)
	}

	if errStdout != nil {
		log.Printf("Failed to capture stdout: %v\n", errStdout)
	}

	if errStderr != nil {
		log.Printf("Failed to capture stderr: %v\n", errStderr)
	}
}

func ExecuteCommandOutput(cmd *exec.Cmd, pathPrefix string, envars []string) string {
	dir := GetCwd()

	cmd.Dir = path.Join(dir, pathPrefix)

	if envars != nil {
		cmd.Env = envars
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Executing command failed: (%v) %s\n", err, out)
	}
	return string(out)
}