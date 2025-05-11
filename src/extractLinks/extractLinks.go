package main

import (
	"github.com/voodooEntity/archivist"
	"github.com/voodooEntity/gits/src/transport"
	"github.com/voodooEntity/go-cyberbrain-plugin-interface/src/interfaces"
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

	contentHtml := input.Children()[0].Children()[0].Value
	domain := input.Value

	links, err := extractLinksFromHTML(contentHtml)
	absoluteLinks := ensureAbsoluteLinks(domain, links)
	linkEntities := []transport.TransportRelation{}
	if nil == err {
		for _, val := range absoluteLinks {
			if "" != val {
				linkEntities = append(linkEntities, transport.TransportRelation{
					Target: transport.TransportEntity{
						ID:             -2,
						Type:           "Link",
						Context:        "extractLinks",
						Value:          val,
						Properties:     make(map[string]string),
						ChildRelations: []transport.TransportRelation{},
					},
				})
			}
		}
	}

	ret := transport.TransportEntity{
		ID:             0,
		Type:           "Content",
		Value:          contentHtml,
		Context:        "loadUrl",
		Properties:     make(map[string]string),
		ChildRelations: linkEntities,
	}

	return []transport.TransportEntity{ret}, nil
}

func extractLinksFromHTML(html string) ([]string, error) {
	linkPattern := `<a\s+(?:[^>]*?\s+)?href=["']([^"']+)["']`

	regex := regexp.MustCompile(linkPattern)
	matches := regex.FindAllStringSubmatch(html, -1)

	var links []string
	for _, match := range matches {
		links = append(links, match[1])
	}

	return links, nil
}

func ensureAbsoluteLinks(baseURL string, links []string) []string {
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
								Value:      "Domain",
								Context:    "System",
								Properties: map[string]string{"Mode": "Set", "Type": "Secondary"},
								ChildRelations: []transport.TransportRelation{
									{
										Target: transport.TransportEntity{
											ID:         -1,
											Type:       "Structure",
											Value:      "Page",
											Context:    "System",
											Properties: map[string]string{"Mode": "Set", "Type": "Secondary"},
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
								},
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
