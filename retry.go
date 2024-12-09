package videosys_starter

import (
	videocollection "github.com/kordar/video-collection"
	"github.com/spf13/cast"
	"strings"
)

func toIntSlice(str string) []int {
	data := make([]int, 0)
	if str == "" {
		return data
	}
	items := strings.Split(str, ",")
	for _, item := range items {
		v := cast.ToInt(item)
		data = append(data, v)
	}
	return data
}

func retryHandle(cfg map[string]string) videocollection.Retry {
	if cfg["retry"] == "default" {
		retryMaxTimes := cast.ToInt(cfg["retry_max_times"])
		waitSeconds := toIntSlice(cfg["retry_wait_seconds"])
		return &videocollection.DefaultRetry{
			MaxTimes:    retryMaxTimes,
			WaitSeconds: waitSeconds,
		}
	}
	return nil
}
