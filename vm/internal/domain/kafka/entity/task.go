package entity

// Task 任务结构体
type Task struct {
	VMID     string `json:"vm_id"`
	ExpireAt int64  `json:"expire_at"`
	TaskType string `json:"task_type"` // "notify" 或 "destroy"
}
