package update

import (
	"bytes"
	"errors"
	"goflutter/internal/flutter"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func getGoModPath() (string, error) {
	cmd := exec.Command("go", "env", "GOMODCACHE")
	buf := bytes.Buffer{}
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	modpath := buf.String()

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
	for _, d := range dirs {
		println(d.Name())
	}
	return "", nil
}

func detectAssetPath() (string, error) {
	log.Info().Msgf("detect asset path...")
	assetpath, err := filepath.Abs("./asset")
	if err != nil {
		println(err.Error())
	}
	if _, err := os.Stat(assetpath); err != nil {
		log.Warn().Msgf("asset path %s not exist,error:%s", assetpath, err.Error())
	} else {
		return assetpath, nil
	}

	assetpath, err = getModAssetPath()
	if err != nil {
		log.Warn().Msgf("asset path %s not exist,error:%s", assetpath, err.Error())
	} else {
		return assetpath, nil
	}
	return "", errors.New("asset is missing")
}

func Update(projpath string) error {

	assetpath, err := detectAssetPath()
	if err != nil {
		log.Error().Msg("asset path is missing,maybe reinstall this project")
		return err
	}

	log.Info().Msgf("asset path:%s", assetpath)
	proj, err := flutter.NewFlutterProject(projpath)
	if err != nil {
		log.Error().Msgf("project %s is not valid flutter project", projpath)
		return err
	}
	log.Info().Msgf("update project ...")
	filepath.Walk(assetpath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(path, assetpath) {
			subpath := strings.ReplaceAll(path, assetpath, "")
			if strings.HasPrefix(subpath, string(filepath.Separator)) {
				proj.UpdateFile(assetpath, subpath[1:])
			}
		}
		return nil
	})

	log.Info().Msgf("build flutter project ...")
	if err := proj.Build(); err != nil {
		log.Error().Msgf("project build error %s", err.Error())
		return err
	}
	log.Info().Msgf("build flutter project finished!!")
	return nil
}
