package yogsot

import (
	"errors"
)

// CreateStack creates group of resources and logically bundles them together.
func (y *YogClient) CreateStack(request CreateStackRequest) (CreateStackResponse, error) {
	_, err := parseTemplate(request.TemplateBody)
	if err != nil {
		return CreateStackResponse{}, errors.New("error while parsing tempalte: " + err.Error())
	}
	// for _, v := range req.Resources {
	// 	fmt.Println("Received the following resources:")
	// 	// fmt.Printf("Name: %s, Type: %s\n", v.Name, v.Type)
	// }
	return CreateStackResponse{}, nil
}

// DeleteStack deletes a given stack.
func (y *YogClient) DeleteStack(request DeleteStackRequest) (DeleteStackResponse, error) {
	return DeleteStackResponse{}, nil
}

// DescribeStack returns information about a created stack.
func (y *YogClient) DescribeStack(request DescribeStackRequest) (DescribeStackResponse, error) {
	return DescribeStackResponse{}, nil
}
