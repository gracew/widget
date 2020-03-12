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

func DeployAPI(apiName string, deployID string, auth model.Auth, customLogic []*model.CustomLogic) error {
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
		"API_NAME=" + apiName,
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
		return errors.Wrapf(err, "failed to launch docker container: %s", string(out))
	}

	// TODO(gracew): clean up the tmp files eventually...
	return nil
}

func DeployCustomLogic(deployID string, customLogic []*model.CustomLogic) error {
	customLogicDir, err := ioutil.TempDir(TMP_DIR, "customLogic-")
	if err != nil {
		return errors.Wrap(err, "failed to create customLogic temp dir")
	}

	// for now we just pick the first language
	ext, err := getExtension(customLogic[0].Language)
	if err != nil {
		return errors.Wrap(err, "could not determine file extension")
	}

	for _, logic := range customLogic {
		if logic.Before != nil {
			writeFileInDir(customLogicDir, "beforeCreate" + ext, *logic.Before)
		}
		if logic.After != nil {
			writeFileInDir(customLogicDir, "afterCreate" + ext, *logic.After)
		}
	}

	image, err := getImage(customLogic[0].Language)
	if err != nil {
		return errors.Wrap(err, "could not determine image")
	}

	cmd := exec.Command("docker",
		"run",
		"-d",
		"-v",
		customLogicDir + ":/app/customLogic",
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

func getExtension(language model.Language) (string, error) {
	// for now we just pick the first language
	switch language {
		case model.LanguageJavascript:
			return ".js", nil
		case model.LanguagePython:
			return ".py", nil
		default:
			return "", errors.New("unknown custom logic language: " + language.String())
	}
}

func getImage(language model.Language) (string, error) {
	// for now we just pick the first language
	switch language {
		case model.LanguageJavascript:
			return "node-runner", nil
		case model.LanguagePython:
			return "python-runner", nil
		default:
			return "", errors.New("unknown custom logic language: " + language.String())
	}
}

func writeTmpFile(input interface{}, prefix string) (string, error) {
	file, err := ioutil.TempFile(TMP_DIR, prefix)
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp file")
	}

	err = json.NewEncoder(file).Encode(input)
	if err != nil {
		return "", errors.Wrap(err, "failed to encode object to file")
	}
	return filepath.Abs(file.Name())
}

func writeFileInDir(dir string, name string, input string) error {
	path := filepath.Join(dir, name)
	err := ioutil.WriteFile(path, []byte(input), 0644)
	if err != nil {
		return err
	}
	return nil
}