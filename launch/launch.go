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
	"github.com/gracew/widget/store"
	"github.com/pkg/errors"
)

// TODO(gracew): change this lol
const TMP_DIR = "/Users/gracew/tmp"

type Launcher struct {
	Store       store.Store
	DeployID    string
	API         model.API
	Auth        model.Auth
	CustomLogic *model.AllCustomLogic
}

func (l Launcher) DeployAPI() error {
	l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepGenerateCode, model.DeployStatusInProgress)
	generated, err := l.generateCode()
	if err != nil {
		l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepGenerateCode, model.DeployStatusFailed)
		return errors.Wrapf(err, "failed to generate model code for API %s", l.API.ID)
	}
	l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepGenerateCode, model.DeployStatusComplete)

	l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepBuildImage, model.DeployStatusInProgress)
	err = l.buildImage(generated)
	if err != nil {
		l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepBuildImage, model.DeployStatusFailed)
		return errors.Wrapf(err, "failed to build docker image for API %s", l.API.ID)
	}
	l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepBuildImage, model.DeployStatusComplete)

	l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepLaunchContainer, model.DeployStatusInProgress)
	err = l.launchContainer()
	if err != nil {
		l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepLaunchContainer, model.DeployStatusFailed)
		return errors.Wrapf(err, "failed to launch docker container for API %s", l.API.ID)
	}
	l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepLaunchContainer, model.DeployStatusComplete)

	return nil
}

func (l Launcher) generateCode() (string, error) {
	f := jen.NewFile("generated")
	fields := []jen.Code{
		jen.Id("tableName").Struct().Tag(map[string]string{"sql": strings.ToLower(l.API.Name)}),
		jen.Id("ID").String().Tag(map[string]string{"json": "id", "sql": "type:uuid,default:gen_random_uuid()"}),
		jen.Id("CreatedBy").String().Tag(map[string]string{"json": "createdBy"}),
		jen.Id("CreatedAt").String().Tag(map[string]string{"json": "createdAt", "sql": "default:now()"}),
	}
	for _, field := range l.API.Fields {
		jenField := jen.Id(strings.Title(field.Name))
		tags := map[string]string{"json": field.Name}
		switch field.Type {
		case model.TypeBoolean:
			jenField.Bool()
			tags["sql"] = ",notnull"
		case model.TypeFloat:
			jenField.Float64()
		case model.TypeInt:
			jenField.Int32()
		case model.TypeString:
			jenField.String()
		default:
			return "", errors.New("unknown field type: " + field.Type.String())
		}
		jenField.Tag(tags)
		fields = append(fields, jenField)
	}
	f.Type().Id("Object").Struct(fields...)

	return fmt.Sprintf("%#v", f), nil
}

func (l Launcher) buildImage(generated string) error {
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
		l.API.ID,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed to build docker image for api %s: %s", l.API.ID, string(out))
	}

	return nil
}

func (l Launcher) launchContainer() error {
	// write API definition to temp files
	apiPath, err := writeTmpFile(l.API, "api-")
	if err != nil {
		return errors.Wrap(err, "failed to write api to temp file")
	}

	authPath, err := writeTmpFile(l.Auth, "auth-")
	if err != nil {
		return errors.Wrap(err, "failed to write auth to temp file")
	}

	customLogicPath, err := writeTmpFile(l.CustomLogic, "customLogic-")
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
		apiPath+":/app/api.json",
		"-v",
		authPath+":/app/auth.json",
		"-v",
		customLogicPath+":/app/customLogic.json",
		"-e",
		"API_NAME="+l.API.Name,
		"--name",
		"widget-proxy",
		"--network",
		"widget-proxy_default",
		l.API.ID, // image name
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed to launch docker container: %s", string(out))
	}

	// TODO(gracew): clean up the tmp files eventually...
	return nil
}

func (l Launcher) DeployCustomLogic() error {
	l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepLaunchCustomLogicContainer, model.DeployStatusInProgress)
	err := l.deployCustomLogic()
	if err != nil {
		l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepLaunchCustomLogicContainer, model.DeployStatusFailed)
		return err
	}
	l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepLaunchCustomLogicContainer, model.DeployStatusComplete)
	return nil
}

func (l Launcher) deployCustomLogic() error {
	customLogicDir, err := ioutil.TempDir(TMP_DIR, "customLogic-")
	if err != nil {
		return errors.Wrap(err, "failed to create customLogic temp dir")
	}

	// for now we just pick the first language
	// TODO(gracew): fix this
	language := l.CustomLogic.Create.Language
	ext, err := getExtension(l.CustomLogic.Create.Language)
	if err != nil {
		return errors.Wrap(err, "could not determine file extension")
	}

	writeCustomLogicFiles(customLogicDir, "create", ext, l.CustomLogic.Create)
	writeCustomLogicFiles(customLogicDir, "delete", ext, l.CustomLogic.Delete)
	for actionName, actionCustomLogic := range l.CustomLogic.Update {
		writeCustomLogicFiles(customLogicDir, actionName, ext, actionCustomLogic)
	}

	image, err := getImage(language)
	if err != nil {
		return errors.Wrap(err, "could not determine image")
	}

	cmd := exec.Command("docker",
		"run",
		"-d",
		"-v",
		customLogicDir+":/app/customLogic",
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

func writeCustomLogicFiles(dir string, label string, ext string, customLogic *model.CustomLogic) {
	if customLogic.Before != nil {
		writeFileInDir(dir, "before"+label+ext, *customLogic.Before)
	}
	if customLogic.After != nil {
		writeFileInDir(dir, "after"+label+ext, *customLogic.After)
	}
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
