package job

import "time"

func NewJob(sleepTime int) Job {
	return &job{
		sleepTime: sleepTime,
	}
}

type Queue chan Job

type Job interface {
	Do()
}

type job struct {
	sleepTime int
}

func (j *job) Do() {
	time.Sleep(time.Millisecond * time.Duration(j.sleepTime))
}
