package unifi

import (
	"net/http"
)

type SiteCurrentChannels struct {
	ChannelsNA       []int       `json:"channels_na"`
	ChannelsNA160    []int       `json:"channels_na_160"`
	ChannelsNA40     []int       `json:"channels_na_40"`
	ChannelsNA40BCM  []int       `json:"channels_na_40_bcm"`
	ChannelsNA80     []int       `json:"channels_na_80"`
	ChannelsNA80BCM  []int       `json:"channels_na_80_bcm"`
	ChannelsNABCM    []int       `json:"channels_na_bcm"`
	ChannelsNADFS    []int       `json:"channels_na_dfs"`
	ChannelsNAIndoor []int       `json:"channels_na_indoor"`
	ChannelsNG       []int       `json:"channels_ng"`
	ChannelsNG40     []int       `json:"channels_ng_40"`
	ChannelsNG40BCM  []int       `json:"channels_ng_40_bcm"`
	ChannelsNGBCM    []int       `json:"channels_ng_bcm"`
	Code             interface{} `json:"code"` // sometimes string or int
	Key              string      `json:"key"`
	Name             string      `json:"name"`
}

type SiteCurrentChannelsResponse struct {
	Meta CommonMeta            `json:"meta"`
	Data []SiteCurrentChannels `json:"data"`
}

func (c *Client) SiteCurrentChannels(site string) (*SiteCurrentChannelsResponse, error) {
	var resp SiteCurrentChannelsResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/current-channel", nil, &resp)
	return &resp, err
}
