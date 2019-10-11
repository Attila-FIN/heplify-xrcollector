package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

// XRPacket holds UDP data and address
type XRPacket struct {
	addr *net.UDPAddr
	data []byte
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Use %s like: %s [option]\n", "heplify-xrcollector 0.4", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&cfg.HttpServerAddress, "hs", "http://127.0.0.1:80", "HTTP server address")
	flag.StringVar(&cfg.CollectorAddress, "xs", ":5060", "XR collector UDP listen address")
	flag.UintVar(&cfg.HepNodeID, "hi", 3333, "HEP ID")
	flag.BoolVar(&cfg.Debug, "debug", false, "Log with debug level")
	flag.Parse()
}

func main() {
	addrXR, err := net.ResolveUDPAddr("udp", cfg.CollectorAddress)
	if err != nil {
		log.Fatalln(err)
	}

	connXR, err := net.ListenUDP("udp", addrXR)
	if err != nil {
		log.Fatalln(err)
	}

	// connHEP, err := net.Dial("udp", cfg.HepServerAddress)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	inXRCh := make(chan XRPacket, 100)
	outXRCh := make(chan XRPacket, 100)
	//outHEPCh := make(chan []byte, 100)
	outHTTPCh := make(chan []byte, 100)

	go recvXR(connXR, inXRCh, outHTTPCh)
	go sendXR(connXR, outXRCh)
	//go sendHEP(connHEP, outHEPCh)
	go sendHTTP(cfg.HttpServerAddress, outHTTPCh)

	for packet := range inXRCh {
		outXRCh <- packet
	}
}
