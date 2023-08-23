package md5

import (
	"crypto/md5"
	"fmt"
	"strings"
)

type EncryptManager struct {
	Salt string
}

func (e *EncryptManager) EncryptPassword(code string) string {
	return Md5Crypt(code, e.Salt)
}

func Md5Crypt(str string, salt ...interface{}) (CryptString string) {
	if l := len(salt); l > 0 {
		slice := make([]string, l+1)
		str = fmt.Sprintf(str+strings.Join(slice, "%v"), salt...)
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
