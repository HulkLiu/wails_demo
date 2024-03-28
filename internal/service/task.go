package service

import (
	"changeme/config"
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"sync"
	"time"
)

// Task 表示一个任务
type Task struct {
	ID          int    `json:"Id" gorm:"column:Id;AUTO_INCREMENT;not null"`
	Name        string `json:"Name" gorm:"column:Name;not null"`
	Description string `json:"Description" gorm:"column:Description;not null"`
	Completed   int    `json:"Completed" gorm:"column:Completed;not null"`
	CreatedAt   string `json:"CreatedAt" gorm:"column:CreatedAt;not null"`
	UpdatedAt   string `json:"UpdatedAt" gorm:"column:UpdatedAt;not null"`
}

// TableName ...
func (t *Task) TableName() string {
	return "task"
}

func NewTaskDB(dbFilePath string) *gorm.DB {
	//isExist := xfile.IsExist(dbFilePath)
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: gormLogger.Discard,
	})
	if err != nil {
		panic("创建数据库失败")
	}
	// 判断表是否存在
	if !db.Migrator().HasTable(&Task{}) {
		// 创建表
		err = db.AutoMigrate(&Task{})
		if err != nil {
			//panic(err)
			log.Printf("%v", err)
			return nil
		}
	}
	return db
}

// TaskList 表示任务列表
type TaskList struct {
	Tasks []Task
	mu    sync.RWMutex // 用于并发控制
}

// TaskManager 表示任务管理器，包含任务列表
type TaskManager struct {
	TaskList TaskList
	TasksDB  *gorm.DB
}

// NewTask 创建一个新任务
func (tm *TaskManager) NewTask(name, description string) (Task, error) {
	tm.TaskList.mu.Lock()
	defer tm.TaskList.mu.Unlock()

	task := Task{
		//ID:          len(tm.TaskList.Tasks) + 1,
		Name:        name,
		Description: description,
		Completed:   0,
		CreatedAt:   time.Now().Format(config.TimeFormat),
		UpdatedAt:   time.Now().Format(config.TimeFormat),
	}

	tm.TaskList.Tasks = append(tm.TaskList.Tasks, task)
	if err = tm.TasksDB.Create(&task).Error; err != nil {
		return Task{}, err
	}
	return task, nil
}

// CompleteTask 标记任务为完成
func (tm *TaskManager) CompleteTask(id int) error {
	tm.TaskList.mu.Lock()
	defer tm.TaskList.mu.Unlock()

	for i, task := range tm.TaskList.Tasks {
		if task.ID == id {
			tm.TaskList.Tasks[i].Completed = 1
			tm.TaskList.Tasks[i].UpdatedAt = time.Now().Format(config.TimeFormat)
			return nil
		}
	}

	return errors.New("task not found")
}

// DeleteTask 删除一个任务
func (tm *TaskManager) DeleteTask(id int) error {
	tm.TaskList.mu.Lock()
	defer tm.TaskList.mu.Unlock()

	for i, task := range tm.TaskList.Tasks {
		if task.ID == id {
			tm.TaskList.Tasks = append(tm.TaskList.Tasks[:i], tm.TaskList.Tasks[i+1:]...)
			return nil
		}
	}

	return errors.New("task not found")
}

// UpdateTask 更新一个任务
func (tm *TaskManager) UpdateTask(id int, name, description string, completed int) error {
	tm.TaskList.mu.Lock()
	defer tm.TaskList.mu.Unlock()

	for i, task := range tm.TaskList.Tasks {
		if task.ID == id {
			tm.TaskList.Tasks[i].Name = name
			tm.TaskList.Tasks[i].Description = description
			tm.TaskList.Tasks[i].UpdatedAt = time.Now().Format(config.TimeFormat)
			tm.TaskList.Tasks[i].Completed = completed
			return nil
		}
	}

	return errors.New("task not found")
}

// FilterTasks 返回符合过滤条件的任务
func (tm *TaskManager) FilterTasks(filterFunc func(Task) int) []Task {
	tm.TaskList.mu.RLock()
	defer tm.TaskList.mu.RUnlock()

	var filteredTasks []Task
	for _, task := range tm.TaskList.Tasks {
		if filterFunc(task) == 1 {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks
}

// GetTasks 返回所有任务
func (tm *TaskManager) GetTasks() []Task {
	tm.TaskList.mu.RLock()
	defer tm.TaskList.mu.RUnlock()

	List := make([]Task, 0)

	tm.TasksDB.Order("CreatedAt DESC").Find(&List)
	return List
	//return tm.TaskList.Tasks
}

// GetTaskByID 返回指定ID的任务
func (tm *TaskManager) GetTaskByID(id int) (Task, error) {
	tm.TaskList.mu.RLock()
	defer tm.TaskList.mu.RUnlock()

	for _, task := range tm.TaskList.Tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return Task{}, errors.New("task not found")
}
