package config

import (
	"errors"
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

var File *goconfig.ConfigFile

const configFile = "/conf/conf.ini"

func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + configFile

	if !fileExists(configPath) {
		panic(errors.New("配置文件不存在"))
	}

	len := len(os.Args)
	if len > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + configFile
		}
	}

	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("读取配置文件出错", err)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func A() {
	fmt.Println("hello world")
}
