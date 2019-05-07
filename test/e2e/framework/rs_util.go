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

package framework

import (
	"fmt"

	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForReadyReplicaSet waits until the replicaset has all of its replicas ready.
func (f *Framework) WaitForReadyReplicaSet(name string) error {
	err := wait.Poll(Poll, pollShortTimeout, func() (bool, error) {
		rs, err := f.ClientSet.ExtensionsV1beta1().ReplicaSets(f.Namespace.Name).Get(name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return *(rs.Spec.Replicas) == rs.Status.Replicas && *(rs.Spec.Replicas) == rs.Status.ReadyReplicas, nil
	})
	if err == wait.ErrWaitTimeout {
		err = fmt.Errorf("replicaset %q never became ready", name)
	}
	return err
}

// WaitForReplicaSetTargetSpecReplicas waits for .spec.replicas of a RS to equal targetReplicaNum
func (f *Framework) WaitForReplicaSetTargetSpecReplicas(replicaSet *extensions.ReplicaSet, targetReplicaNum int32) error {
	desiredGeneration := replicaSet.Generation
	err := wait.PollImmediate(Poll, pollShortTimeout, func() (bool, error) {
		rs, err := f.ClientSet.ExtensionsV1beta1().ReplicaSets(replicaSet.Namespace).Get(replicaSet.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return rs.Status.ObservedGeneration >= desiredGeneration && *rs.Spec.Replicas == targetReplicaNum, nil
	})
	if err == wait.ErrWaitTimeout {
		err = fmt.Errorf("replicaset %q never had desired number of .spec.replicas", replicaSet.Name)
	}
	return err
}

// WaitForReplicaSetDelete waits for the ReplicateSet to be removed
func (f *Framework) WaitForReplicaSetDelete(replicaSet *extensions.ReplicaSet) error {
	replicas := int32(0)
	replicaSet.Spec.Replicas = &replicas
	err := wait.Poll(Poll, pollShortTimeout, func() (bool, error) {
		_, err := f.ClientSet.ExtensionsV1beta1().ReplicaSets(replicaSet.Namespace).Update(replicaSet)
		if err == nil {
			return true, nil
		}
		// Retry only on update conflict.
		if errors.IsConflict(err) {
			return false, nil
		}
		return false, err
	})
	err = wait.Poll(Poll, pollShortTimeout, func() (bool, error) {
		err := f.ClientSet.ExtensionsV1beta1().ReplicaSets(replicaSet.Namespace).Delete(replicaSet.Name, &metav1.DeleteOptions{})
		if err != nil {
			return false, err
		}
		return true, nil
	})
	if err == wait.ErrWaitTimeout {
		err = fmt.Errorf("replicaset %q not deleted", replicaSet.Name)
	}
	return err
}
