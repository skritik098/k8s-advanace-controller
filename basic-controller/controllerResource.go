package controller

import (
	"fmt"

	appsinformer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corelister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// In OOPS prgramming world, we have 3 types of classes files. 1) Main class file 2) Resource File 3) Resource Service file

// In this file, we will be creating the Resource file

type Controller struct { // For controller, at least we need an informer, clientset, lister, queue
	Clientset   kubernetes.Interface
	PodInformer appsinformer.PodInformer
	PodLister   corelister.PodLister
	Queue       workqueue.RateLimitingInterface
}

func NewController(clientset kubernetes.Interface, informer appsinformer.PodInformer) *Controller {
	c := &Controller{
		Clientset:   clientset,
		PodInformer: informer,
		PodLister:   informer.Lister(),
		Queue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "example"),
	}

	c.PodInformer.Informer().AddEventHandler(
		&cache.ResourceEventHandlerFuncs{
			AddFunc:    HandleAdd,
			DeleteFunc: HanldeDelete,
			UpdateFunc: HandleUpdate,
		},
	)
	return c
}

func (c *Controller) Run(ch <-chan struct{}) {
	fmt.Println("Starting the controller")

	c.PodInformer.Informer().Run(ch)

	// Let's first check the cache to be synced Up     <--- This sync of cache is already been handled internally
	// in the build-in Run() method
	/*
			if !cache.WaitForCacheSync(ch, c.PodInformer.Informer().HasSynced) {
				fmt.Println("Syncing the cachec and Waiting for cache to be synced")
				log.Errorf("Timed out waiting for caches to sync")
			}
		fmt.Println("Cache get synced & running the Informer")
	*/
}
