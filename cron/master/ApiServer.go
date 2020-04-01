package master

import (
	"encoding/json"
	"net"
	"net/http"
	"prepare/cron/common"
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
//POST job={"name":"job1","command":"echo hello","cronExpr":"* * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	//任务保存在Etcd中
	//1.解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	//2.取表单中的job字段
	postJob = req.PostForm.Get("job")
	//3.反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	//4.保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	//5.返回正常应答 错误返回错误应答({"errno":0, "msg":"","data":{....}})
	if bytes, err = common.BuiIdResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:

	//6.异常返回应答
	if bytes, err = common.BuiIdResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}

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
