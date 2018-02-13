package yogsot

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

// Droplets are maps of droplets with corresponding ids
type Droplets struct {
	droplets map[string]int
	sync.Mutex
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
		// Droplet doesn't yet have an ID. This will be updated once they are created.
		if err != nil {
			return CreateStackResponse{}, err
		}
		err = d.buildRequest(request.StackName, v)
		if err != nil {
			return CreateStackResponse{}, err
		}
		builtResources = append(builtResources, d)
	}

	y.launchAllDroplets(builtResources)
	y.setupDropletIDsForResources(builtResources)
	y.launchTheRestOfTheResources(builtResources)

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

// launchAllDroplets goes through the resources and launches all the
// droplets concurrently. It uses a semaphore to limit the number
// of concurrent droplet launches. Currently that is hardcoded to 4.
func (y *YogClient) launchAllDroplets(droplets []interface{}) {
	sem := make(chan int, 4)
	var wg sync.WaitGroup
	for _, v := range droplets {
		if d, ok := v.(*Droplet); ok {
			wg.Add(1)
			go func(d *Droplet) {
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

// launchDroplet launches a single droplet
func (y *YogClient) launchDroplet(droplet *Droplet) error {
	droplets.Lock()
	defer droplets.Unlock()
	log.Println("Launching droplet.")
	err := droplet.build(y)
	if err != nil {
		return err
	}
	// if droplet.Droplet != nil {
	droplets.droplets[droplet.Request.Name] = droplet.Droplet.ID
	// }
	return nil
}

// setupDropletIDsForResources for each service there is a different way
// to provide droplet ids to use
func (y *YogClient) setupDropletIDsForResources(resources []interface{}) error {
	for _, v := range resources {
		switch i := v.(type) {
		case *FloatingIP:
			i.setDropletID(droplets.droplets[i.DropletName])
		case *Droplet:
		default:
			s := fmt.Sprintf("unknown type %v", i)
			return errors.New(s)
		}
	}
	return nil
}

func (y *YogClient) launchTheRestOfTheResources(resources []interface{}) error {
	for _, v := range resources {
		if _, ok := v.(*Droplet); ok {
			continue
		}
		err := v.(Resource).build(y)
		if err != nil {
			return err
		}
	}

	return nil
}
