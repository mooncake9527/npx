package files

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	imgext "github.com/shamsher31/goimgext"
)

func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	if err != nil {
		return 0, errors.New(err.Error())
	}
	return len(content), err
}

// GetExt 获取文件后缀
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// CheckExist 检查文件是否存在
func CheckExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CheckPermission 检查文件权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// IsNotExistMkDir 检查文件夹是否存在
// 如果不存在则新建文件夹
func IsNotExistMkDir(src string) error {
	if exist := !CheckExist(src); !exist {
		if err := MkDir(src); err != nil {
			return errors.New(err.Error())
		}
	}
	return nil
}

// MkDir 新建文件夹
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// Open 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer f.Close()
	return f, nil
}

// GetImgType 获取Img文件类型
func GetImgType(p string) (string, error) {
	file, err := os.Open(p)
	if err != nil {
		return "", errors.New(err.Error())
	}
	defer file.Close()
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return "", errors.New(err.Error())
	}
	filetype := http.DetectContentType(buff)
	ext := imgext.Get()
	for i := 0; i < len(ext); i++ {
		if strings.Contains(ext[i], filetype[6:len(filetype)]) {
			return filetype, nil
		}
	}
	return "", errors.New("invalid image type")
}

// GetType 获取文件类型
func GetType(p string) (string, error) {
	file, err := os.Open(p)
	if err != nil {
		return "", errors.New(err.Error())
	}
	defer file.Close()
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return "", errors.New(err.Error())
	}
	filetype := http.DetectContentType(buff)
	return filetype, nil
}

func PathCreate(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// PathExist 判断目录是否存在
func PathExist(addr string) bool {
	s, err := os.Stat(addr)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()
}

func FileCreate(content bytes.Buffer, name string) {
	file, err := os.Create(name)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(file)
	if err != nil {
		log.Println(err)
	}
	_, err = file.WriteString(content.String())
	if err != nil {
		log.Println(err)
	}
}
