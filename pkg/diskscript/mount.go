package diskscript

import (
	"context"
	"github.com/shirou/gopsutil/disk"
	"log"
)

func (conf *YamlConfig) MountScripts(ctx context.Context) {
	log.Println("开始挂载检查")
	for _, c := range conf.Mount {
		diskusage, err := disk.Usage(c.Mount)
		if err != nil {
			log.Fatalf("挂载点%s错误:%s", c.Mount, err)
		} else {
			PercentageOfOccupancy := float64(diskusage.Used) / float64(diskusage.Total) * 100
			if PercentageOfOccupancy > float64(c.Threshold) {
				log.Printf("挂载:%s 限制:%d%% 当前:%d%% 触发脚本\n", c.Mount, c.Threshold, int64(PercentageOfOccupancy))
				ShellCommand(c.Scripts)
				if c.Alert && conf.Alert.Enable {
					c.AlertAdd(ctx, int64(PercentageOfOccupancy))
				}
			} else {
				log.Printf("挂载:%s 限制:%d%% 当前:%d%% 检查通过\n", c.Mount, c.Threshold, int64(PercentageOfOccupancy))
			}
		}
	}
}
