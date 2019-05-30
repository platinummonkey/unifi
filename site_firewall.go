package unifi

import (
	"net/http"
)

type SiteFirewallRule map[string]interface{}

type SiteFirewallRuleResponse struct {
	Meta CommonMeta         `json:"meta"`
	Data []SiteFirewallRule `json:"data"`
}

func (c *Client) SiteFirewallRules(siteID string) (*SiteFirewallRuleResponse, error) {
	var resp SiteFirewallRuleResponse
	err := c.doSiteRequest(http.MethodGet, siteID, "rest/firewallrule", nil, &resp)
	return &resp, err
}

type SiteFirewallGroup map[string]interface{}

type SiteFirewallGroupResponse struct {
	Meta CommonMeta          `json:"meta"`
	Data []SiteFirewallGroup `json:"data"`
}

func (c *Client) SiteFirewallGroups(siteID string) (*SiteFirewallGroupResponse, error) {
	var resp SiteFirewallGroupResponse
	err := c.doSiteRequest(http.MethodGet, siteID, "rest/firewallgroup", nil, &resp)
	return &resp, err
}
