package launch

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gracew/widget/graph/model"
	"github.com/pkg/errors"
)

// TODO(gracew): change this lol
const TMP_DIR = "/Users/gracew/tmp"

func DeployAPI(auth model.Auth) error {
	// serialize the auth object and write it to a temp file
	bytes, err := json.Marshal(auth)
	if err != nil {
		return errors.Wrap(err, "could not marshal auth object")
	}

	authPath, err := writeTmpFile(bytes, "auth-")
	if err != nil {
		return errors.Wrap(err, "failed to write auth to temp file")
	}
	defer os.Remove(authPath)

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
		"--name",
		"widget-proxy",
		"--network",
		"widget-proxy_default",
		"widget-proxy",
	)
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to execute docker container")
	}

	return nil
}

func writeTmpFile(input []byte, prefix string) (string, error) {
	file, err := ioutil.TempFile(TMP_DIR, prefix)
	if err != nil {
		return "", err
	}

	if _, err = file.Write(input); err != nil {
        return "", err
	}
	return filepath.Abs(file.Name())
}