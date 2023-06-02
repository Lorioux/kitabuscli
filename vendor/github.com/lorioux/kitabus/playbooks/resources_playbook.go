package playbooks

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/lorioux/kitabus/globals"
	"github.com/lorioux/kitabus/mappings"
	// "github.com/lorioux/kitabus/mappings"
)

type TFResourceFactory interface {
	// MakeResource () string
	// DoResourceMetadata () string
	// MakeServiceType () string

	/*Check if resource catalog file's parent folders exist in the directory tree.
	  For instance, a directory pattern e.g. [orgnode/[...folder_name]/project_name]/Service_name/
	  where Service_name corresponds to subdomain (e.g. compute.googleapis.com) part "[compute]".
	  Example 1:
	    /orgnode/ 		# org node
		|--Orgpolicy/ 	# service_name e.g orgpolicy.googleapis.com/Policy
		|  |--*.tf
	  Example 2:
	    /Security/ 		# org unit
		|--Orgpolicy/ 	# service_name e.g orgpolicy.googleapis.com/Policy
		|  |--*.tf
	*/
	CheckDirectoryExistsOrCreate(root string, relative string) bool // create the directory IF NOT EXISTS returning the PATH and BOOL

	/*Check if resource catalog file's exists
		 For instance, a resource catalog file correspond to resource (e.g. compute.googlepais.com/[Instance|Network]) part: "Instance or Network"
		 Example 1:
	    	/Projects/
			|--CloudLabs/           # project_name
			|  |--Compute/          # service_name
			|  |  |--Instance.tf    # resource catalog as terraform configuration file *.tf
			|  |  |--Network.tf     # resource catalog as terraform configuration file *.tf
			|  |--Dns/              # service_name e.g dns.googlepais.com/Policy
			|  |  |--Policy.tf      # resource catalog as terraform configuration file *.tf
	*/
	CheckTFResourceTypeFileExistsOrCreate() bool

	/* Construct the TF Google resource and append into the TF Resource Type File
		   For instance, a TF resource construct will have "[name]" and "[parent]" fields provided by Asset listing output.
		   Hence, if we run $(gcloud asset list --[organization|folder|project] ID --filter "assetType = 'compute.googleapis.com/Instance'" --limit 1 --format json)"
		   Example output:
		   # assets.json
		   [
	        {
	          "ancestors": [
		          "projects/*******",
		          "folders/********",
		          "organizations/********"
	          ],
	          "assetType": "compute.googleapis.com/Instance",
	          "name": "//compute.googleapis.com/projects/[project_id]/zones/europe-central2-a/instances/worker",
	        },
			....]
			# Instance.tf
			...
			[truncated]
			resource "google_compute_instance" "worker" {
				name = "projects/[project_id]/zones/europe-central2-a/instances/worker" #(REQUIRED)
				parent = "projects/*******" #(optional)
			}
	*/
	ConstructTFGCPResource() string

	/* Reconstruct TF Google resource based on the import state values
	 */

	ConstructTFResourceTemplate() string

	IsSupportedTFResource(tftype string, gtype string) bool

	ParseComplexTFResourceParameters() any
}

type TFResourceType struct {
	Parent         string `json:"parent"`
	Name           string
	RequiredFields []string `json:"required_fields"` // map[string]interface{}
	TfrType        string   `json:"resource_type"`
	TfrLabel       string   `json:"label"`
	Path           string   `json:"path"`
	Type           string   `json:"type"`
	Mirror         *globals.Asset
	importer       *schema.ResourceImporter
	data           schema.ResourceData
	TfrInputFile   string
}

type ServiceType struct{}

// var KindResourceCounter map[string]int
// var file *os.File;

var fileOerr error

var ancestors string

var TFMetaModuleInstance *TFMetaModule

func init() {

	TFMetaModuleInstance = &TFMetaModule{
		ListOfTFResources: &sync.Map{},
		ListOfTFFilePaths: []string{},
	}
}

/**
* Let resource be asset of type e.g. compute.googlepais.com/Instance, therefore:
* 	1. Let name be string formatted as "projects/{project_id}/zones/{zone_name}/instances/[name]", and
*   2. Let required_fields be a map of key=value, where key is a required fields different to [name and parent]. And
*   3. Let parent be (optional) e.g. projects/{project_id}, so
* So, there should be a parent folder where:
	1. Folder name is "Compute", so it should have a file with:
	2. File name is "Instance.tf"
*/

func (tfr *TFResourceType) SetTFResourceType(s string) {
	tfr.TfrType = s
}

func (tfr *TFResourceType) GetTFResourceType() string {
	return tfr.TfrType
}

func (tfr *TFResourceType) CheckDirectoryExistsOrCreate(p ...string) bool {
	// TODO: Concatenate the root and relative path  then test if exists
	// TODO: Create the tree IF NOT EXIST
	base := path.Join(p...)
	if _, err := os.ReadDir(base); err != nil {
		if err := os.MkdirAll(base, 0777); err == nil {
			log.Printf("Adding path: %v", base)
		} else {
			return false
		}
	}
	return true
}

