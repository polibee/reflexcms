// Package gateways holds the payment gateway plugin surface. Every gateway
// self-registers at init() time; the pay module only talks to the Registry,
// so adding a provider = dropping a new file in this package.
package gateways

import "fmt"

// CreateRequest carries everything a gateway needs to open a payment.
type CreateRequest struct {
	OrderNo   string // our unique order number
	Title     string
	AmountCents int64
	Currency  string // ISO code, e.g. USD / CNY
	ReturnURL string // sync redirect after payment
	NotifyURL string // async webhook we expose
}

// CreateResult points the buyer at the payment (a hosted page URL today).
type CreateResult struct {
	PayURL string
	Raw    map[string]any
}

// Gateway is one payment provider plugin.
type Gateway interface {
	Name() string
	Enabled() bool
	Create(req CreateRequest) (CreateResult, error)
}

var registry = map[string]Gateway{}

// Register plugs a gateway into the pay module. Called from init().
func Register(g Gateway) {
	registry[g.Name()] = g
}

// Get returns a gateway by name.
func Get(name string) (Gateway, error) {
	g, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown payment gateway %q", name)
	}
	return g, nil
}

// All returns every registered gateway.
func All() []Gateway {
	out := make([]Gateway, 0, len(registry))
	for _, g := range registry {
		out = append(out, g)
	}
	return out
}
