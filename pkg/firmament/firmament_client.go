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

package firmament

import (
	"fmt"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// Schedule sends a schedule request to firmament server.
func Schedule(client FirmamentSchedulerClient) *SchedulingDeltas {
	scheduleResp, err := client.Schedule(context.Background(), &ScheduleRequest{})
	if err != nil {
		grpclog.Fatalf("%v.Schedule(_) = _, %v: ", client, err)
	}
	return scheduleResp
}

// TaskCompleted tells firmament server the given task is completed.
func TaskCompleted(client FirmamentSchedulerClient, tuid *TaskUID) {
	tCompletedResp, err := client.TaskCompleted(context.Background(), tuid)
	if err != nil {
		grpclog.Fatalf("%v.TaskCompleted(_) = _, %v: ", client, err)
	}
	switch tCompletedResp.Type {
	case TaskReplyType_TASK_NOT_FOUND:
		glog.Fatalf("Task %d not found", tuid.TaskUid)
	case TaskReplyType_TASK_JOB_NOT_FOUND:
		glog.Fatalf("Task's %d job not found", tuid.TaskUid)
	case TaskReplyType_TASK_COMPLETED_OK:
	default:
		panic(fmt.Sprintf("Unexpected TaskCompleted response %v for task %v", tCompletedResp, tuid.TaskUid))
	}
}

// TaskFailed tells firmament server the given task is failed.
func TaskFailed(client FirmamentSchedulerClient, tuid *TaskUID) {
	tFailedResp, err := client.TaskFailed(context.Background(), tuid)
	if err != nil {
		grpclog.Fatalf("%v.TaskFailed(_) = _, %v: ", client, err)
	}
	switch tFailedResp.Type {
	case TaskReplyType_TASK_NOT_FOUND:
		glog.Fatalf("Task %d not found", tuid.TaskUid)
	case TaskReplyType_TASK_JOB_NOT_FOUND:
		glog.Fatalf("Task's %d job not found", tuid.TaskUid)
	case TaskReplyType_TASK_FAILED_OK:
	default:
		panic(fmt.Sprintf("Unexpected TaskFailed response %v for task %v", tFailedResp, tuid.TaskUid))
	}
}

// TaskRemoved tells firmament server the given task is removed.
func TaskRemoved(client FirmamentSchedulerClient, tuid *TaskUID) {
	tRemovedResp, err := client.TaskRemoved(context.Background(), tuid)
	if err != nil {
		grpclog.Fatalf("%v.TaskRemoved(_) = _, %v: ", client, err)
	}
	switch tRemovedResp.Type {
	case TaskReplyType_TASK_NOT_FOUND:
		glog.Fatalf("Task %d not found", tuid.TaskUid)
	case TaskReplyType_TASK_JOB_NOT_FOUND:
		glog.Fatalf("Task's %d job not found", tuid.TaskUid)
	case TaskReplyType_TASK_REMOVED_OK:
	default:
		panic(fmt.Sprintf("Unexpected TaskRemoved response %v for task %v", tRemovedResp, tuid.TaskUid))
	}
}

// TaskSubmitted tells firmament server the given task is submitted.
func TaskSubmitted(client FirmamentSchedulerClient, td *TaskDescription) {
	tSubmittedResp, err := client.TaskSubmitted(context.Background(), td)
	if err != nil {
		grpclog.Fatalf("%v.TaskSubmitted(_) = _, %v: ", client, err)
	}
	switch tSubmittedResp.Type {
	case TaskReplyType_TASK_ALREADY_SUBMITTED:
		glog.Fatalf("Task (%s,%d) already submitted", td.JobDescriptor.Uuid, td.TaskDescriptor.Uid)
	case TaskReplyType_TASK_STATE_NOT_CREATED:
		glog.Fatalf("Task (%s,%d) not in created state", td.JobDescriptor.Uuid, td.TaskDescriptor.Uid)
	case TaskReplyType_TASK_SUBMITTED_OK:
	default:
		panic(fmt.Sprintf("Unexpected TaskSubmitted response %v for task (%v,%v)", tSubmittedResp, td.JobDescriptor.Uuid, td.TaskDescriptor.Uid))
	}
}

// TaskUpdated tells firmament server the given task is updated.
func TaskUpdated(client FirmamentSchedulerClient, td *TaskDescription) {
	tUpdatedResp, err := client.TaskUpdated(context.Background(), td)
	if err != nil {
		grpclog.Fatalf("%v.TaskUpdated(_) = _, %v: ", client, err)
	}
	switch tUpdatedResp.Type {
	case TaskReplyType_TASK_NOT_FOUND:
		glog.Fatalf("Task (%s,%d) not found", td.JobDescriptor.Uuid, td.TaskDescriptor.Uid)
	case TaskReplyType_TASK_JOB_NOT_FOUND:
		glog.Fatalf("Task's (%s,%d) job not found", td.JobDescriptor.Uuid, td.TaskDescriptor.Uid)
	case TaskReplyType_TASK_UPDATED_OK:
	default:
		panic(fmt.Sprintf("Unexpected TaskUpdated response %v for task (%v,%v)", tUpdatedResp, td.JobDescriptor.Uuid, td.TaskDescriptor.Uid))
	}
}

// NodeAdded tells firmament server the given node is added.
func NodeAdded(client FirmamentSchedulerClient, rtnd *ResourceTopologyNodeDescriptor) {
	nAddedResp, err := client.NodeAdded(context.Background(), rtnd)
	if err != nil {
		grpclog.Fatalf("%v.NodeAdded(_) = _, %v: ", client, err)
	}
	switch nAddedResp.Type {
	case NodeReplyType_NODE_ALREADY_EXISTS:
		glog.Infof("Tried to add existing node %s", rtnd.ResourceDesc.Uuid)
	case NodeReplyType_NODE_ADDED_OK:
	default:
		panic(fmt.Sprintf("Unexpected NodeAdded response %v for node %v", nAddedResp, rtnd.ResourceDesc.Uuid))
	}
}

// NodeFailed tells firmament server the given node is failed.
func NodeFailed(client FirmamentSchedulerClient, ruid *ResourceUID) {
	nFailedResp, err := client.NodeFailed(context.Background(), ruid)
	if err != nil {
		grpclog.Fatalf("%v.NodeFailed(_) = _, %v: ", client, err)
	}
	switch nFailedResp.Type {
	case NodeReplyType_NODE_NOT_FOUND:
		glog.Fatalf("Tried to fail non-existing node %s", ruid.ResourceUid)
	case NodeReplyType_NODE_FAILED_OK:
	default:
		panic(fmt.Sprintf("Unexpected NodeFailed response %v for node %v", nFailedResp, ruid.ResourceUid))
	}
}

// NodeRemoved tells firmament server the given node is removed.
func NodeRemoved(client FirmamentSchedulerClient, ruid *ResourceUID) {
	nRemovedResp, err := client.NodeRemoved(context.Background(), ruid)
	if err != nil {
		grpclog.Fatalf("%v.NodeRemoved(_) = _, %v: ", client, err)
	}
	switch nRemovedResp.Type {
	case NodeReplyType_NODE_NOT_FOUND:
		glog.Fatalf("Tried to remove non-existing node %s", ruid.ResourceUid)
	case NodeReplyType_NODE_REMOVED_OK:
	default:
		panic(fmt.Sprintf("Unexpected NodeRemoved response %v for node %v", nRemovedResp, ruid.ResourceUid))
	}
}

// NodeUpdated tells firmament server the given node is updated.
func NodeUpdated(client FirmamentSchedulerClient, rtnd *ResourceTopologyNodeDescriptor) {
	nUpdatedResp, err := client.NodeUpdated(context.Background(), rtnd)
	if err != nil {
		grpclog.Fatalf("%v.NodeUpdated(_) = _, %v: ", client, err)
	}
	switch nUpdatedResp.Type {
	case NodeReplyType_NODE_NOT_FOUND:
		glog.Fatalf("Tried to updated non-existing node %s", rtnd.ResourceDesc.Uuid)
	case NodeReplyType_NODE_UPDATED_OK:
	default:
		panic(fmt.Sprintf("Unexpected NodeUpdated response %v for node %v", nUpdatedResp, rtnd.ResourceDesc.Uuid))
	}
}

// AddTaskStats sends task status to firmament server.
func AddTaskStats(client FirmamentSchedulerClient, ts *TaskStats) {
	_, err := client.AddTaskStats(context.Background(), ts)
	if err != nil {
		grpclog.Fatalf("%v.AddTaskStats(_) = _, %v: ", client, err)
	}
}

// AddNodeStats sends node status to firmament server.
func AddNodeStats(client FirmamentSchedulerClient, rs *ResourceStats) {
	_, err := client.AddNodeStats(context.Background(), rs)
	if err != nil {
		grpclog.Fatalf("%v.AddNodeStats(_) = _, %v: ", client, err)
	}
}

// Check tests if firmament server is health
func Check(client FirmamentSchedulerClient, req_service *HealthCheckRequest) (bool, error) {
	res, err := client.Check(context.Background(), req_service)
	if err == nil {
		if res.GetStatus() == ServingStatus_SERVING {
			return true, nil
		}
	}
	return false, err
}

// New creates a firmament scheduler client by a remote server address.
// NOTE: it's an insecure connection.
func New(address string) (FirmamentSchedulerClient, *grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		glog.Errorf("Did not connect to Firmament scheduler: %v", err)
		return nil, nil, err
	}
	fc := NewFirmamentSchedulerClient(conn)
	return fc, conn, nil
}

