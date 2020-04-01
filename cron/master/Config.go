package master

import (
	"encoding/json"
	"io/ioutil"
)

//程序配置
type Config struct {
	ApiPort         int      `json:"api_port"`
	ApiReadTimeout  int      `json:"api_read_timeout"`
	ApiWriteTimeout int      `json:"api_write_timeout"`
	EtcdEndpoints   []string `json:"etcd_endpoints"`
	EtcdDialTimeout int      `json:"etcd_dial_timeout"`
}

//配置单例
var (
	//单例模式
	G_config *Config
)

//加载配置
func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)
	//1.把配置文件读出来
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	//2.json反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	//3.赋值单例
	G_config = &conf

	return
}
