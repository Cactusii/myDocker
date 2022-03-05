package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

/*
每个subsystem需要实现该接口
*/

type Subsystem interface {
	Name() string                               // subsystem 名字
	Set(path string, res *ResourceConfig) error // 设置某个cgroup关联该subsystem, 并设置相应的限制
	Apply(path string, pid int) error           // 将进程添加到某cgroup中
	Remove(path string) error                   // 移除cgroup
}

var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
