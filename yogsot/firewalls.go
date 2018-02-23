package yogsot

import (
	"github.com/digitalocean/godo"
)

// Firewall is a struct for firewall creation requests.
type Firewall struct {
	Response             *godo.Response
	Firewall             *godo.Firewall
	Request              *godo.FirewallRequest
	DropletNames         []string
	InboundDropletNames  map[string][]string
	OutboundDropletNames map[string][]string
}

func (fw *Firewall) buildRequest(stackname string, resource map[string]interface{}) error {
	fw.DropletNames = make([]string, 0)
	fw.InboundDropletNames = make(map[string][]string, 0)
	fw.OutboundDropletNames = make(map[string][]string, 0)
	req := &godo.FirewallRequest{}
	for k, v := range resource {
		if k == "Type" {
			continue
		}

		switch k {
		case "Name":
			req.Name = v.(string)
		case "InboundRules":
			inboundRules := []godo.InboundRule{}
			for inKey, value := range v.(map[interface{}]interface{}) {
				iRule, dropletNames := getIRule(value.(map[interface{}]interface{}))
				fw.InboundDropletNames[inKey.(string)] = dropletNames
				inboundRules = append(inboundRules, iRule)
			}
			req.InboundRules = inboundRules
		case "OutboundRules":
			outboundRules := []godo.OutboundRule{}
			for outKey, value := range v.(map[interface{}]interface{}) {
				oRule, dropletNames := getORule(value.(map[interface{}]interface{}))
				fw.OutboundDropletNames[outKey.(string)] = dropletNames
				outboundRules = append(outboundRules, oRule)
			}
			req.OutboundRules = outboundRules
		case "Tags":
		case "DropletIDs":
			for _, value := range v.([]interface{}) {
				switch o := value.(type) {
				case string:
					fw.DropletNames = append(fw.DropletNames, o)
				case int:
					req.DropletIDs = append(req.DropletIDs, o)
				}
			}
		}
	}
	fw.Request = req
	return nil
}

func (fw *Firewall) build(yogClient *YogClient) error {
	context := NewContext()
	firewall, response, err := yogClient.Firewalls.Create(context, fw.Request)
	if err != nil {
		return err
	}
	fw.Response = response
	fw.Firewall = firewall
	return nil
}

func (fw *Firewall) setFirewallDropletIDs(ids []int) {
	fw.Request.DropletIDs = append(fw.Request.DropletIDs, ids...)
}

func (fw *Firewall) setInboundDropletIDs(ids []int) {
}

func (fw *Firewall) setOutboundDropletIDs(ids []int) {
}

func getIRule(v map[interface{}]interface{}) (iRule godo.InboundRule, dropletNames []string) {
	for key, innerV := range v {
		if key == "Protocol" {
			iRule.Protocol = innerV.(string)
		}
		if key == "PortRange" {
			iRule.PortRange = innerV.(string)
		}
		if key == "Sources" {
			dropletNames, iRule.Sources = getSources(innerV.(map[interface{}]interface{}))
		}
	}
	return
}

func getORule(v map[interface{}]interface{}) (oRule godo.OutboundRule, dropletNames []string) {
	for key, innerV := range v {
		if key == "Protocol" {
			oRule.Protocol = innerV.(string)
		}
		if key == "PortRange" {
			oRule.PortRange = innerV.(string)
		}
		if key == "Sources" {
			dropletNames, oRule.Destinations = getDestinations(innerV.(map[interface{}]interface{}))
		}
	}
	return
}

// TODO: Find a way to remove this duplication and tiet coupling to the
// design of the sources struct.
// Maybe submit a PR to remove the un-needed extra Destination struct.
// Which is literally the same as Sources.
func getSources(sources map[interface{}]interface{}) ([]string, *godo.Sources) {
	ret := &godo.Sources{}
	dropletNames, source := getSource(sources)
	ret.Addresses = source["Addresses"].([]string)
	ret.DropletIDs = source["DropletIDs"].([]int)
	ret.LoadBalancerUIDs = source["LoadBalancerUIDs"].([]string)
	ret.Tags = source["Tags"].([]string)
	return dropletNames, ret
}

func getDestinations(destination map[interface{}]interface{}) ([]string, *godo.Destinations) {
	ret := &godo.Destinations{}
	dropletNames, source := getSource(destination)
	ret.Addresses = source["Addresses"].([]string)
	ret.DropletIDs = source["DropletIDs"].([]int)
	ret.LoadBalancerUIDs = source["LoadBalancerUIDs"].([]string)
	ret.Tags = source["Tags"].([]string)
	return dropletNames, ret
}

func getSource(source map[interface{}]interface{}) ([]string, map[string]interface{}) {
	retMap := make(map[string]interface{}, 0)
	retMap["Addresses"] = make([]string, 0)
	retMap["Tags"] = make([]string, 0)
	retMap["DropletIDs"] = make([]int, 0)
	retMap["LoadBalancerUIDs"] = make([]string, 0)
	dropletNames := make([]string, 0)
	for k, v := range source {
		switch k {
		case "Addresses":
			for _, value := range v.([]interface{}) {
				retMap["Addresses"] = append(retMap["Addresses"].([]string), value.(string))
			}
		case "Tags":
			for _, value := range v.([]interface{}) {
				retMap["Tags"] = append(retMap["Tags"].([]string), value.(string))
			}
		case "DropletIDs":
			for _, value := range v.([]interface{}) {
				switch o := value.(type) {
				case string:
					dropletNames = append(dropletNames, o)
				case int:
					retMap["DropletIDs"] = append(retMap["DropletIDs"].([]int), o)
				}
			}
		case "LoadBalancerUIDs":
			for _, value := range v.([]interface{}) {
				retMap["LoadBalancerUIDs"] = append(retMap["LoadBalancerUIDs"].([]string), value.(string))
			}
		}
	}
	return dropletNames, retMap
}
