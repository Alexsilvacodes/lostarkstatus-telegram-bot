package lostarkstatus

import (
	"sort"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Server struct {
	Name   string
	Status string
}

const BaseURL = "https://www.playlostark.com/en-us/support/server-status"

func GetStatus(f func(s []Server)) {
	var servers []Server

	c := colly.NewCollector(
		colly.AllowedDomains("www.playlostark.com"),
	)

	c.OnHTML(".ags-ServerStatus-content-responses-response-server", func(e *colly.HTMLElement) {
		name := e.ChildText(".ags-ServerStatus-content-responses-response-server-name")
		status_class := e.ChildAttr(".ags-ServerStatus-content-responses-response-server-status", "class")
		status := strings.ReplaceAll(
			status_class,
			"ags-ServerStatus-content-responses-response-server-status ags-ServerStatus-content-responses-response-server-status--",
			"",
		)
		server := Server{name, StatusCapitalizedWithIcon(status)}
		servers = append(servers, server)
	})

	c.OnError(func(_ *colly.Response, err error) {
		f([]Server{})
	})

	c.OnScraped(func(r *colly.Response) {
		sort.Slice(servers, func(i, j int) bool {
			return servers[i].Name < servers[j].Name
		})

		f(servers)
	})

	c.Visit(BaseURL)
}

func StatusCapitalizedWithIcon(status string) string {
	switch status {
	case "good":
		return "Good ✅"
	case "busy":
		return "Busy ❌"
	case "full":
		return "Full 💬"
	case "maintenance":
		return "Maintenance 🚧"
	}

	return ""
}
