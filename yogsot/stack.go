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
	// TODO: This has to be a priority / chain of initialization.
	builtResources := []Resource{}
	for _, v := range csi.Resources {
		d, err := buildResource(v["Type"].(string))
		if err != nil {
			return CreateStackResponse{}, err
		}
		err = d.build(request.StackName, v)
		if err != nil {
			return CreateStackResponse{}, err
		}
		builtResources = append(builtResources, d.(Droplet))
	}
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
