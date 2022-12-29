package md5

import (
	"crypto/md5"
	"encoding/hex"
)

const salt = "MEDICINAL_SHARE_TERRACE"

func Hash(text string) string {
	byt := []byte(text + salt)
	hash := md5.New()
	hash.Write(byt)
	return hex.EncodeToString(hash.Sum(nil))
}
