package main

import (
	"errors"
	"github.com/voodooEntity/archivist"
	"github.com/voodooEntity/gits/src/transport"
	"github.com/voodooEntity/go-cyberbrain-plugin-interface/src/interfaces"
	"net/url"
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

	link := input.ChildRelations[0].Target.ChildRelations[0].Target.ChildRelations[0].Target.Value
	linkDomain, err := extractDomain(link)

	if nil != err {
		return []transport.TransportEntity{}, err
	}

	if linkDomain != input.Value {
		return []transport.TransportEntity{}, errors.New("skipping link due to wrong domain")
	}

	ret := transport.TransportEntity{
		Type:       "Domain",
		ID:         0,
		Value:      input.Value,
		Context:    "webcrawl",
		Properties: make(map[string]string),
		ChildRelations: []transport.TransportRelation{
			{
				Target: transport.TransportEntity{
					ID:         0,
					Type:       "Page",
					Value:      link,
					Context:    "addPage",
					Properties: make(map[string]string),
				},
			},
		},
	}

	return []transport.TransportEntity{ret}, nil
}

func extractDomain(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

func (self Plugin) GetConfig() transport.TransportEntity {
	return transport.TransportEntity{
		ID:         -1,
		Type:       "Action",
		Value:      "addPage",
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
														Properties: map[string]string{"Mode": "Set", "Type": "Secondary"},
														ChildRelations: []transport.TransportRelation{
															{
																Target: transport.TransportEntity{
																	ID:         -1,
																	Type:       "Structure",
																	Value:      "Link",
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
