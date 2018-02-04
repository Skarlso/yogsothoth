package yogsot

// Stack represents a collection of DigitalOcean assests.
type Stack struct {
	Name string
}

// CreateStackRequest create stack request.
type CreateStackRequest struct {
	TemplateBody []byte
	StackName    string
}

// CreateStackResponse create stack response.
type CreateStackResponse struct {
}

// DeleteStackRequest delete stack request.
type DeleteStackRequest struct {
}

// DeleteStackResponse delete stack response.
type DeleteStackResponse struct {
}

// DescribeStackRequest describe stack request.
type DescribeStackRequest struct {
}

// DescribeStackResponse describe stack response.
type DescribeStackResponse struct {
}

// Parameter are the variables of the stack.
type Parameter struct {
	Default     string `yaml:"Default"`
	Type        string `yaml:"Type"`
	Description string `yaml:"Description"`
}

// // Resource is parsed as a map[string][interface{}] so that it can be
// // determined what kind of resource it is and parse it further accordingly.
// type Resource struct {
// }

type Resource interface {
	build(data map[string]interface{}) (interface{}, error)
}

type createStackInput struct {
	Parameters map[string]Parameter              `yaml:"Parameters"`
	Resources  map[string]map[string]interface{} `yaml:"Resources"`
}