// AddTaskStats sends task status to firmament server.
func AddTaskInfo(client FirmamentSchedulerClient, ts *TaskInfo) {
	_, err := client.AddTaskInfo(context.Background(), ts)
	if err != nil {
		grpclog.Fatalf("%v.AddTaskInfo(_) = _, %v: ", client, err)
	}
}

// QueueAdded tells firmament server the given queue is added.
func QueueAdded(client FirmamentSchedulerClient, qd *QueueDescriptor) {
	qAddedResp, err := client.QueueAdded(context.Background(), qd)
	glog.Infof("\n Queue response is %v", qAddedResp)
	if err != nil {
		grpclog.Fatalf("%v.QueueAdded(_) = _, %v: ", client, err)
	}
	switch qAddedResp.Type {
	case QueueReplyType_QUEUE_ALREADY_ADDED:
		glog.Infof("Tried to add existing queue %s", qd.Uid)
	case QueueReplyType_QUEUE_ADDED_OK:
	default:
		panic(fmt.Sprintf("Unexpected QueueAdded response %v for queue %v", qAddedResp, qd.Uid))
	}
}

// QueueRemoved tells firmament server the given queue is removed.
func QueueRemoved(client FirmamentSchedulerClient, qd *QueueDescriptor) {
	qRemovedResp, err := client.QueueRemoved(context.Background(), qd)
	if err != nil {
		grpclog.Fatalf("%v.QueueRemoved(_) = _, %v: ", client, err)
	}
	switch qRemovedResp.Type {
	case QueueReplyType_QUEUE_NOT_FOUND:
		glog.Fatalf("Tried to remove non-existing queue %s", qd.Uid)
	case QueueReplyType_QUEUE_REMOVED_OK:
	default:
		panic(fmt.Sprintf("Unexpected QueueRemoved response %v for queue %v", qRemovedResp, qd.Uid))
	}
}

