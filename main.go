package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"sync"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/format/rtmp"
)

// Config stores the configuration for the RTMP server
type Config struct {
	FirstStreamKey  string `json:"FirstStreamKey"`
	SecondStreamKey string `json:"SecondStreamKey"`
	OutputStreamURL string `json:"OutputStreamURL"`
}

var config Config
var lock sync.Mutex
var primaryConn *rtmp.Conn
var backupConn *rtmp.Conn

func main() {
	loadConfig()

	format.RegisterAll()

	server := &rtmp.Server{}

	server.HandlePublish = func(conn *rtmp.Conn) {
		streamKey := conn.URL.Path
		log.Printf("Stream started: %s", streamKey)

		lock.Lock()
		defer lock.Unlock()

		if streamKey != config.FirstStreamKey && streamKey != config.SecondStreamKey {
			log.Println("Invalid stream key. Rejecting connection.")
			conn.Close()
			return
		}

		if streamKey == config.FirstStreamKey {
			if primaryConn != nil {
				log.Println("Primary stream reconnected. Stopping backup stream.")
				backupConn.Close()
			}
			primaryConn = conn
			if backupConn != nil {
				backupConn.Close()
				backupConn = nil
			}
		} else {
			backupConn = conn
			if primaryConn != nil {
				log.Println("Backup stream connected. Stopping backup.")
				conn.Close()
				return
			}
		}

		go pushToExternalRTMP(conn, streamKey)
	}

	log.Printf("Starting RTMP server on :1935")
	log.Fatal(server.ListenAndServe())
}

func pushToExternalRTMP(conn *rtmp.Conn, streamKey string) {
	dstURL := config.OutputStreamURL

	dst, err := rtmp.Dial(dstURL)
	if err != nil {
		log.Printf("Failed to connect to external server: %v", err)
		conn.Close()
		return
	}
	defer dst.Close()

	streams, err := conn.Streams()
	if err != nil {
		log.Fatalf("Failed to get streams: %v", err)
		return
	}

	err = dst.WriteHeader(streams)
	if err != nil {
		log.Fatalf("Failed to write header: %v", err)
		return
	}
	defer dst.WriteTrailer()

	for {
		packet, err := conn.ReadPacket()
		if err != nil {
			log.Printf("Error reading packet: %v", err)
			if err == io.EOF {
				log.Printf("Stream stopped: %s", streamKey)
			}
			break
		}
		if err = dst.WritePacket(packet); err != nil {
			log.Printf("Failed to write packet to external server: %v", err)
			break
		}
	}

	lock.Lock()
	defer lock.Unlock()

	if primaryConn == conn {
		primaryConn = nil // Reset primary connection reference
		log.Printf("Primary stream disconnected: %s", streamKey)
	} else if backupConn == conn {
		backupConn = nil // Reset backup connection reference
		log.Printf("Backup stream disconnected: %s", streamKey)
	}
}

func loadConfig() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}
