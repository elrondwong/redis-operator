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
// Code generated by lister-gen. DO NOT EDIT.

package v1beta2

import (
	v1beta2 "github.com/elrondwong/redis-operator/api/v1beta2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// RedisSentinelLister helps list RedisSentinels.
// All objects returned here must be treated as read-only.
type RedisSentinelLister interface {
	// List lists all RedisSentinels in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.RedisSentinel, err error)
	// RedisSentinels returns an object that can list and get RedisSentinels.
	RedisSentinels(namespace string) RedisSentinelNamespaceLister
	RedisSentinelListerExpansion
}

// redisSentinelLister implements the RedisSentinelLister interface.
type redisSentinelLister struct {
	indexer cache.Indexer
}

// NewRedisSentinelLister returns a new RedisSentinelLister.
func NewRedisSentinelLister(indexer cache.Indexer) RedisSentinelLister {
	return &redisSentinelLister{indexer: indexer}
}

// List lists all RedisSentinels in the indexer.
func (s *redisSentinelLister) List(selector labels.Selector) (ret []*v1beta2.RedisSentinel, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.RedisSentinel))
	})
	return ret, err
}

// RedisSentinels returns an object that can list and get RedisSentinels.
func (s *redisSentinelLister) RedisSentinels(namespace string) RedisSentinelNamespaceLister {
	return redisSentinelNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RedisSentinelNamespaceLister helps list and get RedisSentinels.
// All objects returned here must be treated as read-only.
type RedisSentinelNamespaceLister interface {
	// List lists all RedisSentinels in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.RedisSentinel, err error)
	// Get retrieves the RedisSentinel from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta2.RedisSentinel, error)
	RedisSentinelNamespaceListerExpansion
}

// redisSentinelNamespaceLister implements the RedisSentinelNamespaceLister
// interface.
type redisSentinelNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all RedisSentinels in the indexer for a given namespace.
func (s redisSentinelNamespaceLister) List(selector labels.Selector) (ret []*v1beta2.RedisSentinel, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.RedisSentinel))
	})
	return ret, err
}

// Get retrieves the RedisSentinel from the indexer for a given namespace and name.
func (s redisSentinelNamespaceLister) Get(name string) (*v1beta2.RedisSentinel, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta2.Resource("redissentinel"), name)
	}
	return obj.(*v1beta2.RedisSentinel), nil
}
