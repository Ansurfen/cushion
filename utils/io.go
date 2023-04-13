package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

func NewConf(name, _type, dir string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigName(name)
	conf.SetConfigType(_type)
	conf.AddConfigPath(dir)
	if err := conf.ReadInConfig(); err != nil {
		panic(err)
	}
	return conf
}

func NewConfFromPath(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)
	if err := conf.ReadInConfig(); err != nil {
		panic(err)
	}
	return conf
}

func Unzip(zipPath, dstDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if err := unzipFile(file, dstDir); err != nil {
			return err
		}
	}
	return nil
}

func unzipFile(file *zip.File, dstDir string) error {
	filePath := path.Join(dstDir, file.Name)
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, rc)
	return err
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	return nil
}

func FetchFile(src, dst string) (int64, error) {
	file, err := os.OpenFile(dst, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	res, err := http.Get(src)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	n, err := io.Copy(file, res.Body)
	return n, err
}

func Mkdirs(path string) error {
	return os.MkdirAll(path, 0777)
}

func SafeMkdirs(path string) error {
	if ok, err := PathIsExist(path); err != nil {
		return err
	} else if !ok {
		if err := Mkdirs(path); err != nil {
			return err
		}
	}
	return nil
}

func SafeBatchMkdirs(dirs []string) error {
	for _, dir := range dirs {
		if err := SafeMkdirs(dir); err != nil {
			return err
		}
	}
	return nil
}

func Exec(arg ...string) ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		out, err := exec.Command("cmd", append([]string{"/C"}, arg...)...).CombinedOutput()
		if err != nil {
			return out, err
		}
		return out, nil
	case "linux":
		out, err := exec.Command("/bin/bash", append([]string{"/C"}, arg...)...).CombinedOutput()
		if err != nil {
			return out, err
		}
		return out, nil
	case "darwin":
	default:
	}
	return []byte(""), nil
}

func ExecStr(args string) ([]byte, error) {
	return Exec(strings.Fields(args)...)
}

func WriteFile(file string, data []byte) error {
	fp, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fp.Close()
	fp.Write(data)
	return nil
}

func SafeWriteFile(file string, data []byte) error {
	if ok, err := PathIsExist(file); err != nil {
		return err
	} else if !ok {
		if err := WriteFile(file, data); err != nil {
			return err
		}
	}
	return nil
}

func ReadStraemFromFile(file string) ([]byte, error) {
	fp, err := os.Open(file)
	if err != nil {
		return []byte(""), err
	}
	defer fp.Close()
	raw, err := ioutil.ReadAll(fp)
	if err != nil {
		return []byte(""), err
	}
	return raw, nil
}

func ReadLineFromFile(file string, filter func(string) string) ([]byte, error) {
	fp, err := os.Open(file)
	if err != nil {
		return []byte(""), err
	}
	defer fp.Close()
	fileScanner := bufio.NewScanner(fp)
	var ret []byte
	for fileScanner.Scan() {
		ret = append(ret, []byte(filter(fileScanner.Text()))...)
	}
	if err := fileScanner.Err(); err != nil {
		return []byte(""), err
	}
	return ret, nil
}

func ReadLineFromString(str string, filter func(string) string) ([]byte, error) {
	scanner := bufio.NewScanner(strings.NewReader(str))
	scanner.Split(bufio.ScanLines)
	var ret []byte
	for scanner.Scan() {
		ret = append(ret, []byte(filter(scanner.Text()))...)
	}
	if err := scanner.Err(); err != nil {
		return []byte(""), err
	}
	return ret, nil
}

func PathIsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Filename(fullpath string) string {
	filename := path.Base(fullpath)
	ext := path.Ext(filename)
	return filename[:len(filename)-len(ext)]
}
