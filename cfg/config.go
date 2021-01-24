package cfg

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ServeDir    string           `json:"serve_dir"`
	BaiduAIConf AppKeySecretConf `json:"baidu_ai_conf"`
}

type AppKeySecretConf struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

var _cfg *Config

func LoadConfig(fn string) (*Config, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	_cfg = &Config{}
	err = json.Unmarshal(data, _cfg)
	return _cfg, err
}

func Get() *Config {
	return _cfg
}
