package workers

import (
	"Raven/internals/evaluator"
)

type Response struct {
	Result string
	Err    error
}

type Job struct {
	Query        string
	ResponseChan chan Response
}

type WorkerPool struct {
	JobQueue    chan Job
	WorkerCount int
}

func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
	wp := &WorkerPool{
		JobQueue:    make(chan Job, queueSize),
		WorkerCount: workerCount,
	}
	wp.Start()
	return wp
}

func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.WorkerCount; i++ {
		workerID := i
		go func(id int) {
			for job := range wp.JobQueue {
				res := executeJob(job, id)
				job.ResponseChan <- res
			}
		}(workerID)
	}
}

func executeJob(job Job, workerID int) Response {
	result, err := evaluator.EvaluateQuery(job.Query)
	if err != nil {
		return Response{Err: err}
	}
	return Response{Result: result}
}
