package playbooks

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"
	"time"

	// "time"

	"github.com/lorioux/kitabus/globals"
	"github.com/lorioux/kitabus/mappings"
)

type TFImportModulesFactory interface {
	/*Check  TF importing resource meta file EXISTS OR CREATE NEW.
	  For instance, a TF importing resource meta file contains a reference of service folders as modules"
	  Example 1:
	    /Projects/
	    |--CloudLabs/           # project_name
	    |  |--meta_modules.tf 	# Terraform modules for imports
	  Example 2:
	    /[orgnode | orgunit]/   # org node
	    |--meta_modules.tf      # Terraform importing resource modules
	*/

	CheckTFMetaModulesFileExistsOrCreate() bool

	/*Construct a meta module for importing resource type. For instance. meta modules have source field
	    pointing to the resources relative path.
	    Example 1:
		    /Projects/
		    |--CloudLabs/
		    |  |--Compute/
		    |  |  |--Instance.tf
	         |--meta_modules.tf   # content example
	             [truncated]
	             module "compute" {
	                source = "./compute"
	             }
	             [truncated]
	*/

	ConstructTFMetaModuleTemplate() string
}

/*
Also it can have a `resources block, optional` as:

	resources = list(object({
	 resource_type = string  # e.g. google_compute_instance, google_compue_disk
	 identifier    = string  # e.g. workernode, workernode-disk
	 resource_id   = string  # e.g. projects/-/zones/-/instances/workernode
	}))
*/

var ListOfModulesTFResources = new(sync.Map)

type TFMetaModule struct {
	// source string
	//resources []interface{}
	Scope                      string
	Path                       string
	ListOfTFResources          *sync.Map
	ListOfTFFilePaths          []string
	ModuleBaseAutoVarFilesPath map[string]string
	ModuleServiceType          string
	ModuleResourceNamesByType  map[string][]string
	ModuleResourceMirrorByName map[string]*TFResourceType
}

type Resource struct {
	Type      string                   `json:"type"`
	Name      string                   `json:"name"`
	Instances []map[string]interface{} `json:"instances"`
}

type TFModuleInstance struct {
	source    string
	resources map[string][]Resource
}

func (tfmi *TFModuleInstance) ConstructTFRBaseValuesTemplate(tfrp string, tfrs any) {
	var filePath string
	instances := []map[string]any{}

	if strings.Contains(tfrp, ".ACTIVE") {
		filePath = strings.Split(tfrp, ".ACTIVE")[0]
	} else {
		filePath = strings.Split(tfrp, ".META")[0]
	}

	rs := tfrs.([]*TFResourceType)

	filePath = path.Join(filePath, ".BASE")
	if err := os.MkdirAll(filePath, 0777); err != nil {
		log.Println(err)
	}

	filePath = path.Join(filePath, rs[0].TfrInputFile)
	//for each tf resouce create an instance object to map to Instances
	for _, tfr := range rs {

		instance := map[string]any{
			"index_key": tfr.Name,
		}
		attributes := map[string]any{}
		var data map[string]any
		rt, _ := json.Marshal(tfr)

		if err := json.Unmarshal(rt, &data); err != nil {
			log.Println("FAILED TO UNMARSHAL")
		}

		for _, attr := range tfr.RequiredFields {
			attributes[attr] = data[attr]

		}
		instance["attributes"] = attributes
		instances = append(instances, instance)
	}

	tfmi.resources = map[string][]Resource{}
	tfmi.resources["resources"] = []Resource{}
	tfmi.resources["resources"] = append(tfmi.resources["resources"], Resource{
		Type:      rs[0].TfrType,
		Name:      rs[0].TfrLabel,
		Instances: instances,
	})

	// tfmi.resources[0]
	file, _ := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0777)
	if data, err := json.MarshalIndent(tfmi.resources, "", "\t"); err == nil {
		if err := os.WriteFile(filePath, data, 0777); err != nil {
			log.Println("ERROR FAILED TO CREATE:", err, string(data), filePath)
		}
		file.Close()
	}

}

