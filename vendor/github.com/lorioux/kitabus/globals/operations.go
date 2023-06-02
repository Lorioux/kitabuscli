package globals

import (
	"context"
	"log"
	"strings"

	cam "google.golang.org/api/cloudasset/v1"
)

func (a *Asset) PickResourceName() string {
	name := strings.Split(a.Name, "/")
	return name[len(name)-1]
}

func (a *Asset) PickResourceId() string {
	id := strings.SplitN(a.Name, "/", 4)
	return id[len(id)-1]
}

func (a *Asset) PickParentId() string {

	switch a.PickResourceType() {
	case "Organization":
		return a.Organization
	// case "Folder" :
	// 	strings.Split(a.Parent, "/")
	// case "Project":
	// 	return a.Project
	default:
		id := strings.Split(a.Parent, "/")
		return id[len(id)-1]
	}
}

func (a *Asset) PickResourceType() string {
	rType := strings.Split(a.AssetType, "/")
	return rType[len(rType)-1] //strings.ToTitle(aType)
}

func (a *Asset) PickParentType() string {
	pType := strings.Split(a.ParentAssetType, "/")
	return pType[len(pType)-1] //strings.ToTitle(aType)
}

func (a *Asset) SetData(data any) {
	a.Data = data.(map[string]any)
}

func (a *Asset) PickServiceType() string {
	sType := strings.Split(a.AssetType, ".")[0]

	return strings.ToLower(sType)
}

func (a *Asset) SetSupported(s bool) {
	switch {
	case s:
		a.IsSupported = s
	default:
		log.Printf("UNSUPPORTED RSTYPE: %v", a.AssetType)
	}
}

func InitService() (*cam.Service, error) {
	return cam.NewService(GetContext())
}

func InitContext() context.Context {
	return context.Background()
}
