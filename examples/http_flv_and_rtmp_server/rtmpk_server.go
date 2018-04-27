package main

import (
	"sync"
	"io"
	"net/http"
	"github.com/LiveStudioSolution/joy4/format"
	"github.com/LiveStudioSolution/joy4/av/avutil"
	"github.com/LiveStudioSolution/joy4/av/pubsub"
	"github.com/LiveStudioSolution/joy4/format/rtmp"
	"github.com/LiveStudioSolution/joy4/format/rtmpk"
	"fmt"
	"flag"
)

func init() {
	format.RegisterAll()
}

type writeFlusherk struct {
	httpflusher http.Flusher
	io.Writer
}

func (self writeFlusherk) Flush() error {
	self.httpflusher.Flush()
	return nil
}

func main() {
	protocol := flag.String("protocol", "rtmpk", "rtmp or rtmpk")
	addr := flag.String("addr", ":1936", " server bind address, default :1936")
	flag.Parse()
	switch *protocol {
	case "rtmp":
		runRtmpServer(*addr)
	case "rtmpk":
		runRtmpKServer(*addr)
	default:
		fmt.Errorf("unkonown protocol %s", protocol)
	}
}

func runRtmpKServer(addr string) {
	server := &rtmpk.Server{}
	//server.Addr = "0.0.0.0:1936"
	server.Addr = addr
	rtmpk.Debug = true

	l := &sync.RWMutex{}
	type Channel struct {
		que *pubsub.Queue
	}
	channels := map[string]*Channel{}

	server.HandlePlay = func(conn *rtmpk.Conn) {
		l.RLock()
		ch := channels[conn.URL.Path]
		l.RUnlock()

		if ch != nil {
			fmt.Printf("HandlePublish rtmpk.Conn %+v", conn)
			cursor := ch.que.Latest()
			avutil.CopyFile(conn, cursor)
		}
	}

	server.HandlePublish = func(conn *rtmpk.Conn) {
		streams, _ := conn.Streams()

		l.Lock()
		ch := channels[conn.URL.Path]
		if ch == nil {
			ch = &Channel{}
			ch.que = pubsub.NewQueue()
			ch.que.WriteHeader(streams)
			channels[conn.URL.Path] = ch
		} else {
			ch = nil
		}
		l.Unlock()
		if ch == nil {
			return
		}
		fmt.Printf("HandlePublish rtmpk.Conn %+v", conn)
		avutil.CopyPackets(ch.que, conn)

		l.Lock()
		delete(channels, conn.URL.Path)
		l.Unlock()
		ch.que.Close()
	}

	server.ListenAndServe()

	// ffmpeg -re -i movie.flv -c copy -f flv rtmpk://localhost/movie
}

func runRtmpServer(addr string) {
	server := &rtmp.Server{}
	//server.Addr = "0.0.0.0:1936"
	server.Addr = addr
	rtmp.Debug = true

	l := &sync.RWMutex{}
	type Channel struct {
		que *pubsub.Queue
	}
	channels := map[string]*Channel{}

	server.HandlePlay = func(conn *rtmp.Conn) {
		l.RLock()
		ch := channels[conn.URL.Path]
		l.RUnlock()

		if ch != nil {
			fmt.Printf("HandlePublish rtmp.Conn %+v", conn)
			cursor := ch.que.Latest()
			avutil.CopyFile(conn, cursor)
		}
	}

	server.HandlePublish = func(conn *rtmp.Conn) {
		streams, _ := conn.Streams()

		l.Lock()
		ch := channels[conn.URL.Path]
		if ch == nil {
			ch = &Channel{}
			ch.que = pubsub.NewQueue()
			ch.que.WriteHeader(streams)
			channels[conn.URL.Path] = ch
		} else {
			ch = nil
		}
		l.Unlock()
		if ch == nil {
			return
		}
		fmt.Printf("HandlePublish rtmp.Conn %+v", conn)
		avutil.CopyPackets(ch.que, conn)

		l.Lock()
		delete(channels, conn.URL.Path)
		l.Unlock()
		ch.que.Close()
	}

	server.ListenAndServe()

	// ffmpeg -re -i movie.flv -c copy -f flv rtmp://localhost/movie
}