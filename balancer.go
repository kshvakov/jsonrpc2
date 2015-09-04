package jsonrpc2

import (
	"sync"
	"time"
)

type Discovery interface {
	Get() ([]string, error)
}

func newBalancer(discovery Discovery) *balancer {

	b := balancer{
		discovery: discovery,
		mutex:     &sync.Mutex{},
	}

	if addresses, err := discovery.Get(); err == nil {

		b.addresses = addresses
	}

	go b.watch()

	return &b
}

type balancer struct {
	discovery   Discovery
	addresses   []string
	nextAddress int
	mutex       *sync.Mutex
}

func (b *balancer) len() int {

	b.mutex.Lock()

	defer b.mutex.Unlock()

	return len(b.addresses)
}

func (b *balancer) next() (string, error) {

	b.mutex.Lock()

	defer func() {

		b.nextAddress++
		b.mutex.Unlock()

	}()

	if len(b.addresses) == 0 {

		return "", &ErrorNoLiveUpstreams{}
	}

	if b.nextAddress > len(b.addresses)-1 {

		b.nextAddress = 0
	}

	return b.addresses[b.nextAddress], nil
}

func (b *balancer) watch() {

	for {

		<-time.After(time.Second)

		if addresses, err := b.discovery.Get(); err == nil {

			b.mutex.Lock()

			b.addresses = addresses

			b.mutex.Unlock()
		}
	}
}
