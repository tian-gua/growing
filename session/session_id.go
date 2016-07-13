package gsession

import (
	"io"
	"crypto/rand"
	"encoding/base64"
)


//生成sessionid
func getSessionId() string {
	b := make([]byte, 32)
	//取32个随机数
	_, err := io.ReadFull(rand.Reader, b);
	if err != nil {
		panic(err)
	}
	//base64加密32个随机数
	return base64.URLEncoding.EncodeToString(b)
}
