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
	expectedDropletNames := []string{"FwDroplet"}
	if !reflect.DeepEqual(fw.DropletNames, expectedDropletNames) {
		t.Fatalf("Expected %+v did not equal actual %+v", fw.DropletNames, expectedDropletNames)
	}
	expectedInboundDropletNames := make(map[string][]string, 0)
	expectedInboundDropletNames["Inbound1"] = []string{"MyDroplet"}
	expectedOutboundDropletNames := make(map[string][]string, 0)
	expectedOutboundDropletNames["Outbound1"] = []string{"NotMyDroplet"}
	if !reflect.DeepEqual(fw.InboundDropletNames, expectedInboundDropletNames) {
		t.Fatalf("Expected %+v did not equal actual %+v", fw.InboundDropletNames, expectedInboundDropletNames)
	}
	if !reflect.DeepEqual(fw.OutboundDropletNames, expectedOutboundDropletNames) {
		t.Fatalf("Expected %+v did not equal actual %+v", fw.OutboundDropletNames, expectedOutboundDropletNames)
	}
}

func TestSettingDropletIdsForIndividualRules(t *testing.T) {
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
	d := Droplets{
		droplets: make(map[string]int, 0),
	}
	d.droplets["MyDroplet"] = 999
	d.droplets["NotMyDroplet"] = 998
	d.droplets["FwDroplet"] = 997
	fids := make([]int, 0)
	for _, v := range fw.DropletNames {
		fids = append(fids, d.GetID(v))
	}
	fw.setFirewallDropletIDs(fids)
	for inboundName, names := range fw.InboundDropletNames {
		var inIds []int
		for _, name := range names {
			inIds = append(inIds, d.GetID(name))
		}
		fw.setInboundDropletIDs(fw.InboundRequestsForName[inboundName], inIds)
	}
	for outboundName, names := range fw.OutboundDropletNames {
		var outIds []int
		for _, name := range names {
			outIds = append(outIds, d.GetID(name))
		}
		fw.setOutboundDropletIDs(fw.OutboundRequestsForName[outboundName], outIds)
	}
	expectedInIDs := []int{1234, 999}
	expectedOutIDs := []int{4321, 998}
	expectedFwIDs := []int{12, 997}
	if !reflect.DeepEqual(expectedFwIDs, fw.Request.DropletIDs) {
		t.Fatalf("Expected %+v did not equal actual %+v", expectedFwIDs, fw.Request.DropletIDs)
	}
	if !reflect.DeepEqual(expectedInIDs, fw.Request.InboundRules[0].Sources.DropletIDs) {
		t.Fatalf("Expected %+v did not equal actual %+v", expectedInIDs, fw.Request.InboundRules[0].Sources.DropletIDs)
	}
	if !reflect.DeepEqual(expectedOutIDs, fw.Request.OutboundRules[0].Destinations.DropletIDs) {
		t.Fatalf("Expected %+v did not equal actual %+v", expectedOutIDs, fw.Request.OutboundRules[0].Destinations.DropletIDs)
	}
}
