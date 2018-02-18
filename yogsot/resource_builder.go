package yogsot

import "errors"

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
	default:
		return nil, errors.New("unknown resource type: " + T.String())
	}
}
