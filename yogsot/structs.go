package yogsot

// Stack represents a collection of DigitalOcean assests.
type Stack struct {
}

// CreateStackRequest create stack request.
type CreateStackRequest struct {
	TemplateBody string
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

type createStackInput struct {
}
