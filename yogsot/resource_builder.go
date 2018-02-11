package yogsot

import "errors"

func buildResource(T Service) (Resource, error) {
	switch T {
	case DROPLET:
		return &Droplet{Priority: 0}, nil
	case FLOATINGIP:
		return &FloatingIP{Priority: 1}, nil
	default:
		return nil, errors.New("unknown resource type: " + T.String())
	}
}
