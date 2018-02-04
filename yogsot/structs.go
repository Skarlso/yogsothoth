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
	Name      string
	Error     error
	Resources []Resource
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

// Resource defines a resource which is able to build itself.
type Resource interface {
	build(string, map[string]interface{}) error
}

type createStackInput struct {
	Parameters map[string]Parameter              `yaml:"Parameters"`
	Resources  map[string]map[string]interface{} `yaml:"Resources"`
}
