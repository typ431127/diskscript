package main

import (
	"diskscript/pkg/diskscript"
	"flag"
)

var configfile string

func init() {
	flag.StringVar(&configfile, "conf", "config.yml", "指定配置文件")
}

func main() {
	flag.Parse()
	diskscript.Execute(configfile)
}