func (tfm *TFMetaModule) ConstructTFMetaModuleTemplate(root string) string {
	const meta = (`{{block "main" . -}}{{$paths := .ModuleBaseAutoVarFilesPath}}{{$names := .ModuleResourceNamesByType}}{{- "\n" -}}
  module "{{- .ModuleServiceType -}}" { {{- "\n\t" -}}
    source = "./{{.ModuleServiceType}}" {{- "\n\t" -}}
    {{- range $rtype, $rnames := $names}}{{$rtype_suffix := "_resources" -}}
    {{- $rtype}} = {  {{- "\n\t\t" -}}
      names = [{{range $nx := $rnames}}"{{$nx}}",{{end}}] {{- "\n\t\t" -}}
      base_states = "{{index $paths $rtype}}"
    }{{end}} {{- "\n" -}}
  }{{- "\n\n\t" -}}{{- end}}`)
	// resource_base_inputs = {
	//{{$rtype_suffix -}}
	//   {{range $rtype, $inpath := .ModuleBaseAutoVarFilesPath}}
	//     {{$rt}} = {{$inpath}} {{"\n\t"}}
	//   {{end}}{{"\n\t"}}}
	/*
	 {{range $tpe, $path := .ModuleBaseAutoVarFilesPath}}
	      {{if eq $tpe $rtype}}base_states = {{$path}}{{else}}{{end}}
	      {{end}}
	*/

	module, err := template.New(strings.ToLower(tfm.ModuleServiceType)).Parse(meta)
	// resource, err := resource.ParseFiles("./playbooks/.meta_template")
	if err != nil {
		log.Println("FAILED TEMPLATE PARSE:", err)
	}
	file, err := os.OpenFile(root+"/modules.tf", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Panicf("FAILED TO OPEN FILE: %v TO APPEND TEMPLATE. %v", "tfr.Path", err)
	}
	defer file.Close()

	if err := module.Execute(file, tfm); err == nil {
		log.Printf("ADDED TF MODULE TEMPLATE FOR: %v --- INTO %v", tfm.ModuleServiceType, root)
	} else {
		log.Println("ERROR FAILED TO ADD TF MODULE: ", err, root)
	}

	return ""
}

func (tfm *TFMetaModule) CreateModules() {

	muw := new(sync.WaitGroup)
	mux := new(sync.RWMutex)

	globals.ListOfAllTFResourcesPerPath.Range(func(key, value any) bool {
		muw.Add(1)
		tfm.ModuleResourceNamesByType = map[string][]string{}
		tfm.ModuleBaseAutoVarFilesPath = map[string]string{}
		tfm.ModuleResourceMirrorByName = map[string]*TFResourceType{}
		var modulePath string
		rsPath := path.Dir(key.(string))

		tfm.Scope = key.(string)

		if strings.Contains(rsPath, ".ACTIVE") {
			modulePath = strings.SplitAfter(rsPath, ".ACTIVE")[0]
		}
		if strings.Contains(rsPath, ".META") {
			modulePath = strings.SplitAfter(rsPath, ".META")[0]
		}

		terraProviders := path.Join(globals.OrgRootPath, "providers.tf")
		// terraBinary := path.Join(os.Getenv("PWD"),".terraform")

		go func(tfm *TFMetaModule, p string, muw *sync.WaitGroup, mx *sync.RWMutex) {
			mx.RLock()

			defer mx.RUnlock()
			defer muw.Done()
			defer tfm.ExecuteTFResourceImportState()
			assets := value.([]*globals.Asset)
			resources := []*TFResourceType{}
			// var provider_template = `

			// provider "google" {
			//   project = "%v"
			// }
			// `
			tfm.Path = modulePath
			if file, err := os.OpenFile(modulePath+"/providers.tf", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777); err == nil {
				globals.GlobalCMDRunner(globals.OrgRootPath, []any{false}, "cp", terraProviders, modulePath)
				providerSym := path.Join(modulePath, ".terraform")
				if err := os.MkdirAll(providerSym, 0777); err != nil {
					log.Println("FAILED: ", err)
				}

				if _, err := os.Readlink(providerSym + "/providers"); err != nil {
					if err := os.Symlink(path.Join(globals.OrgRootPath, ".terraform/providers/"), providerSym+"/providers"); err != nil {
						log.Println("FAILED: ", err)
					}
				}

				// mprovider := fmt.Sprintf(provider_template, assets[0].PickParentId() )

				// if _, err := file.Write([]byte(mprovider)); err != nil {
				//   log.Println(err)
				// }
				file.Close()

			} else {
				log.Fatal("ERROR:", err)
			}

			for _, mirror := range assets {
				resource := &TFResourceType{
					Mirror: mirror,
				}
				resource.SetRequiredFields(p)
				resources = append(resources, resource)
				resource.ConstructTFResourceTemplate() //.InvokeTFResourceImporter() //

				tfm.ModuleServiceType = mirror.PickServiceType()
				// tfm.ModuleResourceNamesByType = map[string][]string{
				//   mirror.PickResourceType() : []string{},
				// }

				tfm.ModuleResourceNamesByType[resource.TfrType] = append(tfm.ModuleResourceNamesByType[resource.TfrType], resource.Name)
				tfm.ModuleResourceMirrorByName[resource.Name] = resource
				tfm.ModuleBaseAutoVarFilesPath[resource.TfrType] = fmt.Sprint("../", ".BASE", "/", resource.TfrInputFile)
			}

			instances := TFModuleInstance{}
			instances.ConstructTFRBaseValuesTemplate(rsPath, resources)
			tfm.ConstructTFMetaModuleTemplate(modulePath)
			globals.GlobalCMDRunner(modulePath, []any{false}, "terraform", "init", "-lock=false")
			// states := make(mappings.LocalTFStateContent)
			// states.HandleTFStateDeconstructs(resources[0].TfrType, tfm.Scope, nil)

		}(tfm, key.(string), muw, mux)

		muw.Wait()
		return true
	})
}

