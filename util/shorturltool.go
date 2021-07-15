package util

import (
	"encoding/base64"
)

func GenUrlID(url, expireAt string) string {
	data := url+expireAt
	urlEncode := base64.StdEncoding.EncodeToString([]byte(data))
	start := urlEncode[0:5]
	end := urlEncode[len(urlEncode)-5:]
	return start+end
}
