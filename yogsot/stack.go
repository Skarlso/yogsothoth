package yogsot

import (
	"errors"
	"fmt"
	"sync"
)

const (
	// FailedToParseYaml would indicate a syntax error in the build yaml file.
	FailedToParseYaml = iota
	// NoTypeForResource would indicate that the yaml is missing type information
	// for a specified resource.
	NoTypeForResource
	// UnknownResourceType indicates that resource defined in the yaml isn't known
	// to this library or hasn't been implemented yet.
	UnknownResourceType
	// FailedToBuildRequest indicates that there was a problem creating an object representation
	// of the yaml configuration for a given resource. The request here means that it tried to
	// construct a request for the godo library and failed in doing so. Probable cause could be
	// that there was a missing element in the yaml file that the request required to be there.
	// Consult the attached error message for more information.
	FailedToBuildRequest
	// ErrorLaunchingDroplets indicates that there was a problem launching a droplet or several droplets.
	// Consult the attached error message for more information.
	ErrorLaunchingDroplets
	// ErrorSettingUpDropletIDSForResources happens when Yog tries to setup IDs for the resources that
	// required them in the YAML through reference. Probable cause could be a missing link or invalid name
	// of the resource.
	ErrorSettingUpDropletIDSForResources
	// ErrorLaunchingResource indicates a problem while trying to launch any other resource
	// on DigitalOcean than a droplet.
	ErrorLaunchingResource
)

// CreateStackRequest create stack request.
type CreateStackRequest struct {
	TemplateBody []byte
	StackName    string
}

// CreateStackResponse create stack response.
type CreateStackResponse struct {
	Name      string
	Error     error
	Resources []Resource
}

// DeleteStackRequest delete stack request.
type DeleteStackRequest struct {
}

// DeleteStackResponse delete stack response.
type DeleteStackResponse struct {
	DeletedResources []Resource
}

// DescribeStackRequest describe stack request.
type DescribeStackRequest struct {
}

// DescribeStackResponse describe stack response.
type DescribeStackResponse struct {
}

// DropletError is an error that contains information about
// droplet launch faiulre such as, the name of the droplet
// and the failure reason.
type DropletError struct {
	DropletName string
	Error       error
	Message     string
}

// YogError is an error which accumulates multiple errors
// with different contexts. Like droplet start errors or
// errors with parsing yaml containing as much information
// as possible.
type YogError struct {
	Code          int
	Stackname     string
	DropletErrors []DropletError
	Error         error
	Message       string
}

// Droplets are maps of droplets with corresponding ids
type Droplets struct {
	droplets map[string]int
	sync.RWMutex
}

var droplets = Droplets{
	droplets: make(map[string]int, 0),
}

// AddDroplet locks the droplets and inserts a new id
func (d *Droplets) AddDroplet(droplet string, id int) {
	d.Lock()
	defer d.Unlock()
	d.droplets[droplet] = id
}

// GetID locks and returns the ID of a droplet
func (d *Droplets) GetID(droplet string) int {
	d.RLock()
	defer d.RUnlock()
	return d.droplets[droplet]
}

// CreateStack creates group of resources and logically bundles them together.
func (y *YogClient) CreateStack(request CreateStackRequest) (CreateStackResponse, YogError) {
	csi, err := parseTemplate(request.TemplateBody)
	if err != nil {
		ye := YogError{
			Code:          FailedToParseYaml,
			Message:       "error while parsing template",
			Error:         err,
			DropletErrors: []DropletError{},
			Stackname:     request.StackName,
		}
		return CreateStackResponse{}, ye
	}

	response := CreateStackResponse{Name: request.StackName, Error: nil}
	builtResources := []Resource{}
	for k, v := range csi.Resources {
		var s Service
		if _, ok := v["Type"]; !ok {
			message := fmt.Sprintf("no 'Type' provided for resource '%s'", k)
			ye := YogError{
				Code:          NoTypeForResource,
				Message:       message,
				Error:         errors.New(message),
				DropletErrors: []DropletError{},
				Stackname:     request.StackName,
			}
			return CreateStackResponse{}, ye
		}
		d, err := buildResource(s.Service(v["Type"].(string)))
		// Droplet doesn't yet have an ID. This will be updated once they are created.
		if err != nil {
			ye := YogError{
				Code:          UnknownResourceType,
				Message:       "",
				Error:         err,
				DropletErrors: []DropletError{},
				Stackname:     request.StackName,
			}
			return CreateStackResponse{}, ye
		}
		err = d.buildRequest(request.StackName, v)
		if err != nil {
			ye := YogError{
				Code:          FailedToBuildRequest,
				Message:       "",
				Error:         err,
				DropletErrors: []DropletError{},
				Stackname:     request.StackName,
			}
			return CreateStackResponse{}, ye
		}
		builtResources = append(builtResources, d)
	}

	de := y.launchAllDroplets(builtResources)
	if len(de) != 0 {
		ye := YogError{
			Code:          ErrorLaunchingDroplets,
			Stackname:     request.StackName,
			Error:         errors.New("error launching droplets. please see DropletErrors for more detail"),
			Message:       "",
			DropletErrors: de,
		}
		return CreateStackResponse{}, ye
	}
	err = y.setupDropletIDsForResources(builtResources)
	if err != nil {
		ye := YogError{
			Code:          ErrorSettingUpDropletIDSForResources,
			Stackname:     request.StackName,
			Error:         err,
			Message:       "error setting up droplet ids for resource",
			DropletErrors: []DropletError{},
		}
		return CreateStackResponse{}, ye
	}
	err = y.launchTheRestOfTheResources(builtResources)
	if err != nil {
		ye := YogError{
			Code:          ErrorLaunchingResource,
			Stackname:     request.StackName,
			Error:         err,
			Message:       "error while trying to launch the rest of the resources",
			DropletErrors: []DropletError{},
		}
		return CreateStackResponse{}, ye
	}

	response.Resources = builtResources
	return response, YogError{}
}

