package yogsot

import (
	"reflect"
	"testing"

	"github.com/digitalocean/godo"
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
    RedirectHttpToHttps: true
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

func TestLoadBalancerHealthChk(t *testing.T) {
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
    RedirectHttpToHttps: true
    HealthCheck:
      Protocol: proto
      Port: 1234
      Path: /health
      CheckIntervalSeconds: 15
      ResponseTimeoutSeconds: 150
      HealthyThreshold: 2
      UnhealthyThreshold: 4
    Type: LoadBalancer
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
	if lb.Request.HealthCheck == nil {
		t.Fatalf("lb healtcheck was nil")
	}
	expected := &godo.HealthCheck{}
	expected.Path = "/health"
	expected.Port = 1234
	expected.Protocol = "proto"
	expected.CheckIntervalSeconds = 15
	expected.ResponseTimeoutSeconds = 150
	expected.HealthyThreshold = 2
	expected.UnhealthyThreshold = 4
	if !reflect.DeepEqual(expected, lb.Request.HealthCheck) {
		t.Fatalf("expected healthcheck didn't match actual. actual: %+v", lb.Request.HealthCheck)
	}
}
