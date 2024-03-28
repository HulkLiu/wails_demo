package internal

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/process"
)

//func main() {
//	err := cmdDockerStart()
//	if err != nil {
//		log.Printf("err:%v", err)
//	}
//}

func CmdDockerStart() error {
	// 启动Docker Desktop
	cmd := exec.Command("C:\\Program Files\\Docker\\Docker\\Docker Desktop.exe")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Docker Desktop started successfully.")

	// 监听Docker Desktop启动完成
	go func() {

		for {
			if isDockerDesktopRunning() {
				DockerDesktopStarted <- true
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()
	go startContainer()
	// 设置超时时间
	return nil
}

func startContainer() error {
	for {
		// 等待Docker Desktop启动完成或超时
		select {
		case <-DockerDesktopStarted:
			// 启动容器命令
			cmd := exec.Command("cmd.exe", "/C", "docker", "start", ContainerID)
			_, err := cmd.CombinedOutput()
			if err != nil {
				return err
			}
			DockerDesktopStarted2 <- true
			break
			//fmt.Printf("%s", res)
		case <-time.After(10 * time.Second):
			return fmt.Errorf("timeout waiting for Docker Desktop to start")
		}
	}
}

// 使用gopsutil包来检查Docker Desktop是否正在运行
func isDockerDesktopRunning() bool {
	processes, err := process.Processes()
	if err != nil {
		log.Fatal(err)
	}
	var n int
	for _, p := range processes {
		name, _ := p.Name()

		if strings.Contains(name, "Docker") {
			n++
			log.Printf("%v => Pid:%v, Name:%v", n, p.Pid, name)
			if n == 5 {
				return true
			}
		}
		/*
			2023/12/19 22:03:47 1 => Pid:5372, Name:Docker Desktop.exe
			2023/12/19 22:03:47 2 => Pid:14248, Name:Docker Desktop.exe
			2023/12/19 22:03:48 3 => Pid:18412, Name:Docker Desktop.exe
			2023/12/19 22:03:48 4 => Pid:19792, Name:Docker Desktop.exe
			2023/12/19 22:03:48 5 => Pid:23024, Name:Docker Desktop.exe
		*/
	}

	return false
}
