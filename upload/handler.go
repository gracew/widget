package upload

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gracew/widget/graph/model"
	"github.com/gracew/widget/store"
	"github.com/pkg/errors"
)

type Handler struct {
	Store store.Store
}

func (h Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.TempFile(os.TempDir(), "tar-")
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()
	io.Copy(file, r.Body)

	vars := mux.Vars(r)
	go h.untarAndBuild(vars["id"], file)
}

func (h Handler) untarAndBuild(apiID string, file *os.File) {
	dest, err := ioutil.TempDir(os.TempDir(), "tar-")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	file.Seek(0, 0)
	err = Untar(dest, file)
	if err != nil {
		panic(err)
	}

	err = h.BuildCustomLogicImage(apiID, dest)
	if err != nil {
		panic(err)
	}
}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
// adapted from https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func Untar(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return errors.Wrap(err, "failed to create new gzip reader")
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func (h Handler) BuildCustomLogicImage(apiID string, tmpDir string) error {
	err := copyFile(fmt.Sprintf("./upload/%s.dockerfile", "node-runner"), path.Join(tmpDir, "Dockerfile"))
	if err != nil {
		return errors.Wrap(err, "failed to copy dockerfile")
	}

	imageName := apiID + "-custom"
	cmd := exec.Command(
		"docker",
		"build",
		tmpDir,
		"-t",
		apiID+"-custom",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed to build custom logic docker image for api %s: %s", apiID, string(out))
	}

	err = h.saveCustomLogic(apiID, imageName, tmpDir)
	if err != nil {
		return errors.Wrap(err, "failed to save custom logic image name in database")
	}

	return nil
}

func (h Handler) saveCustomLogic(apiID string, imageName string, tmpDir string) error {
	infos, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return errors.Wrap(err, "failed to list dir "+tmpDir)
	}

	customLogic := model.AllCustomLogic{APIID: apiID, ImageURL: imageName, Update: make(map[string]*model.CustomLogic)}
	// TODO(gracew): improve data model
	placeholder := ""
	for _, info := range infos {
		if isBeforeCreate(info.Name()) {
			customLogic.Create = &model.CustomLogic{
				Before: &placeholder,
			}
		} else if isAfterCreate(info.Name()) {
			customLogic.Create = &model.CustomLogic{
				After: &placeholder,
			}
		} else if isBeforeDelete(info.Name()) {
			customLogic.Delete = &model.CustomLogic{
				Before: &placeholder,
			}
		} else if isAfterDelete(info.Name()) {
			customLogic.Delete = &model.CustomLogic{
				After: &placeholder,
			}
		} else if strings.HasPrefix(info.Name(), "before") {
			i := strings.Index(info.Name(), ".")
			if i >= 0 {
				action := info.Name()[len("before"):i]
				customLogic.Update[action] = &model.CustomLogic{
					Before: &placeholder,
				}
			}
		} else if strings.HasPrefix(info.Name(), "after") {
			i := strings.Index(info.Name(), ".")
			if i >= 0 {
				action := info.Name()[len("after"):i]
				customLogic.Update[action] = &model.CustomLogic{
					Before: &placeholder,
				}
			}
		}
	}
	err = h.Store.SaveCustomLogic(&customLogic)
	if err != nil {
		return errors.Wrap(err, "failed to save custom logic")
	}
	return nil
}

func isBeforeCreate(filename string) bool {
	return strings.HasPrefix(filename, "beforecreate.")
}

func isAfterCreate(filename string) bool {
	return strings.HasPrefix(filename, "aftercreate.")
}

func isBeforeDelete(filename string) bool {
	return strings.HasPrefix(filename, "beforedelete.")
}

func isAfterDelete(filename string) bool {
	return strings.HasPrefix(filename, "afterdelete.")
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
