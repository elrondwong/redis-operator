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

// FakeRedises implements RedisInterface
type FakeRedises struct {
	Fake *FakeV1beta2
	ns   string
}

var redisResource = v1beta2.SchemeGroupVersion.WithResource("redis")

var redisesKind = v1beta2.SchemeGroupVersion.WithKind("Redis")

// Get takes name of the redis, and returns the corresponding redis object, and an error if there is any.
func (c *FakeRedises) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta2.Redis, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(redisesResource, c.ns, name), &v1beta2.Redis{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.Redis), err
}

// List takes label and field selectors, and returns the list of Redises that match those selectors.
func (c *FakeRedises) List(ctx context.Context, opts v1.ListOptions) (result *v1beta2.RedisList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(redisesResource, redisesKind, c.ns, opts), &v1beta2.RedisList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta2.RedisList{ListMeta: obj.(*v1beta2.RedisList).ListMeta}
	for _, item := range obj.(*v1beta2.RedisList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested redises.
func (c *FakeRedises) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(redisesResource, c.ns, opts))

}

// Create takes the representation of a redis and creates it.  Returns the server's representation of the redis, and an error, if there is any.
func (c *FakeRedises) Create(ctx context.Context, redis *v1beta2.Redis, opts v1.CreateOptions) (result *v1beta2.Redis, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(redisesResource, c.ns, redis), &v1beta2.Redis{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.Redis), err
}

// Update takes the representation of a redis and updates it. Returns the server's representation of the redis, and an error, if there is any.
func (c *FakeRedises) Update(ctx context.Context, redis *v1beta2.Redis, opts v1.UpdateOptions) (result *v1beta2.Redis, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(redisesResource, c.ns, redis), &v1beta2.Redis{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.Redis), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRedises) UpdateStatus(ctx context.Context, redis *v1beta2.Redis, opts v1.UpdateOptions) (*v1beta2.Redis, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(redisesResource, "status", c.ns, redis), &v1beta2.Redis{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.Redis), err
}

// Delete takes name of the redis and deletes it. Returns an error if one occurs.
func (c *FakeRedises) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(redisesResource, c.ns, name, opts), &v1beta2.Redis{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRedises) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(redisesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta2.RedisList{})
	return err
}

// Patch applies the patch and returns the patched redis.
func (c *FakeRedises) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta2.Redis, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(redisesResource, c.ns, name, pt, data, subresources...), &v1beta2.Redis{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.Redis), err
}
