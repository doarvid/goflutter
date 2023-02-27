package update

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/doarvid/goflutter/internal/flutter"

	"log"
)

func getGoModPath() (string, error) {
	cmd := exec.Command("go", "env", "GOMODCACHE")
	buf := bytes.Buffer{}
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	modpath := buf.String()
	println(modpath)
	modpath = strings.TrimSpace(modpath)
	modpath = strings.ReplaceAll(modpath, "\\\\", "\\")
	return modpath, nil
}

func getModAssetPath() (string, error) {
	modroot, err := getGoModPath()
	if err != nil {
		return "", err
	}

	user_proj_path := path.Join(modroot, "github.com", "doarvid")

	dirs, err := os.ReadDir(user_proj_path)
	if err != nil {
		return "", err
	}
	var projects []string
	for _, d := range dirs {
		if strings.HasPrefix(d.Name(), "goflutter") {
			projects = append(projects, path.Join(user_proj_path, d.Name(), "asset"))
		}
	}
	if len(projects) == 0 {
		return "", errors.New("no go mod")
	}
	sort.Slice(projects, func(i, j int) bool {
		return strings.Compare(projects[i], projects[j]) > 0
	})
	return projects[0], nil
}

func detectAssetPath() (string, error) {
	log.Printf("detect asset path...")
	assetpath, err := filepath.Abs("./asset")
	if err != nil {
		println(err.Error())
	}
	if _, err := os.Stat(assetpath); err != nil {
		log.Printf("asset path %s not exist,error:%s", assetpath, err.Error())
	} else {
		return assetpath, nil
	}

	assetpath, err = getModAssetPath()
	if err != nil {
		log.Printf("asset path %s not exist,error:%s", assetpath, err.Error())
	} else {
		return assetpath, nil
	}
	return "", errors.New("asset is missing")
}

func buildfloders(dir string, top bool) ([]string, error) {
	dir_, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer dir_.Close()
	fis, err := dir_.Readdir(0)
	if err != nil {
		return nil, err
	}

	filenames := []string{}
	for _, fi := range fis {
		subpath := path.Join(dir, fi.Name())
		if fi.IsDir() {
			fns, err := buildfloders(subpath, false)
			if err == nil {
				filenames = append(filenames, fns...)
			}
			continue
		}
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		filenames = append(filenames, fi.Name())
	}
	if top {
		return filenames, nil
	}
	var paths []string
	for _, fn := range filenames {
		paths = append(paths, path.Join(path.Base(dir_.Name()), fn))
	}
	return paths, nil
}

func Update(projpath string) error {

	assetpath, err := detectAssetPath()
	if err != nil {
		log.Printf("asset path is missing,maybe reinstall this project")
		return err
	}

	log.Printf("asset path:%s", assetpath)
	proj, err := flutter.NewFlutterProject(projpath)
	if err != nil {
		log.Printf("project %s is not valid flutter project", projpath)
		return err
	}
	log.Printf("update project ...")

	files, err := buildfloders(assetpath, true)
	if err != nil {
		log.Printf("get asset error :%s", err.Error())
		return err
	}
	for _, f := range files {
		proj.UpdateFile(assetpath, f)
	}
	log.Printf("build flutter project ...")
	if err := proj.Build(); err != nil {
		log.Printf("project build error %s", err.Error())
		return err
	}

	log.Printf("fix template ...")
	proj.UpdateFileSym("flutter/runner.go", "runnerlib", "-l"+proj.Name())
	log.Printf("build go app...")
	proj.BuildGoApp(true)
	log.Printf("finished")
	return nil
}
