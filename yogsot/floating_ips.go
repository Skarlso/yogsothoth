package yogsot

import "github.com/digitalocean/godo"

// FloatingIP is a struct which creates a floating ip create request.
type FloatingIP struct {
	Response   *godo.Response
	FloatingIP *godo.FloatingIP
	Request    *godo.FloatingIPCreateRequest
	Priority   int
}

func (fip *FloatingIP) build(yogClient *YogClient) error {
	return nil
}

func (fip *FloatingIP) buildRequest(stackname string, resource map[string]interface{}) error {
	return nil
}
