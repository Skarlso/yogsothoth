package yogsot

import (
	"errors"
	"log"
	"sync"
)

// Droplets are maps of droplets with corresponding ids
type Droplets struct {
	droplets map[string]int
	*sync.Mutex
}

var droplets = Droplets{
	droplets: make(map[string]int, 0),
}

// CreateStack creates group of resources and logically bundles them together.
func (y *YogClient) CreateStack(request CreateStackRequest) (CreateStackResponse, error) {
	csi, err := parseTemplate(request.TemplateBody)
	if err != nil {
		return CreateStackResponse{}, errors.New("error while parsing tempalte: " + err.Error())
	}

	response := CreateStackResponse{Name: request.StackName, Error: nil}
	builtResources := []interface{}{}
	for _, v := range csi.Resources {
		var s Service
		d, err := buildResource(s.Service(v["Type"].(string)))
		if err != nil {
			return CreateStackResponse{}, err
		}
		builtResources = append(builtResources, d)

		// r, err := d.build(request.StackName, v, y)
		// if err != nil {
		// 	return CreateStackResponse{}, err
		// }
	}

	// There can be many droplet assigned to many services.
	// Need a way ~!Ref~ to tie a droplet to a service.
	// Once located, create the droplet and save its ID.
	// ID Is saved by map[string]int -> {Droplet: 1}
	// Later on !Ref: Droplet will contain the name of the resource.
	// This will be easy because dropletIds will have the name like:
	// dropletIds[refName] -> 1
	y.launchAllDroplets(builtResources)
	response.Resources = builtResources
	return response, nil
}

// DeleteStack deletes a given stack.
func (y *YogClient) DeleteStack(request DeleteStackRequest) (DeleteStackResponse, error) {
	return DeleteStackResponse{}, nil
}

// DescribeStack returns information about a created stack.
func (y *YogClient) DescribeStack(request DescribeStackRequest) (DescribeStackResponse, error) {
	return DescribeStackResponse{}, nil
}

func (y *YogClient) launchAllDroplets(droplets []interface{}) {
	sem := make(chan int, 4)
	var wg sync.WaitGroup
	for _, v := range droplets {
		if d, ok := v.(Droplet); ok {
			wg.Add(1)
			go func(d Droplet) {
				defer wg.Done()
				sem <- 1
				if err := y.launchDroplet(d); err != nil {
					log.Fatal("Error while launching droplet: ", d.Droplet.Name)
				}
				<-sem
			}(d)
		}
	}
	wg.Wait()
}

func (y *YogClient) launchDroplet(droplet Droplet) error {
	droplets.Lock()
	defer droplets.Unlock()
	log.Println("Launching droplet.")
	err := droplet.build("stackname", y)
	if err != nil {
		return err
	}
	droplets.droplets[droplet.Droplet.Name] = droplet.Droplet.ID
	return nil
}
