package vo

// VMSpec 虚拟机规格值对象
type VMSpec struct {
	CPUCores    int
	MemoryGB    float64
	OSType      string
	OSVersion   string
}

// NewVMSpec 创建新的虚拟机规格
func NewVMSpec(cpuCores int, memoryGB float64, osType, osVersion string) *VMSpec {
	return &VMSpec{
		CPUCores:    cpuCores,
		MemoryGB:    memoryGB,
		OSType:      osType,
		OSVersion:   osVersion,
	}
}

// Validate 验证虚拟机规格是否合法
func (s *VMSpec) Validate() bool {
	return s.CPUCores > 0 && s.CPUCores <= 8 && 
		s.MemoryGB > 0 && s.MemoryGB <= 16 && 
		s.OSType != "" && s.OSVersion != ""
}