package batchclient

import (
	"fmt"
	"github.com/golang/glog"
	kbv1 "github.com/kubernetes-sigs/poseidon/pkg/apis/scheduling/v1alpha1"
	"github.com/kubernetes-sigs/poseidon/pkg/client/clientset/versioned"
	kbinfo "github.com/kubernetes-sigs/poseidon/pkg/client/informers/externalversions"
	"github.com/kubernetes-sigs/poseidon/pkg/firmament"
	"github.com/kubernetes-sigs/poseidon/pkg/k8sclient"
	//v1 "k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/runtime
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	"k8s.io/client-go/tools/cache"

	"time"
)

// NewQueueWatcher initializes a QueueWatcher.
func NewQueueWatcher(client versioned.Interface, fc firmament.FirmamentSchedulerClient) (*QueueWatcher, kbinfo.SharedInformerFactory) {
	glog.V(2).Info("Starting QueueWatcher...")
	QueueToQD = make(map[string]*firmament.QueueDescriptor)
	kbinformer := kbinfo.NewSharedInformerFactory(client, 0)

	// create informer for Queue information
	queueInformer := kbinformer.Scheduling().V1alpha1().Queues().Informer()

	queueWatcher := &QueueWatcher{
		clientset: client,
		fc:        fc,
		//informer: queueInformer,
		lister:     kbinformer.Scheduling().V1alpha1().Queues().Lister(),
		qWorkQueue: k8sclient.NewKeyedQueue(),
		qSynced:    queueInformer.HasSynced,
	}

	/*scheme.AddToScheme(scheme.Scheme)
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	*/
	queueInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				glog.Errorf("AddFunc: error getting key %v", err)
			}
			queueWatcher.enqueueQueueAddition(key, obj)
		},
		UpdateFunc: func(old, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err != nil {
				glog.Errorf("UpdateFunc: error getting key %v", err)
			}
			queueWatcher.enqueueQueueUpdate(key, old, new)
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				glog.Errorf("DeleteFunc: error getting key %v", err)
			}
			queueWatcher.enqueueQueueDeletion(key, obj)
		},
	})

	glog.V(2).Info("Values in queueinformer2...+%v", queueInformer)

	//queueWatcher.controller = queueInformer
	//glog.V(2).Info("Values in queueinformer3...+%v",queueInformer)
	//glog.Info("++++ Inside : QW set for q watcher +%v +%v",queueWatcher.temp, kbinformer.Scheduling().V1alpha1().Queues().Informer().AddEventHandler)
	return queueWatcher, kbinformer

}

func (qw *QueueWatcher) enqueueQueueAddition(key interface{}, obj interface{}) {
	q, ok := obj.(*kbv1.Queue)
	if !ok {
		glog.Errorf("Cannot convert to *kbv1.Queue: %v", obj)
		return
	}

	addedQ := qw.parseQueue(q, QueueAdded)
	qw.qWorkQueue.Add(key, addedQ)

	return
}

func (qw *QueueWatcher) enqueueQueueUpdate(old interface{}, new interface{}, obj interface{}) {

	q, ok := obj.(*kbv1.Queue)
	if !ok {
		glog.Errorf("Cannot convert to *kbv1.Queue: %v", obj)
		return
	}

	updateQ := qw.parseQueue(q, QueueUpdated)
	qw.enqueueQueueDeletion(old, updateQ)
	qw.enqueueQueueAddition(new, updateQ)

}

func (qw *QueueWatcher) enqueueQueueDeletion(key interface{}, obj interface{}) {
	q, ok := obj.(*kbv1.Queue)
	if !ok {
		glog.Errorf("Cannot convert to *kbv1.Queue: %v", obj)
		return
	}
	deletedQ := qw.parseQueue(q, QueueDeleted)
	qw.qWorkQueue.Add(key, deletedQ)
	glog.Info("enqueueQueueDeletion: Deleted queue ", deletedQ.Name)
}

func (qw *QueueWatcher) parseQueue(q *kbv1.Queue, phase QueuePhase) *QueueProportion {

	queue := &QueueProportion{
		Name:   q.Name,
		Weight: q.Spec.Weight,
		State:  phase,
	}

	return queue
}

// Run starts queue watcher.
func (qw *QueueWatcher) Run(stopCh <-chan struct{}, qWorkers int) {
	defer utilruntime.HandleCrash()
	// The workers can stop when we are done.
	defer qw.qWorkQueue.ShutDown()

	glog.Info("Getting queue updates...")

	if !cache.WaitForCacheSync(stopCh, qw.qSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	glog.Info("Starting queue watching workers")
	for i := 0; i < qWorkers; i++ {
		go wait.Until(qw.qWorker, time.Second, stopCh)
	}

	glog.Info("Started worker threads")
	<-stopCh
	glog.Info("Shutting down worker threads")

}

func (qw *QueueWatcher) qWorker() {
	for {
		func() {
			key, items, quit := qw.qWorkQueue.Get()
			if quit {
				return
			}
			for _, item := range items {
				q := item.(*QueueProportion)
				switch q.State {
				case QueueAdded:
					QueueMux.Lock()
					qd := qw.Get(q)
					_, ok := QueueToQD[q.Name]
					if ok {
						glog.Infof("Queue %s already exists", q.Name)
						QueueMux.Unlock()
						continue
					}
					QueueToQD[q.Name] = qd
					QueueMux.Unlock()
					firmament.QueueAdded(qw.fc, qd)
					glog.Infof(" in QueueAdded")
				case QueueDeleted:
					QueueMux.RLock()
					qd, ok := QueueToQD[q.Name]
					QueueMux.RUnlock()
					if !ok {
						glog.Fatalf("Queue %s does not exist", qd.Name)
					}
					QueueMux.Lock()
					delete(QueueToQD, q.Name)
					QueueMux.Unlock()
					firmament.QueueRemoved(qw.fc, qd)
					glog.Info(QueueToQD, " in QueueDeleted")
				case QueueUpdated:
					QueueMux.RLock()
					qd, ok := QueueToQD[q.Name]
					if !ok {
						glog.Fatalf("Queue %s does not exist", qd.Name)
					}
					qw.updateQueueDescriptor(q, qd)
					QueueMux.RUnlock()
					firmament.QueueUpdated(qw.fc, qd)
					glog.Info(QueueToQD, " in QueueUpdated")

				default:
					glog.Fatalf("Unexpected queue %s phase %s", q.Name, q.State)
				}
			}
			defer qw.qWorkQueue.Done(key)
		}()
	}
}

func (qw *QueueWatcher) Get(proportion *QueueProportion) *firmament.QueueDescriptor {

	q := &firmament.QueueDescriptor{
		Name:   proportion.Name,
		Weight: proportion.Weight,
	}

	return q
}

func (qw *QueueWatcher) updateQueueDescriptor(queue *QueueProportion, qd *firmament.QueueDescriptor) {

	queue.Name = qd.Name
	queue.Weight = qd.Weight

}

func init() {
	//Initialize to default queue
	qw := &QueueWatcher{}
	queue := &QueueProportion{Name: "default", Weight: 1}
	qw.Get(queue)
}
