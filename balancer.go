package jsonrpc2

type balancer struct {
	discovery Discovery
}

func (b *balancer) Next() (string, error) {

	return "", nil
}
