package jsonrpc2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testDiscovery struct {
	addresses []string
	err       error
}

func (t *testDiscovery) Get() ([]string, error) {

	if t.err != nil {

		return nil, t.err
	}

	return t.addresses, nil
}

func TestBalancerNext(t *testing.T) {

	addresses := []string{"a", "b", "c"}

	b := newBalancer(&testDiscovery{addresses: addresses})

	if assert.True(t, b.len() == 3) {

		for _, address := range addresses {

			if a, err := b.next(); assert.NoError(t, err) {

				assert.Equal(t, address, a)
			}
		}

		for _, address := range addresses {

			if a, err := b.next(); assert.NoError(t, err) {

				assert.Equal(t, address, a)
			}
		}
	}
}

func TestBalancerErrorNoLiveUpstreams(t *testing.T) {

	b := newBalancer(&testDiscovery{})

	if assert.True(t, b.len() == 0) {

		_, err := b.next()

		if assert.Error(t, err) {

			_, ok := err.(*ErrorNoLiveUpstreams)

			assert.True(t, ok)
		}
	}
}
