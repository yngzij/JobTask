package job

import (
	"sync"
)

type JobInterface interface {
	Prepare()
	Do()
}

func NewJob() *Job {
	return &Job{
		maxWorkers: 5,
	}
}

type Job struct {
	prepareQ   chan JobInterface
	JobQ       chan JobInterface
	jobs       []JobInterface
	mux        sync.Mutex
	maxWorkers int
}

func (j *Job) Empty() bool {
	j.mux.Lock()
	defer j.mux.Unlock()
	return len(j.jobs) == 0
}

func (j *Job) SetMaxWorkers(max int) {
	j.maxWorkers = max
}

func (j *Job) PrepareJob() {
	for _, job := range j.jobs {
		j.prepareQ <- job
	}
	close(j.prepareQ)
}

func (j *Job) DoJob() {
	for _, job := range j.jobs {
		j.JobQ <- job
	}
	close(j.JobQ)
}

func (j *Job) Run() {
	j.prepareQ = make(chan JobInterface, 0)
	j.JobQ = make(chan JobInterface, 0)

	j.RunPrepare()
	j.RunTask()
}

func (j *Job) RunTask() {
	go j.DoJob()
	wg := sync.WaitGroup{}
	ch := make(chan interface{}, j.maxWorkers)
	for job := range j.JobQ {
		wg.Add(1)
		ch <- struct{}{}
		go func(wg *sync.WaitGroup, job JobInterface, j *Job) {
			job.Do()
			wg.Done()
			<-ch
		}(&wg, job, j)
	}
	wg.Wait()
}

func (j *Job) RunPrepare() sync.WaitGroup {
	go j.PrepareJob()
	wg := sync.WaitGroup{}
	ch := make(chan interface{}, j.maxWorkers)
	for job := range j.prepareQ {
		wg.Add(1)
		ch <- struct{}{}
		go func(wg *sync.WaitGroup, job JobInterface) {
			job.Prepare()
			<-ch
			wg.Done()
		}(&wg, job)
	}
	wg.Wait()
	close(ch)
	return wg
}

func (j *Job) PushJob(jobs ...JobInterface) {
	for _, job := range jobs {
		j.jobs = append(j.jobs, job)
	}
}
