package main

import (
	"log"
	"spider/config"
	"spider/scheduler"
	"spider/tasks"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建调度器
	sched := scheduler.NewScheduler()

	// 初始化并添加任务
	for _, taskCfg := range cfg.Tasks {
		var task tasks.Task
		// 根据任务类型选择具体任务实现
		// 这里假设所有任务都是 Task，可以根据实际需求扩展
		task = tasks.NewTask()
		if err := task.Init(taskCfg); err != nil {
			log.Printf("Failed to init task %s: %v", taskCfg.Name, err)
			continue
		}
		if err := sched.AddTask(taskCfg.Schedule, task.Execute); err != nil {
			log.Printf("Failed to add task %s to scheduler: %v", taskCfg.Name, err)
			continue
		}
	}

	// 启动调度器
	sched.Start()

	// 阻塞主线程
	select {}
}
