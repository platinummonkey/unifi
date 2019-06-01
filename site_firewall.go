package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// SiteFirewallRule defines a site firewall rule
type SiteFirewallRule map[string]interface{}

// SiteFirewallRuleResponse contains the firewall rules response
type SiteFirewallRuleResponse struct {
	Meta CommonMeta         `json:"meta"`
	Data []SiteFirewallRule `json:"data"`
}

// SiteFirewallRules queries the site firewall rules
// site - the site to query
func (c *Client) SiteFirewallRules(site string) (*SiteFirewallRuleResponse, error) {
	var resp SiteFirewallRuleResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/firewallrule", nil, &resp)
	return &resp, err
}

// SiteFirewallGroup defines the site firewall group
type SiteFirewallGroup map[string]interface{}

// SiteFirewallGroupResponse contains the firewall groups response
type SiteFirewallGroupResponse struct {
	Meta CommonMeta          `json:"meta"`
	Data []SiteFirewallGroup `json:"data"`
}

// SiteFirewallGroups will list firewall groups
// site - the site to query
// groupID - filter on the associated group, if zero-value it returns all for the entire site.
func (c *Client) SiteFirewallGroups(site string, groupID string) (*SiteFirewallGroupResponse, error) {
	extPath := "rest/firewallgroup"
	if groupID != "" {
		extPath = extPath + "/" + strings.TrimSpace(groupID)
	}

	var resp SiteFirewallGroupResponse
	err := c.doSiteRequest(http.MethodGet, site, extPath, nil, &resp)
	return &resp, err
}

// FirewallGroupType defines firewall group types
type FirewallGroupType string

// The only supported firewall group types
const (
	FirewallGroupTypeAddressGroup     FirewallGroupType = "address-group"
	FirewallGroupTypeIPV6AddressGroup FirewallGroupType = "ipv6-address-group"
	FirewallGroupTypePortGroup        FirewallGroupType = "port-group"
)

// IsValid returns true if it's a valid firewall group type.
// there are only a few valid types
func (g FirewallGroupType) IsValid() bool {
	switch g {
	case FirewallGroupTypeAddressGroup, FirewallGroupTypeIPV6AddressGroup, FirewallGroupTypePortGroup:
		return true
	default:
		return false
	}
}

// FirewallGroupMembers defines the firewall group members configuration
type FirewallGroupMembers struct {
	IPV4Addresses []string
	IPV6Addresses []string
	Ports         []int
}

// CreateFirewallGroup will create a new firewall group
// site - the site to modify
// name - the name of the firewall group
// groupType - the type of firewall group
// groupMembers - the firewall group member configuration
func (c *Client) CreateFirewallGroup(site string, name string, groupType FirewallGroupType, groupMembers FirewallGroupMembers) (*GenericResponse, error) {
	if !groupType.IsValid() {
		return nil, fmt.Errorf("invalid groupType specified: %s", groupType)
	}

	members := make([]interface{}, 0)
	switch groupType {
	case FirewallGroupTypeAddressGroup:
		for _, a := range groupMembers.IPV4Addresses {
			members = append(members, a)
		}
	case FirewallGroupTypeIPV6AddressGroup:
		for _, a := range groupMembers.IPV6Addresses {
			members = append(members, a)
		}
	case FirewallGroupTypePortGroup:
		for _, a := range groupMembers.Ports {
			members = append(members, a)
		}
	}

	payload := map[string]interface{}{
		"name":          name,
		"group_type":    string(groupType),
		"group_members": members,
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "rest/firewallgroup", bytes.NewReader(data), &resp)
	return &resp, err
}

// UpdateFirewallGroup will update an existing firewall group
// site - the site to modify
// siteID - the ID of the site
// groupID - the ID of the firewall group
// name - the name of the firewall group
// groupType - the type of firewall group, note you cannot change a group type
// groupMembers - the firewall group member configuration
func (c *Client) UpdateFirewallGroup(site string, siteID string, groupID string, name string, groupType FirewallGroupType, groupMembers FirewallGroupMembers) (*GenericResponse, error) {
	if !groupType.IsValid() {
		return nil, fmt.Errorf("invalid groupType specified: %s", groupType)
	}

	members := make([]interface{}, 0)
	switch groupType {
	case FirewallGroupTypeAddressGroup:
		for _, a := range groupMembers.IPV4Addresses {
			members = append(members, a)
		}
	case FirewallGroupTypeIPV6AddressGroup:
		for _, a := range groupMembers.IPV6Addresses {
			members = append(members, a)
		}
	case FirewallGroupTypePortGroup:
		for _, a := range groupMembers.Ports {
			members = append(members, a)
		}
	}

	payload := map[string]interface{}{
		"_id":           groupID,
		"name":          name,
		"group_type":    string(groupType),
		"group_members": members,
		"site_id":       siteID,
	}

	data, _ := json.Marshal(payload)

	extPath := fmt.Sprintf("rest/firewallgroup/%s", strings.TrimSpace(groupID))

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, extPath, bytes.NewReader(data), &resp)
	return &resp, err
}

// DeleteFirewallGroup will delete an existing firewall group
// site - the site to modify
// siteID - the ID of the site
// groupID - the ID of the firewall group
func (c *Client) DeleteFirewallGroup(site string, groupID string) (*GenericResponse, error) {
	extPath := fmt.Sprintf("rest/firewallgroup/%s", strings.TrimSpace(groupID))

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodDelete, site, extPath, nil, &resp)
	return &resp, err
}
