package videosys_starter

import (
	"github.com/kordar/gocron"
)

type Worker interface {
	Heartbeat(options *WorkerOption)
	JobInfo(options *WorkerOption)
}

type FetchMethod struct {
	Name  string      `json:"name"`
	JobId string      `json:"job_id"`
	Param interface{} `json:"param"`
}

type WorkerOption struct {
	NodeId             string `json:"node_id"`
	NodeHost           string `json:"node_host"`
	Feign              string `json:"feign"`
	FeignHost          string `json:"feign_host"`
	FeignTrace         string `json:"feign_trace"`
	FeignDebug         string `json:"feign_debug"`
	FeignTimeout       int    `json:"feign_timeout"`
	FeignRetryCount    int    `json:"feign_retry_count"`
	FeignRetryWaitTime int    `json:"feign_retry_wait_time"`

	HeartbeatSpec string `json:"heartbeat_spec"`
	HeartbeatUrl  string `json:"heartbeat_url"`
	JobInfoSpec   string `json:"job_info_spec"`
	JobInfoUrl    string `json:"job_info_url"`
}

type WorkerHeartbeatSchedule struct {
	options *WorkerOption
	Worker  Worker
	*gocron.BaseSchedule
}

func NewWorkerHeartbeatSchedule(options *WorkerOption, worker Worker) *WorkerHeartbeatSchedule {
	return &WorkerHeartbeatSchedule{options, worker, &gocron.BaseSchedule{}}
}

func (h *WorkerHeartbeatSchedule) GetId() string {
	return "#worker-heartbeat"
}

func (h *WorkerHeartbeatSchedule) GetSpec() string {
	if h.options.HeartbeatSpec != "" {
		return h.options.HeartbeatSpec
	} else {
		return h.BaseSchedule.GetSpec()
	}
}

func (h *WorkerHeartbeatSchedule) Execute() {
	h.Worker.Heartbeat(h.options)
}

// ------------ job info ----------------------

type WorkerJobInfoSchedule struct {
	options *WorkerOption
	Worker  Worker
	*gocron.BaseSchedule
}

func NewWorkerJobInfoSchedule(options *WorkerOption, worker Worker) *WorkerJobInfoSchedule {
	return &WorkerJobInfoSchedule{options, worker, &gocron.BaseSchedule{}}
}

func (h *WorkerJobInfoSchedule) GetId() string {
	return "#worker-job-info"
}

func (h *WorkerJobInfoSchedule) GetSpec() string {
	if h.options.JobInfoSpec != "" {
		return h.options.JobInfoSpec
	} else {
		return h.BaseSchedule.GetSpec()
	}
}

func (h *WorkerJobInfoSchedule) Execute() {
	h.Worker.JobInfo(h.options)
}
