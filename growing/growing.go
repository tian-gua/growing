package growing

import "net/http"


//启动http服务器
func start(ip, port string) {

	err := http.ListenAndServe(ip + port, nil)
	if err != nil {
		panic(err)
	}

}
