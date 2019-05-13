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

// NewPodGroupWatcher initializes a PgWatcher.
func NewPodGroupWatcher(client versioned.Interface, fc firmament.FirmamentSchedulerClient) (*PgWatcher, kbinfo.SharedInformerFactory) {
	glog.V(2).Info("Starting PgWatcher...")
	PgToPGD = make(map[string]*firmament.PodGroupDescriptor)
	kbinformer := kbinfo.NewSharedInformerFactory(client, 0)

	// create informer for Pg information
	pgInformer := kbinformer.Scheduling().V1alpha1().PodGroups().Informer()

	pgWatcher := &PgWatcher{
		clientset: client,
		fc:        fc,
		//informer: queueInformer,
		lister:      kbinformer.Scheduling().V1alpha1().PodGroups().Lister(),
		pgWorkQueue: k8sclient.NewKeyedQueue(),
		pgSynced:    pgInformer.HasSynced,
	}

	pgInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				glog.Errorf("AddFunc: error getting key %v", err)
			}
			pgWatcher.enqueuePgAddition(key, obj)
		},
		UpdateFunc: func(old, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err != nil {
				glog.Errorf("UpdateFunc: error getting key %v", err)
			}
			pgWatcher.enqueuePgUpdate(key, old, new)
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				glog.Errorf("DeleteFunc: error getting key %v", err)
			}
			pgWatcher.enqueuePgDeletion(key, obj)
		},
	})

	return pgWatcher, kbinformer

}

func (pgw *PgWatcher) enqueuePgAddition(key interface{}, obj interface{}) {
	pg, ok := obj.(*kbv1.PodGroup)
	if !ok {
		glog.Errorf("Cannot convert to *kbv1.PgGroup: %v", obj)
		return
	}
	pg.Status.Phase = "Added"
	addedPg := pgw.parsePg(pg, PGAdded)
	//x:=	pgw.lister.PodGroups(addedQ.Name)
	pgw.pgWorkQueue.Add(key, addedPg)
	glog.Info("enqueuePodGroupAdition: Added queue ", addedPg.Name)
	return
}

func (pgw *PgWatcher) enqueuePgUpdate(old interface{}, new interface{}, obj interface{}) {

	pg, ok := obj.(*kbv1.PodGroup)
	if !ok {
		glog.Errorf("Cannot convert to *kbv1.PodGroup: %v", obj)
		return
	}
	updatePg := pgw.parsePg(pg, PGUpdated)
	pgw.enqueuePgDeletion(old, updatePg)
	pgw.enqueuePgAddition(new, updatePg)

}

func (pgw *PgWatcher) enqueuePgDeletion(key interface{}, obj interface{}) {
	pg, ok := obj.(*kbv1.PodGroup)
	if !ok {
		glog.Errorf("Cannot convert to *kbv1.PodGroup: %v", obj)
		return
	}
	deletedPg := pgw.parsePg(pg, PGDeleted)
	pgw.pgWorkQueue.Add(key, deletedPg)
	glog.Info("enqueuePgDeletion: Deleted pod ", deletedPg.Name)
}

func (pgw *PgWatcher) parsePg(pg *kbv1.PodGroup, phase PodGroupPhase) *PodGroup {
	podGroup := &PodGroup{
		Name:              pg.Name,
		MinMember:         pg.Spec.MinMember,
		Queue:             pg.Spec.Queue,
		PriorityClassName: pg.Spec.PriorityClassName,
		State:             phase,
	}

	return podGroup
}

// Run starts queue watcher.
func (pgw *PgWatcher) Run(stopCh <-chan struct{}, qWorkers int) {
	defer utilruntime.HandleCrash()

	// The workers can stop when we are done.
	defer pgw.pgWorkQueue.ShutDown()
	defer glog.Info("Shutting down PgWatcher")
	glog.Info("Getting pod group updates...")
	//go pgw.controller.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, pgw.pgSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	glog.Info("Starting pod group watching workers")
	for i := 0; i < qWorkers; i++ {
		go wait.Until(pgw.pgWorker, time.Second, stopCh)
	}

	<-stopCh
	glog.Info("Stopping pod group watcher")
}

func (pgw *PgWatcher) pgWorker() {
	for {
		func() {
			key, items, quit := pgw.pgWorkQueue.Get()
			if quit {
				return
			}
			for _, item := range items {
				pg := item.(*PodGroup)
				switch pg.State {
				case PGAdded:
					PgMux.Lock()
					pgd := pgw.Get(pg)
					_, ok := PgToPGD[pg.Name]
					if ok {
						glog.Infof("Pod Group %s already exists", pg.Name)
						PgMux.Unlock()
						continue
					}
					PgToPGD[pg.Name] = pgd
					PgMux.Unlock()
					firmament.PodGroupAdded(pgw.fc, pgd)
					glog.Info(PgToPGD, " in PodGroupAdded")
				case PGDeleted:
					PgMux.RLock()
					pgd, ok := PgToPGD[pg.Name]
					PgMux.RUnlock()
					if !ok {
						glog.Fatalf("Pod Group %s does not exist", pg.Name)
					}
					PgMux.Lock()
					delete(PgToPGD, pg.Name)
					PgMux.Unlock()
					firmament.PodGroupRemoved(pgw.fc, pgd)
					glog.Info(PgToPGD, " in PodGroupDeleted")
				case PGUpdated:
					QueueMux.RLock()
					pgd, ok := PgToPGD[pg.Name]
					if !ok {
						glog.Fatalf("Queue %s does not exist", pgd.Name)
					}
					pgw.updatePgDescriptor(pg, pgd)
					QueueMux.RUnlock()
					firmament.PodGroupUpdated(pgw.fc, pgd)
					glog.Info(QueueToQD, " in PodGroupUpdated")
				default:
					glog.Fatalf("Unexpected pod group %s phase %s", pg.Name, pg.State)
				}
			}
			defer pgw.pgWorkQueue.Done(key)
		}()
	}
}

func (pgw *PgWatcher) Get(pg *PodGroup) *firmament.PodGroupDescriptor {

	q := &firmament.PodGroupDescriptor{
		Name:      pg.Name,
		MinMember: pg.MinMember,
		Queue:     pg.Queue,
	}

	return q
}

func (pgw *PgWatcher) updatePgDescriptor(pg *PodGroup, qd *firmament.PodGroupDescriptor) {

	pg.Name = qd.Name
	pg.MinMember = qd.MinMember
	pg.Queue = qd.Queue

}