func (tfm *TFMetaModule) RegisterTFResources(tfr *TFResourceType) {
	if val, ok := globals.ListOfAllTFResourcesPerPath.Load(tfr.Path); ok {
		val = append(val.([]*globals.Asset), tfr.Mirror)
		globals.ListOfAllTFResourcesPerPath.Swap(tfr.Path, val)
	} else {
		globals.ListOfAllTFResourcesPerPath.Store(tfr.Path, []*globals.Asset{tfr.Mirror})
		// tfm.RegisterTFResourcePath(tfr.Path)
		// log.Printf("FAILED TO STORE")
	}
}

/*META Index of modules stored in .kitabus */
func (tfm *TFMetaModule) WriteMetaModulesIndex() {
	var data = map[string][]*globals.Asset{}
	globals.ListOfAllTFResourcesPerPath.Range(func(key, value any) bool {
		r := value.([]*globals.Asset)
		data[key.(string)] = r
		return true
	})

	// file , err := os.OpenFile(".kitabus", os.O_CREATE|os.O_WRONLY, 0777)
	// if err != nil {
	//   log.Println("ERROR: Failed to open \".kitabus\" file", err)
	// }
	// defer file.Close()
	if d, err := json.MarshalIndent(data, "", "\t"); err == nil {
		os.Remove(".kitabus")
		if err := os.WriteFile(".kitabus", d, 0444); err != nil {
			log.Println(err)
		}
	} else {
		log.Println("ERROR: Failed to write content into the \".kitabus\" file", err)
	}
}

func (tfm *TFMetaModule) OverrideMetaModulesOrRemove(mpath string, oper string) {
	operator := func(p *string, o string) {
		switch o {
		case "modules[override]":
			if err := os.Remove(*p); err != nil {
				log.Println("ERROR:", err)
			} else {
				log.Println("REMOVED:", *p)
			}
		}
	}
	// log.Println("HEERE",3)
	override := func(mpath string) bool {
		// mpath = key.(string)

		if strings.Contains(mpath, ".ACTIVE") {
			mpath = strings.SplitAfter(mpath, ".ACTIVE")[0] //+ "/" + "modules.tf"
		}
		if strings.Contains(mpath, ".META") {
			mpath = strings.SplitAfter(mpath, ".META")[0] //+ "/" + "modules.tf"
		}
		log.Println("REMOVE", mpath)
		go operator(&mpath, oper)

		return true
	}

	if strings.EqualFold(mpath, "") {
		data, _ := os.ReadFile(".kitabus")
		var paths any
		if err := json.Unmarshal(data, &paths); err == nil {
			for p, _ := range paths.(map[string]any) {
				override(p)
			}
		} else {
			log.Println(err)
		}
	} else {
		go operator(&mpath, oper)
	}
}

func (tfm *TFMetaModule) ExecuteTFResourceImportState() {
	mw := new(sync.WaitGroup)
	mw.Add(1)
	go func(tfm *TFMetaModule, mw *sync.WaitGroup) {

		for name, rs := range tfm.ModuleResourceMirrorByName {

			nameTF := name

			typeTF := rs.TfrType
			lableTF := rs.TfrLabel
			// parentTF := rs.Mirror.PickParentId()
			addressTF := fmt.Sprintf("module.%s.%s.%s[\"%s\"]", tfm.ModuleServiceType, typeTF, lableTF, nameTF)
			// idTF := fmt.Sprintf("%s/%s",parentTF,rs.Mirror.Name)
			globals.GlobalCMDRunner(tfm.Path, []any{true}, "terraform", "import", addressTF, rs.Mirror.PickResourceId())
			time.Sleep(5 * time.Second)
			// break

		}
		mw.Done()
	}(tfm, mw)
	// Update base input values
	mw.Wait()
	states := &mappings.LocalTFStateContent{}
	states.SetTFStatesFilePath(tfm.Path + "/" + "terraform.tfstate")
	states.ReadTFStatesFile()
	for tp, _ := range tfm.ModuleResourceNamesByType {
		states.UpdateBaseValuesFile(tp)

		states.HandleTFStateDeconstructs(tp, tfm.Scope, nil)
	}
}
