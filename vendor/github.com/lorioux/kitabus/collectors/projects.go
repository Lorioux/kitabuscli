package collectors

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	crm "google.golang.org/api/cloudresourcemanager/v1"


	"github.com/lorioux/kitabus/globals"
)



type Project struct {
	Identifier map[string][]string
	Parents map[string]string 	// {parentid:[ parentname, level index]}
	Assets 	map[string][]string // {asset_type: [names...]}
}

var ActiveProjects []map[string]string
var LifecycleState any

// list all projects
func GetProjects() []map[string]string {
	if LifecycleState == nil {
		LifecycleState = "ACTIVE"
	}
	CrmService, err := crm.NewService(globals.GetContext())
	if err == nil {

		response, err := CrmService.Projects.List().Filter(fmt.Sprintf("lifecycleState=%s", LifecycleState)).Do()

		if err != nil {
			log.Fatal("Error: ", err)
		}

		for i := 0; i < len(response.Projects)-1; i++ {
			projects, _ := response.Projects[i].MarshalJSON()
			var holder map[string]string
			if err := json.Unmarshal(projects, &holder); err != nil {
				// fmt.Println(holder)
				ActiveProjects = append(ActiveProjects, holder)
			}
		}
		// fmt.Println(ActiveProjects)
	} else {
		log.Fatal(err)
	}
	return ActiveProjects
}

func NewProjectInstance() Project {
	p := Project{}
	p.Parents = map[string]string{}
	p.Assets = map[string][]string{}
	p.Identifier = map[string][]string{}
	return p
}

func (p *Project) SetName(s string){
	id := strings.Split(s, "/")[1]
	if iden := p.Identifier[id]; iden == nil {
		// log.Print(s)
		p.Identifier[id] = FetchResourceNameById(s).([]string)
	}
}


func (p *Project) SetParents(pa []interface{}) {
	for _, pId := range pa {
		name := FetchResourceNameById(pId.(string))
		if anc := p.Parents[pId.(string)]; anc == "" {
			p.Parents[pId.(string)] = name.(string)
		}
	}
}

func (p *Project) AddAssets(t string, n string) {
	// get asset list and append and a new asset for specific asset type
	nn := strings.Split(n, "/")[0:]
	nx :=  nn[len(nn)-1]
	if assetType := p.Assets[t]; assetType == nil {
		p.Assets[t] = []string{nx}
	} else if assetType != nil {
		p.Assets[t] = append(p.Assets[t], nx)
	}
}