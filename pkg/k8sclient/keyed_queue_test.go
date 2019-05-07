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

package k8sclient

import (
	"reflect"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	fakeQueue := NewKeyedQueue()
	var testDatas = []struct {
		key   interface{}
		value interface{}
	}{
		{"Item1", "Value1"},
		{"Item2", "Value2"},
		{"Item2", "Value22"},
		{"Item3", "Value2"},
		{10.12, 99},
		{"Item2", "Value23"},
	}

	var testResult = []struct {
		key   interface{}
		value []interface{}
	}{
		{"Item1", []interface{}{"Value1"}},
		{"Item2", []interface{}{"Value2", "Value22", "Value23"}},
		{"Item3", []interface{}{"Value2"}},
		{10.12, []interface{}{99}},
	}

	for _, testValue := range testDatas {
		fakeQueue.Add(testValue.key, testValue.value)

	}
	fakeQueue.Done(testDatas[0].key)
	for _, testValue := range testResult {
		key, value, _ := fakeQueue.Get()
		if !(reflect.DeepEqual(key, testValue.key) && reflect.DeepEqual(value, testValue.value)) {

			t.Error("expected ", testValue.key, testValue.value, "got ", key, value)
		}
	}
}

func TestNotDone(t *testing.T) {
	fakeQueue := NewKeyedQueue()
	var testDatas = []struct {
		key   interface{}
		value interface{}
	}{
		{"Item1", "Value1"},
		{"Item2", "Value2"},
		{"Item2", "Value22"},
		{"Item3", "Value2"},
		{10.12, 99},
		{"Item2", "Value23"},
	}

	var testResult = []struct {
		key   interface{}
		value []interface{}
	}{
		{"Item1", []interface{}{"Value1"}},
		{"Item3", []interface{}{"Value2"}},
		{10.12, []interface{}{99}},
	}

	// Get() moves the key to a temporary processing queue,
	// Add() following a Get() on the same key will not be available in the queue
	// till the key from the previous Get() call is processed.
	fakeQueue.Add(testDatas[1].key, testDatas[1].value)
	fakeQueue.Get()

	for _, testValue := range testDatas {
		fakeQueue.Add(testValue.key, testValue.value)
	}

	for _, testValue := range testResult {
		key, value, _ := fakeQueue.Get()
		if reflect.DeepEqual(key, testValue.key) {
			if !reflect.DeepEqual(value, testValue.value) {
				t.Error("expected ", testValue.key, testValue.value, "got ", key, value)
			}
		} else {

			t.Error("expected ", testValue.key, testValue.value, "got ", key, value)
		}
	}
}

func TestDone(t *testing.T) {
	fakeQueue := NewKeyedQueue()
	var testDatas = []struct {
		key   interface{}
		value interface{}
	}{
		{"Item1", "Value1"},
		{"Item2", "Value2"},
		{"Item2", "Value22"},
		{"Item3", "Value2"},
		{10.12, 99},
		{"Item2", "Value23"},
	}

	var testResult = []struct {
		key   interface{}
		value []interface{}
	}{
		{"Item1", []interface{}{"Value1"}},
		{"Item3", []interface{}{"Value2"}},
		{10.12, []interface{}{99}},
		{"Item2", []interface{}{"Value2", "Value22", "Value23"}},
	}

	fakeQueue.Add(testDatas[1].key, testDatas[1].value)
	fakeQueue.Get()
	for _, testValue := range testDatas {
		fakeQueue.Add(testValue.key, testValue.value)

	}
	fakeQueue.Done(testDatas[1].key)
	//fakeQueue.Get()
	for _, testValue := range testResult {
		key, value, _ := fakeQueue.Get()
		if reflect.DeepEqual(key, testValue.key) {
			if !reflect.DeepEqual(value, testValue.value) {
				t.Error("expected ", testValue.key, testValue.value, "got ", key, value)
			}
		} else {

			t.Error("expected ", testValue.key, testValue.value, "got ", key, value)
		}
	}
}

func TestShutDown(t *testing.T) {
	fakeQueue := NewKeyedQueue()
	var testDatas = []struct {
		key   interface{}
		value interface{}
	}{
		{"Item1", "Value1"},
		{"Item2", "Value2"},
		{"Item2", "Value22"},
		{"Item3", "Value2"},
		{10.12, 99},
		{"Item2", "Value23"},
	}

	var testResult = []struct {
		key   interface{}
		value []interface{}
	}{
		{"Item1", []interface{}{"Value1"}},
		{"Item2", []interface{}{"Value2"}},
	}

	// Add() followed by Shutdown on queue will not be processed
	// this Add will be ignored.

	for index, testValue := range testDatas {
		if index == 2 {
			fakeQueue.ShutDown()
		}
		fakeQueue.Add(testValue.key, testValue.value)
	}
	for _, testValue := range testResult {
		key, value, _ := fakeQueue.Get()
		if reflect.DeepEqual(key, testValue.key) {
			if !reflect.DeepEqual(value, testValue.value) {
				t.Error("expected ", testValue.key, testValue.value, "got ", key, value)
			}
		} else {
			t.Error("expected ", testValue.key, testValue.value, "got ", key, value)
		}
	}
}

// Get after shutdown on empty queue
func TestGetAfterShutDown(t *testing.T) {
	fakeQueue := NewKeyedQueue()
	// Shutting down on nil queue
	fakeQueue.ShutDown()
	key, value, down := fakeQueue.Get()
	if key != nil && value != nil && down != true {
		t.Error("expected ", nil, nil, true, "got ", key, value, down)
	}
}

// Get on an Empty Queue
func TestGetOnEmptyQueue(t *testing.T) {
	fakeQueue := NewKeyedQueue()
	var testResult = []struct {
		key   interface{}
		value []interface{}
	}{
		{"Item1", []interface{}{"Value1"}},
	}
	go func() {
		key, value, _ := fakeQueue.Get()
		//blocks
		if !reflect.DeepEqual(key, testResult[0].key) && !reflect.DeepEqual(value, testResult[0].value) {
			t.Error("expected Item1,value1 got ", key, value)
		}
	}()
	time.Sleep(1 * time.Second)
	fakeQueue.Add("Item1", "Value1")
}

//get after shutdown on empty queue
func TestShuttingDown(t *testing.T) {
	fakeQueue := NewKeyedQueue()
	//shutting down on nil queue
	fakeQueue.ShutDown()
	key, value, down := fakeQueue.Get()
	if down != fakeQueue.ShuttingDown() && key != nil && value != nil {
		t.Error("expected ", nil, nil, true, "got ", key, value, down)
	}
}
