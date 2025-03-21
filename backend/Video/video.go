package video

import (
	"fmt"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type Video struct {
	Input  string
	Output string
}

func (v *Video) Segment() bool {
	err := ffmpeg_go.Input(v.Input).Output(v.Output+"/"+"segment_%04d.ts", ffmpeg_go.KwArgs{
		"segment_list":         v.Output + "/" + "playlist.m3u8",
		"segment_time":         "10",
        "c": "copy",
        "map": "0",
        "f": "segment",
	}).Run()

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