// DeleteStack deletes a given stack.
func (y *YogClient) DeleteStack(request DeleteStackRequest) (DeleteStackResponse, YogError) {
	ret := DeleteStackResponse{}

	return ret, YogError{}
}

// DescribeStack returns information about a created stack.
func (y *YogClient) DescribeStack(request DescribeStackRequest) (DescribeStackResponse, YogError) {
	return DescribeStackResponse{}, YogError{}
}

// launchAllDroplets goes through the resources and launches all the
// droplets concurrently. It uses a semaphore to limit the number
// of concurrent droplet launches. Currently that is hardcoded to 4.
func (y *YogClient) launchAllDroplets(droplets []Resource) []DropletError {
	dropletErrors := make([]DropletError, 0)
	sem := make(chan int, 4)
	var wg sync.WaitGroup
	var wl *sync.Mutex
	for _, v := range droplets {
		if d, ok := v.(*Droplet); ok {
			wg.Add(1)
			go func(d *Droplet, wl *sync.Mutex) {
				defer wg.Done()
				sem <- 1
				if err := y.launchDroplet(d); err != nil {
					de := DropletError{
						DropletName: d.Droplet.Name,
						Error:       err,
						Message:     "error while launching droplet",
					}
					wl.Lock()
					dropletErrors = append(dropletErrors, de)
					// This isn't deferred because of the defer wg.Done above
					wl.Unlock()
				}
				<-sem
			}(d, wl)
		}
	}
	wg.Wait()
	return dropletErrors
}

// launchDroplet launches a single droplet
func (y *YogClient) launchDroplet(droplet *Droplet) error {
	err := droplet.build(y)
	if err != nil {
		return err
	}
	droplets.AddDroplet(droplet.Request.Name, droplet.Droplet.ID)
	return nil
}

// setupDropletIDsForResources for each service there is a different way
// to provide droplet ids to use
func (y *YogClient) setupDropletIDsForResources(resources []Resource) error {
	for _, v := range resources {
		switch i := v.(type) {
		case *FloatingIP:
			i.setDropletID(droplets.GetID(i.DropletName))
		case *LoadBalancer:
			var ids []int
			for _, v := range i.DropletNames {
				ids = append(ids, droplets.GetID(v))
			}
			i.addDropletIDs(ids)
		case *Firewall:
			var fids []int
			for _, v := range i.DropletNames {
				fids = append(fids, droplets.GetID(v))
			}
			i.setFirewallDropletIDs(fids)
			for inboundName, names := range i.InboundDropletNames {
				var inIds []int
				for _, name := range names {
					inIds = append(inIds, droplets.GetID(name))
				}
				i.setInboundDropletIDs(i.InboundRequestsForName[inboundName], inIds)
			}

			for outboundName, names := range i.OutboundDropletNames {
				var outIds []int
				for _, name := range names {
					outIds = append(outIds, droplets.GetID(name))
				}
				i.setOutboundDropletIDs(i.OutboundRequestsForName[outboundName], outIds)
			}
		case *Domain:
		case *Droplet:
		default:
			s := fmt.Sprintf("unknown type %v", i)
			return errors.New(s)
		}
	}
	return nil
}

func (y *YogClient) launchTheRestOfTheResources(resources []Resource) error {
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
