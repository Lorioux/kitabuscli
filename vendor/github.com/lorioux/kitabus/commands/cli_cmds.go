package commands

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	"google.golang.org/api/cloudasset/v1"
	"google.golang.org/api/googleapi"

	"github.com/lorioux/kitabus/globals"
	"github.com/lorioux/kitabus/playbooks"
)

var (
	scope             string
	atypes            string
	withRelationships bool
	relations         []string

	GenerateTfTemplates bool
	generate_cmd        string
	// search_cmd bool
	import_cmd string

	CLICMD *exec.Cmd
)

// var GCLOUD_ASSET_CMD []string

// const GCLOUD_ASSET_SEARCH_CMD = "search-all-resources"

// // const GCLOUD_ASSET_FILTER = `"state = ACTIVE"`
// const GCLOUD_ASSET_ORDER_BY = `parentFullResourceName DESC,assetType DESC`
// const GCLOUD_ASSET_OUT_FMT = `json(name,displayName,assetType,parentAssetType,parentFullResourceName,project,folders,organization)`

// var GCLOUD_ASSET_OUT_PATH = path.Join(os.Getenv("PWD"), ".catalogue")
// var GCLOUD_ASSET_SEARCH_SCOPE = ""
// var OutFile *os.File

const (
	generate_cmd_usage = `
Usage: kitabus --generate [option] [--scope|kind|with-relations|relationships|path]
[options]: 
	[catalogue]	--path "."	"Search resources in the --scope and --kind list. The catalogue will be Generate in the --path (default '.')" 
	[orgtree]	--path "."	"Generate organization directory tree relative to --path (default '.')"
	[modules]	--path "."	"Generate terraform module templates and append to the existing based on the catalogue and relative to --path (default '.')"
	[modules[override]] --path "."	"Generate terraform module templates , overriding the existing, based on the catalogue and relative to --path (default '.')"
`
	import_cmd_usage = `
Usage: kitabus --import [options] [--path]
[options]: 
	[resource] --path "."	"Search resources in the --scope and --kind list. The catalogue will be Generate in the --path (default '.')" 
`
)

func init() {
	flag.StringVar(
		&generate_cmd,
		"generate",
		"",
		``+generate_cmd_usage)

	flag.StringVar(
		&import_cmd,
		"import",
		"",
		``+import_cmd_usage)

	flag.StringVar(&scope, "scope", os.Getenv("GOOGLE_CLOUD_PROJECT"),
		`"The scope of the assets to list. Example: projects/{ID|NUMBER}, folders/{NUMBER}, organizations/{NUMBER}"`)

	flag.StringVar(&atypes, "kind", ".*", "A comma-separated list of asset types to list, or a empty to match all asset types.")

	flag.StringVar(&globals.OrgRootPath, "path", os.Getenv("PWD"), "Example: kitabus --Generate templates --[path].\nA command to crete resources terraform templates.")

}

func CmdExecute() {
	flag.Parse()

	if strings.EqualFold(generate_cmd, "") && strings.EqualFold(import_cmd, "") {
		log.Print(generate_cmd_usage, import_cmd_usage)
	} else if !strings.EqualFold(generate_cmd, "") {
		switch generate_cmd {
		case "catalogue":
			CmdGenerateCatalogus()
		case "orgtree":
			CmdGenerateOrgTreeIndex()
		case "modules":
			CmdGenerateResourceTFTemplates()
		case "modules[override]":
			playbooks.TFMetaModuleInstance.OverrideMetaModulesOrRemove("", generate_cmd)
			CmdGenerateResourceTFTemplates()
		case "auto":
			CmdExecuteForDebugger()
		default:
			log.Print(generate_cmd_usage)
		}
	} else if !strings.EqualFold(import_cmd, "") {
		switch import_cmd {
		case "resources":
			CmdImportResourceTFState()
		default:
			log.Print(import_cmd_usage)
		}
	}
}

func CmdExecuteForDebugger() {
	// playbooks.TFMetaModuleInstance.OverrideMetaModulesOrRemove(orgroot,"modules[override]")
	// CmdDiscoverAllResourcesWithData()
	CmdGenerateCatalogus()
	CmdGenerateOrgTreeIndex()
	CmdGenerateResourceTFTemplates()
	// CmdImportResourceTFState()
}

