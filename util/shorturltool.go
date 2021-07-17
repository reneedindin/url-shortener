package util

import (
	"encoding/base64"
)

func GenUrlID(url, expireAt string) string {
	data := url+expireAt
	urlEncode := base64.StdEncoding.EncodeToString([]byte(data))
	return urlEncode
}
