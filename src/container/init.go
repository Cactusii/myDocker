package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

/*

此方法是在容器中运行的，在此方法运行之前，容器已经创建，并有了一个init进程

此方法主要是挂在proc文件系统，以便可以检查当前进程的资源情况，其中Flags分别表示：
MS_NOEXEC: 在本文件系统中不允许运行其他程序；
MS_NOSUID: 在本文件系统运行程序时，不允许set-user-ID/set-group-ID；
MS_NODEV: 所有mount的系统的默认参数。

最后调用syscall.Exec()，其实对应了execve(filename, argv)系统调用，功能是执行filename进程，
并覆盖当前进程，包括镜像、数据、堆栈、pid等信息。也就是说此进程会替换掉init进程称为新的init进程。
*/

func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("command %s", command)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}
