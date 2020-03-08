package launch

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"github.com/gracew/widget/graph/model"
	"github.com/pkg/errors"
)

// TODO(gracew): change this lol
const TMP_DIR = "/Users/gracew/tmp"

func DeployAPI(deployID string, auth model.Auth, customLogic []*model.CustomLogic) error {
	// write auth and customLogic objects to temp files
	authPath, err := writeTmpFile(auth, "auth-")
	if err != nil {
		return errors.Wrap(err, "failed to write auth to temp file")
	}

	customLogicPath, err := writeTmpFile(customLogic, "customLogic-")
	if err != nil {
		return errors.Wrap(err, "failed to write customLogic to temp file")
	}

	// launch the docker container
	// TODO(gracew): replace with k8s + a proper client lib
	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"-p",
		// TODO(gracew): this is bad
		"8081:8080",
		"-v",
		authPath + ":/app/auth.json",
		"-v",
		customLogicPath + ":/app/customLogic.json",
		"-e",
		"DEPLOY_ID=" + deployID,
		"--name",
		"widget-proxy",
		"--network",
		"widget-proxy_default",
		"widget-proxy",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		print(string(out))
		return errors.Wrapf(err, "failed to launch docker container: %s", string(out))
	}

	// TODO(gracew): clean up the tmp files eventually...
	return nil
}

func DeployCustomLogic(deployID string, customLogic []*model.CustomLogic) error {
	customLogicPath, err := writeTmpFile(customLogic, "customLogic-")
	if err != nil {
		return errors.Wrap(err, "failed to write customLogic to temp file")
	}


	var image string
	// for now we just pick the first language
	switch customLogic[0].Language {
		case model.LanguageJavascript:
			image = "node-runner"
		case model.LanguagePython:
			image = "python-runner"
		default:
			return errors.New("unknown custom logic language: " + customLogic[0].Language.String())
	}

	cmd := exec.Command("docker",
		"run",
		"-d",
		"-v",
		customLogicPath + ":/app/customLogic.json",
		"--name",
		"custom-logic",
		"--network",
		"widget-proxy_default",
		image,
	)

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to launch custom logic container")
	}

	return nil
}

func writeTmpFile(input interface{}, prefix string) (string, error) {
	file, err := ioutil.TempFile("/tmp", prefix)
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp file")
	}

	err = json.NewEncoder(file).Encode(input)
	if err != nil {
		return "", errors.Wrap(err, "failed to encode object to file")
	}
	return filepath.Abs(file.Name())
}

func writeTmpFileBytes(input []byte, prefix string) (string, error) {
	file, err := ioutil.TempFile(TMP_DIR, prefix)
	if err != nil {
		return "", err
	}

	if _, err = file.Write(input); err != nil {
        return "", err
	}
	return filepath.Abs(file.Name())
}