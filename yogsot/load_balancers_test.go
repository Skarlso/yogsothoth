package yogsot

import (
	"testing"
)

func TestLoadBalancerForwardingRules(t *testing.T) {
	template := []byte(`
Parameters:
  StackName:
    Description: The name of the stack to deploy
    Type: String
    Default: FurnaceStack
  Port:
    Description: Test port
    Type: Number
    Default: 80

Resources:
  LoadBalancer:
    Name: TestBalancer
    Algorithm: random
    Region: nyc3
    Tag: BalancerTest
    StickySessions: asdf
    RedirectHttpToHttps: true
    HealthCheck: asdf
    Type: LoadBalancer
    ForwardingRules:
      ForwardingRule1:
        EntryProtocol: garbage
      ForwardingRule2:
        EntryPort: 1234
    DropletIDs:
      - MyDroplet1
      - MyDroplet2`)
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	lb := LoadBalancer{}
	err = lb.buildRequest("TestStack", response.Resources["LoadBalancer"])
	if err != nil {
		t.Fatal("expected error to be nil. was: ", err)
	}
	if len(lb.Request.ForwardingRules) < 2 {
		t.Fatalf("lb forwarding rules count should have been greater than 2. was: %d", len(lb.Request.ForwardingRules))
	}
}
