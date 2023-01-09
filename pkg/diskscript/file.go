package diskscript

import (
	"context"
	"log"
	"math"
)

func (conf *YamlConfig) FileScripts(ctx context.Context) {
	log.Println("开始文件检查")
	for _, c := range conf.File {
		if !IsFile(c.File) {
			log.Printf("文件:%s检查失败不是文件\n", c.File)
			continue
		}
		filesize, err := FileSize(c.File)
		if err != nil {
			log.Printf("文件检查失败:%s\n", err)
		} else {
			filesizeformat := UnitConvert(filesize)
			if filesize > int64(math.Round(UnitParse(c.Threshold))) {
				log.Printf("文件:%s 阈值:%s 当前占用:%s 触发脚本\n", c.File, c.Threshold, filesizeformat)
				ShellCommand(c.Scripts)
				if c.Alert && conf.Alert.Enable {
					c.AlertAdd(ctx, filesizeformat)
				}
			} else {
				log.Printf("文件:%s 阈值:%s 当前占用:%s 检查通过\n", c.File, c.Threshold, filesizeformat)
			}
		}
	}
}
