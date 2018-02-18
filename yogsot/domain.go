package yogsot

import (
	"errors"

	"github.com/digitalocean/godo"
)

// Domain is a struct which creates a Domain create request.
type Domain struct {
	Response *godo.Response
	Domain   *godo.Domain
	Request  *godo.DomainCreateRequest
}

func (d *Domain) build(yogClient *YogClient) error {
	context := NewContext()
	domain, response, err := yogClient.Domains.Create(context, d.Request)
	if err != nil {
		return err
	}
	d.Response = response
	d.Domain = domain
	return nil
}

func (d *Domain) buildRequest(stackname string, resource map[string]interface{}) error {
	req := &godo.DomainCreateRequest{}
	if v, ok := resource["Name"]; ok {
		req.Name = v.(string)
	} else {
		return errors.New("missing 'Name' key")
	}
	if v, ok := resource["IPAddress"]; ok {
		req.IPAddress = v.(string)
	} else {
		return errors.New("missing 'IPAddress' key")
	}
	d.Request = req
	return nil
}
