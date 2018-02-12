package yogsot

import (
	"errors"
	"reflect"

	"github.com/digitalocean/godo"
)

// Droplet is a struct that builds a droplet request.
type Droplet struct {
	Response *godo.Response
	Droplet  *godo.Droplet
	Request  *godo.DropletCreateRequest
}

func (d *Droplet) buildRequest(stackname string, resource map[string]interface{}) error {
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
					Fingerprint: fingerprint.(string),
				})
				id++
			}
			req.SSHKeys = keys
			continue
		}

		if k == "Volumes" {
			names := v.([]interface{})
			volumes := []godo.DropletCreateVolume{}
			for _, name := range names {
				volumes = append(volumes, godo.DropletCreateVolume{
					ID:   name.(string),
					Name: name.(string),
				})
			}
			req.Volumes = volumes
			continue
		}

		ref := reflect.ValueOf(req)
		val := reflect.Indirect(ref).FieldByName(k)
		if val == reflect.ValueOf(nil) {
			return errors.New("field not found: " + k)
		}
		val.Set(reflect.ValueOf(v))
	}
	req.Tags = append(req.Tags, stackname)
	d.Request = req
	return nil
}

func (d *Droplet) build(yogClient *YogClient) error {
	context := NewContext()
	droplet, response, err := yogClient.Droplets.Create(context, d.Request)
	if err != nil {
		return err
	}
	d.Response = response
	d.Droplet = droplet
	return nil
}
