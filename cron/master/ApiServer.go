package master

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

//任务的http接口
type ApiServer struct {
	httpServer *http.Server
}

//单例对象
var (
	G_apiServer *ApiServer
)

//保存任务接口
func handleJobSave(w http.ResponseWriter, r *http.Request) {

}

//初始化http服务
func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	//启动TCP监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	//创建一个http服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	//赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	//启动服务端
	go httpServer.Serve(listener)

	return
}