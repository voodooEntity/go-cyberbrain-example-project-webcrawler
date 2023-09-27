package main

import (
	"github.com/voodooEntity/archivist"
	"github.com/voodooEntity/gits/src/transport"
	"github.com/voodooEntity/go-cyberbrain-plugin-interface/src/interfaces"
	"net/url"
	"regexp"
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

	links, err := extractLinksFromHTML2(input.Value)
	absoluteLinks := ensureAbsoluteLinks2(input.Properties["domain"], links)
	archivist.Info("absolute link", absoluteLinks)
	ret := []transport.TransportRelation{}
	if nil == err {
		for _, val := range absoluteLinks {
			if "" != val {
				ret = append(ret, transport.TransportRelation{
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

	input.ChildRelations = ret
	input.Properties = map[string]string{"domain": input.Properties["domain"]}
	//archivist.InfoF("Return data of extractLinks %+v", input)
	return []transport.TransportEntity{input}, nil
}

func extractLinksFromHTML2(html string) ([]string, error) {
	linkPattern := `<a\s+(?:[^>]*?\s+)?href=["']([^"']+)["']`

	regex := regexp.MustCompile(linkPattern)
	matches := regex.FindAllStringSubmatch(html, -1)

	var links []string
	for _, match := range matches {
		links = append(links, match[1])
	}

	return links, nil
}

func ensureAbsoluteLinks2(baseURL string, links []string) []string {
	baseURL = "https://" + baseURL
	var absoluteLinks []string

	for _, link := range links {
		// Check if the link starts with "http://", "https://", or "//"
		if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") || strings.HasPrefix(link, "//") {
			absoluteLinks = append(absoluteLinks, link)
		} else if strings.HasPrefix(link, "mailto:") {
			// Skip mailto links
			continue
		} else {
			// Construct an absolute URL by combining the base URL and the relative link
			absoluteURL := baseURL + "/" + link
			absoluteLinks = append(absoluteLinks, absoluteURL)
		}
	}

	return absoluteLinks
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
		cLink := link
		// Parse the link
		linkURL, err := url.Parse(cLink)
		if err != nil || strings.HasPrefix(link, "mailto:") {
			// Handle the error appropriately
			//archivist.ErrorF("Error parsing link URL: %v\n", err)
			continue // Skip this link if there was an error parsing it
		}
		//archivist.Info("LinkURL", linkURL)
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
