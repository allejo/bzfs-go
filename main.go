package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

const ReplayFileName = "20170701-1926-fun.rec"

type Duration struct {
	days  int
	hours int
	mins  int
	secs  int
	usecs int
}

func calcDuration(timestamp int64) (time Duration) {
	secs := timestamp / 1000000

	time.days = int(secs / (24 * 60 * 60))
	secs %= 24 * 60 * 60

	time.hours = int(secs / (60 * 60))
	secs %= 60 * 60

	time.mins = int(secs / 60)
	secs %= 60

	time.secs = int(secs)
	time.usecs = int(timestamp % 1000000)

	return
}

func main() {
	dat, _ := ioutil.ReadFile(ReplayFileName)
	buf := bytes.NewBuffer(dat)

	header := loadReplayHeader(buf)
	timeStruct := calcDuration(header.filetime)

	p, err := loadReplayPacket(buf)

	fmt.Printf("magic:     0x%04X\n", header.magic)
	fmt.Printf("replay:    version %d\n", header.version)
	fmt.Printf("offset:    %d\n", header.offset)
	fmt.Printf("length:    %-d days, %d hours, %d minutes, %d seconds, %d usecs\n", timeStruct.days, timeStruct.hours, timeStruct.mins, timeStruct.secs, timeStruct.usecs)
	fmt.Printf("start:     %s\n", time.Unix(p.timestamp/1000000, 0))
	fmt.Printf("end:       %s\n", time.Unix((p.timestamp+header.filetime)/1000000, 0))
	fmt.Printf("author:    %s (%s)\n", header.callSign, header.motto)
	fmt.Printf("bzfs:      bzfs-%s\n", header.appVersion)
	fmt.Printf("protocol:  %.8s\n", header.ServerVersion)
	fmt.Printf("flagSize:  %d\n", header.flagsSize)
	fmt.Printf("worldSize: %d\n", header.worldSize)
	fmt.Printf("worldHash: %s\n", header.realHash)
	fmt.Printf("\n")

	for err == nil {
		var packet interface{}

		switch p.mode {
		case RealPacket:
			fmt.Printf("Real Packet\n")
			break

		case StatePacket:
			fmt.Printf("State Packet\n")
			packet = packetDataToStruct(p.code, p.data)
			break

		case UpdatePacket:
			fmt.Printf("Update Packet\n")
			break
		}

		marshaled, _ := json.Marshal(packet)
		fmt.Println(string(marshaled))

		// Move on to the next packet
		p, err = loadReplayPacket(buf)
	}
}