// QueueUpdated tells firmament server the given queue is updated.
func QueueUpdated(client FirmamentSchedulerClient, qd *QueueDescriptor) {
	qUpdatedResp, err := client.QueueUpdated(context.Background(), qd)
	if err != nil {
		grpclog.Fatalf("%v.QueueUpdated(_) = _, %v: ", client, err)
	}
	switch qUpdatedResp.Type {
	case QueueReplyType_QUEUE_NOT_FOUND:
		glog.Fatalf("Tried to updated non-existing queue %s", qd.Uid)
	case QueueReplyType_QUEUE_UPDATED_OK:
	default:
		panic(fmt.Sprintf("Unexpected QueueUpdated response %v for queue %v", qUpdatedResp, qd.Uid))
	}
}

// PodGroupAdded tells firmament server the given pod group is added.
func PodGroupAdded(client FirmamentSchedulerClient, pgd *PodGroupDescriptor) {
	pgAddedResp, err := client.PodGroupAdded(context.Background(), pgd)
	glog.Infof("\n Pod Group response is %v", pgAddedResp)
	if err != nil {
		grpclog.Fatalf("%v.QueueAdded(_) = _, %v: ", client, err)
	}
	switch pgAddedResp.Type {
	case PodGroupReplyType_PODGROUP_ALREADY_ADDED:
		glog.Infof("Tried to add existing pod group %s", pgd.Uid)
	case PodGroupReplyType_PODGROUP_ADDED_OK:
	default:
		panic(fmt.Sprintf("Unexpected PodGroupAdded response %v for pod group %v", pgAddedResp, pgd.Uid))
	}
}

// PodGroupRemoved tells firmament server the given pod group is removed.
func PodGroupRemoved(client FirmamentSchedulerClient, pgd *PodGroupDescriptor) {
	pgRemovedResp, err := client.PodGroupRemoved(context.Background(), pgd)
	if err != nil {
		grpclog.Fatalf("%v.PodGroupRemoved(_) = _, %v: ", client, err)
	}
	switch pgRemovedResp.Type {
	case PodGroupReplyType_PODGROUP_NOT_FOUND:
		glog.Fatalf("Tried to remove non-existing pod group %s", pgd.Uid)
	case PodGroupReplyType_PODGROUP_REMOVED_OK:
	default:
		panic(fmt.Sprintf("Unexpected PodGroupRemoved response %v for pod group %v", pgRemovedResp, pgd.Uid))
	}
}

// PodGroupUpdated tells firmament server the given pod group is updated.
func PodGroupUpdated(client FirmamentSchedulerClient, pgd *PodGroupDescriptor) {
	pgUpdatedResp, err := client.PodGroupUpdated(context.Background(), pgd)
	if err != nil {
		grpclog.Fatalf("%v.PodGroupUpdated(_) = _, %v: ", client, err)
	}
	switch pgUpdatedResp.Type {
	case PodGroupReplyType_PODGROUP_NOT_FOUND:
		glog.Fatalf("Tried to updated non-existing pod group %s", pgd.Uid)
	case PodGroupReplyType_PODGROUP_UPDATED_OK:
	default:
		panic(fmt.Sprintf("Unexpected PodGroupUpdated response %v for pod group %v", pgUpdatedResp, pgd.Uid))
	}
}
