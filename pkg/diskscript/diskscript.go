package diskscript

import (
	"context"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

func Execute(configfile string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("读取配置文件: ", configfile)
	file, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Fatalf("yaml文件读取失败%s", err)
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		log.Fatalf("yaml文件解析失败: %v", err)
	}
	AlertmanagerURL, err = url.Parse(Conf.Alert.Url)
	if err != nil {
		log.Fatal(err)
	}
	Conf.MountScripts(ctx)
	Conf.DirectoryScripts(ctx)
	Conf.FileScripts(ctx)
}

func init() {
	Hostname, _ = os.Hostname()
}