func (tfr *TFResourceType) CheckTFResourceTypeFileExistsOrCreate(p ...string) bool {
	//TODO: For each asset type create a file in the directory tree
	filepath := path.Join(p...) + ".tf"
	if file, err := os.OpenFile(filepath, os.O_CREATE|os.O_SYNC, 0777); err == nil {
		tfr.SetTFResourceFilePath(filepath)
		return file.Close() == nil
	}
	return false
}

func (tfr *TFResourceType) SetRequiredFields(p ...string) *TFResourceType {
	tfr.TfrType = globals.Reversor.(*mappings.ResourceReversor).MatchTFRNameByCloudRSType(tfr.Mirror.AssetType).(string)

	// if tfr.TfrType =="" {
	// 	// log.Printf("REQUIRED.....: %v", tfr.Mirror.AssetType)
	// 	return nil
	// }

	tfr.Mirror.SetSupported(true)

	// log.Print("HERE WALK")
	tfr.Name = tfr.Mirror.DisplayName

	if rq, err := mappings.PickResourceRequiredFieldsByTFRName(tfr.TfrType); err == nil {
		tfr.RequiredFields = rq
	} else {
		return nil
	}
	// go tfr.SetTFResourceType(tfrType)
	// KindResourceCounter[tfr.TfrType] += 1
	// count := KindResourceCounter[tfr.TfrType]
	tfr.Parent = tfr.Mirror.PickParentId()

	tfr.TfrLabel = strings.ToLower(tfr.Mirror.PickResourceType()) //fmt.Sprint(tfr.TfrType,"_",count)
	tfr.TfrInputFile = fmt.Sprintf(".%s", tfr.TfrLabel)
	// go mappings.TFImportState(tfrType.(string), tfr.TfrLabel, tfr.Name)
	// log.Printf("RSTYPE: %v --- TF: %v", tfr.Mirror.AssetType, tfr.TfrType)
	tfr.SetTFResourceFilePath(p...)
	// go tfr.SetTFResourceImporter()
	return tfr
}

func (tfr *TFResourceType) ConstructTFResourceTemplate() any {

	// tfrType := mappings.MatchTFRNameByCloudRSType()
	const form = `{{block "main" .}}{{"\n" -}}
	locals { {{$label := .TfrLabel}}{{$rtype := .TfrType}}{{$name := .Name}}{{$parent := .Parent}}{{$hasname := false}}{{$hasparent := false}}{{"\n"}} {{"\n\t" -}}
		{{$rtype}}    = jsondecode(file("${var.{{$rtype}}["base_states"]}"))["resources"][0] {{- "\n\t" -}}
		{{$label}}_attrs = {for resource in local.{{$rtype}}.instances : (resource["index_key"]) => (resource["attributes"])}{{"\n" -}}
	}{{- "\n\n" -}}
	variable "{{$rtype}}" {
	type = object({
		names = list(string)
		base_states = string
	}){{"\n"}} }{{- "\n\n" -}}
	resource "{{$rtype}}" "{{$label}}" {
	for_each     = { for name in var.{{$rtype}}["names"] : name => local.{{$label}}_attrs["${name}"] }
	# required fields{{"\n\t"}}
    {{- range $i, $req := .RequiredFields -}}
	{{$isname := ne $req "name"}}{{$isparent := ne $req "parent"}}
	{{- if eq $req "project" -}}
		{{- $req}} = "projects/{{$parent -}}"{{"\n\t"}}
	{{- else -}}
		{{- $req}} = each.value["{{$req -}}"]{{"\n\t"}}{{end}}{{end -}}
    {{"\n"}}}{{"\n\n"}}{{- end}}`

	resource, err := template.New(strings.ToLower(tfr.TfrType)).Parse(form)
	// resource, err := resource.ParseFiles("./playbooks/.meta_template")
	if err != nil {
		log.Printf("FAILED TEMPLATE PARSE: %v", err)
	}

	file, err := os.OpenFile(tfr.Path, os.O_WRONLY, 0777)
	if err != nil {
		log.Printf("FAILED TO OPEN FILE %v TO APPEND TEMPLATE. %v", tfr.Path, err)
	}
	defer file.Close()

	if err := resource.Execute(file, tfr); err == nil {
		log.Printf("ADDING TF TEMPLATE FOR: %v --- INTO %v", tfr.TfrType, tfr.Path)
		// log.Panic(err)
	}

	return "Done!"
}

func (tfr *TFResourceType) ConstructTFResourceLocalValuesFile() any {
	// Construct TF Local Variable Template
	// Get the resource file path and add a JSON File
	rspath := path.Dir(tfr.Path)
	if file, err := os.OpenFile(path.Join(rspath, tfr.Mirror.PickResourceName()), os.O_CREATE|os.O_APPEND, 0777); err == nil {
		os.WriteFile(rspath, []byte{}, 0777)
		file.Close()
	}

	const local = `{{block "main .}}
	local {
		{{range .}}
	}{{"\n\n"}}`
	return nil
}

func (tfr *TFResourceType) SetTFResourceFilePath(p ...string) {
	tfr.Path = path.Join(p...)
}
