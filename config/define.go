/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package config

import (
	"github.com/fengyuqin/kungfu/v2/treaty"
)

type Config struct {
	Discover  DiscoverConf              `json:"discover"`
	Rpc       RpcConf                   `json:"rpc"`
	Stores    StoresConf                `json:"stores"`
	Connector ConnectorConf             `json:"connector"`
	Servers   map[string]*treaty.Server `json:"servers"`
	Domains   map[string]string         `json:"domains"`
	Ssl       SslConf                   `json:"ssl"`
}

type DiscoverConf struct {
	UseType      string   `json:"use_type"`
	DialTimeout  int      `json:"dial_timeout"`
	Endpoints    []string `json:"endpoints"`
	ServerPrefix string   `json:"server_prefix"`
	DataPrefix   string   `json:"data_prefix"`
}

type RpcConf struct {
	UseType     string   `json:"use_type" mapstructure:"use_type"`
	DialTimeout int      `json:"dial_timeout" mapstructure:"dial_timeout"`
	Endpoints   []string `json:"endpoints" mapstructure:"endpoints"`
	DebugMsg    bool     `json:"debug_msg" mapstructure:"debug_msg"`
	Prefix      string   `json:"prefix" mapstructure:"prefix"`
}

type StoresConf struct {
	UseType     string   `json:"use_type"`
	DialTimeout int      `json:"dial_timeout"`
	Endpoints   []string `json:"endpoints"`
	Password    string   `json:"password"`
	DB          int      `json:"db"`
	Prefix      string   `json:"prefix"`
}

type ConnectorConf struct {
	UseType           string `json:"use_type"`            //使用的协议
	UseWebsocket      bool   `json:"use_websocket"`       //是否使用websocket
	WebsocketPath     string `json:"websocket_path"`      //websocket路径
	UseSerializer     string `json:"use_serializer"`      //使用的协议
	ProtoPath         string `json:"proto_path"`          //protobuf位置
	HeartbeatInterval int    `json:"heartbeat_interval"`  //心跳间隔
	Version           string `json:"version"`             //当前tcpserver版本号
	MaxPacketSize     int32  `json:"max_packet_size"`     //都需数据包的最大值
	MaxConn           int    `json:"max_conn"`            //当前服务器主机允许的最大链接个数
	WorkerPoolSize    int    `json:"worker_pool_size"`    //业务工作Worker池的数量
	MaxWorkerTaskLen  int32  `json:"max_worker_task_len"` //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen     int32  `json:"max_msg_chan_len"`    //SendBuffMsg发送消息的缓冲最大长度
	LogDir            string `json:"log_dir"`             //日志所在文件夹 默认"./log"
	LogFile           string `json:"log_file"`            //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
	LogDebugClose     bool   `json:"log_debug_close"`     //是否关闭Debug日志级别调试信息 默认false  -- 默认打开debug信息
	TokenKey          string `json:"token_key"`           //token生成key
}

type SslConf struct {
	PowerOn  bool   `json:"power_on"`  //是否是用ssl
	CertFile string `json:"cert_file"` //证书文件地址
	KeyFile  string `json:"key_file"`  //key文件地址
}
type TecentOBS struct {
	SecretId   string `json:"secret_id"`   //秘钥ID
	SecretKey  string `json:"secret_key"`  //秘钥key
	ServiceUrl string `json:"service_url"` //服务地址
	BulletUrl  string `json:"bullet_url"`  //存储桶地址
}

type TecentSms struct {
	SecretId  string `json:"secret_id" mapstructure:"secret_id"`   //秘钥ID
	SecretKey string `json:"secret_key" mapstructure:"secret_key"` //秘钥key
	EndPoint  string `json:"endpoint" mapstructure:"endpoint"`     //服务地址
	Region    string `json:"region" mapstructure:"region"`         //地域
	SdkAppid  string `json:"sdk_appid" mapstructure:"sdk_appid"`   //sdk appid
	SignName  string `json:"sign_name" mapstructure:"sign_name"`   //签名内容
}