// func ExecuteByRelationShips(rel ...string) {
// 	// relationships := flag.String("relationships", "", "A comma-separated list of asset types to list. Set only when content-type \"RELATIONSHIP\" is specified.")
// }

func CmdGenerateCatalogus() {
	CmdDiscoverAllResourcesWithOutData()
	data, _ := json.MarshalIndent(globals.ListOfAllTFResourcesWithData, "", "\t")
	os.Remove(globals.OrgRootPath + "/" + ".catalogus")
	if err := os.WriteFile(path.Join(globals.OrgRootPath, ".catalogus"), data, 0444); err != nil {
		log.Println(err)
	}
}

func CmdDiscoverAllResourcesWithOutData() error {
	searchx := globals.CamService.V1.SearchAllResources(scope)
	searchx.ReadMask("name,display_name,asset_type,project,folders,organization,parent_asset_type,parent_full_resource_name,state,relationships")
	// search.Query("relationships:INSTANCE_TO_INSTANCEGROUP")
	// search.Fields([]googleapi.Field{"assets",}...)
	// FetchResourceData(*scope, assType...)
	searchx.AssetTypes(strings.Split(atypes, ",")...)
	searchx.OrderBy("AssetType")
	// log.Println(asst)

	if err := searchx.Pages(globals.GetContext(), func(sarr *cloudasset.SearchAllResourcesResponse) error {
		if sarr == nil {
			return *new(error)
		}
		results := sarr.Results
		for _, res := range results {
			// data, _ := resource.MarshalJSON()

			name := strings.SplitN(res.Name, "/", 4)[3]
			// dname := strings.Split(res.DisplayName, "/")
			asset := &globals.Asset{
				Name:            res.Name,
				DisplayName:     res.DisplayName,
				Ancestors:       res.Folders,
				Project:         res.Project,
				Parent:          res.ParentFullResourceName,
				ParentAssetType: res.ParentAssetType,
				Organization:    res.Organization,
				// Data: nil,
			}
			if strings.EqualFold(res.AssetType, "compute.googleapis.com/Project") && strings.Contains(asset.ParentAssetType, "Project") {
				asset.AssetType = "cloudresourcemanager.googleapis.com/Project"
			} else {
				asset.AssetType = res.AssetType
			}
			time.Sleep(1 * time.Second)
			globals.MapOfAllTFResourcesWithData.LoadOrStore(name, asset)
			globals.ListOfAllTFResourcesWithData = append(globals.ListOfAllTFResourcesWithData, asset)
			// log.Println(res.State,name , dname[len(dname)-1])
			// globals.Glossary[asset.PickResourceId()] = asset.
		}
		time.Sleep(1 * time.Second)
		return nil

	}); err == nil {
		// time.Sleep(2 * time.Second)
		// CmdListAllResourcesWithData()
		return err
	} else {
		log.Println(err)
		return err
	}
}

func CmdListAllResourcesWithData() error {
	search := globals.CamService.Assets.List(scope).ContentType("RESOURCE").AssetTypes(strings.Split(atypes, ",")...).Fields([]googleapi.Field{
		"assets(name,assetType,ancestors,relatedAsset,resource/data)",
	}...)

	if err := search.Pages(globals.GetContext(), func(lar *cloudasset.ListAssetsResponse) error {
		if lar == nil {
			log.Println(lar.ServerResponse)
			return *new(error)
		}

		assets := lar.Assets
		for _, res := range assets {
			name := strings.SplitN(res.Name, "/", 4)[3]
			var data, _ = res.Resource.Data.MarshalJSON()

			var outdata map[string]any
			json.Unmarshal(data, &outdata)
			if as, ok := globals.MapOfAllTFResourcesWithData.Load(name); ok {
				asx := as.(*globals.Asset)
				asx.Data = outdata
				globals.MapOfAllTFResourcesWithData.Store(name, asx)
				globals.ListOfAllTFResourcesWithData = append(globals.ListOfAllTFResourcesWithData, asx)
			} else {
				// asx := as.(*globals.Asset)
				// globals.ListOfAllTFResourcesWithData = append(globals.ListOfAllTFResourcesWithData, asx)
				// log.Println("ERROR:", "Asset not found!", name)
			}
		}

		return nil
	}); err != nil {
		log.Println(err)
		return err
	} else {
		return err
	}
}

