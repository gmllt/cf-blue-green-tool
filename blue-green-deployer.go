package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"io"
	"os"
)

type BlueGreenDeploy struct {
	Connection plugin.CliConnection
	Out        io.Writer
	ErrorFunc  ErrorHandler
}

func (p *BlueGreenDeploy) Setup(connection plugin.CliConnection) {
	p.Connection = connection
}

func (p *BlueGreenDeploy) pushNewApp(manifest Manifest) () {
	manifest.SuffixApp("new")
	var manifestFilePath = "manifest-new.yml"
	_ = manifest.GenerateFile(manifestFilePath)
	fmt.Println("cf push -f "+manifestFilePath)
	_ = os.Remove(manifestFilePath)
}

func (p *BlueGreenDeploy) MapNewApp(newManifest Manifest, manifest Manifest) () {
	for NewAppIndex, NewApp := range newManifest.Applications {
		var OldApp = manifest.Applications[NewAppIndex]
		for _, OldRoute := range manifest.GetRoutes(OldApp.Name) {
			fmt.Println("cf map-route " + NewApp.Name + " -n " + OldRoute.Host + " " + OldRoute.Domain)
		}
	}
}

func (p *BlueGreenDeploy) RemoveOldApp(oldManifest Manifest) () {
	for _, NewApp := range oldManifest.Applications {
		fmt.Println("cf delete -f " + NewApp.Name)
	}
}

func (p *BlueGreenDeploy) RemoveOldRoute(manifest Manifest) {
	for _, App := range manifest.Applications {
		for _, oldRoute := range manifest.GetRoutes(App.Name) {
			fmt.Println("cf unmap-route "+App.Name+" -n "+oldRoute.Host+" "+oldRoute.Domain)
			fmt.Println("cf delete-route -f -n "+oldRoute.Host+" "+oldRoute.Domain)
		}
	}
}

func (p *BlueGreenDeploy) UnMapNewRoute(newManifest Manifest) () {
	for _, NewApp := range newManifest.Applications {
		for _, OldRoute := range newManifest.GetRoutes(NewApp.Name) {
			fmt.Println("cf unmap-route " + NewApp.Name + " -n " + OldRoute.Host + " " + OldRoute.Domain)
			fmt.Println("cf delete-route -n " + OldRoute.Host + " " + OldRoute.Domain)
		}
	}
}

func (p *BlueGreenDeploy) RenameApp(newManifest Manifest, manifest Manifest) () {
	for AppIndex, App := range newManifest.Applications {
		var CurrentApp = manifest.Applications[AppIndex]
		fmt.Println("cf rename " + CurrentApp.Name + " " + App.Name)
	}
}

func (p *BlueGreenDeploy) MapOldApp(oldManifest Manifest) () {
	for _, OldApp := range oldManifest.Applications {
		for _, OldRoute := range oldManifest.GetRoutes(OldApp.Name) {
			fmt.Println("cf map-route " + OldApp.Name + " -n " + OldRoute.Host + " " + OldRoute.Domain)
		}
	}
}

func (p *BlueGreenDeploy) UnMapOldApp(oldManifest Manifest, manifest Manifest) () {
	for OldAppIndex, OldApp := range oldManifest.Applications {
		var App = manifest.Applications[OldAppIndex]
		for _, OldRoute := range manifest.GetRoutes(App.Name) {
			fmt.Println("cf unmap-route " + OldApp.Name + " -n " + OldRoute.Host + " " + OldRoute.Domain)
		}
	}
}
