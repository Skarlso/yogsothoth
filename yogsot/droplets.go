package yogsot

import "github.com/digitalocean/godo"

// Droplet A single DO droplet representation.
type Droplet struct {
	*godo.Droplet
}

func (d *Droplet) createDroplet(request *godo.DropletCreateRequest) {

}
