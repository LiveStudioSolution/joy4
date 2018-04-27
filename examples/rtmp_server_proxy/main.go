package main

import (
	"github.com/LiveStudioSolution/joy4/format"
	"github.com/LiveStudioSolution/joy4/av/avutil"
	"github.com/LiveStudioSolution/joy4/format/rtmp"
	"fmt"
	"strings"
)

func init() {
	format.RegisterAll()
}

func main() {
	server := &rtmp.Server{}
	server.Addr ="0.0.0.0:1936"

	server.HandlePlay = func(conn *rtmp.Conn) {
		segs := strings.Split(conn.URL.Path, "/")
		//url := fmt.Sprintf("%s://%s", segs[1], strings.Join(segs[2:], "/"))
		url := fmt.Sprintf("rtmp://127.0.0.1/%s/%s", segs[3], segs[4])
		fmt.Printf("proxy to ")
		//url := "rtmp://127.0.0.1/live/lzh1155"
		src, _ := avutil.Open(url)
		src, _ := avutil.OpenK(url)
		avutil.CopyFile(conn, src)
		// srs --> proxy --> ffplay
	}

	// joy -rtmp -> proxy -rtmp-> srs|nginx -->
	// publish -->

	server.ListenAndServe()

	// ffplay rtmp://localhost/rtsp/192.168.1.1/camera1
	// ffplay rtmp://localhost/rtmp/live.hkstv.hk.lxdns.com/live/hks
	// joy listen 1935
	// proxy listen 1936
	// publish 127.0.0.1 1935
	// ffplay rtmp://127.0.01:1936/live/lzh1155
}
