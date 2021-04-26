/*
	Description : 任务队列
	Author : ManGe
	Version : v0.1
	Date : 2021-04-24
*/

package gathertool

import (
	"fmt"
	"log"
	"sync"
)

type TodoQueue interface {
	Add(task *Task)    //向队列中添加元素
	Poll()   *Task  //移除队列中最前面的元素
	Clear()  bool   //清空队列
	Size()  int     //获取队列的元素个数
	IsEmpty() bool  //判断队列是否是空
	Print() // 打印
}

// 队列
type Queue struct {
	mux sync.RWMutex
	list []*Task
}

// 任务对象
type Task struct {
	Url string
	Context map[string]interface{}
	Urls []*ReqUrl // 多步骤使用
}

// 单个请求地址对象
type ReqUrl struct {
	Url string
	Method  string
	Params  map[string]interface{}
}

// NewQueue 新建一个队列
func NewQueue() TodoQueue {
	list := make([]*Task, 0)
	return &Queue{list: list, mux: sync.RWMutex{}}
}

// Add 向队列中添加元素
func (q *Queue) Add(task *Task) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.list = append(q.list,task)
}

// Poll 移除队列中最前面的额元素
func (q *Queue) Poll() *Task {
	q.mux.RLock()
	defer q.mux.RUnlock()
	if q.IsEmpty() {
		fmt.Println("queue is empty!")
		return &Task{}
	}

	first := q.list[0]
	q.list = q.list[1:]
	return first
}

func (q *Queue) Clear() bool {
	if q.IsEmpty() {
		fmt.Println("queue is empty!")
		return false
	}
	for i:=0 ; i< q.Size() ; i++ {
		q.list[i].Url = ""
	}
	q.list = nil
	return true
}

func (q *Queue) Size() int {
	return len(q.list)
}

func (q *Queue) IsEmpty() bool {
	if len(q.list) == 0 {
		return true
	}
	return false
}

func (q *Queue) Print() {
	log.Println(q.list)
}