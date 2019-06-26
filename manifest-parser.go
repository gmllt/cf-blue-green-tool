package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"strings"
)

type Manifest struct {
	Applications []struct {
		Name      string `yaml:"name,omitempty"`
		Memory    string `yaml:"memory,omitempty"`
		DiskQuota string `yaml:"disk_quota,omitempty"`
		Instances int    `yaml:"instances,omitempty"`
		Command   string `yaml:"command,omitempty"`
		Docker    struct {
			Image    string `yaml:"image,omitempty"`
			Username string `yaml:"username,omitempty"`
		} `yaml:"docker,omitempty"`
		HealthCheckHttpEndpoint string   `yaml:"health-check-http-endpoint,omitempty"`
		HealthCheckType         string   `yaml:"health-check-type,omitempty"`
		NoRoute                 bool     `yaml:"no-route,omitempty"`
		RandomRoute             bool     `yaml:"random-route,omitempty"`
		Stack                   string   `yaml:"stack,omitempty"`
		Timeout                 int      `yaml:"timeout,omitempty"`
		Domain                  string   `yaml:"domain,omitempty"`
		Domains                 []string `yaml:"domains,omitempty"`
		Host                    string   `yaml:"host,omitempty"`
		Hosts                   []string `yaml:"hosts,omitempty"`
		NoHostname              bool     `yaml:"no-hostname,omitempty"`
		Routes                  []struct {
			Route string `yaml:"route"`
		} `yaml:"routes,omitempty"`
		Services   []string               `yaml:"services,omitempty"`
		Path       string                 `yaml:"path,omitempty"`
		BuildPack  string                 `yaml:"buildpack,omitempty"`
		BuildPacks []string               `yaml:"buildpacks,omitempty"`
		Env        map[string]interface{} `yaml:"env,omitempty"`
	} `yaml:"applications,omitempty"`
}

type Route struct {
	Host   string
	Domain string
}

type ErrorHandler func(string, error)

type ManifestParse struct {
	Connection plugin.CliConnection
	Out        io.Writer
	ErrorFunc  ErrorHandler
}

func (p *ManifestParse) parseByte(Byte []byte) (Manifest, error) {
	var manifest Manifest
	var err = yaml.Unmarshal(Byte, &manifest)
	if err != nil {
		return manifest, err
	}
	return manifest, err
}

func (p *ManifestParse) parseFile(manifestFile string) (Manifest, error) {
	var manifest Manifest
	yamlFile, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return manifest, err
	}
	return p.parseByte(yamlFile)
}

func (m *Manifest) GetRoutes(AppName string) []Route {
	var NewRoutes []Route
	for _, App := range m.Applications {
		if App.Name == AppName {
			for _, CurrentRoute := range App.Routes {
				var Parts = strings.Split(CurrentRoute.Route, ".")
				var NewRoute Route
				NewRoute.Host = Parts[0]
				NewRoute.Domain = strings.Join(Parts[:0+copy(Parts[0:], Parts[1:])], ".")
				NewRoutes = append(NewRoutes, NewRoute)
			}
		}
	}
	return NewRoutes
}

func (m *Manifest) GenerateFile(FileName string) (error) {
	Byte, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(FileName, Byte, 0644)
	return err
}

func (m *Manifest) SuffixApp(suffix string) () {
	for AppIndex, App := range m.Applications {
		for RouteIndex, CurrentRoute := range m.GetRoutes(App.Name) {
			m.Applications[AppIndex].Routes[RouteIndex].Route = CurrentRoute.Host + "-" + suffix + "." + CurrentRoute.Domain
		}
		m.Applications[AppIndex].Name = App.Name + "-" + suffix
	}
}
