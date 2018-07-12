package yogsot

import "errors"

// Service represents a DigitalOcean asset that can be created
type Service int

const (
	// DROPLET is the main building block on which all other services rely
	DROPLET Service = iota
	// FLOATINGIP is a static IP which can be assigned to a droplet
	FLOATINGIP
	// FIREWALL ability to restrict network access to and from a droplet
	FIREWALL
	// IMAGE may either be a: snapshot, backup, application image
	IMAGE
	// LOADBALANCER distribute traffic across multiple droplet
	LOADBALANCER
	// DOMAIN is names purchased from a domain name registrar
	DOMAIN
)

func (s Service) String() string {
	switch s {
	case DROPLET:
		return "Droplet"
	case FLOATINGIP:
		return "Floating IP"
	case FIREWALL:
		return "FireWall"
	case IMAGE:
		return "Image"
	case LOADBALANCER:
		return "Load Balancer"
	case DOMAIN:
		return "Domain"
	default:
		return "Unknown Type"
	}
}

// Service creates a service out of a String definition of a service
func (s Service) Service(T string) Service {
	switch T {
	case "Droplet":
		return DROPLET
	case "FloatingIP":
		return FLOATINGIP
	case "Firewall":
		return FIREWALL
	case "Image":
		return IMAGE
	case "LoadBalancer":
		return LOADBALANCER
	case "Domain":
		return DOMAIN
	default:
		return 999
	}
}

// Resource defines a resource which is able to build itself.
type Resource interface {
	build(*YogClient) error
	buildRequest(string, map[string]interface{}) error
}

func buildResource(T Service) (Resource, error) {
	switch T {
	case DROPLET:
		return new(Droplet), nil
	case FLOATINGIP:
		return new(FloatingIP), nil
	case LOADBALANCER:
		return new(LoadBalancer), nil
	case DOMAIN:
		return new(Domain), nil
	case FIREWALL:
		return new(Firewall), nil
	default:
		return nil, errors.New("unknown resource type: " + T.String())
	}
}
