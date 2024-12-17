package cryptos

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/pkg/errors"
	"io"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// MD5
func MD5(data []byte) string {
	h := md5.Sum(data)
	return hex.EncodeToString(h[:])
}

// SHA256 sha256
func SHA256(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// GenPwd 生成密码
func GenPwd(pwd string) (enPwd string, err error) {
	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost); err != nil {
		err = errors.New(err.Error())
		return
	} else {
		enPwd = string(hash)
	}
	return
}

// CompPwd 验证密码
func CompPwd(hashPwd, srcPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(srcPwd)); err != nil {
		return false
	}
	return true
}

// MD5File /*
func MD5File(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
		return ""
	}
	hash := md5.New()
	_, _ = io.Copy(hash, f)
	MD5Str := hex.EncodeToString(hash.Sum(nil))
	_ = f.Close()
	return MD5Str
}
