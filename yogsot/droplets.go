package yogsot

import (
	"errors"
	"reflect"

	"github.com/digitalocean/godo"
)

// Droplet is a struct that builds a droplet request.
type Droplet struct {
	Response *godo.Response
}

func (d Droplet) buildRequest(stackname string, resource map[string]interface{}) (*godo.DropletCreateRequest, error) {
	req := &godo.DropletCreateRequest{}
	for k, v := range resource {
		if k == "Type" {
			continue
		}
		if k == "Image" {
			req.Image = godo.DropletCreateImage{
				Slug: v.(map[interface{}]interface{})["Slug"].(string),
			}
			continue
		}
		if k == "SSHKeys" {
			fingerprints := v.([]interface{})
			keys := []godo.DropletCreateSSHKey{}
			id := 0
			for _, fingerprint := range fingerprints {
				keys = append(keys, godo.DropletCreateSSHKey{
					ID:          id,
					Fingerprint: fingerprint.(map[interface{}]interface{})["Fingerprint"].(string),
				})
				id++
			}
			req.SSHKeys = keys
			continue
		}

		ref := reflect.ValueOf(req)
		val := reflect.Indirect(ref).FieldByName(k)
		if val == reflect.ValueOf(nil) {
			return req, errors.New("field not found: " + k)
		}
		val.Set(reflect.ValueOf(v))
	}
	req.Tags = []string{stackname}
	return req, nil
}

func (d Droplet) build(stackname string, resource map[string]interface{}) error {
	_, err := d.buildRequest(stackname, resource)
	if err != nil {
		return err
	}
	// build droplet with client here
	return nil
}

func createDroplet(request *godo.DropletCreateRequest) {

}
