/*
Copyright 2020 Opstree Solutions.

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

package fake

import (
	"context"

	v1beta2 "github.com/elrondwong/redis-operator/api/v1beta2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRedisReplications implements RedisReplicationInterface
type FakeRedisReplications struct {
	Fake *FakeV1beta2
	ns   string
}

var redisreplicationsResource = v1beta2.SchemeGroupVersion.WithResource("redisreplications")

var redisreplicationsKind = v1beta2.SchemeGroupVersion.WithKind("RedisReplication")

// Get takes name of the redisReplication, and returns the corresponding redisReplication object, and an error if there is any.
func (c *FakeRedisReplications) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta2.RedisReplication, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(redisreplicationsResource, c.ns, name), &v1beta2.RedisReplication{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.RedisReplication), err
}

// List takes label and field selectors, and returns the list of RedisReplications that match those selectors.
func (c *FakeRedisReplications) List(ctx context.Context, opts v1.ListOptions) (result *v1beta2.RedisReplicationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(redisreplicationsResource, redisreplicationsKind, c.ns, opts), &v1beta2.RedisReplicationList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta2.RedisReplicationList{ListMeta: obj.(*v1beta2.RedisReplicationList).ListMeta}
	for _, item := range obj.(*v1beta2.RedisReplicationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested redisReplications.
func (c *FakeRedisReplications) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(redisreplicationsResource, c.ns, opts))

}

// Create takes the representation of a redisReplication and creates it.  Returns the server's representation of the redisReplication, and an error, if there is any.
func (c *FakeRedisReplications) Create(ctx context.Context, redisReplication *v1beta2.RedisReplication, opts v1.CreateOptions) (result *v1beta2.RedisReplication, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(redisreplicationsResource, c.ns, redisReplication), &v1beta2.RedisReplication{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.RedisReplication), err
}

// Update takes the representation of a redisReplication and updates it. Returns the server's representation of the redisReplication, and an error, if there is any.
func (c *FakeRedisReplications) Update(ctx context.Context, redisReplication *v1beta2.RedisReplication, opts v1.UpdateOptions) (result *v1beta2.RedisReplication, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(redisreplicationsResource, c.ns, redisReplication), &v1beta2.RedisReplication{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.RedisReplication), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRedisReplications) UpdateStatus(ctx context.Context, redisReplication *v1beta2.RedisReplication, opts v1.UpdateOptions) (*v1beta2.RedisReplication, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(redisreplicationsResource, "status", c.ns, redisReplication), &v1beta2.RedisReplication{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.RedisReplication), err
}

// Delete takes name of the redisReplication and deletes it. Returns an error if one occurs.
func (c *FakeRedisReplications) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(redisreplicationsResource, c.ns, name, opts), &v1beta2.RedisReplication{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRedisReplications) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(redisreplicationsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta2.RedisReplicationList{})
	return err
}

// Patch applies the patch and returns the patched redisReplication.
func (c *FakeRedisReplications) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta2.RedisReplication, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(redisreplicationsResource, c.ns, name, pt, data, subresources...), &v1beta2.RedisReplication{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.RedisReplication), err
}
