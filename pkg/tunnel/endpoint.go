package tunnel

import (
	"fmt"
	"strconv"
	"strings"
)

type Endpoint struct {
	User string
	Host string
	Port int
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

func NewEndpoint(s string) *Endpoint {

	ep := &Endpoint{
		Host: s,
	}

	if parts := strings.Split(ep.Host, "@"); len(parts) > 1 {
		ep.User = parts[0]
		ep.Host = parts[1]
	}

	if parts := strings.Split(ep.Host, ":"); len(parts) > 1 {
		ep.Host = parts[0]
		ep.Port, _ = strconv.Atoi(parts[1])
	}

	return ep
}
