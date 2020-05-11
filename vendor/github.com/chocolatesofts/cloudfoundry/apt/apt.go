package apt
import (
	"fmt"
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"os"
	
	"github.com/cloudfoundry/libbuildpack"
)

const(
	Aptname = "apt"
)
type Stager interface{
	LinkDirectoryInDepDir(string, string) error
	WriteProfileD(string, string) error
	CacheDir() string
	DepDir() string
	BuildDir() string
}

type Command interface {
	Execute(string, io.Writer, io.Writer, string, ...string) error
	Output(string, string, ...string) (string, error)
}

type Supplier struct{
	Stager         Stager
	Command        Command
	Setup 			bool
	Update 			bool
}

func AptSetup(s *Supplier,installDir string) (error,[]string,[]string,string){
	cacheDir:=s.Stager.CacheDir()
	aptCacheDir:=filepath.Join(cacheDir, "apt", "cache")
	stateDir := filepath.Join(cacheDir, "apt", "state")
	preferences := filepath.Join(cacheDir, "apt", "etc", "preferences")
	archiveDir :=filepath.Join(aptCacheDir, "archives")
	sourcelist:=filepath.Join(cacheDir,"apt","sources","sources.list")
	rootDir := "/etc/apt"
	aptSources:= filepath.Join(rootDir, "sources.list")
	options:=[]string{
			"-o", "debug::nolocking=true",
			"-o", "dir::cache=" + aptCacheDir,
			"-o", "dir::etc::sourcelist=" + sourcelist,
			"-o", "dir::state=" + stateDir,
			"-o", "Dir::Etc::preferences=" + preferences}

	doptions := append(options, "-y", "--allow-downgrades", "--allow-remove-essential", "--allow-change-held-packages", "-d", "install", "--reinstall")

	if (s.Setup == false){
		if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
			return err,nil,nil,""
		}
		if err := os.MkdirAll(aptCacheDir, os.ModePerm); err != nil {
			return err,nil,nil,""
		}
		if err := os.MkdirAll(stateDir, os.ModePerm); err != nil {
			return err,nil,nil,""
		}
		if err := os.MkdirAll(installDir, os.ModePerm); err != nil {
			return err,nil,nil,""
		}
		if err := os.MkdirAll(archiveDir, os.ModePerm); err != nil {
			return err,nil,nil,""
		}
		aptPrefs := filepath.Join(rootDir, "preferences")
		if exists, err := libbuildpack.FileExists(aptPrefs); err != nil {
			return err,nil,nil,""
		} else if exists {
			if err := libbuildpack.CopyFile(aptPrefs, preferences); err != nil {
				return err,nil,nil,""
			}
		} else {
			dirPath := filepath.Dir(preferences)
			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				return err,nil,nil,""
			}
		}
		if err := libbuildpack.CopyFile(aptSources, sourcelist); err != nil {
			return err,nil,nil,""
		}
	}
	s.Setup=true
	return nil,options,doptions,archiveDir
}

func AptUpdate(s *Supplier,options []string) error {
	if(s.Update==true){
		return nil
	}
	uargs := append(options, "update")	
	var errBuff bytes.Buffer
	if err := s.Command.Execute("/", &errBuff, &errBuff, "apt-get", uargs...); err != nil {
		return fmt.Errorf("failed to apt-get update %s\n\n%s", errBuff.String(), err)
	}
	s.Update=true
	return nil;
}

func InstallPackages(s *Supplier,archiveDir string,installDir string) error{
	files, err := filepath.Glob(filepath.Join(archiveDir, "*.deb"))
	if err != nil {
		return err
	}
	for _, file := range files {
		output, err := s.Command.Output("/", "dpkg", "-x", file, installDir)
	if err != nil {
		return fmt.Errorf("failed to install pkg %s\n\n%s\n\n%s", file, output, err.Error())
	}
	}
	return nil;
}

func LinkPackages(s *Supplier,installDir string) error {
	for _, dirs := range [][]string{
		{"usr/bin", "bin"},
		{"usr/lib", "lib"},
		{"usr/lib/i386-linux-gnu", "lib"},
		{"usr/lib/x86_64-linux-gnu", "lib"},
		{"lib/x86_64-linux-gnu", "lib"},
		{"usr/include", "include"},
	} {
		dest := filepath.Join(installDir, dirs[0])
		if exists, err := libbuildpack.FileExists(dest); err != nil {
			return err
		} else if exists {
			if err := s.Stager.LinkDirectoryInDepDir(dest, dirs[1]); err != nil {
				return err
			}
		}
	}
	for _, dirs := range [][]string{
		{"usr/lib/i386-linux-gnu/pkgconfig", "pkgconfig"},
		{"usr/lib/x86_64-linux-gnu/pkgconfig", "pkgconfig"},
		{"usr/lib/pkgconfig", "pkgconfig"},
		} {
			dest := filepath.Join(installDir, dirs[0])
			if exists, err := libbuildpack.FileExists(dest); err != nil {
				return err
			}else if exists {
				files, err := ioutil.ReadDir(dest)
				if err != nil {
					return err
				}
			destDir := filepath.Join(s.Stager.DepDir(), dirs[1])
			if err := os.MkdirAll(destDir, 0755); err != nil {
				return err
			}
			for _, file := range files {
				//TODO: better way to copy a file?
				contents, err := ioutil.ReadFile(filepath.Join(dest, file.Name()))
				if err != nil {
					return err
				}
				newContents := strings.Replace(string(contents[:]), "prefix=/usr\n", "prefix="+filepath.Join(installDir, "usr")+"\n", -1)
				err = ioutil.WriteFile(filepath.Join(destDir, file.Name()), []byte(newContents), 0666)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil;
}

func InstallDir(s *Supplier,pkg string)string{
	installDir := filepath.Join(s.Stager.DepDir(),Aptname,pkg)
	return installDir
}