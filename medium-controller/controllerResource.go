package mediumcontroller

import (
	"fmt"
	"time"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformer "k8s.io/client-go/informers/core/v1"
	corelister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

// Sample controller construct need data objects as: 1) Pod-HasSynced 2) Resource Lister 3) WorkQueue

type Controller struct {
	HasSynced cache.InformerSynced
	Lister    corelister.PodLister
	Queue     workqueue.RateLimitingInterface
}

// Define method to create a controller or creating a contructor

func NewController(pod appsinformer.PodInformer) *Controller {
	c := &Controller{
		HasSynced: pod.Informer().HasSynced,
		Lister:    pod.Lister(),
		Queue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "medium-controller"), // Second argument is the controller name
	}

	// Register Event Handler function with the informer
	pod.Informer().AddEventHandler(
		&cache.ResourceEventHandlerFuncs{
			AddFunc:    c.HandleAdd, // As these are mthods not simple functions hence need to call them "c." instead of calling them directly
			UpdateFunc: c.HandleUpdate,
			DeleteFunc: c.HanldeDelete,
		},
	)
	return c
}

// Now we need define the Run method to start the controller

func (c *Controller) Run(threads int, stopch <-chan struct{}) { // Arguments are no of worker threads and Channel for input to stop controller or not
	// don't let panics crash the process
	defer utilruntime.HandleCrash()

	// make sure the work queue is shutdown which will trigger workers to end
	defer c.Queue.ShutDown()

	klog.Info("Starting medium Controller")

	// Wait for cache to get synced up or wait for your secondary caches to fill before starting your work

	if !cache.WaitForCacheSync(stopch, c.HasSynced) {
		klog.Info("Waiting for cache to get synced up")
		return
	}

	// start up your worker threads based on threadiness.  Some controllers
	// have multiple kinds of workers
	for i := 0; i < threads; i++ {
		// runWorker will loop until "something bad" happens.  The .Until will
		// then rekick the worker after one second
		go wait.Until(c.runWorker, time.Second, stopch)
	}

	// wait until we're told to stop
	<-stopch
	klog.Infof("Shutting down medium controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// This will execute the method continuously in an infinite loop until the method processNextItem() returns false
	}
}

func (c *Controller) processNextItem() bool { // Here return type bool is must, because we are using this return as conditional for the for loop function
	// pull the next work item from queue.  It should be a key we use to lookup
	// something in a cache
	item, quit := c.Queue.Get()
	if quit {
		return false
	}

	defer c.Queue.Done(item) // Here we saying that we have done/completed processing on this key in the queue

	// This key is of type namespace/resource-name, and now we need to use this to get the object details from lister

	key, err := cache.MetaNamespaceKeyFunc(item) // It will return the namespace & name as key from that API resource obj

	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid cache item: %s", key))
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		//		return nil
	}

	// Now we have the namespace & resource name, next we would need to get the resource details
	pod, err := c.Lister.Pods(ns).Get(name)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("error getting resource from lister %s", err.Error()))
	}

	fmt.Printf("Pod is %s", pod.Name)
	return true
}
