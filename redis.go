/*
	Description : redis 相关方法
	Author : ManGe
	Version : v0.1
	Date : 2021-04-30
*/

package gathertool

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net"
	"sync"
	"time"
)

type Rds struct {
	SSHUser string
	SSHPassword string
	SSHAddr string
	RedisHost string
	RedisPost string
	RedisPassword string

	// redis DB
	RedisDB int

	// 单个连接
	Conn redis.Conn

	//	最大闲置数，用于redis连接池
	RedisMaxIdle int

	//	最大连接数
	RedisMaxActive int

	//	单条连接Timeout
	RedisIdleTimeoutSec int

	// 连接池
	Pool *redis.Pool
}

type SSHConnInfo struct {
	SSHUser string
	SSHPassword string
	SSHAddr string
}

func NewSSHInfo( addr, user, password string) *SSHConnInfo {
	return &SSHConnInfo{
		SSHUser : user,
		SSHPassword : password,
		SSHAddr : addr,
	}
}

func NewRedis(host, port, password string, db int, vs ...interface{}) (*Rds) {
	var sshConnInfo SSHConnInfo

	for _,v := range vs{
		log.Println("v = ", v)
		switch vv := v.(type) {
		case *SSHConnInfo:
			sshConnInfo = *vv
		case SSHConnInfo:
			sshConnInfo = vv
		}
	}

	return &Rds{
		SSHUser: sshConnInfo.SSHUser,
		SSHPassword : sshConnInfo.SSHPassword,
		SSHAddr: sshConnInfo.SSHAddr,
		RedisHost:host,
		RedisPost:port,
		RedisPassword:password,
		RedisDB:db,
	}
}

func NewRedisPool(host, port, password string, db, maxIdle, maxActive, idleTimeoutSec int, vs ...interface{}) (*Rds) {
	var sshConnInfo SSHConnInfo

	for _,v := range vs{
		switch vv := v.(type) {
		case *SSHConnInfo:
			sshConnInfo = *vv
		case SSHConnInfo:
			sshConnInfo = vv
		}
	}

	return &Rds{
		SSHUser: sshConnInfo.SSHUser,
		SSHPassword : sshConnInfo.SSHPassword,
		SSHAddr: sshConnInfo.SSHAddr,
		RedisHost:host,
		RedisPost:port,
		RedisPassword:password,
		RedisDB:db,
		RedisMaxIdle:maxIdle,
		RedisMaxActive:maxActive,
		RedisIdleTimeoutSec:idleTimeoutSec,
	}
}


// 	redis连接
func (r *Rds) RedisConn() (err error) {
	host := fmt.Sprintf("%s:%s", r.RedisHost, r.RedisPost)

	if r.SSHPassword != "" && r.SSHUser != "" && r.SSHAddr != ""{
		if sshClient, err := SSHClient(r.SSHUser, r.SSHPassword, r.SSHAddr); err == nil{
			var conn net.Conn
			conn, err = sshClient.Dial("tcp", host)
			r.Conn = redis.NewConn(conn, -1, -1)
		}
	}else{
		r.Conn, err = redis.Dial("tcp", host)
	}

	if err != nil{
		return
	}

	if r.Conn == nil{
		err = errors.New("redis conn is null")
		return
	}

	if r.RedisPassword != "" {
		if _, authErr := r.Conn.Do("AUTH", r.RedisPassword); authErr != nil {
			err = fmt.Errorf("redis auth password error: %s", authErr)
			return
		}
	}

	if r.RedisDB < 1{
		r.RedisDB = 0
	}
	_, err = r.Conn.Do("select", fmt.Sprintf("%d", r.RedisDB))
	return
}


// RPool 连接池连接
// 返回redis连接池  *redis.Pool.Get() 获取redis连接
func (r *Rds) RedisPool() error {
	host := fmt.Sprintf("%s:%s", r.RedisHost, r.RedisPost)

	r.Pool = &redis.Pool{
		MaxIdle:     r.RedisMaxIdle,
		MaxActive:   r.RedisMaxActive,
		IdleTimeout: time.Duration(r.RedisIdleTimeoutSec) * time.Second,
		Dial: func() (redis.Conn, error) {
			var (
				c redis.Conn
				err error
			)
			if r.SSHPassword != "" && r.SSHUser != "" && r.SSHAddr != ""{
				sshClient, err := SSHClient(r.SSHUser, r.SSHPassword, r.SSHAddr)
				if  sshClient != nil {
					var conn net.Conn
					conn, err = sshClient.Dial("tcp", host)
					c = redis.NewConn(conn, -1, -1)
				}
				if err != nil{
					return nil, err
				}
			}else{
				c, err = redis.Dial("tcp", host)
			}

			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}

			if c == nil {
				return nil, fmt.Errorf("redis connection is null.")
			}

			//验证redis密码
			if r.RedisPassword != "" {
				if _, authErr := c.Do("AUTH", r.RedisPassword); authErr != nil {
					return nil, fmt.Errorf("redis auth password error: %s", authErr)
				}
			}

			_, err = c.Do("select", fmt.Sprintf("%d", r.RedisDB))
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
	return nil
}

func (r *Rds) GetConn() redis.Conn{
	if r.Conn != nil{
		return r.Conn
	}
	rc := r.Pool.Get()
	if rc != nil{
		return rc
	}
	return nil
}

func (r *Rds) SelectDB(dbNumber int) error {
	rc := r.GetConn()
	if rc == nil{
		return errors.New("redis conn is nil")
	}
	_, err := rc.Do("select", fmt.Sprintf("%d", dbNumber))
	return  err
}


// Del key
func RedisDELKeys(rds *Rds, keys string, jobNumber int){

	rds.RedisMaxActive = rds.RedisMaxActive+jobNumber*2
	rds.RedisMaxIdle = rds.RedisMaxIdle+jobNumber*2
	rds.RedisPool()

	log.Println(rds.Pool.MaxActive, rds.Pool.MaxIdle)

	conn := rds.Pool.Get()
	queue := NewQueue()
	res, err := redis.Strings(conn.Do("keys", keys))
	if err != nil {
		log.Println("GET error", err.Error())
	}
	conn.Close()

	for _, v := range res {
		queue.Add(&Task{Url: v})
	}
	allNumber := queue.Size()
	log.Println("allNumber = ", allNumber)

	var wg sync.WaitGroup
	for job:=0;job<jobNumber;job++{
		wg.Add(1)
		go func(i int){
			defer wg.Done()

			log.Println("启动第",i ,"个任务")
			for {
				if queue.IsEmpty() || queue.Size() < 2 {
					break
				}
				task := queue.Poll()
				log.Println("第",i,"个任务取的值： ", task.Url)
				c := rds.Pool.Get()
				s,err := redis.Int64(c.Do("DEL", task.Url))
				if err != nil || s == 0 {
					log.Println("redis command:  err : ", err)
					queue.Add(task)
				}else{
					log.Println("删除成功 ！！！")
				}
				c.Close()
				log.Println(fmt.Sprintf("[进度] %d/%d  %f %%", allNumber - queue.Size(),
					allNumber, (float64(allNumber - queue.Size())/float64(allNumber))*100))
			}

			log.Println("第",i ,"个任务结束！！")
		}(job)
	}
	wg.Wait()
	log.Println("执行完成！！！")

}
