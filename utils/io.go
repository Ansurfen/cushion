package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// OpenConf unmarshal file which located in disk to memory according to name, type, dir
func OpenConf(name, _type, dir string, opts ...viper.Option) *viper.Viper {
	conf := viper.NewWithOptions(opts...)
	conf.SetConfigName(name)
	conf.SetConfigType(_type)
	conf.AddConfigPath(dir)
	return conf
}

// OpenConfFromPath unmarshal file which located in disk to memory according to path
func OpenConfFromPath(path string, opts ...viper.Option) *viper.Viper {
	conf := viper.NewWithOptions(opts...)
	conf.SetConfigFile(path)
	return conf
}

// Zip to compress zip of source to specify path
func Zip(zipPath string, paths ...string) error {
	if err := os.MkdirAll(filepath.Dir(zipPath), os.ModePerm); err != nil {
		return err
	}
	archive, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()
	// traverse the file or directory
	for _, srcPath := range paths {
		// remove the trailing path separator if path is a directory
		srcPath = strings.TrimSuffix(srcPath, string(os.PathSeparator))

		// visit all the files or directories in the tree
		err = filepath.Walk(srcPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// create a local file header
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// set compression
			header.Method = zip.Deflate

			// set relative path of a file as the header name
			header.Name, err = filepath.Rel(filepath.Dir(srcPath), path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				header.Name += string(os.PathSeparator)
			}

			// create writer for the file header and save content of the file
			headerWriter, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(headerWriter, f)
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Unzip unzip zip of source to specify path
func Unzip(src, dst string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if err := unzipFile(file, dst); err != nil {
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

// MoveFile move file from src to dst, like mv or move command
func MoveFile(src, dst string) error {
	inputFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(dst)
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

// FetchFile fetch file from remote source to local destination
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

// Mkdirs recurse to create path
func Mkdirs(path string) error {
	return os.MkdirAll(path, 0777)
}

// SafeMkdirs recurse to create path when path isn't exist
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

// SafeBatchMkdirs recurse to create dirs when path isn't exist
func SafeBatchMkdirs(dirs []string) error {
	for _, dir := range dirs {
		if err := SafeMkdirs(dir); err != nil {
			return err
		}
	}
	return nil
}

// Exec automatically fit in os enviroment to execute command.
// windows 10+ -> powershell, others -> cmd;
// linux, darwin -> /bin/bash
func Exec(arg ...string) ([]byte, error) {
	switch CurPlatform.OS {
	case "windows":
		switch CurPlatform.Ver {
		case "10", "11":
			out, err := exec.Command("powershell", arg...).CombinedOutput()
			if err != nil {
				return out, err
			}
			return out, nil
		default:
			out, err := exec.Command("cmd", append([]string{"/C"}, arg...)...).CombinedOutput()
			if err != nil {
				return out, err
			}
			return out, nil
		}
	case "linux", "darwin":
		out, err := exec.Command("/bin/bash", append([]string{"/C"}, arg...)...).CombinedOutput()
		if err != nil {
			return out, err
		}
		return out, nil
	default:
	}
	return []byte(""), nil
}

// ExecStr automatically split string to string arrary, then call Exec to execute
func ExecStr(args string) ([]byte, error) {
	return Exec(strings.Fields(args)...)
}

// WriteFile write data or create file to write data according to file
func WriteFile(file string, data []byte) error {
	fp, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fp.Close()
	fp.Write(data)
	return nil
}

// WriteFile write data or create file to write data according to file when file isn't exist
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

// ReadStraemFromFile return total data from specify file
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

// ReadStraemFromFile return data to be filter from specify file
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

// ReadStraemFromFile return data to be filter from string
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

// PathIsExist judge whether path exist. If exist, return true.
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

// Filename returns the last element name of fullpath.
func Filename(fullpath string) string {
	filename := path.Base(fullpath)
	ext := path.Ext(filename)
	return filename[:len(filename)-len(ext)]
}
