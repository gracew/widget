package launch

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/gracew/widget/graph/model"
	"github.com/pkg/errors"
)

// TODO(gracew): change this lol
const TMP_DIR = "/Users/gracew/tmp"

func DeployAPI(api model.API, auth model.Auth, customLogic []*model.CustomLogic) error {
	generated, err := generateCode(api)
	if err != nil {
		return errors.Wrapf(err, "failed to generate model code for API %s", api.ID)
	}

	err = buildImage(api.ID, generated)
	if err != nil {
		return errors.Wrapf(err, "failed to build docker image for API %s", api.ID)
	}

	err = launchContainer(api, auth, customLogic)
	if err != nil {
		return errors.Wrapf(err, "failed to launch docker container for API %s", api.ID)
	}
	return err
}

func generateCode(api model.API) (string, error) {
	f := jen.NewFile("generated")
	fields := []jen.Code{
		jen.Id("ID").String().Tag(map[string]string{"json": "id", "sql": "type:uuid,default:gen_random_uuid()"}),
		jen.Id("CreatedBy").String().Tag(map[string]string{"json": "createdBy"}),
	}
	for _, field := range api.Definition.Fields {
		jenField := jen.Id(strings.Title(field.Name))
		switch field.Type {
			case model.TypeBoolean:
				jenField.Bool()
			case model.TypeFloat:
				jenField.Float64()
			case model.TypeInt:
				jenField.Int32()
			case model.TypeString:
				jenField.String()
			default:
				return "", errors.New("unknown field type: " + field.Type.String())
		}
		jenField.Tag(map[string]string{"json": field.Name})
		fields = append(fields, jenField)
	}
	f.Type().Id("Object").Struct(fields...)
	return fmt.Sprintf("%#v", f), nil
}

func buildImage(apiID string, generated string) error {
	tmpDir, err := ioutil.TempDir(TMP_DIR, "build-image-")
	if err != nil {
		return errors.Wrap(err, "failed to create temp dir")
	}

	err = ioutil.WriteFile(path.Join(tmpDir, "model.go"), []byte(generated), 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write generated code to temp dir")
	}

	err = copyFile("./launch/Dockerfile", path.Join(tmpDir, "Dockerfile"))
	if err != nil {
		return errors.Wrap(err, "failed to copy dockerfile")
	}

	cmd := exec.Command(
		"docker",
		"build",
		tmpDir,
		"-t",
		apiID,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed to build docker image for api %s: %s", apiID, string(out))
	}

	return nil
}

func launchContainer(api model.API, auth model.Auth, customLogic []*model.CustomLogic) error {
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
		"API_NAME=" + api.Name,
		"--name",
		"widget-proxy",
		"--network",
		"widget-proxy_default",
		api.ID, // image name
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

func copyFile(srcPath string, destPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return errors.Wrap(err, "failed to open source file")
	}
	defer src.Close()

	dest, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open dest file")
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return errors.Wrap(err, "failed to copy file")
	}
	return nil
  }