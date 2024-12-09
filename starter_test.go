package videosys_starter_test

import (
	goframeworkvideosys "github.com/kordar/goframework-videosys"
	logger "github.com/kordar/gologger"
	videocollection "github.com/kordar/video-collection"
	videosysstarter "github.com/kordar/videosys-starter"
	"github.com/xfrr/goffmpeg/transcoder"
	"testing"
	"time"
)

var cfg = map[string]interface{}{
	"id":                     "xxx",
	"ffmpeg_input_path":      "rtsp://admin:a1234567@192.168.10.67:554/h264/ch1/sub/av_stream",
	"output_type":            "0",
	"ffmpeg_raw_input_args":  "-re",
	"ffmpeg_raw_output_args": "-r 1 -q:v 2",
}

func TestNewStreamModule(t *testing.T) {
	module := videosysstarter.NewStreamModule("test", func(moduleName string, itemId string, configuration videocollection.Configuration, retry videocollection.Retry, item map[string]string) {
		collection := &videocollection.FFmpegCollection{
			Running: func(value transcoder.Progress, cfg *videocollection.Configuration, collect *videocollection.FFmpegCollection) {
			},
			ExecPipe: func(buff []byte, cfg *videocollection.Configuration) {
				logger.Infof("==========>>>%v", len(buff))
			},
		}
		_ = goframeworkvideosys.AddStreamInstance(itemId, collection, &configuration, retry)
	})
	module.Load(cfg)

	//list := goframeworkvideosys.ConfigList()
	//logger.Infof("----------%+v", list)

	go func() {
		time.Sleep(10 * time.Second)
		goframeworkvideosys.Stop("xxx")
		list := goframeworkvideosys.ConfigList()
		logger.Infof("----------%+v", list)
	}()

	time.Sleep(100 * time.Second)
}
