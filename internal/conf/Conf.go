package conf

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func Conf() Config {
	file, err := ioutil.ReadFile("./configs/conf.yml")
	if err != nil {
		panic("Err read file configuration: " + err.Error())
	}
	conf := Config{}
	yaml.Unmarshal(file, &conf)
	return conf
}
