package spthread

import "sync"

type TaskThreadList struct {
	TaskList map[string]struct {
		// 任务名称
		TaskName string
		// 任务状态
		TaskStatus string
		// 任务执行时间
		TaskTime string
	}
	Lock *sync.Mutex
}

type ThreadWaitList struct {
	TaskName string
	Lock     *sync.Mutex
}
