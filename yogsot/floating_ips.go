package yogsot

import "github.com/digitalocean/godo"

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
	req.Region = resource["Region"].(string)
	fip.Request = req
	fip.DropletName = resource["DropletID"].(string)
	return nil
}

func (fip *FloatingIP) setDropletID(ID int) {
	fip.Request.DropletID = ID
}
