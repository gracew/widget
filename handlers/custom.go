package handlers

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
)

// TODO(gracew): change this lol
const TMP_DIR = "/Users/gracew/tmp"

func Execute(input []byte, customLogic string) ([]byte, error) {
	// write the input and customLogic to temp files
	inputPath, err := writeTmpFile(input)
	if err != nil {
		log.Print("failed to write input", err);
		return nil, err
	}

	customLogicPath, err := writeTmpFile([]byte(customLogic))
	if err != nil {
		log.Print("failed to write custom logic", err);
		return nil, err
	}

	outputDirPath, err := ioutil.TempDir(TMP_DIR, "output-")
	if err != nil {
		log.Print("failed to create output dir", err);
		return nil, err
	}

	// run the docker container
	cmd := exec.Command(
		"docker",
		"run",
		"-v",
		customLogicPath + ":/app/customLogic.js",
		"-v",
		inputPath + ":/app/input.json",
		"-v",
		outputDirPath + ":/app/output",
		"node-runner",
	)
	err = cmd.Run()
	if err != nil {
		log.Print("failed to execute docker container", err);
		return nil, err
	}

	// read the output
	output, err := ioutil.ReadFile(filepath.Join(outputDirPath, "output.json"))
	if err != nil {
		log.Print("failed to read output");
		return nil, err
	}
	// TODO(gracew): close all the files...
	return output, nil
}

func writeTmpFile(input []byte) (string, error) {
	file, err := ioutil.TempFile(TMP_DIR, "input-")
	if err != nil {
		return "", err
	}

	if _, err = file.Write(input); err != nil {
        return "", err
	}
	return filepath.Abs(file.Name())
}