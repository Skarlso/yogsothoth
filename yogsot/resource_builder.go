package yogsot

import "errors"

func buildResource(T string) (Resource, error) {
	switch T {
	case "Droplet":
		return Droplet{}, nil
	default:
		return nil, errors.New("unknown resource type: " + T)
	}
}
