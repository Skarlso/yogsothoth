package yogsot

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/digitalocean/godo"
)

func TestSimpleFirewallRules(t *testing.T) {
	template, err := ioutil.ReadFile("./fixtures/firewall_test_TestSimpleFirewallRules.yaml")
	if err != nil {
		t.Fatal("unexpected error while opening fixture: ", err)
	}
	response, err := parseTemplate(template)
	if err != nil {
		t.Fatal("failed with error: ", err)
	}
	fw := Firewall{}
	err = fw.buildRequest("TestStack", response.Resources["FireWall"])
	if err != nil {
		t.Fatal("expected error to be nil. was: ", err)
	}
	expectedInbound := godo.InboundRule{
		PortRange: "2345:2345",
		Protocol:  "ProtocolValue",
		Sources: &godo.Sources{
			Addresses:        []string{"skarlso.io", "nagios.skarlso.io"},
			Tags:             []string{"Multiple", "Tags"},
			DropletIDs:       []int{1234},
			LoadBalancerUIDs: []string{"UID1", "UID2"},
		},
	}
	expectedOutbound := godo.OutboundRule{
		PortRange: "1234:2345",
		Protocol:  "ProtocolValue2",
		Destinations: &godo.Destinations{
			Addresses:        []string{"skarlso.io", "nagios.skarlso.io"},
			Tags:             []string{"Outbound", "Tags"},
			DropletIDs:       []int{4321},
			LoadBalancerUIDs: []string{"UID3", "UID4"},
		},
	}
	if !reflect.DeepEqual(fw.Request.InboundRules[0], expectedInbound) {
		t.Fatalf("Expected %+v did not equal actual %+v", expectedInbound, fw.Request.InboundRules[0])
	}
	if !reflect.DeepEqual(fw.Request.OutboundRules[0], expectedOutbound) {
		t.Fatalf("Expected %+v did not equal actual %+v", expectedOutbound, fw.Request.OutboundRules[0])
	}
	expectedDropletNames := []string{"MyDroplet", "NotMyDroplet"}
	if !reflect.DeepEqual(fw.DropletNames, expectedDropletNames) {
		t.Fatalf("Expected %+v did not equal actual %+v", fw.DropletNames, expectedDropletNames)
	}
}
