package main

import (
	"flag"
	"fmt"
	"time"

	//	contRes "controller/basic-controller" // Importing the module from a different package

	medContRes "controller/medium-controller"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kuberconfig := flag.String("kubeconfig", "/Users/kritiksachdeva/.kube/config", "Locati on to k8s configuration")
	config, err := clientcmd.BuildConfigFromFlags("", *kuberconfig)

	clientcmd.BuildConfigFromFlags()
	if err != nil {
		// Handle error reading the kubeconfig
		fmt.Printf("error %s reading kubeconfig", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Error %s reading config from service account", err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error %s reading clientset", err.Error())
	}
	Informer := informers.NewSharedInformerFactory(clientset, 10*time.Second)

	dynClientset, dynErr := dynamic.NewForConfig(config)

	if dynErr != nil {
		fmt.Printf("Error %s reading dynamiclientset", dynErr.Error())
	}

	DynamicInformer := dynamicinformer.NewDynamicSharedInformerFactory(dynClientset, 10*time.Second)

	dynClientset.Resource()
	//	controller := contRes.NewController(clientset, Informer.Core().V1().Pods())

	controller := medContRes.NewController(Informer.Core().V1().Pods())
	ch := make(chan struct{})
	defer close(ch)

	//Informer.Start(ch)

	// controller.Run(ch2)
	controller.Run(1, ch)
	<-ch // Waiting for input from the channel and until then keep an hold onto the main goroutine function

	fmt.Println(controller)

}
