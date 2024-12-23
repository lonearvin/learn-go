package main

import (
	"fmt"
	"math/rand"
)

type Job struct {
	Id      int
	RandNum int
}

type Result struct {
	// 这里必须传入对象实例
	job *Job
	// 求和
	sum int
}

func main() {
	jobChan := make(chan *Job, 128)
	resultChan := make(chan *Result, 128)

	createPool(64, jobChan, resultChan)
	// 打开打印的协程
	go func(ResultChan chan *Result) {
		for result := range resultChan {
			fmt.Printf("job id:%v randnum:%v result:%v\n", result.job.Id,
				result.job.RandNum, result.sum)
		}
	}(resultChan)

	var id int
	for {
		id++
		// 生成随机数
		rNum := rand.Int()
		job := &Job{
			Id:      id,
			RandNum: rNum,
		}
		jobChan <- job
	}
}

func createPool(num int, job chan *Job, res chan *Result) {
	// 根据开的协程的个数去跑
	for i := 0; i < num; i++ {
		go func(job chan *Job, res chan *Result) {
			// 执行计算
			for j := range job {
				// 随机数接过来
				rNum := j.RandNum
				var sum int
				for rNum != 0 {
					tmp := rNum % 10
					sum += tmp
					rNum /= 10
				}
				r := &Result{
					job: j,
					sum: sum,
				}
				res <- r
			}
		}(job, res)
	}
}
