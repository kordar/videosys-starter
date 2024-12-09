package videosys_starter

import (
	"github.com/go-resty/resty/v2"
	goframeworkcron "github.com/kordar/goframework-cron"
	"github.com/kordar/goframework-gorm-mysql"
	goframeworkresty "github.com/kordar/goframework-resty"
	logger "github.com/kordar/gologger"
	"github.com/kordar/goresty"
	"github.com/spf13/cast"
	"time"
)

type StreamNodeInfoModule struct {
	name string
	load func(id string, option *WorkerOption, cfg map[string]string)
}

func NewStreamNodeInfoModule(name string, load func(id string, option *WorkerOption, cfg map[string]string)) *StreamNodeInfoModule {
	return &StreamNodeInfoModule{name, load}
}

func (m StreamNodeInfoModule) Name() string {
	return m.name
}

func (m StreamNodeInfoModule) Load(value interface{}) {

	cfg := cast.ToStringMapString(value)

	if cfg["id"] == "" {
		logger.Warnf("[%s] please configure the attribute id.", m.Name())
		return
	}

	id := cfg["id"]
	option := &WorkerOption{
		NodeId:             cfg["node_id"],
		NodeHost:           cfg["node_host"],
		Feign:              cfg["feign"],
		FeignHost:          cfg["feign_host"],
		FeignTrace:         cfg["feign_trace"],
		FeignDebug:         cfg["feign_debug"],
		FeignTimeout:       cast.ToInt(cfg["feign_timeout"]),
		FeignRetryCount:    cast.ToInt(cfg["feign_retry_count"]),
		FeignRetryWaitTime: cast.ToInt(cfg["feign_retry_wait_time"]),
		HeartbeatSpec:      cfg["heartbeat_spec"],
		HeartbeatUrl:       cfg["heartbeat_url"],
		JobInfoSpec:        cfg["job_info_spec"],
		JobInfoUrl:         cfg["job_info_url"],
	}

	if cfg["type"] == "custom" && m.load != nil {
		m.load(id, option, cfg)
		return
	}

	if cfg["type"] == "polling_db" {
		_ = goframeworkcron.AddGocronInstance(id, nil, nil)
		if cfg["db"] == "" {
			logger.Errorf("[%s] the database instance must.", m.Name())
			return
		}
		exist := goframework_gorm_mysql.HasMysqlInstance(cfg["db"])
		if !exist {
			logger.Errorf("[%s] the database '%s' instance must.", m.Name(), cfg["db"])
			return
		}
		db := goframework_gorm_mysql.GetMysqlDB(cfg["db"])
		worker := NewLocalPollingWorker(db)
		goframeworkcron.AddJob(id, NewWorkerHeartbeatSchedule(option, worker))
		goframeworkcron.AddJob(id, NewWorkerJobInfoSchedule(option, worker))
		return
	}

	if cfg["type"] == "polling_api" {
		_ = goframeworkcron.AddGocronInstance(id, nil, nil)
		feign := getFeign(m.Name(), option)
		worker := NewAPIPollingWorker(feign)
		goframeworkcron.AddJob(id, NewWorkerHeartbeatSchedule(option, worker))
		goframeworkcron.AddJob(id, NewWorkerJobInfoSchedule(option, worker))
		return
	}

}

func (m StreamNodeInfoModule) Close() {
}

func getFeign(name string, options *WorkerOption) *goresty.Feign {
	if options.NodeId == "" {
		logger.Fatalf("[%s] you must configure the parameter \"node_id\"", name)
	}

	if options.FeignHost == "" {
		logger.Fatalf("[%s] you must configure the parameter \"feign_host\"", name)
	}

	if options.Feign != "" && goframeworkresty.HasFeignInstance(options.Feign) {
		return goframeworkresty.GetFeignClient(options.Feign)
	}

	feign := goresty.NewFeign(nil)
	feign.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		if options.FeignTrace == "enable" {
			request.EnableTrace()
		}
		return nil
	})

	feign.Options(func(client *resty.Client) {
		client.SetBaseURL(options.FeignHost)
		if options.FeignDebug == "enable" {
			client.SetDebug(true)
		}
		if options.FeignTimeout != 0 {
			remoteTimeout := cast.ToDuration(options.FeignTimeout)
			client.SetTimeout(time.Second * remoteTimeout)
		}
		if options.FeignRetryCount != 0 {
			client.SetRetryCount(options.FeignRetryCount)
		}
		if options.FeignRetryWaitTime != 0 {
			remoteRetryWaitTime := cast.ToDuration(options.FeignRetryWaitTime)
			client.SetRetryWaitTime(time.Second * remoteRetryWaitTime)
		}
	})

	feign.OnError(func(request *resty.Request, err error) {
		logger.Errorf("[%s] request err = %+v", name, err)
	})

	return feign
}
