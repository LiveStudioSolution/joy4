package format

import (
	"github.com/LiveStudioSolution/joy4/format/mp4"
	"github.com/LiveStudioSolution/joy4/format/ts"
	"github.com/LiveStudioSolution/joy4/format/rtmp"
	"github.com/LiveStudioSolution/joy4/format/rtsp"
	"github.com/LiveStudioSolution/joy4/format/flv"
	"github.com/LiveStudioSolution/joy4/format/aac"
	"github.com/LiveStudioSolution/joy4/av/avutil"
	"github.com/LiveStudioSolution/joy4/format/rtmpk"
)

func RegisterAll() {
	avutil.DefaultHandlers.Add(mp4.Handler)
	avutil.DefaultHandlers.Add(ts.Handler)
	avutil.DefaultHandlers.Add(rtmp.Handler)
	avutil.DefaultHandlers.Add(rtmpk.Handler)
	avutil.DefaultHandlers.Add(rtsp.Handler)
	avutil.DefaultHandlers.Add(flv.Handler)
	avutil.DefaultHandlers.Add(aac.Handler)
}

