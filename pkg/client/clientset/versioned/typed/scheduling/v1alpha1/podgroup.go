/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/kubernetes-sigs/poseidon/pkg/apis/scheduling/v1alpha1"
	scheme "github.com/kubernetes-sigs/poseidon/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PodGroupsGetter has a method to return a PodGroupInterface.
// A group's client should implement this interface.
type PodGroupsGetter interface {
	PodGroups(namespace string) PodGroupInterface
}

// PodGroupInterface has methods to work with PodGroup resources.
type PodGroupInterface interface {
	Create(*v1alpha1.PodGroup) (*v1alpha1.PodGroup, error)
	Update(*v1alpha1.PodGroup) (*v1alpha1.PodGroup, error)
	UpdateStatus(*v1alpha1.PodGroup) (*v1alpha1.PodGroup, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.PodGroup, error)
	List(opts v1.ListOptions) (*v1alpha1.PodGroupList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PodGroup, err error)
	PodGroupExpansion
}

// podGroups implements PodGroupInterface
type podGroups struct {
	client rest.Interface
	ns     string
}

// newPodGroups returns a PodGroups
func newPodGroups(c *SchedulingV1alpha1Client, namespace string) *podGroups {
	return &podGroups{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the podGroup, and returns the corresponding podGroup object, and an error if there is any.
func (c *podGroups) Get(name string, options v1.GetOptions) (result *v1alpha1.PodGroup, err error) {
	result = &v1alpha1.PodGroup{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podgroups").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PodGroups that match those selectors.
func (c *podGroups) List(opts v1.ListOptions) (result *v1alpha1.PodGroupList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.PodGroupList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podgroups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested podGroups.
func (c *podGroups) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("podgroups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a podGroup and creates it.  Returns the server's representation of the podGroup, and an error, if there is any.
func (c *podGroups) Create(podGroup *v1alpha1.PodGroup) (result *v1alpha1.PodGroup, err error) {
	result = &v1alpha1.PodGroup{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("podgroups").
		Body(podGroup).
		Do().
		Into(result)
	return
}

// Update takes the representation of a podGroup and updates it. Returns the server's representation of the podGroup, and an error, if there is any.
func (c *podGroups) Update(podGroup *v1alpha1.PodGroup) (result *v1alpha1.PodGroup, err error) {
	result = &v1alpha1.PodGroup{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podgroups").
		Name(podGroup.Name).
		Body(podGroup).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *podGroups) UpdateStatus(podGroup *v1alpha1.PodGroup) (result *v1alpha1.PodGroup, err error) {
	result = &v1alpha1.PodGroup{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podgroups").
		Name(podGroup.Name).
		SubResource("status").
		Body(podGroup).
		Do().
		Into(result)
	return
}

// Delete takes name of the podGroup and deletes it. Returns an error if one occurs.
func (c *podGroups) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podgroups").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *podGroups) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podgroups").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched podGroup.
func (c *podGroups) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PodGroup, err error) {
	result = &v1alpha1.PodGroup{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("podgroups").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
