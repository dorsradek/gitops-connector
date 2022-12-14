package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/fluxcd/pkg/runtime/events"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	var (
		dryRun        bool
		labelSelector string
	)

	flags := flag.NewFlagSet("args", flag.ExitOnError)

	flags.BoolVar(&dryRun, "dry-run", false, "Dry-run")
	flags.StringVar(&labelSelector, "label-selector", "", "Label selector")

	if err := flags.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	flags.VisitAll(func(f *flag.Flag) {
		log.Infof("%s: %s\n", f.Name, f.Value)
	})

	config, err := rest.InClusterConfig()
	if err != nil {
		kubeConfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			panic(err)
		}
		log.Infoln("Running from outside of the cluster")
	} else {
		log.Infoln("Running inside the cluster")
	}

	_, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	log.Infoln("Hello world1")
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8081", nil)
}

func hello(w http.ResponseWriter, req *http.Request) {
	log.Infoln("hello")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var e events.Event
	err = json.Unmarshal(body, &e)
	if err != nil {
		panic(err)
	}
	log.Infoln(e.Message)
	log.Infoln(e.Reason)
	log.Infoln(e.Timestamp)
}
