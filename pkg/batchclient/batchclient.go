/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package batchclient

import (
	"github.com/kubernetes-sigs/poseidon/pkg/client/clientset/versioned"
	"github.com/kubernetes-sigs/poseidon/pkg/firmament"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/golang/glog"
	config2 "github.com/kubernetes-sigs/poseidon/pkg/config"
	//kbver "github.com/kubernetes-sigs/poseidon/pkg/client/clientset/versioned"
)

var ClientSet versioned.Clientset

// GetClientConfig returns a kubeconfig object which to be passed to a Kubernetes client on initialization.
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

// New initializes a firmament and Kubernetes client and starts watching Pod and Node.
func New(kubeConfig string, firmamentAddress string) {

	glog.Info("First Line: batch newclient called")
	config, err := GetClientConfig(kubeConfig)
	if err != nil {
		glog.Fatalf("Failed to load client config: %v", err)
	}
	config.QPS = config2.GetQPS()
	config.Burst = config2.GetBurst()

	ClientSet, err := versioned.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create connection: %v", err)
	}

	fc, conn, err := firmament.New(firmamentAddress)
	if err != nil {
		glog.Fatalf("Failed to connect to Firmament: %v", err)
	}
	defer conn.Close()

	glog.Info("batch newclient called")
	stopCh := make(chan struct{})

	qw, qInformer := NewQueueWatcher(ClientSet, fc)

	glog.Info("++++QW set for q watcher +%v +%v", qInformer.Scheduling().V1alpha1().Queues().Informer().AddEventHandler)

	go qInformer.Start(stopCh)

	go qw.Run(stopCh, 10)

	pgw, pgInformer := NewPodGroupWatcher(ClientSet, fc)

	go pgInformer.Start(stopCh)

	go pgw.Run(stopCh, 10)

	<-stopCh

}

func init() {

	glog.Info("batchclient init called")
}
