package yogsot

import "github.com/digitalocean/godo"

// FloatingIP is a struct which creates a floating ip create request.
type FloatingIP struct {
	Response   *godo.Response
	FloatingIP *godo.FloatingIP
	Request    *godo.FloatingIPCreateRequest
	Priority   int
}

func (fip *FloatingIP) build(stackname string, yogClient *YogClient) error {
	return nil
}
