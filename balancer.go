package jsonrpc2

type balancer struct{}

func (b *balancer) Next() (string, bool) {

	return "", false
}
