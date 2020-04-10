package launch

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"github.com/gracew/widget/graph/model"
	"github.com/pkg/errors"
)

func (l Launcher) DeployCustomLogic() error {
	l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepLaunchCustomLogicContainer, model.DeployStatusInProgress)
	var err error
	// TODO(gracew): fix this
	if l.CustomLogic.ImageURL != "" {
		err = l.deployCustomLogicImage()
	} else {
		err = l.deployCustomLogic()
	}
	if err != nil {
		l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepLaunchCustomLogicContainer, model.DeployStatusFailed)
		return err
	}
	l.Store.SaveDeployStepStatus(l.DeployID, model.DeployStepLaunchCustomLogicContainer, model.DeployStatusComplete)
	return nil
}

func (l Launcher) deployCustomLogicImage() error {
	cmd := exec.Command("docker",
		"run",
		"-d",
		"--name",
		"custom-logic",
		"--network",
		"widget-proxy_default",
		l.API.ID+"-custom",
	)

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to launch custom logic container")
	}

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
