// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

// +build !windows

package main

import (
	"bytes"
	"os"
	"runtime/pprof"

	"github.com/DataDog/datadog-agent/cmd/agent/app"
	log "github.com/cihub/seelog"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			defer log.Flush()
			buf := new(bytes.Buffer)
			pprof.Lookup("goroutine").WriteTo(buf, 1)
			log.Errorf("Whoopsies! Agent panicked ಥ_ಥ - stacktraces: %s", buf)
			os.Exit(-1)
		}
	}()
	// Invoke the Agent
	if err := app.AgentCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
