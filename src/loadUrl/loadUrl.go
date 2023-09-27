package main

import (
	"errors"
	"github.com/voodooEntity/archivist"
	"github.com/voodooEntity/gits/src/transport"
	"github.com/voodooEntity/go-cyberbrain-plugin-interface/src/interfaces"
	"io/ioutil"
	"net/http"
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
	var currUrl string
	var checkDomain string
	if "alpha" == requirement {
		currUrl = input.Value
		checkDomain = input.Properties["domain"]
	} else {
		currUrl = input.Value
	}

	currUrlDomain, err := extractDomain(currUrl)
	if nil != err {
		return []transport.TransportEntity{}, errors.New("cant extract domain from page url")
	}
	if "alpha" == requirement && currUrlDomain != checkDomain {
		return []transport.TransportEntity{}, errors.New("skipping url due to wrong domain")
	}

	pageContent, err := loadPage(currUrl)
	if nil != err {
		return []transport.TransportEntity{}, err
	}

	ret := []transport.TransportEntity{
		{
			ID:         0,
			Type:       "Domain",
			Context:    "loadUrl",
			Value:      currUrlDomain,
			Properties: make(map[string]string),
			ChildRelations: []transport.TransportRelation{
				{
					Target: transport.TransportEntity{
						ID:         0,
						Type:       "Page",
						Context:    "loadUrl",
						Value:      currUrl,
						Properties: make(map[string]string),
						ChildRelations: []transport.TransportRelation{
							{
								Target: transport.TransportEntity{
									ID:         0,
									Type:       "Content",
									Context:    "loadUrl",
									Value:      pageContent,
									Properties: map[string]string{"domain": currUrlDomain},
								},
							},
						},
					},
				},
			},
		},
	}
	return ret, nil
}

func loadPage(currUrl string) (string, error) {
	archivist.Debug("loading url", currUrl)
	req, err := http.NewRequest(http.MethodGet, currUrl, nil)
	if err != nil {
		archivist.DebugF("client: could not create request: %s\n", err)
		return "", errors.New("client: could not create request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		archivist.DebugF("client: error making http request: %s\n", err)
		return "", errors.New("client: error making http request")
	}

	archivist.DebugF("client: got response!\n")
	archivist.DebugF("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		archivist.DebugF("client: could not read response body: %s\n", err)
		return "", errors.New("client: could not read response body")
	}

	archivist.DebugF("client: response body: %s\n", resBody)

	if !isTextOrHTMLContentType(res.Header.Get("Content-Type")) {
		archivist.DebugF("client: response content is not plain text or HTML")
		return "", errors.New("client: response content is not plain text or HTML")
	}
	return string(resBody[:]), nil
}

func isTextOrHTMLContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "text/html") || strings.HasPrefix(contentType, "application/xhtml+xml")
}

func extractDomain(inputURL string) (string, error) {
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "https://" + inputURL
	}

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	parts := strings.Split(parsedURL.Hostname(), ".")
	if len(parts) >= 2 {
		// If there are at least two parts, return the last two joined by a dot
		return parts[len(parts)-2] + "." + parts[len(parts)-1], nil
	}

	return "", errors.New("Invalid URL format")
}

func (self Plugin) GetConfig() transport.TransportEntity {
	return transport.TransportEntity{
		ID:         -1,
		Type:       "Action",
		Value:      "loadUrl",
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
								ID:             -1,
								Type:           "Structure",
								Value:          "Link",
								Context:        "System",
								Properties:     map[string]string{"Mode": "Set", "Type": "Primary"},
								ChildRelations: []transport.TransportRelation{},
							},
						},
					},
				},
			},
			{
				Target: transport.TransportEntity{
					ID:         -1,
					Type:       "Dependency",
					Value:      "beta",
					Context:    "System",
					Properties: make(map[string]string),
					ChildRelations: []transport.TransportRelation{
						{
							Target: transport.TransportEntity{
								ID:             -1,
								Type:           "Structure",
								Value:          "InitUrl",
								Context:        "System",
								Properties:     map[string]string{"Mode": "Set", "Type": "Primary"},
								ChildRelations: []transport.TransportRelation{},
							},
						},
					},
				},
			},
			{
				Target: transport.TransportEntity{
					Type:       "Category",
					Value:      "webcrawler",
					Properties: make(map[string]string),
					Context:    "System",
				},
			},
		},
	}
}
