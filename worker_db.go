package videosys_starter

import (
	goframeworkvideosys "github.com/kordar/goframework-videosys"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type LocalPollingWorker struct {
	db *gorm.DB
}

func NewLocalPollingWorker(db *gorm.DB) *LocalPollingWorker {
	return &LocalPollingWorker{db}
}

func (l LocalPollingWorker) Heartbeat(options *WorkerOption) {
	info := NodeInfo{
		NodeId:      options.NodeId,
		NodeHost:    options.NodeHost,
		RefreshTime: time.Now(),
	}
	l.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}, {Name: "node_host"}},
		DoUpdates: clause.AssignmentColumns([]string{"refresh_time"}),
	}).Create(&info)
}

func (l LocalPollingWorker) JobInfo(options *WorkerOption) {
	list := goframeworkvideosys.ConfigList()
	data := make([]NodeJobInfo, 0)
	for _, vo := range list {
		info := ConvertToNodeInfo(vo)
		info.NodeId = options.NodeId
		info.NodeHost = options.NodeHost
		data = append(data, info)
	}
	if len(data) == 0 {
		return
	}
	l.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}, {Name: "node_host"}, {Name: "name"}},
		UpdateAll: true,
	}).Create(&data)
}
