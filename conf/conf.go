package conf

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	atomicConf atomic.Value
)

func Init(confPath string) {
	config, err := loadConfig(confPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	atomicConf.Store(config)
	go refreshConfig(confPath)
}

func InitParams(confParams map[string]string) {
	config := conf{}
	config.Server.Addr = confParams["addr"]
	if _, exist := confParams["maxhttptime"]; exist {
		if err := config.Server.MaxHTTPTime.UnmarshalText([]byte(confParams["maxhttptime"])); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	} else {
		config.Server.MaxHTTPTime.Duration = time.Second * 5
	}
	cs := confParams["corpsecret"]
	config.WeChatConf = map[string]*WeChatConf{}
	config.WeChatConf[cs] = &WeChatConf{
		CorpId:     confParams["corpid"],
		CorpSecret: confParams["corpsecret"],
		MediaId:    confParams["mediaid"],
	}
	config.WeChatConf[cs].AgentId, _ = strconv.ParseInt(confParams["agentid"], 10, 64)
	config.WeChatConf[cs].EnableDuplicateCheck, _ = strconv.Atoi(confParams["enableduplicatecheck"])
	config.WeChatConf[cs].DuplicateCheckInterval, _ = strconv.Atoi(confParams["duplicatecheckinterval"])
	atomicConf.Store(config)
}

func GetConfig() conf {
	return atomicConf.Load().(conf)
}

func loadConfig(confPath string) (conf, error) {
	var config conf
	_, err := toml.DecodeFile(confPath, &config)
	if err != nil {
		return config, errors.New(fmt.Sprintf("config file err: %v", err))
	}
	return config, nil
}

func refreshConfig(confPath string) {
	for {
		time.Sleep(time.Second)
		config, err := loadConfig(confPath)
		if err == nil {
			atomicConf.Store(config)
		} else {
			log.Println(err)
		}
	}
}

type conf struct {
	Server     Server
	WeChatConf map[string]*WeChatConf
}

type WeChatConf struct {
	CorpId                 string
	CorpSecret             string
	AgentId                int64
	MediaId                string
	EnableDuplicateCheck   int
	DuplicateCheckInterval int
}

type Server struct {
	Addr        string
	MaxHTTPTime duration
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
