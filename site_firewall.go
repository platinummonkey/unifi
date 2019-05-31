package unifi

import (
	"net/http"
)

type SiteFirewallRule map[string]interface{}

type SiteFirewallRuleResponse struct {
	Meta CommonMeta         `json:"meta"`
	Data []SiteFirewallRule `json:"data"`
}

func (c *Client) SiteFirewallRules(site string) (*SiteFirewallRuleResponse, error) {
	var resp SiteFirewallRuleResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/firewallrule", nil, &resp)
	return &resp, err
}

type SiteFirewallGroup map[string]interface{}

type SiteFirewallGroupResponse struct {
	Meta CommonMeta          `json:"meta"`
	Data []SiteFirewallGroup `json:"data"`
}

func (c *Client) SiteFirewallGroups(site string) (*SiteFirewallGroupResponse, error) {
	var resp SiteFirewallGroupResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/firewallgroup", nil, &resp)
	return &resp, err
}