func CmdGenerateOrgTreeIndex() {
	data, err := os.ReadFile(path.Join(globals.OrgRootPath, ".catalogus"))
	defer playbooks.TFMetaModuleInstance.WriteMetaModulesIndex()
	if err == nil {
		if err := json.Unmarshal(data, &globals.ListOfAllTFResourcesWithData); err != nil {
			log.Println("ERROR: ", "Failed to read the \".catalogus\" file. \n\tPlease, ensure the file exists by running: kitabus --Generate inventory --[params,..] ")
		} else {
			playbooks.ExecutePlayBook()
		}
	}
}

func CmdGenerateResourceTFTemplates() {
	go InitTFProviders()
	file, err := os.ReadFile(globals.OrgRootPath + "/" + ".kitabus")
	if err != nil {
		log.Println("ERROR:", "Please run kitabus --generate orgtree --path ", err, ".kitabus")
	}
	var data map[string][]*globals.Asset
	if err := json.Unmarshal(file, &data); err == nil {
		//
		globals.ListOfAllTFResourcesPerPath = new(sync.Map)
		// go func (map[string][]*globals.Asset)  {
		for k, v := range data {
			globals.ListOfAllTFResourcesPerPath.LoadOrStore(k, v)
		}
		// }(data)
	} else {
		log.Println(err)
	}
	// log.Println(data)
	playbooks.TFMetaModuleInstance.CreateModules()
}

func CmdGenerateResourceTFModules() {}
func CmdImportResourceTFState()     {}

// func CMDBuilder(opts any) {

// 	GCLOUD_ASSET_CMD = []string{
// 		"gcloud", "asset",
// 		GCLOUD_ASSET_SEARCH_CMD,
// 		"--scope", GCLOUD_ASSET_SEARCH_SCOPE,
// 		// "--filter", GCLOUD_ASSET_FILTER,
// 		"--order-by", GCLOUD_ASSET_ORDER_BY,
// 		"--format", GCLOUD_ASSET_OUT_FMT,
// 		// GCLOUD_ASSET_OUT_PATH,
// 	}
// }

func LocalCmdRunner(pwd string, pipe []any, cmdp ...string) any {

	CLICMD = &exec.Cmd{
		//   Stdout: file,
		Stderr:    os.Stderr,
		WaitDelay: 10 * time.Second,
	}
	CLICMD.Dir = pwd
	CLICMD.Path = "/usr/bin/" + cmdp[0]
	CLICMD.Args = cmdp

	switch pipe[0] {
	case "file":
		log.Println("RUNNING: ", CLICMD.Args, pwd)
		globals.OutFile, _ = os.OpenFile(pipe[1].(string), os.O_CREATE|os.O_WRONLY, 0777)
		CLICMD.Stdout = globals.OutFile
	case false:
		CLICMD.Stdout = nil
	default:
		if pipe[0].(bool) {
			CLICMD.Stdout = os.Stdout
		}

	}
	CLICMD.Run()
	CLICMD.Wait()
	if CLICMD.ProcessState.ExitCode() == 0 {
		log.Println("SUCCEEDED: ", CLICMD.Args, pwd)
	}

	if pipe[0] == "file" {
		globals.OutFile.Close()
		log.Println("CHECK YOUR FILE: ", pipe[1])
		time.Sleep(15 * time.Second)
	}

	return CLICMD.ProcessState
}

func InitTFProviders() {
	data, _ := globals.MetaTFMappingsFile.ReadFile("assets/providers.template")
	if err := os.WriteFile(globals.OrgRootPath+"/providers.tf", data, 0777); err != nil {
		log.Println("FAILED TO CREATE THE TF PROVIDERS FILE", err)
		os.Exit(1)
	} else {
		globals.GlobalCMDRunner(globals.OrgRootPath, []any{false}, "terraform", "init", "-reconfigure")
	}
}
