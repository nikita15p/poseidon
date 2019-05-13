package batchclient

import (
	"github.com/kubernetes-sigs/poseidon/pkg/client/clientset/versioned"
	"github.com/kubernetes-sigs/poseidon/pkg/client/listers/scheduling/v1alpha1"
	"github.com/kubernetes-sigs/poseidon/pkg/firmament"
	"github.com/kubernetes-sigs/poseidon/pkg/k8sclient"
	"k8s.io/client-go/tools/cache"
	"sync"
)

// NodePhase represents a node phase.
type QueuePhase string

// QueueMux is used to guard access to the queue maps.
var QueueMux sync.RWMutex

// PgMux is used to guard access to the pod group maps
var PgMux sync.RWMutex

// QueuetoQD maps queue name to firmament queue descriptor.
var QueueToQD map[string]*firmament.QueueDescriptor

//PgToPGD maps pod groeps to firmament pod group descriptor
var PgToPGD map[string]*firmament.PodGroupDescriptor

const (
	// QueueAdded represents a queue added phase.
	QueueAdded QueuePhase = "Added"
	// QueueDeleted represents a queue deleted phase.
	QueueDeleted QueuePhase = "Deleted"
	// QueueUpdated represents a queue updated phase.
	QueueUpdated QueuePhase = "Updated"
)

// PodGroupPhase represents a pod phase.
type PodGroupPhase string

const (
	// PGAdded represents a pod group added phase
	PGAdded PodGroupPhase = "Added"
	//PGDeleted represents a pod group deleted phase
	PGDeleted PodGroupPhase = "Deleted"
	// PGUpdated irepresents a pod group updated phase
	PGUpdated PodGroupPhase = "Updated"
)

type PgWatcher struct {
	clientset   versioned.Interface
	pgWorkQueue k8sclient.Queue
	//informer cache.Controller
	lister v1alpha1.PodGroupLister
	//recorder     record.EventRecorder
	fc       firmament.FirmamentSchedulerClient
	pgSynced cache.InformerSynced
	//mutex sync.Mutex

}

type PodGroup struct {
	Name              string
	MinMember         int32
	Queue             string
	PriorityClassName string
	State             PodGroupPhase
}

type QueueWatcher struct {
	clientset  versioned.Interface
	qWorkQueue k8sclient.Queue
	//informer cache.Controller
	lister v1alpha1.QueueLister
	//recorder     record.EventRecorder
	fc      firmament.FirmamentSchedulerClient
	qSynced cache.InformerSynced
	//mutex sync.Mutex

}

type QueueProportion struct {
	Name   string
	Weight int32
	State  QueuePhase
}
