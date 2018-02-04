package yogsot

func buildResource(stackname string, T string, resource map[string]interface{}) Resource {

	switch T {
	case "Droplet":
		return Droplet{}
	}
	return nil
}
