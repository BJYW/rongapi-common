package config

var LoggerConfig loggerConfig

type loggerConfig struct {
	Url    string `goblet:"url,http://1130158802826311.mns.cn-beijing-internal.aliyuncs.com"`
	Queue  string `goblet:"queue,logger"`
	Key    string `goblet:"key,SzBTLQdeMWnF7fSA"`
	Secret string `goblet:"secret,RvfSHgrpYceGUnkLGaQ8SYaslCdOIz"`
}


