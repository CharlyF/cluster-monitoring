package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/CharlyF/cluster-monitoring/pkg/aggregator"
	"github.com/CharlyF/cluster-monitoring/util"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	podInformer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	podLister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type PodController struct {
	podLister podLister.PodLister
	podListerSync cache.InformerSynced
	queue workqueue.RateLimitingInterface
	clientSet kubernetes.Interface
}

func NewPodController(client kubernetes.Interface,inf podInformer.PodInformer) (*PodController, error) {
	p := &PodController{
		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultItemBasedRateLimiter(), "pod"),
	}
	p.clientSet = client

	inf.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: p.addPod,
			UpdateFunc: p.updatePod,
			DeleteFunc: p.deletePod,
		})
	p.podLister = inf.Lister()
	p.podListerSync = inf.Informer().HasSynced

	return p, nil
}

func (p *PodController) Run(stop <- chan struct{}) {
	defer p.queue.ShutDown()
	fmt.Printf("Starting runner for PodController \n")
	defer fmt.Printf("Stopping PodController \n")
	if !cache.WaitForCacheSync(stop, p.podListerSync) {
		return
	}
	go wait.Until(p.worker, time.Second, stop)
	<-stop

}

func (p *PodController) addPod(obj interface{}) {
	newPod, ok := obj.(*v1.Pod)
	if !ok {
		fmt.Printf("not OK adding %v \n", obj)
		return
	}
	//fmt.Printf("adding pod %v\n", newPod.Name)
	p.enqueue(newPod)
}

func (p *PodController) updatePod(oldObj interface{}, newObj interface{})  {
	newPod, ok := newObj.(*v1.Pod)
	if !ok {
		return
	}
	oldPod, ok := oldObj.(*v1.Pod)
	if !ok {
		return
	}
	fmt.Printf("updating pod %s to %s \n", oldPod.Name, newPod.Name)
	// remove old from cache ?
	p.enqueue(newPod)
}
func (p *PodController) deletePod(obj interface{}) {
	deadPod, ok := obj.(*v1.Pod)
	if !ok {
		return
	}
	// remove from the DB ?
	podIp := deadPod.Status.PodIP
	fmt.Printf("deleting %v", deadPod.Name)
	if podIp != ""{
		util.Cache.Delete(podIp)
	}
	p.queue.Done(deadPod)
}

func (p *PodController) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	p.queue.AddRateLimited(key)
}

func (p *PodController) worker() {
	for p.processNext() {
	}
}

func (p *PodController) processNext() bool{
	key, quit := p.queue.Get()
	if quit {
		return false
	}
	defer p.queue.Done(key)

	pod, cached := p.isPodCached(key)
    //ctnMeta := p.getContainerMeta(pod.Status.ContainerStatuses[0].ContainerID)
	// add ctnMeta to the cache.
	//fmt.Println(ctnMeta)
	if !cached && pod != nil {
		// Cache the rigth values.
		err := addToCache(pod)
		if err != nil {
			fmt.Errorf(err.Error())
			return false
		}
	}
	return true
}

func addToCache(pod *v1.Pod) error{
	pm := &aggregator.PodMetadata{}
	pm.Name = pod.Name
	pm.UID = string(pod.UID)
	pmJson, err := json.Marshal(pm)
	if err != nil {
		return err
	}
	return util.Cache.Add(pod.Status.PodIP, string(pmJson), 0)
}

func (p *PodController) isPodCached(key interface{}) (*v1.Pod, bool){

	ns, name, err := cache.SplitMetaNamespaceKey(key.(string))
	if err != nil {
		return nil, false
	}
	fmt.Printf("processing %s/%s \n", ns, name)
	// check if in cache here
	// if not, fetch from informer cache, get data
	// IP: { Docker: {"dockerLabels": ["foo:bar"]}, Kubernetes: {"replicaset": "foofoo"}}
	pod, err := p.podLister.Pods(ns).Get(name)
	if err != nil {
		return nil, false
	}

	fmt.Printf("client-go cache pod: \n- Name: %s \n- IP: %s", pod.Name, pod.Status.PodIP)

	podCached, cached := util.Cache.Get(pod.Status.PodIP)
	if !cached {
		return pod, false
	}
	model := &aggregator.PodMetadata{}
	err = json.Unmarshal([]byte(podCached.(string)), model)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Pod %s retrieved from cache: %s \n", model.Name, cached)

	return pod,cached
}

//func (p *PodController) getContainerMeta(cId string) *types.ContainerJSONBase {
//	ctn, err := p.dockerUtil.Inspect(cId, false)
//	if err != nil {
//		return nil
//	}
//	return ctn.ContainerJSONBase
//}
