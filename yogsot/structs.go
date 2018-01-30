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

type Parameter struct {
	Default     string `yaml:"Default"`
	Type        string `yaml:"Type"`
	Description string `yaml:"Description"`
}

type Resource struct {
	Name string `yaml:"Name"`
	Type string `yaml:"Type"`
}

type createStackInput struct {
	Parameters map[string]Parameter `yaml:"Parameters"`
	Resources  map[string]Resource  `yaml:"Resources"`
}
