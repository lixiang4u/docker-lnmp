package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ToJson(v any) string {
	bs, _ := json.MarshalIndent(v, "", "\t")
	return string(bs)
}

func StringMd5(s string) string {
	var md5Hash = md5.New()
	_, err := io.WriteString(md5Hash, s)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}

func StringHash(s string) string {
	return StringMd5(s)[:8]
}

func AppDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return filepath.Dir(filepath.Dir(dir))
}
