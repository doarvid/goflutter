package flutter

import (
	"errors"
	"fmt"
	"goflutter/pkg/binlookup"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/rs/zerolog/log"
)

var (
	flutterBinLookup = binlookup.BinLookup{
		Name:                "flutter",
		InstallInstructions: "Please install flutter,https://flutter.dev/docs/get-started/install",
	}
)

type FlutterProject struct {
	path string
}

func NewFlutterProject(path string) (*FlutterProject, error) {
	proj := &FlutterProject{path: path}
	if !proj.IsFlutterProject() {
		return nil, errors.New(path + " is not flutter project directory")
	}
	return proj, nil
}

func (f *FlutterProject) isProjectFileExists(file string) bool {
	fullpath := path.Join(f.path, file)
	_, err := os.Stat(fullpath)
	if err != nil {
		log.Warn().Msgf("%s stat error:%s", fullpath, err.Error())
	}
	return err == nil
}
func (f *FlutterProject) IsFlutterProject() bool {
	if !f.isProjectFileExists("pubspec.yaml") {
		return false
	}
	if !f.isProjectFileExists("windows") {
		return false
	}
	if !f.isProjectFileExists("lib") {
		return false
	}

	return true
}

func (f *FlutterProject) UpdateFile(rootdir string, file string) error {
	srcfile := path.Join(rootdir, file)
	dstfile := path.Join(f.path, file)
	log.Info().Msgf("%s => %s", srcfile, dstfile)
	if _, err := os.Stat(dstfile); err == nil {
		os.Remove(dstfile)
	}
	data, err := ioutil.ReadFile(srcfile)
	if err != nil {
		return fmt.Errorf("source file %s is missing,err:%s", srcfile, err.Error())
	}

	if err := ioutil.WriteFile(dstfile, data, 0777); err != nil {
		return fmt.Errorf("dst file %s write err:%s", srcfile, err.Error())
	}
	return nil
}

func (f *FlutterProject) Build() error {
	cmd := exec.Command("flutter", "build", "windows")
	d, err := cmd.Output()
	if err != nil {
		return err
	}
	println(string(d))
	return nil
}
