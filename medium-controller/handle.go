package mediumcontroller

import (
	"fmt"

	"k8s.io/client-go/tools/cache"
)

// Make these as methods of the controller for which we add them in the Queue

func (c *Controller) HandleAdd(obj interface{}) {
	fmt.Println("Add was called")
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err == nil {
		c.Queue.Add(key)
	}
}

func (c *Controller) HandleUpdate(oldobj, newobj interface{}) {
	fmt.Println("Update was called")
	key, err := cache.MetaNamespaceKeyFunc(newobj) // Here we are adding the key in the queue which is of the format namespace-resource
	if err == nil {
		c.Queue.Add(key)
	}
}

func (c *Controller) HanldeDelete(obj interface{}) {
	fmt.Println("Delete was called")
	// IndexerInformer uses a delta nodeQueue, therefore for deletes we have to use this
	// key function.
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err == nil {
		c.Queue.Add(key)
	}
}
