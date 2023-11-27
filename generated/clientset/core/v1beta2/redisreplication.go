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

// RedisReplicationLister helps list RedisReplications.
// All objects returned here must be treated as read-only.
type RedisReplicationLister interface {
	// List lists all RedisReplications in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.RedisReplication, err error)
	// RedisReplications returns an object that can list and get RedisReplications.
	RedisReplications(namespace string) RedisReplicationNamespaceLister
	RedisReplicationListerExpansion
}

// redisReplicationLister implements the RedisReplicationLister interface.
type redisReplicationLister struct {
	indexer cache.Indexer
}

// NewRedisReplicationLister returns a new RedisReplicationLister.
func NewRedisReplicationLister(indexer cache.Indexer) RedisReplicationLister {
	return &redisReplicationLister{indexer: indexer}
}

// List lists all RedisReplications in the indexer.
func (s *redisReplicationLister) List(selector labels.Selector) (ret []*v1beta2.RedisReplication, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.RedisReplication))
	})
	return ret, err
}

// RedisReplications returns an object that can list and get RedisReplications.
func (s *redisReplicationLister) RedisReplications(namespace string) RedisReplicationNamespaceLister {
	return redisReplicationNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RedisReplicationNamespaceLister helps list and get RedisReplications.
// All objects returned here must be treated as read-only.
type RedisReplicationNamespaceLister interface {
	// List lists all RedisReplications in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.RedisReplication, err error)
	// Get retrieves the RedisReplication from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta2.RedisReplication, error)
	RedisReplicationNamespaceListerExpansion
}

// redisReplicationNamespaceLister implements the RedisReplicationNamespaceLister
// interface.
type redisReplicationNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all RedisReplications in the indexer for a given namespace.
func (s redisReplicationNamespaceLister) List(selector labels.Selector) (ret []*v1beta2.RedisReplication, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.RedisReplication))
	})
	return ret, err
}

// Get retrieves the RedisReplication from the indexer for a given namespace and name.
func (s redisReplicationNamespaceLister) Get(name string) (*v1beta2.RedisReplication, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta2.Resource("redisreplication"), name)
	}
	return obj.(*v1beta2.RedisReplication), nil
}
