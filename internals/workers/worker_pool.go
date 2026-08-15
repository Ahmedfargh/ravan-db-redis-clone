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
	evaluator   *evaluator.Evaluator
}

func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
	wp := &WorkerPool{
		JobQueue:    make(chan Job, queueSize),
		WorkerCount: workerCount,
		evaluator:   evaluator.NewEvaluator(nil),
	}
	wp.Start()
	return wp
}

func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.WorkerCount; i++ {
		workerID := i
		go func(id int) {
			for job := range wp.JobQueue {
				res := wp.executeJob(job, id)
				job.ResponseChan <- res
			}
		}(workerID)
	}
}

func (wp *WorkerPool) executeJob(job Job, workerID int) Response {
	result, err := wp.evaluator.EvaluateQuery(job.Query)
	if err != nil {
		return Response{Err: err}
	}
	return Response{Result: result}
}
