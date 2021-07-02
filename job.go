/*
	Description : 并发工作任务
	Author : ManGe
	Version : v0.1
	Date : 2021-04-28
*/

package gathertool

import (
	"log"
	"net/http"
	"runtime"
	"sync"
)

//TODO:  StartJob 开始运行并发， 取task 的 method
func StartJob(){}

// 设置遇到错误执行 Retry 事件
type Err2Retry bool

// StartJobGet 并发执行Get,直到队列任务为空
// @jobNumber 并发数，
// @queue 全局队列，
// @client 单个并发任务的client，
// @SucceedFunc 成功方法，
// @ RetryFunc重试方法，
// @FailedFunc 失败方法
func StartJobGet(jobNumber int, queue TodoQueue, vs ...interface{}){

	var (
		client *http.Client
		succeed SucceedFunc
		retry RetryFunc
		failed FailedFunc
		err2Retry Err2Retry
	)

	runtime.GOMAXPROCS(runtime.NumCPU())

	for _,v := range vs{
		switch vv := v.(type) {
		case *http.Client:
			client = vv
		case SucceedFunc:
			succeed = vv
		case FailedFunc:
			failed = vv
		case RetryFunc:
			retry = vv
		case Err2Retry:
			err2Retry = vv
			}
	}

	var wg sync.WaitGroup
	for job:=0;job<jobNumber;job++{
		wg.Add(1)
		go func(i int){
			log.Println("启动第",i ,"个任务")
			defer wg.Done()
			for {
				if queue.IsEmpty(){
					break
				}
				task := queue.Poll()
				log.Println("第",i,"个任务取的值： ", task)
				ctx := NewGet(task.Url, task)
				if client != nil {
					ctx.Client = client
				}
				if succeed != nil {
					ctx.SetSucceedFunc(succeed)
				}
				if retry != nil {
					ctx.SetRetryFunc(retry)
				}
				if failed != nil {
					ctx.SetFailedFunc(failed)
				}
				if err2Retry {
					ctx.OpenErr2Retry()
				}

				switch task.Type {
				case "","do":
					ctx.Do()
				case "upload":
					if task.SavePath == ""{
						task.SavePath = task.SaveDir + task.FileName
					}
					ctx.Upload(task.SavePath)
				default:
					ctx.Do()
				}

			}
			log.Println("第",i ,"个任务结束！！")
		}(job)
	}
	wg.Wait()
	log.Println("执行完成！！！")
}


// StartJobPost 开始运行并发Post
func StartJobPost(jobNumber int, queue TodoQueue, vs ...interface{}){
	var (
		client *http.Client
		succeed SucceedFunc
		retry RetryFunc
		failed FailedFunc
		err2Retry Err2Retry
	)

	for _,v := range vs{
		switch vv := v.(type) {
		case *http.Client:
			loger("have Client")
			client = vv
		case SucceedFunc:
			succeed = vv
		case FailedFunc:
			failed = vv
		case RetryFunc:
			retry = vv
		case Err2Retry:
			err2Retry = vv
		}
	}

	var wg sync.WaitGroup
	for job:=0;job<jobNumber;job++{
		wg.Add(1)
		go func(i int){
			log.Println("启动第",i ,"个任务")
			defer wg.Done()
			for {
				if queue.IsEmpty(){
					break
				}
				task := queue.Poll()
				log.Println("第",i,"个任务取的值： ", task, task.HeaderMap)
				ctx := NewPost(task.Url, []byte(task.JsonParam), "application/json;", task)
				if client != nil {
					ctx.Client = client
				}
				if succeed != nil {
					ctx.SetSucceedFunc(succeed)
				}
				if retry != nil {
					ctx.SetRetryFunc(retry)
				}
				if failed != nil {
					ctx.SetFailedFunc(failed)
				}
				if err2Retry {
					ctx.OpenErr2Retry()
				}

				switch task.Type {
				case "","do":
					ctx.Do()
				case "upload":
					if task.SavePath == ""{
						task.SavePath = task.SaveDir + task.FileName
					}
					ctx.Upload(task.SavePath)
				default:
					ctx.Do()
				}

			}
			log.Println("第",i ,"个任务结束！！")
		}(job)
	}
	wg.Wait()
	log.Println("执行完成！！！")
}
