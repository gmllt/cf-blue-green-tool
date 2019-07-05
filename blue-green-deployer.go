package main

import (
	"io"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/cli/plugin"
)

// BlueGreenDeploy is a simple structure to execute cf-cli functions
type BlueGreenDeploy struct {
	Connection plugin.CliConnection
	Out        io.Writer
	ErrorFunc  ErrorHandler
}

// Setup : Setup object
func (p *BlueGreenDeploy) Setup(connection plugin.CliConnection) {
	p.Connection = connection
}

// pushNewApp : push new app
func (p *BlueGreenDeploy) pushNewApp(manifest Manifest, manifestFile string) {
	var DirPath = filepath.Dir(manifestFile)
	var manifestFilePath = DirPath + "/manifest-new.yml"
	_ = manifest.GenerateFile(manifestFilePath)
	//fmt.Println("cf push -f " + manifestFilePath)
	if _, err := p.Connection.CliCommand("push", "-f", manifestFilePath); err != nil {
		p.ErrorFunc("Could not push new app.", err)
	}
	_ = os.Remove(manifestFilePath)
}

// MapNewApp : map live routes to -new app
func (p *BlueGreenDeploy) MapNewApp(newManifest Manifest, manifest Manifest) {
	for NewAppIndex, NewApp := range newManifest.Applications {
		var OldApp = manifest.Applications[NewAppIndex]
		for _, OldRoute := range manifest.GetRoutes(OldApp.Name) {
			p.MapRoute(NewApp.Name, OldRoute.Host, OldRoute.Domain)
		}
	}
}

// MapRoute : map route by host and domain to app
func (p *BlueGreenDeploy) MapRoute(AppName string, Host string, Domain string) {
	//fmt.Println("cf map-route " + AppName + " -n " + Host + " " + Domain)
	if _, err := p.Connection.CliCommand("map-route", AppName, "-n", Host, Domain); err != nil {
		p.MapDomain(AppName, Host+"."+Domain)
	}

}

// CheckApp : check if an app exist
func (p *BlueGreenDeploy) CheckApps(manifest Manifest) bool {
	for _, App := range manifest.Applications {
		if _, err := p.Connection.CliCommand("app", App.Name); err != nil {
			return false
		}
	}
	return true
}

// MapDomain : map domain to app
func (p *BlueGreenDeploy) MapDomain(AppName string, Domain string) {
	//fmt.Println("cf map-route " + AppName + " " + Domain)
	if _, err := p.Connection.CliCommand("map-route", AppName, Domain); err != nil {
		p.ErrorFunc("Could not map as route or domain.", err)
	}
}

// UnMapRoute : unmap route by host and domain from app
func (p *BlueGreenDeploy) UnMapRoute(AppName string, Host string, Domain string) {
	//fmt.Println("cf unmap-route " + AppName + " -n " + Host + " " + Domain)
	if _, err := p.Connection.CliCommand("unmap-route", AppName, "-n", Host, Domain); err != nil {
		p.UnMapDomain(AppName, Host+"."+Domain)
	}
}

// UnMapDomain : unmap domain from app
func (p *BlueGreenDeploy) UnMapDomain(AppName string, Domain string) {
	//fmt.Println("cf unmap-route " + AppName + " " + Domain)
	p.Connection.CliCommand("unmap-route", AppName, Domain)
}

// DeleteRoute : delete route by host and domain
func (p *BlueGreenDeploy) DeleteRoute(Host string, Domain string) {
	//fmt.Println("cf delete-route -f -n " + Host + " " + Domain)
	p.Connection.CliCommand("delete-route", "-f", "-n", Host, Domain)
}

// RemoveOldApp : delete old app
func (p *BlueGreenDeploy) RemoveOldApp(oldManifest Manifest) {
	for _, NewApp := range oldManifest.Applications {
		//fmt.Println("cf delete -f " + NewApp.Name)
		p.Connection.CliCommand("delete", "-f", NewApp.Name)
	}
}

// RemoveOldRoute : remove routes from app
func (p *BlueGreenDeploy) RemoveOldRoute(manifest Manifest) {
	for _, App := range manifest.Applications {
		for _, oldRoute := range manifest.GetRoutes(App.Name) {
			p.UnMapRoute(App.Name, oldRoute.Host, oldRoute.Domain)
			p.DeleteRoute(oldRoute.Host, oldRoute.Domain)
		}
	}
}

// UnMapNewRoute : unmap -new route from -new app
func (p *BlueGreenDeploy) UnMapNewRoute(newManifest Manifest) {
	for _, NewApp := range newManifest.Applications {
		for _, OldRoute := range newManifest.GetRoutes(NewApp.Name) {
			p.UnMapRoute(NewApp.Name, OldRoute.Host, OldRoute.Domain)
			p.DeleteRoute(OldRoute.Host, OldRoute.Domain)
		}
	}
}

// RenameApp : rename second to first
func (p *BlueGreenDeploy) RenameApp(newManifest Manifest, manifest Manifest) {
	for AppIndex, App := range newManifest.Applications {
		var CurrentApp = manifest.Applications[AppIndex]
		//fmt.Println("cf rename " + CurrentApp.Name + " " + App.Name)
		if _, err := p.Connection.CliCommand("rename", CurrentApp.Name, App.Name); err != nil {
			p.ErrorFunc("Could not map as route or domain.", err)
		}
	}
}

// MapOldApp : map -old routes to -old app
func (p *BlueGreenDeploy) MapOldApp(oldManifest Manifest) {
	for _, OldApp := range oldManifest.Applications {
		for _, OldRoute := range oldManifest.GetRoutes(OldApp.Name) {
			p.MapRoute(OldApp.Name, OldRoute.Host, OldRoute.Domain)
		}
	}
}

// UnMapOldApp : unmap -old routes from -old app
func (p *BlueGreenDeploy) UnMapOldApp(oldManifest Manifest, manifest Manifest) {
	for OldAppIndex, OldApp := range oldManifest.Applications {
		var App = manifest.Applications[OldAppIndex]
		for _, OldRoute := range manifest.GetRoutes(App.Name) {
			p.UnMapRoute(OldApp.Name, OldRoute.Host, OldRoute.Domain)
		}
	}
}
