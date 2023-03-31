package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func DockerLogFormat(s string) []string {
	var strList []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		// "2023-03-30T09:30:23.339998298Z  High performance, minimalist Go web framework\r",
		if len(line) >= 30 && line[10:11] == "T" && line[29:30] == "Z" {
			line = fmt.Sprintf("%s %s%s", line[0:10], line[11:19], line[30:])
		}
		strList = append(strList, line)
	}
	return strList
}
