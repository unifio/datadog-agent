// +build !windows

package main

import (
	_ "expvar"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/DataDog/datadog-agent/cmd/agent/app"
	"github.com/DataDog/datadog-agent/pkg/util/k8s"
)

func main() {
	// go_expvar server
	go http.ListenAndServe("127.0.0.1:5000", http.DefaultServeMux)

	kubeUtil, err := k8s.NewKubeUtil()
	if err != nil {
		fmt.Println("Err NewKubeUtil: ", err)
		os.Exit(-1)
	}

	data, err := k8s.PerformKubeletQuery("192.168.99.100:10255/pods")
	if err != nil {
		fmt.Println("Err GetNodeInfo: ", err)
	} else {
		fmt.Printf("pods: %s", data)
	}

	podList, err := kubeUtil.GetLocalPodList()
	if err != nil {
		fmt.Println("Err GetNodeInfo: ", err)
	} else {
		fmt.Printf("podList: %+v", podList)
	}

	ip, name, err := kubeUtil.GetNodeInfo()
	if err != nil {
		fmt.Println("Err GetNodeInfo: ", err)
	} else {
		fmt.Println("ip: ", ip, "; name: ", name)
	}

	os.Exit(0)

	// Invoke the Agent
	if err := app.AgentCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
