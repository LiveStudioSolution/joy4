package main

import (
	"sync"
	"io"
	"net/http"
	"github.com/LiveStudioSolution/joy4/format"
	"github.com/LiveStudioSolution/joy4/av/avutil"
	"github.com/LiveStudioSolution/joy4/av/pubsub"
	//"github.com/LiveStudioSolution/joy4/format/rtmp"
	rtmp "github.com/LiveStudioSolution/joy4/format/rtmpk"
	"github.com/LiveStudioSolution/joy4/format/flv"
	"fmt"
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
	server := &rtmp.Server{}
	server.Addr = "0.0.0.0:1936"
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
			fmt.Printf("HandlePublish rtmpk.Conn %+v", conn)
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
		fmt.Printf("HandlePublish rtmpk.Conn %+v", conn)
		avutil.CopyPackets(ch.que, conn)

		l.Lock()
		delete(channels, conn.URL.Path)
		l.Unlock()
		ch.que.Close()
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		l.RLock()
		ch := channels[r.URL.Path]
		l.RUnlock()

		if ch != nil {
			w.Header().Set("Content-Type", "video/x-flv")
			w.Header().Set("Transfer-Encoding", "chunked")		
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(200)
			flusher := w.(http.Flusher)
			flusher.Flush()

			muxer := flv.NewMuxerWriteFlusher(writeFlusherk{httpflusher: flusher, Writer: w})
			cursor := ch.que.Latest()

			avutil.CopyFile(muxer, cursor)
		} else {
			http.NotFound(w, r)
		}
	})

	go http.ListenAndServe(":8089", nil)

	server.ListenAndServe()

	// ffmpeg -re -i movie.flv -c copy -f flv rtmp://localhost/movie
	// ffmpeg -f avfoundation -i "0:0" .... -f flv rtmp://localhost/screen
	// ffplay http://localhost:8089/movie
	// ffplay http://localhost:8089/screen
}
