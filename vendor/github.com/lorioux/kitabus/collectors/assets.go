package collectors

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	rsm "cloud.google.com/go/resourcemanager/apiv3"
	rpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	cam "google.golang.org/api/cloudasset/v1"

	"github.com/lorioux/kitabus/globals"
)

var AllProjectsAssets map[string]Project

// get all projects assets
func GetAssetsByProjectId(c chan []map[string]Project, id ...string) any {

	camService, _ := cam.NewService(globals.GetContext())
	AllProjectsAssets = map[string]Project{}
	if id == nil {
		for _, project := range ActiveProjects {
			// get the project ID
			pid := project["projectId"]

			// if !strings.Contains(pid, "cloudlabs") {
			// 	continue
			// }
			ListProjectAssets(pid, camService)
		}
	} else {
		for _, i := range id {
			ListProjectAssets(i, camService)
		}
	}
	return AllProjectsAssets
}

func ListProjectAssets(id string, camService *cam.Service) {
	// list all assets in the project
	call := camService.Assets.List(fmt.Sprintf("projects/%s", id))
	call.ContentType("RESOURCE")
	var holder map[string]interface{}
	var project = NewProjectInstance()
	// start the call
	if response, err := call.Do(); err == nil {
		// go func() {
		for index, assets := range response.Assets {

			resolve, _ := assets.MarshalJSON()
			// log.Println(resolve)

			if err := json.Unmarshal(resolve, &holder); err != nil {
				log.Printf("[WARNING] %s: %v", err, assets)
			}

			// log.Println(holder)
			parent := holder["ancestors"].([]interface{})[0]

			key := parent.(string)
			project.SetName(key)

			if index == 0 {
				project.SetParents(holder["ancestors"].([]interface{})[1:])
				// ca <- []string{holder["name"].(string)}
			}
			// Add assets
			project.AddAssets(holder["assetType"].(string), holder["name"].(string))
			// Ancestors[key] = append(Ancestors[key], holder["assetType"])
		}
		// }()
	} else {
		log.Fatal("Failed to retrieve assets in the project: ", id, "\n", err)
	}
	AllProjectsAssets[id] = Project(project)
}

func FetchResourceNameById(s string) any {
	var holder any
	// var q = cam.V1QueryAssetsCall{}
	if strings.Contains(s, "projects") {
		// Retrieve the project display name
		// call := CrmService.Projects .Get(s)
		req := &rpb.GetProjectRequest{Name: s}
		cli, err := rsm.NewProjectsClient(globals.GetContext())
		if err != nil {
			log.Fatal(err)
		}
		if project, err := cli.GetProject(globals.GetContext(), req); err == nil {
			holder = project.DisplayName //[]string{project.DisplayName, project.ProjectId,}
		}
		// log.Printf("Project: %v", holder)
		defer cli.Close()
	}

	if strings.Contains(s, "folders") {
		req := &rpb.GetFolderRequest{Name: s}
		cli, err := rsm.NewFoldersClient(globals.GetContext())
		if err != nil {
			log.Fatalf("[ERROR] Folders client failure... %v", err)
		}
		defer cli.Close()
		if folder, err := cli.GetFolder(globals.GetContext(), req); err == nil {
			holder = folder.DisplayName

		}
	}

	if strings.Contains(s, "organizations") {
		req := &rpb.GetOrganizationRequest{Name: s}
		cli, err := rsm.NewOrganizationsClient(globals.GetContext())
		if err != nil {
			log.Fatal("[ERROR] Folders client failure...")
		}
		defer cli.Close()
		if org, err := cli.GetOrganization(globals.GetContext(), req); err == nil {
			holder = org.DisplayName

		}
	}
	if holder == nil {
		log.Panicf("FAILED TO RETRIEVE DISPLAY NAME FOR ID: %v", s)
	}
	return holder
}
