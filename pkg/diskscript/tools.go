package diskscript

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func FileSize(path string) (int64, error) {
	var size int64
	info, err := os.Stat(path)
	if err != nil {
		return size, err
	}
	if os.IsNotExist(err) {
		log.Printf("文件:%s不存在", path)
	}
	return info.Size(), nil
}
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func ExecCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf(":%s 执行错误: %s", cmd, err)
	}
	if Conf.Console {
		log.Printf("\n%s", stdoutStderr)
	}
}

func ShellCommand(c []string) {
	for _, shell := range c {
		ExecCommand(shell)
		log.Printf("执行脚本:%s 完成\n", shell)
	}
}
func UnitCheck(size string) bool {
	unit := size[len(size)-1:]
	for _, un := range []string{"k", "m", "g", "t"} {
		if unit == un {
			return true
		}
	}
	return false
}
func UnitParse(size string) float64 {
	size = strings.ToLower(size)
	unit := size[len(size)-1:]
	var sizefloat64 float64
	if UnitCheck(size) == false {
		log.Fatalln("单位错误,只能选择 k m g t 当前:", size)
	}
	sizefloat64, _ = strconv.ParseFloat(strings.Split(size, unit)[0], 64)
	switch {
	case unit == "k":
		sizefloat64 = sizefloat64 * float64(1000)
	case unit == "m":
		sizefloat64 = sizefloat64 * float64(1000*1000)
	case unit == "g":
		sizefloat64 = sizefloat64 * float64(1000*1000*1000)
	case unit == "t":
		sizefloat64 = sizefloat64 * float64(1000*1000*1000*1000)
	}
	return sizefloat64

}

func UnitConvert(fileSize int64) (size string) {
	if fileSize < 1000 {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1000 * 1000) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1000))
	} else if fileSize < (1000 * 1000 * 1000) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1000*1000))
	} else if fileSize < (1000 * 1000 * 1000 * 1000) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1000*1000*1000))
	} else if fileSize < (1000 * 1000 * 1000 * 1000 * 1000) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1000*1000*1000*1000))
	} else {
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1000*1000*1000*1000*1000))
	}
}
