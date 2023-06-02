package globals

import (
	"context"
	"os"
	"sync"

	cam "google.golang.org/api/cloudasset/v1"
	crm "google.golang.org/api/cloudresourcemanager/v1"
	// "github.com/briandowns/spinner"
)

var (
	GetContext                   func() context.Context
	CrmService                   crm.Service
	CamService                   *cam.Service
	MapOfAllTFResourcesWithData  *sync.Map
	ListOfAllTFResourcesWithData []*Asset
	ListOfAllTFResourcesPerPath  *sync.Map
	Glossary                     map[string]string
	// CmdProgress *spinner.Spinner
	Reversor any

	Control *sync.Map
	OutFile *os.File

	GlobalCMDRunner func(pwd string, pipe []any, cmdp ...string) any
)

type (
	Ancestors []string
	Asset     struct {
		Ancestors       Ancestors      `json:"folders"`
		AssetType       string         `json:"assetType"`
		Name            string         `json:"name"`
		DisplayName     string         `json:"displayName"`
		ParentAssetType string         `json:"parentAssetType"`
		Project         string         `json:"project"`
		Organization    string         `json:"organization"`
		Parent          string         `json:"parentFullResourceName"`
		Data            map[string]any `json:"data"`
		// updateTime string `json:update_time`
		IsSupported bool `json:"supported"`
	}
)

// var OrgPrimaryResourceTypes = []string{"Project", "Folder", "TagKey", "TagBinding", "TagKey", "TagValue"}
var parentPathMap = map[string][2]any{}
var UserFriendlyNames *sync.Map

var OrgRootPath string
