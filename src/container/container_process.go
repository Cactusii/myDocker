package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	// 在创建子进程之前先创建用于父子进程通信的管道, write管道给父进程，read给子进程
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	/*
		将读文件句柄传递给子进程, 之后转到init.go readUserCommand()，子进程启动后第一件事就是用读文件句柄读取管道中的命令
		另外这是一个文件类型对象，不能通过字符参数的形式传递，因此使用了command的cmd.ExtraFiles方法，
		这个属性的意思时会外带着这个文件句柄去创建子进程。一般一个进程默认会有三个文件描述符，
		分别时标准输入、标准输出、标准错误，这三个是创建进程的时候就会默认带着，那么外带的这个文件描述符就成了第四个，
		可以通过/proc/self/fd查看进程拥有的文件描述符。
	*/
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe() // 返回两个文件句柄
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

/*
一般用户创建进程时会指定一个命令，就是容器创建后执行的指令，例如/bin/bash，
但是考虑到如果用户传送的指令过长，或者其中带有一些特殊字符。
在这里用到了进程间通信机制，父进程通过进程间通信将命令传递给子进程

无名管道通常是具有血缘关系的进程之间的通信，管道就是一个文件，
但和普通的文件不同的是管道的大小有限制，当存满之后再往里面写数据将会阻塞，
同样如果管道当中没有数据，读数据的时候就会阻塞
*/
