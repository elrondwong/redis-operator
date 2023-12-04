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
// Code generated by informer-gen. DO NOT EDIT.

package v1beta2

import (
	"context"
	time "time"

	apiv1beta2 "github.com/elrondwong/redis-operator/api/v1beta2"
	internalinterfaces "github.com/elrondwong/redis-operator/generated/clientset/externalversions/internalinterfaces"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	clientset "github.com/elrondwong/redis-operator/generated/clientset/versioned"
	v1beta2 "github.com/elrondwong/redis-operator/generated/listers/core/v1beta2"
)

// RedisSentinelInformer provides access to a shared informer and lister for
// RedisSentinels.
type RedisSentinelInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta2.RedisSentinelLister
}

type redisSentinelInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewRedisSentinelInformer constructs a new informer for RedisSentinel type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewRedisSentinelInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredRedisSentinelInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredRedisSentinelInformer constructs a new informer for RedisSentinel type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredRedisSentinelInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1beta2().RedisSentinels(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1beta2().RedisSentinels(namespace).Watch(context.TODO(), options)
			},
		},
		&apiv1beta2.RedisSentinel{},
		resyncPeriod,
		indexers,
	)
}

func (f *redisSentinelInformer) defaultInformer(client clientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredRedisSentinelInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *redisSentinelInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apiv1beta2.RedisSentinel{}, f.defaultInformer)
}

func (f *redisSentinelInformer) Lister() v1beta2.RedisSentinelLister {
	return v1beta2.NewRedisSentinelLister(f.Informer().GetIndexer())
}