package yogsot

import (
	"errors"
)

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
	// When that is done, create the rest of the services belonging to that droplet.
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
