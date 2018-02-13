package yogsot

import (
	"errors"

	"github.com/digitalocean/godo"
)

// FloatingIP is a struct which creates a floating ip create request.
type FloatingIP struct {
	Response    *godo.Response
	FloatingIP  *godo.FloatingIP
	Request     *godo.FloatingIPCreateRequest
	DropletName string
}

func (fip *FloatingIP) build(yogClient *YogClient) error {
	context := NewContext()
	floatingIP, response, err := yogClient.FloatingIPs.Create(context, fip.Request)
	if err != nil {
		return err
	}
	fip.Response = response
	fip.FloatingIP = floatingIP
	return nil
}

func (fip *FloatingIP) buildRequest(stackname string, resource map[string]interface{}) error {
	req := &godo.FloatingIPCreateRequest{}
	if v, ok := resource["Region"]; ok {
		req.Region = v.(string)
	} else {
		return errors.New("missing 'Region' key")
	}
	if v, ok := resource["DropletID"]; ok {
		if id, ok := v.(string); ok {
			fip.DropletName = id
			req.DropletID = -1
		} else if id, ok := v.(int); ok {
			req.DropletID = id
		}
	} else {
		return errors.New("missing DropletID key; set reference to a droplet name or set id to use")
	}
	fip.Request = req
	return nil
}

func (fip *FloatingIP) setDropletID(ID int) {
	if fip.Request.DropletID == -1 {
		fip.Request.DropletID = ID
	}
}
