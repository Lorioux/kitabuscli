package playbooks

import (
	"log"
	"net/url"
	"path"
	"sync"

	"github.com/lorioux/kitabus/collectors"
	"github.com/lorioux/kitabus/globals"
	"github.com/lorioux/kitabus/mappings"
)

func ExecutePlayBook() {

	// file, err := os.ReadFile(assetsFile)
	// defer file.Close()
	// var broker map[string]interface{}

	// if err == nil {
	// 	if err := json.Unmarshal(file, &content); err != nil {
	// 		log.Println(err)
	// 	} else {
	// ConstructResourceHierarchy()
	// UserFriendlyNames = new(sync.Map)
	// go func() {
	// 	for _, con := range content {
	// 		UserFriendlyNames.LoadOrStore(con.PickResourceName(), []any{con.DisplayName, con.PickParentId()})
	// 	}
	// }()

	walkThroughOrgTreeResources()
	// 	}
	// } else {
	// 	log.Fatalln("Error: ", err)
	// }
	// defer mappings.UpdateUnsupportedResourceMap()
}

func walkThroughOrgTreeResources() {

	// mux := new(sync.RWMutex)
	mapPaths := &map[string][2]any{}
	muw := new(sync.WaitGroup)

	for index, con := range globals.ListOfAllTFResourcesWithData {
		/// Check if resource type is in the allowed list.
		if globals.Reversor.(*mappings.ResourceReversor).MatchTFRNameByCloudRSType(con.AssetType) == nil {
			continue
		}
		// else if !(con.PickResourceType() == "Instance") {
		// 	continue
		// }

		resource := &TFResourceType{Mirror: con}
		resource.Mirror.IsSupported = true
		globals.Control.Store(index, []any{false, nil})
		muw.Add(2)
		go resource.walkOrgTreeInBackground(index, muw)
		go resource.makeUserFriendlyResourcePath(index, muw, *mapPaths)
		muw.Wait()
	}
}

func (tfr *TFResourceType) walkOrgTreeInBackground(
	// rstypes []string,
	index int,
	muw *sync.WaitGroup) {

	defer muw.Done()
	if name, err := url.QueryUnescape(tfr.Mirror.Name); err == nil {
		tfr.Mirror.Name = name
	}
	if parent, err := url.QueryUnescape(tfr.Mirror.PickParentId()); err == nil {
		tfr.Parent = parent
	}
	go tfr.SetRequiredFields()

}

func (tfr *TFResourceType) makeUserFriendlyResourcePath(
	index int,
	muw *sync.WaitGroup,
	mapPaths map[string][2]any) {

	defer muw.Done()
	var pathTemp string

	work := func(c [2]any, a string) {
		switch c {
		case [2]any{}:
			nameT := collectors.FetchResourceNameById(a)
			pathTemp = path.Join(nameT.(string), pathTemp)
			c = [2]any{true, pathTemp}
		default:
			pathTemp = path.Join(c[1].(string), pathTemp)
			c = [2]any{true, pathTemp}
			// continue
		}
	}

	lineage := func() {
		ancestors := tfr.Mirror.Ancestors
		for _, a := range ancestors {
			work(mapPaths[a], a)
		}
	}

	// Only execute the step if path is not yet in the globals.Control at the
	// specified index.
	contx, _ := globals.Control.Load(index)
	cx := contx.([]any)
	if !cx[0].(bool) && cx[1] == nil {
		// work(mapPaths[tfr.Mirror.Organization], tfr.Mirror.Organization)
		switch tfr.Mirror.PickParentType() {
		case "Project":
			work(mapPaths[tfr.Mirror.Project], tfr.Mirror.Project)
			lineage()
			pathTemp = path.Join(pathTemp, ".ACTIVE", tfr.Mirror.PickServiceType())
		case "Folder":
			lineage()
			pathTemp = path.Join(pathTemp, ".META", tfr.Mirror.PickServiceType())

		case "Organization":
			pathTemp = path.Join(pathTemp, ".META", tfr.Mirror.PickServiceType())
		default:
			var name = tfr.Mirror.PickResourceType()
			if name == "Project" {
				work(mapPaths[tfr.Mirror.PickResourceName()], tfr.Mirror.PickResourceName())
				lineage()
				pathTemp = path.Join(pathTemp, ".META", tfr.Mirror.PickServiceType())
			} else {
				if tfr.Mirror.Project != "" {
					work(mapPaths[tfr.Mirror.Project], tfr.Mirror.Project)
					lineage()
					pathTemp = path.Join(pathTemp, ".ACTIVE", tfr.Mirror.PickServiceType())
				}
			}
		}
	}

	if pathTemp != "" {

		work(mapPaths[tfr.Mirror.Organization], tfr.Mirror.Organization)
		if _, ok := globals.Control.Swap(index, []any{true, pathTemp, tfr.Mirror.PickResourceType()}); !ok {
		}

		// log.Println(globals.Control.Load(index))
		muw.Add(1)
		go tfr.ConstructResourceHierarchy(index, muw)
	}
}

func (rsc *TFResourceType) ConstructResourceHierarchy(
	index int,
	muw *sync.WaitGroup) {

	defer muw.Done()
	cm, _ := globals.Control.Load(index)
	counter := 1
	if len(cm.([]any)) < 3 /*Minimum length of */ {
		log.Println("FAILED:---", counter, "---", index)
		return
	}

	cx := cm.([]any)

	if rsc.CheckDirectoryExistsOrCreate(globals.OrgRootPath, cx[1].(string)) {
		// log.Printf("RH CONSTRUCT: %v --- DIR: %v => %v", index, cm[1], cm[3])
		if !rsc.CheckTFResourceTypeFileExistsOrCreate(globals.OrgRootPath, cx[1].(string), cx[2].(string)) {
			log.Panicf("FAILED TO ADD THE RESOURCE TYPE %s", cx[2])
		} else {
			log.Println("TF File added successful:---", index, "---", cx[1], cx[2])
		}

		if rsc.Mirror.IsSupported { //&& rsc.GetTFResourceType() != ""
			// go rsc.ConstructTFResourceTemplate()
			go TFMetaModuleInstance.RegisterTFResources(rsc)
		}
	}
}
