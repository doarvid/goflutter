package flutter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/akavel/rsrc/rsrc"
	"github.com/doarvid/goflutter/pkg/binlookup"
	"github.com/spf13/viper"

	"log"
)

var (
	flutterBinLookup = binlookup.BinLookup{
		Name:                "flutter",
		InstallInstructions: "Please install flutter,https://flutter.dev/docs/get-started/install",
	}
)

type FlutterProject struct {
	path   string
	config *viper.Viper
}

func NewFlutterProject(path string) (*FlutterProject, error) {
	proj := &FlutterProject{path: path}
	if !proj.IsFlutterProject() {
		return nil, errors.New(path + " is not flutter project directory")
	}
	if err := proj.init(); err != nil {
		return nil, errors.New("project init error")
	}
	return proj, nil
}

func (f *FlutterProject) isProjectFileExists(file string) bool {
	fullpath := path.Join(f.path, file)
	_, err := os.Stat(fullpath)
	if err != nil {
		log.Printf("%s stat error:%s", fullpath, err.Error())
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

func (f *FlutterProject) init() error {
	config := path.Join(f.path, "pubspec.yaml")
	f.config = viper.New()
	f.config.SetConfigFile(config)
	return f.config.ReadInConfig()
}

func (f *FlutterProject) UpdateFile(rootdir string, file string) error {
	srcfile := path.Join(rootdir, file)
	dstfile := path.Join(f.path, file)
	log.Printf("%s => %s", srcfile, dstfile)
	if _, err := os.Stat(dstfile); err == nil {
		os.Remove(dstfile)
	}
	dirpath := path.Dir(dstfile)
	os.MkdirAll(dirpath, 0777)
	data, err := ioutil.ReadFile(srcfile)
	if err != nil {
		return fmt.Errorf("source file %s is missing,err:%s", srcfile, err.Error())
	}

	ctx := string(data)
	newctx := strings.ReplaceAll(ctx, fmt.Sprintf("{{%s}}", "projectname"), f.Name())

	if err := ioutil.WriteFile(dstfile, []byte(newctx), 0777); err != nil {
		return fmt.Errorf("dst file %s write err:%s", srcfile, err.Error())
	}
	return nil
}

func (f *FlutterProject) Name() string {
	return f.config.GetString("name")
}
func (f *FlutterProject) Build() error {
	cmd := exec.Command("flutter", "build", "windows")
	cmd.Dir = f.path
	d, err := cmd.Output()
	if err != nil {
		return err
	}
	println(string(d))
	return nil
}
func (f *FlutterProject) BuildGoApp(gui bool) error {
	rsrcfile := path.Join(f.path, "goflutter.syso")
	err := rsrc.Embed(rsrcfile, "amd64", path.Join(f.path, "goflutter.manifest"), "")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	dstfile := path.Join(f.path, fmt.Sprintf("/build/windows/runner/Release/%s.exe", f.Name()))
	cmdargs := []string{"go", "build", "-ldflags", "-H windowsgui", "-o", dstfile}
	if !gui {
		cmdargs = []string{"go", "build", "-o", dstfile}
	}
	cmd := exec.Command(cmdargs[0], cmdargs[1:]...)
	cmd.Dir = f.path
	d, err := cmd.Output()
	if err != nil {
		return err
	}

	builder := strings.Builder{}
	for _, c := range cmdargs {
		if strings.Contains(c, " ") {
			builder.WriteString("\"")
			builder.WriteString(c)
			builder.WriteString("\"")
		} else {
			builder.WriteString(c)
		}
		builder.WriteString(" ")
	}
	log.Println(builder.String())
	if len(d) > 0 {
		println(string(d))
	}
	return nil
}

func (f *FlutterProject) UpdateFileSym(filename string, symbol string, val string) error {
	fp := path.Join(f.path, filename)
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	log.Printf("update %s symbol %s value %s", fp, symbol, val)
	ctx := string(data)
	newctx := strings.ReplaceAll(ctx, fmt.Sprintf("{{%s}}", symbol), val)
	err = ioutil.WriteFile(fp, []byte(newctx), 0777)
	if err != nil {
		return err
	}
	return nil
}
