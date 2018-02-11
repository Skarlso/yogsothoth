package yogsot

// Service represents a DigitalOcean asset that can be created
type Service int

const (
	// DROPLET is the main building block on which all other services rely
	DROPLET Service = iota
	// FLOATINGIP is a static IP which can be assigned to a droplet
	FLOATINGIP
	// FIREWALL ability to restrict network access to and from a droplet
	FIREWALL
	// IMAGE may either be a: snapshot, backup, application image
	IMAGE
	// LOADBALANCER distribute traffic across multiple droplet
	LOADBALANCER
	// DOMAINRECORD is names purchased from a domain name registrar
	DOMAINRECORD
)

func (s Service) String() string {
	switch s {
	case DROPLET:
		return "Droplet"
	case FLOATINGIP:
		return "Floating IP"
	case FIREWALL:
		return "FireWall"
	case IMAGE:
		return "Image"
	case LOADBALANCER:
		return "Load Balancer"
	case DOMAINRECORD:
		return "Domain Record"
	default:
		return "Unkown Type"
	}
}

// Service creates a service out of a String definition of a service
func (s Service) Service(T string) Service {
	switch T {
	case "Droplet":
		return DROPLET
	case "FloatingIP":
		return FLOATINGIP
	case "Firewall":
		return FIREWALL
	case "Image":
		return IMAGE
	case "LoadBalancer":
		return LOADBALANCER
	case "DomainRecord":
		return DOMAINRECORD
	default:
		return 999
	}
}

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
	Resources []interface{}
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
	build(string, *YogClient) error
}

type createStackInput struct {
	Parameters map[string]Parameter              `yaml:"Parameters"`
	Resources  map[string]map[string]interface{} `yaml:"Resources"`
}
