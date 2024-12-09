package videosys_starter

import (
	goframeworkvideosys "github.com/kordar/goframework-videosys"
	logger "github.com/kordar/gologger"
	"github.com/kordar/goresty"
	"time"
)

type APIPollingWorker struct {
	feigns []*goresty.Feign
}

func NewAPIPollingWorker(feign ...*goresty.Feign) *APIPollingWorker {
	return &APIPollingWorker{feigns: feign}
}

func (A APIPollingWorker) Heartbeat(options *WorkerOption) {

	body := map[string]interface{}{
		"node_id":   options.NodeId,
		"node_host": options.NodeHost,
		"time":      time.Now().Format("2006-01-02 15:04:05"),
	}

	for _, feign := range A.feigns {
		fetchMethod := FetchMethod{}
		_, err := feign.Request().SetBody(body).SetResult(&fetchMethod).Get(options.HeartbeatUrl)
		if err != nil {
			logger.Warnf("[Videosys-APIPollingWorker] request %s err: %v", options.HeartbeatUrl, err)
			return
		}

		logger.Infof("[Videosys-APIPollingWorker] response data is %+v", fetchMethod)

		if fetchMethod.Name == "stop-job" {
			goframeworkvideosys.Stop(fetchMethod.JobId)
		}

		if fetchMethod.Name == "start-job" {
			_ = goframeworkvideosys.Start(fetchMethod.JobId)
		}
	}

}

func (A APIPollingWorker) JobInfo(options *WorkerOption) {
	list := goframeworkvideosys.ConfigList()
	for _, feign := range A.feigns {
		_, err := feign.Request().SetBody(list).Post(options.JobInfoUrl)
		if err != nil {
			logger.Warnf("[Videosys-APIPollingWorker] request %s err: %v", options.JobInfoUrl, err)
			return
		}
	}
}
