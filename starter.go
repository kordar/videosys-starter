package videosys_starter

import (
	goframeworkvideosys "github.com/kordar/goframework-videosys"
	logger "github.com/kordar/gologger"
	videocollection "github.com/kordar/video-collection"
	"github.com/spf13/cast"
	"strings"
)

type Load func(moduleName string, itemId string, configuration videocollection.Configuration, retry videocollection.Retry, item map[string]string)

type StreamModule struct {
	name string
	load Load
}

func NewStreamModule(name string, load Load) *StreamModule {
	return &StreamModule{name, load}
}

func (m StreamModule) Name() string {
	return m.name
}

func (m StreamModule) _load(id string, cfg map[string]string) {
	if id == "" {
		logger.Fatalf("[%s] the attribute id cannot be empty.", m.Name())
		return
	}

	if m.load == nil {
		logger.Fatalf("[%s] the load function required.", m.Name())
		return
	}

	configuration := videocollection.Configuration{
		Input:              videocollection.Input(cast.ToInt(cfg["input"])),
		Output:             videocollection.Output(cast.ToInt(cfg["output"])),
		OutputType:         videocollection.OutputFormat(cast.ToInt(cfg["output_type"])),
		FFmpegInputPath:    cfg["ffmpeg_input_path"],
		FFmpegOutputPath:   cfg["ffmpeg_output_path"],
		FFmpegPipeBuffSize: cast.ToInt(cfg["ffmpeg_pipe_buff_size"]),
	}

	if cfg["name"] != "" {
		configuration.Name = cfg["name"]
	} else {
		configuration.Name = id
	}

	if cfg["ffmpeg_raw_input_args"] != "" {
		configuration.FFmpegRawInputArgs = strings.Split(cfg["ffmpeg_raw_input_args"], " ")
	}

	if cfg["ffmpeg_raw_output_args"] != "" {
		configuration.FFmpegRawOutputArgs = strings.Split(cfg["ffmpeg_raw_output_args"], " ")
	}

	retry := retryHandle(cfg)
	// TODO please insert instance of videocollection.Collection into goframeworkvideosys
	m.load(m.name, id, configuration, retry, cfg)

	go func() {
		if err := goframeworkvideosys.Start(id); err == nil {
			logger.Infof("[%s] start stream '%s' successfully", m.Name(), id)
		} else {
			logger.Warnf("[%s] start stream '%s' failed", m.Name(), id)
		}
	}()

	logger.Infof("[%s] loading module '%s' finished", m.Name(), id)
}

func (m StreamModule) Load(value interface{}) {

	items := cast.ToStringMap(value)
	if items["id"] != nil {
		id := cast.ToString(items["id"])
		m._load(id, cast.ToStringMapString(value))
		return
	}

	for key, item := range items {
		m._load(key, cast.ToStringMapString(item))
	}

}

func (m StreamModule) Close() {
}
