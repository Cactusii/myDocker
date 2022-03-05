//package main
//
//import (
//	"fmt"
//	"io/ioutil"
//	"os"
//	"os/exec"
//	"path"
//	"strconv"
//	"syscall"
//)
//
//// 挂载memory subsystem的hierarchy的根目录位置
//const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"
//
//func main() {
//
//	/*
//	创建一个进程
//	 */
//	if os.Args[0] == "/proc/self/exe" {
//		fmt.Printf("current pid %d", syscall.Getpid())
//		fmt.Println()
//		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
//		cmd.SysProcAttr = &syscall.SysProcAttr{}
//		cmd.Stdin = os.Stdin
//		cmd.Stdout = os.Stdout
//		cmd.Stderr = os.Stderr
//		if err := cmd.Run(); err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//	}
//	cmd := exec.Command("/proc/self/exe")
//
//	/*
//	创建Namespace
//	 */
//	cmd.SysProcAttr = &syscall.SysProcAttr{
//		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
//			syscall.CLONE_NEWNS,
//	}
//	cmd.Stdin = os.Stdin
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//
//	/*
//	创建cgroups
//	 */
//	if err := cmd.Start(); err != nil {
//		fmt.Println("ERROR", err)
//		os.Exit(1)
//	} else {
//		// 得到fork出来的进程映射在外部命名空间的pid
//		fmt.Printf("%v", cmd.Process.Pid)
//
//		// 在系统默认创建挂在了memory sybsystem的Hierarchy上创建cgroup
//		os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit"), 0755)
//
//		// 将容器进程加入到这个cgroup中，也就是添加到
//		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "tasks"),
//			[]byte(strconv.Itoa(cmd.Process.Pid)), 0644)
//		// 限制cgroup进程使用
//		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "memory.limit_in_bytes"),
//			[]byte("100m"), 0644)
//	}
//	cmd.Process.Wait()
//}
//
//



package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = `mydocker`

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}