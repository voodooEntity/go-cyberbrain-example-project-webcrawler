package main

import (
	"github.com/voodooEntity/archivist"
	"github.com/voodooEntity/gits/src/transport"
	"github.com/voodooEntity/go-cyberbrain-plugin-interface/src/interfaces"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

// define a type to can bind our methods on to export
type Plugin struct{}

// export the interfaces struct so it can be accessed via plugin Lookup
var Export Plugin

func (self Plugin) New() interfaces.PluginInterface {
	return Plugin{}
}

// Execute method mandatory
func (self Plugin) Execute(input transport.TransportEntity, requirement string, context string) ([]transport.TransportEntity, error) {
	archivist.DebugF("Plugin executed with input %+v", input)

	links, err := extractLinksFromHTML(input.Value)
	absoluteLinks := ensureAbsoluteLinks(input.Properties["domain"], links)
	archivist.Debug("absolute link", absoluteLinks)
	if nil == err {
		for _, val := range absoluteLinks {
			if "" != val {
				input.ChildRelations = append(input.ChildRelations, transport.TransportRelation{
					Target: transport.TransportEntity{
						ID:             0,
						Type:           "Link",
						Context:        "loadUrl",
						Value:          val,
						Properties:     map[string]string{"domain": input.Properties["domain"]},
						ChildRelations: []transport.TransportRelation{},
					},
				})
			}
		}
	}

	input.Properties = make(map[string]string)
	return []transport.TransportEntity{input}, nil
}

func extractLinksFromHTML(inputHTML string) ([]string, error) {
	links := []string{}

	// Parse the inputHTML string
	doc, err := html.Parse(strings.NewReader(inputHTML))
	if err != nil {
		return links, err
	}

	// Define a recursive function to traverse the HTML nodes
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}

	// Start the traversal
	findLinks(doc)

	return links, nil
}

func ensureAbsoluteLinks(baseURL string, links []string) []string {
	baseURL = "https://" + baseURL
	var absoluteLinks []string

	// Parse the base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		// Handle the error appropriately
		archivist.ErrorF("Error parsing base URL: %v\n", err)
		return absoluteLinks
	}

	for _, link := range links {
		// Parse the link
		linkURL, err := url.Parse(link)
		if err != nil || strings.HasPrefix(link, "mailto:") {
			// Handle the error appropriately
			archivist.ErrorF("Error parsing link URL: %v\n", err)
			continue // Skip this link if there was an error parsing it
		}
		archivist.Info("LinkURL", linkURL)
		// Check if the link is relative
		if linkURL.Scheme == "" && linkURL.Host == "" {
			// Construct an absolute URL by combining the base URL and the relative link
			absoluteURL := base.ResolveReference(linkURL)
			absoluteLinks = append(absoluteLinks, absoluteURL.String())
		} else {
			// The link is already absolute
			absoluteLinks = append(absoluteLinks, link)
		}
	}

	return absoluteLinks
}

func (self Plugin) GetConfig() transport.TransportEntity {
	return transport.TransportEntity{
		ID:         -1,
		Type:       "Action",
		Value:      "extractLinks",
		Context:    "System",
		Properties: make(map[string]string),
		ChildRelations: []transport.TransportRelation{
			{
				Target: transport.TransportEntity{
					ID:         -1,
					Type:       "Dependency",
					Value:      "alpha",
					Context:    "System",
					Properties: make(map[string]string),
					ChildRelations: []transport.TransportRelation{
						{
							Target: transport.TransportEntity{
								ID:         -1,
								Type:       "Structure",
								Value:      "Content",
								Context:    "System",
								Properties: map[string]string{"Mode": "Set", "Type": "Primary"},
							},
						},
					},
				},
			},
			{
				Target: transport.TransportEntity{
					ID:         0,
					Type:       "Category",
					Value:      "webcrawler",
					Properties: make(map[string]string),
					Context:    "System",
				},
			},
		},
	}
}
