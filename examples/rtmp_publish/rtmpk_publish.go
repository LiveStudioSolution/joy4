package main

import (
	"github.com/LiveStudioSolution/joy4/av/pktque"
	"github.com/LiveStudioSolution/joy4/format"
	"github.com/LiveStudioSolution/joy4/av/avutil"
	//"github.com/LiveStudioSolution/joy4/format/rtmp"
	rtmp "github.com/LiveStudioSolution/joy4/format/rtmpk"
	"flag"
	"fmt"
	"os"
)

func init() {
	format.RegisterAll()
}

// as same as: ffmpeg -re -i projectindex.flv -c copy -f flv rtmp://localhost:1936/app/publish

var srcUrl *string
var dstUrl *string

func parseFlags(){
	srcUrl = flag.String("srcUrl", "", "input rtmp/file url")
	dstUrl = flag.String("dstUrl", "", "publish to rtmp/rtmpk url")
	flag.Parse()
}

func main() {
	parseFlags()
	fmt.Printf("srcUrl = %s, dstUrl =%s\n", *srcUrl, *dstUrl)
	if *srcUrl == "" || *dstUrl == ""{
		fmt.Fprintf(os.Stderr,"invalid argument\n")
		flag.Usage()
		return
	}
	//file, _ := avutil.Open("/Users/flonly/Music/mp4samples/long.flv")
	src, err := avutil.Open(*srcUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr,"open source error %v\n", err)
		return
	}
	//file, _ := avutil.Open("projectindex.flv")
	//dst, _ := rtmp.Dial("rtmpk://localhost:1936/live/lzh1155")
	dst, err := rtmp.Dial(*dstUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr,"open dest error %v\n", err)
		return
	}

	// conn, _ := avutil.Create("rtmp://localhost:1936/app/publish")
	demuxer := &pktque.FilterDemuxer{Demuxer: src, Filter: &pktque.Walltime{}}
	avutil.CopyFile(dst, demuxer)

	src.Close()
	dst.Close()
}

