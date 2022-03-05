package cgroups

import (
	"github.com/sirupsen/logrus"
	"myDocker/cgroups/subsystems"
)

/*
创建带资源限制的容器时:
CgroupManager 在配置容器资源限制时，首先会初始化 Subsystem 的实例，
然后遍历调用 Subsystem 实例的 Set 方法，创建和配置不同 Subsystem 挂载的 hierarchy 中的 cgroup，
最后再通过调用 Subsystem 实例将容器的进程分别加入到那些 cgroup 中，实现容器的 资源限制.
*/

type CgroupManager struct {
	// cgroup在hierarchy中的一个相对路径 相当于创建的cgroup目录相对于root cgroup目录的路径
	Path string
	// 资源配置
	//Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// 将进程pid加入到这个cgroup中

func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

// 设置cgroup资源限制

func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

//释放cgroup

func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
