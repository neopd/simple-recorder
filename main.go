package main

import (
	"errors"
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/nats-io/nats.go"
)

type Packet struct {
	id   int
	size int
	data []byte
}

var rxq chan Packet
var done chan bool

const MAX_RXQ_COUNT int = 5

// This example shows how to connect to a RTSP server
// and read all tracks on a path.

// format: area.%d.cam.%d.%d
func suject2id(subject *string) (area int, cam int, ch int, err error) {
	slices := strings.Split(*subject, ".")
	if len(slices) != 5 {
		err = errors.New("subject should have 5 fields seperated by '.'")
		return
	}
	if slices[0] != "area" {
		err = errors.New("subject need 'area' keywoard at the first")
		return
	}
	if slices[2] != "cam" {
		err = errors.New("subject need 'cam' keywoard at the third")
		return
	}
	area, err = strconv.Atoi(slices[1])
	if err != nil {
		return
	}

	cam, err = strconv.Atoi(slices[3])
	if err != nil {
		return
	}

	ch, err = strconv.Atoi(slices[4])
	if err != nil {
		return
	}

	err = nil
	return
}

var rxCount int

func onNATSMessage(m *nats.Msg) {
	area, cam, ch, err := suject2id(&m.Subject)
	if err != nil {
		log.Fatal(err)
	}

	if rxCount%5000 == 0 {
		log.Printf("Got data len: %d from area.%d.cam.%d.%d [%s] total count=%d", len(m.Data), area, cam, ch, m.Subject, rxCount)
	}
	rxCount++
}

func main() {

	var rxqCount = flag.Int("txq", MAX_RXQ_COUNT, "Max rx-q count")
	var natsAddr = flag.String("nats", "172.20.10.120:4222", "NATS Address")
	var subject = flag.String("subject", "area.*.cam.*.*", "NATS Subject")
	flag.Parse()

	log.Printf("%s/subject=%s", *natsAddr, *subject)

	rxq = make(chan Packet, *rxqCount)
	done = make(chan bool)

	nc, err := nats.Connect(*natsAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	if _, err := nc.Subscribe(*subject, onNATSMessage); err != nil {
		log.Fatal(err)
	}

	<-done

}
