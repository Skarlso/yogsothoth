package yogsot

func buildResource(T string) Resource {
	switch T {
	case "Droplet":
		return Droplet{}
	}
	return nil
}
