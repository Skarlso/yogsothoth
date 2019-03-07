package yogsot

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/digitalocean/godo"
)

func TestLoadBalancerForwardingRules(t *testing.T) {
	template, err := ioutil.ReadFile("./testdata/load_test_TestLoadBalancerForwardingRules.yaml")
	if err != nil {
		t.Fatal("unexpected error while opening fixture: ", err)
	}
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
	template, err := ioutil.ReadFile("./testdata/load_test_TestLoadBalancerHealthChk.yaml")
	if err != nil {
		t.Fatal("unexpected error while opening fixture: ", err)
	}
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

func TestLoadBalancerInvalidTypeForValue(t *testing.T) {
	template, err := ioutil.ReadFile("./testdata/load_test_TestLoadBalancerInvalidTypeForValue.yaml")
	if err != nil {
		t.Fatal("unexpected error while opening fixture: ", err)
	}
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	lb := LoadBalancer{}
	err = lb.buildRequest("TestStack", response.Resources["LoadBalancer"])
	if err == nil {
		t.Fatal("expected error to be not nil")
	}
}
