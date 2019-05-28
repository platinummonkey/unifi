package main

import (
	"log"
	"os"
	"strings"

	"github.com/platinummonkey/unifi"
)

func hr() {
	log.Printf("---------------------------------------------\n")
}

func main() {
	baseURL := os.Getenv("UNIFI_BASEURL") // like `https://192.168.123.456:7890`
	disableCertCheck := strings.HasPrefix(strings.ToLower(os.Getenv("UNIFI_DISABLE_CERT_CHECK")), "t")
	username := os.Getenv("UNIFI_USER")
	pass := os.Getenv("UNIFI_PASSWORD")

	c, _ := unifi.NewClient(baseURL, disableCertCheck)
	err := c.Login(username, pass, false) // set to true for long running sessions, for this example it's short.
	if err != nil {
		log.Printf("login error: %v\n", err)
	} else {
		log.Printf("login successful\n")
	}
	hr()

	status, err := c.ControllerStatus()
	if err != nil {
		log.Printf("status error: %v", err)
	} else {
		log.Printf("status: %v\n", status)
	}
	hr()

	self, err := c.Self()
	if err != nil {
		log.Printf("self error: %v\n", err)
	} else {
		log.Printf("self: %v\n", self)
	}
	hr()

	siteAdmins, err := c.SiteAdmins()
	if err != nil {
		log.Printf("site-admins error: %v\n", err)
	} else {
		log.Printf("site-admins: %v\n", siteAdmins)
	}
	hr()

	sites, err := c.AvailableSites()
	if err != nil {
		log.Printf("available sites error: %v\n", err)
	} else {
		log.Printf("available sites: %v\n", sites)
	}
	hr()

	sitesVerbose, err := c.AvailableSitesVerbose()
	if err != nil {
		log.Printf("available sites verbose error: %v\n", err)
	} else {
		log.Printf("available sites verbose: %v\n", sitesVerbose)
	}
	hr()
	siteID := sites.Data[0].Name

	// get the site health
	siteHealth, err := c.SiteHealth(siteID)
	if err != nil {
		log.Printf("site-%s-health error: %v\n", siteID, err)
	} else {
		log.Printf("site-%s-health: %v\n", siteID, siteHealth)
	}
	hr()

	// get events for first site
	siteEvents, err := c.SiteEvents(siteID, 0, 0, 10, unifi.EventSortOrderTimeDescending)
	if err != nil {
		log.Printf("site-%s-events error: %v\n", siteID, err)
	} else {
		log.Printf("site-%s-events: %v\n", siteID, siteEvents)
	}
	hr()

	// get alarms
	siteAlarms, err := c.SiteAlarms(siteID, 0, 0, 10, unifi.EventSortOrderTimeDescending)
	if err != nil {
		log.Printf("site-%s-alarms error: %v\n", siteID, err)
	} else {
		log.Printf("site-%s-alarms: %v\n", siteID, siteAlarms)
	}
	hr()

	// active clients
	activeClients, err := c.SiteActiveClients(siteID)
	if err != nil {
		log.Printf("site-%s-active-clients error: %v\n", siteID, err)
	} else {
		log.Printf("site-%s-active-clients: %v\n", siteID, activeClients)
	}
	hr()

	// logout
	err = c.Logout()
	if err != nil {
		log.Printf("logout error: %v\n", err)
	} else {
		log.Printf("logout successful\n")
	}

}
