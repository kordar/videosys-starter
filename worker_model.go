package videosys_starter

import (
	"encoding/json"
	videocollection "github.com/kordar/video-collection"
	"time"
)

type NodeInfo struct {
	NodeId      string    `json:"node_id" gorm:"column:node_id"`
	NodeHost    string    `json:"node_host" gorm:"column:node_host"`
	RefreshTime time.Time `json:"refresh_time" gorm:"column:refresh_time"`
}

type NodeJobInfo struct {
	NodeId              string `json:"node_id" gorm:"column:node_id"`
	NodeHost            string `json:"node_host" gorm:"column:node_host"`
	Name                string `json:"name" gorm:"column:name"`
	Input               int    `json:"input" gorm:"column:input"`
	InputLabel          string `json:"input_label" gorm:"column:input_label"`
	Output              int    `json:"output" gorm:"column:output"`
	OutputLabel         string `json:"output_label" gorm:"column:output_label"`
	OutputType          int    `json:"output_type" gorm:"column:output_type"`
	OutputTypeLabel     string `json:"output_type_label" gorm:"column:output_type_label"`
	RetryTime           string `json:"retry_time" gorm:"column:retry_time"`
	RetryCount          int    `json:"retry_count" gorm:"column:retry_count"`
	RetryStatus         int    `json:"retry_status" gorm:"column:retry_status"`
	RetryStatusLabel    string `json:"retry_status_label" gorm:"column:retry_status_label"`
	ProgressStatus      int    `json:"progress_status" gorm:"column:progress_status"`
	ProgressStatusLabel string `json:"progress_status_label" gorm:"column:progress_status_label"`
	Err                 string `json:"err" gorm:"column:err_info"`
	Data                string `json:"data" gorm:"column:data"`
}

func ConvertToNodeInfo(vo *videocollection.ConfigurationVO) NodeJobInfo {
	m := map[string]interface{}{
		"ffmpeg_input_path":      vo.FFmpegInputPath,
		"ffmpeg_output_path":     vo.FFmpegOutputPath,
		"ffmpeg_raw_input_args":  vo.FFmpegRawInputArgs,
		"ffmpeg_raw_output_args": vo.FFmpegRawOutputArgs,
		"ffmpeg_pipe_buff_size":  vo.FFmpegPipeBuffSize,
	}
	marshal, _ := json.Marshal(m)

	return NodeJobInfo{
		Name:                vo.Name,
		Input:               int(vo.Input),
		InputLabel:          vo.InputLabel,
		Output:              int(vo.Output),
		OutputLabel:         vo.OutputLabel,
		OutputType:          int(vo.OutputType),
		OutputTypeLabel:     vo.OutputTypeLabel,
		RetryTime:           vo.RetryTime,
		RetryCount:          vo.RetryCount,
		RetryStatus:         int(vo.RetryStatus),
		RetryStatusLabel:    vo.RetryStatusLabel,
		ProgressStatus:      int(vo.ProgressStatus),
		ProgressStatusLabel: vo.ProgressStatusLabel,
		Err:                 vo.Err,
		Data:                string(marshal),
	}
}
