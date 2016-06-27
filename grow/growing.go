package grow

import (
	"net/http"
	"strconv"
)


//启动http服务器
func Start(port int) {

	err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		panic(err)
	}

}


