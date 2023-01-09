package diskscript

import (
	"context"
	"log"
	"math"
)

func (conf *YamlConfig) DirectoryScripts(ctx context.Context) {
	log.Println("开始目录检查")
	for _, c := range conf.Directory {
		if !IsDir(c.Directory) {
			log.Printf("目录:%s检查失败不是目录\n", c.Directory)
			continue
		}
		directorysize, err := DirSize(c.Directory)
		if err != nil {
			log.Printf("目录检查失败:%s\n", err)
		} else {
			filesizeformat := UnitConvert(directorysize)
			if directorysize > int64(math.Round(UnitParse(c.Threshold))) {
				log.Printf("目录:%s 阈值:%s 当前占用:%s 触发脚本", c.Directory, c.Threshold, filesizeformat)
				ShellCommand(c.Scripts)
				if c.Alert && conf.Alert.Enable {
					c.AlertAdd(ctx, filesizeformat)
				}
			} else {
				log.Printf("目录:%s 阈值:%s 当前占用:%s 检查通过", c.Directory, c.Threshold, filesizeformat)
			}
		}
	}
}
