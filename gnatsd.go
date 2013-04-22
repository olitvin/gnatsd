// Copyright 2012 Apcera Inc. All rights reserved.

package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/apcera/gnatsd/server"
)

func main() {
	// logging setup
	server.LogSetup()

	opts := server.Options{}

	var debugAndTrace bool
	var configFile string

	// Parse flags
	flag.IntVar(&opts.Port, "port", server.DEFAULT_PORT, "Port to listen on.")
	flag.IntVar(&opts.Port, "p", server.DEFAULT_PORT, "Port to listen on.")
	flag.StringVar(&opts.Host, "host", server.DEFAULT_HOST, "Network host to listen on.")
	flag.StringVar(&opts.Host, "h", server.DEFAULT_HOST, "Network host to listen on.")
	flag.StringVar(&opts.Host, "net", server.DEFAULT_HOST, "Network host to listen on.")
	flag.BoolVar(&opts.Debug, "D", false, "Enable Debug logging.")
	flag.BoolVar(&opts.Debug, "debug", false, "Enable Debug logging.")
	flag.BoolVar(&opts.Trace, "V", false, "Enable Trace logging.")
	flag.BoolVar(&opts.Trace, "trace", false, "Enable Trace logging.")
	flag.BoolVar(&debugAndTrace, "DV", false, "Enable Debug and Trace logging.")
	flag.StringVar(&opts.Username, "user", "", "Username required for connection.")
	flag.StringVar(&opts.Password, "pass", "", "Password required for connection.")
	flag.StringVar(&opts.Authorization, "auth", "", "Authorization token required for connection.")
	flag.IntVar(&opts.HttpPort, "m", 0, "HTTP Port for /varz, /connz endpoints.")
	flag.IntVar(&opts.HttpPort, "http_port", 0, "HTTP Port for /varz, /connz endpoints.")
	flag.StringVar(&configFile, "c", "", "Configuration file.")
	flag.StringVar(&configFile, "config", "", "Configuration file.")

	flag.Parse()

	if debugAndTrace {
		opts.Trace, opts.Debug = true, true
	}

	var fileOpts *Options

	// Parse config if given
	if configFile != "" {
		fileOpts = server.processConfigFile(configFile)
	}

	// Create the server with appropriate options.
	s := server.New(server.mergeOptions(fileOpts, &opts))

	// Start up the http server if needed.
	if opts.HttpPort != 0 {
		s.StartHTTPMonitoring()
	}

	// Profiler
	go func() {
		log.Println(http.ListenAndServe("localhost:6062", nil))
	}()

	// Wait for clients.
	s.AcceptLoop()
}
