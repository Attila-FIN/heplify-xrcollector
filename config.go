package main

var cfg config

type config struct {
	HttpServerAddress string
	CollectorAddress  string
	HepNodeID         uint
	Debug             bool
}
