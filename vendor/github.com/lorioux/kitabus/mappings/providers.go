package mappings

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ProviderGA "github.com/lorioux/google-provider/google"
	"github.com/lorioux/kitabus/globals"
)

type (
	LocalTFStateContent map[string]any
	resourceWalker      struct {
		r             schema.Schema
		attrs         map[string][]any
		attr          string
		data          map[string]any
		schemas       map[string]*schema.Schema
		borders       []int
		ladder        int
		keysTemplate  []any
		attrsTemplate strings.Builder
		body          any
	}
)

const (
	TypeInvalid int = iota
	BlockStart
	BlockEnd
	MapStart
	MapEnd
	ListStart
	ListEnd
	SetStart
	SetEnd
	BaseType
)

var (
	ProviderRequiredParams []string
	ResourceRequiredParams *sync.Map
	Provider               *schema.Provider
	ResourceMap            map[string]*schema.Resource
	walkoverTopAttrs       func(ancestors []any, attr any, s *schema.Schema)
	walkoverElemAttrs      func(a []any, i any, e any)
	// Cfg *meta.Config

	blockDelimiters     = []string{"\t{\n", "\n}\n"}
	mapDelimiters       = []string{"\t{\n", "\n}\n"}
	setDelimiters       = []string{"\t[\n", "\n]\n"}
	baseTags            = []string{"\n", "\n"}
	anyBody         any = nil
)

func init() {
	ProviderRequiredParams = []string{}
	ResourceRequiredParams = new(sync.Map)
	Provider = ProviderGA.Provider()
	ResourceMap = Provider.ResourcesMap
	// walkoverTopAttrs = ResourceWalker.WalkOverTopAttrs
	// walkoverElemAttrs = walkOverElemAttrs[*any]
}

func ExecuteCallProvider() {
	for param := range Provider.Schema {
		if !Provider.Schema[param].Optional {
			ProviderRequiredParams = append(ProviderRequiredParams, param)
		}
	}
}

func (rw *resourceWalker) handleDelimiter(k any, d any, dt any, dc any, dk ...any) {

	if d == dt {
		switch dt {
		case BlockStart, ListStart:
			rw.body = fmt.Sprint(k, blockDelimiters[0])

		case MapStart:
			// anyBody = fmt.Sprint(dk...)
			rw.body = fmt.Sprint(k, " = ", mapDelimiters[0])
		case SetStart:
			rw.body = fmt.Sprint(k, setDelimiters[0])
		}
		rw.keysTemplate = append(rw.keysTemplate, k)
	} else {
		rw.body = fmt.Sprint(dc)
	}
}

func (rw *resourceWalker) bordersBuilder(k any, v any, delimiter any, s *schema.Schema) {

	border := len(rw.borders) - 1
	switch delimiter {
	case BlockStart, BlockEnd:
		rw.handleDelimiter(k, delimiter, BlockStart, blockDelimiters[1]) //block...
	case MapStart, MapEnd:
		// var block = []any{k, " = ", mapDelimiters[0]}
		rw.handleDelimiter(k, delimiter, MapStart, mapDelimiters[1]) //block...
	case ListStart, ListEnd:
		// var block = []any{k, blockDelimiters[0]}
		rw.handleDelimiter(k, delimiter, ListStart, blockDelimiters[1]) //block...
	case SetStart, SetEnd:
		// var block = []any{k, setDelimiters[0]}
		rw.handleDelimiter(k, delimiter, SetStart, setDelimiters[1]) //block...
	default:
		border = rw.borders[border]
		if border == MapStart || border == ListStart {
			rw.body = fmt.Sprint(k, " = ", v, "\n")
		}
		if border == SetStart {
			rw.body = fmt.Sprint(v, ",")
		}
		if border == BaseType {
			rw.body = fmt.Sprint(baseTags[0], k, " = ", v, baseTags[1])
			rw.keysTemplate = append(rw.keysTemplate, k)
		}
	}

	rw.attrsTemplate.WriteString(rw.body.(string))
	// return block.String()
}

func (rw *resourceWalker) handle(attr any, s *schema.Schema) {

	// Handle primitive types.
	switch s.Type {
	case schema.TypeString, schema.TypeBool, schema.TypeFloat, schema.TypeInt:
		// ancestors = append(ancestors, attr)
		if len(rw.borders) == 0 {
			rw.borders = append(rw.borders, BaseType)
		}
		if attr == "" {
			rw.bordersBuilder(rw.attr, nil, BaseType, s)
		} else {
			rw.bordersBuilder(attr, nil, BaseType, s)
		}

	// Handle list, set, and map types.
	case schema.TypeSet:
		rw.borders = append(rw.borders, SetStart)
		// ancestors = append(ancestors, " ", "{\n")
		rw.bordersBuilder(attr, nil, SetStart, s)
		rw.WalkOverElemAttrs(attr, s.Elem)
		rw.bordersBuilder(nil, nil, SetEnd, s)

	case schema.TypeList:
		rw.borders = append(rw.borders, ListStart)
		// ancestors = append(ancestors, " ", "{\n")
		rw.bordersBuilder(attr, nil, ListStart, s)
		rw.WalkOverElemAttrs("", s.Elem)
		rw.bordersBuilder(nil, nil, ListEnd, s)
	case schema.TypeMap:
		rw.borders = append(rw.borders, MapStart)
		// ancestors = append(ancestors, " ", "{\n")
		rw.bordersBuilder(attr, nil, MapStart, s)
		rw.WalkOverElemAttrs("", s.Elem)
		rw.bordersBuilder(nil, nil, MapEnd, s)
	}
}

func (rw *resourceWalker) WalkOverTopAttrs(attr any, s *schema.Schema) {
	rw.ladder += 1
	if s.Required {
		rw.handle(attr, s)
	} else if s.Computed && s.Optional { //rw.ladder <= 2 && (s.Type == schema.TypeList)
		// Optinal attributes path
		rw.handle(attr, s)
	}

	// rw.attrs[rw.attr] = append(rw.attrs[rw.attr], rw.keysTemplate)
	// return fmt.Sprint(rw.T)
}

func (rw *resourceWalker) WalkOverElemAttrs(k any, e any) {
	switch e.(type) {
	case *schema.Schema:
		rw.WalkOverTopAttrs(k, e.(*schema.Schema))
	case *schema.Resource:
		xe := e.(*schema.Resource)
		for k, sc := range xe.Schema {
			rw.WalkOverTopAttrs(k, sc)
		}
	}
}

func (rw *resourceWalker) WalkOverDataTopKeys(ancestors any, attr any, s any) {
	// ancestors = ancestors.(string) + attr.(string)
	// Handle primitive types.
	switch s.(type) {
	case string, bool, float64, int32:
		// key := fmt.Sprint(ancestors...)
		ancestors = fmt.Sprintf("%v = %#v", ancestors.(string), s)
		rw.attrs[rw.attr] = append(rw.attrs[rw.attr], ancestors)

	case []any:
		rw.WalkOverDataChildAttrs(ancestors, 0, s.([]any)[0])

	case map[string]any:
		rw.WalkOverDataChildAttrs(ancestors, 0, s.(map[string]any)[attr.(string)])

	default:
		rw.attrs[rw.attr] = append(rw.attrs[rw.attr], attr)
	}
}

func (rw *resourceWalker) WalkOverDataChildAttrs(a any, index any, e any) {
	a = a.(string) + ".0."
	switch e.(type) {
	case []any:
		for index, e := range e.([]any) {
			// a = a.(string) + ".0."
			rw.WalkOverDataTopKeys(a, index, e)
		}
	case map[string]any:
		xe := e.(map[string]any)
		for k, sc := range xe {
			// a = a.(string) + ".0."
			rw.WalkOverDataTopKeys(a, k, sc)
		}
	}
}

func PickResourceRequiredFieldsByTFRName(tfrName string) ([]string, error) {
	// tfrName := MatchTFRNameByCloudRSType()

	resource := ResourceMap[tfrName]
	var fields []string
	complexFields := map[string]any{}
	// time.Sleep(1 * time.Second)
	if resource == nil {
		return nil, errors.New("NIL")
	}

	for param_name, param := range resource.Schema {
		if param.Required {
			fields = append(fields, param_name)
			if param.Type == schema.TypeList {
				complexFields[param_name] = param.Elem
			}
		} else if strings.EqualFold(param_name, "parent") && param.Optional {
			fields = append(fields, param_name)
		} else if strings.EqualFold(param_name, "project") && param.Optional {
			fields = append(fields, param_name)
		}

	}

	if fields == nil {
		return nil, nil
	}
	// log.Printf("[TF_RESOURCE_REF]: %v ---- [REQUIRED_FIELDS] : %v", tfrName, fields)
	ResourceRequiredParams.LoadOrStore(tfrName, fields)
	return fields, nil
}

func (lts *LocalTFStateContent) HandleTFStateDeconstructs(rtype string, scope string, resources map[string][]*globals.Asset) {
	// for _, r := range resources{
	// a := cty.Object(r.Data.(any))
	// rtype := MatchTFRNameByCloudRSType(r[0].AssetType).(string)
	res := ResourceMap[rtype]

	rwalker := &resourceWalker{
		attrs: map[string][]any{},
		// data: r[0].Data,
		schemas: res.Schema,
	}

	// file, err := os.OpenFile(scope, os.O_APPEND|os.O_WRONLY, 0777)
	// defer file.Close()
	for k, sc := range rwalker.schemas {
		// log.Printf("%v---%T --- %v", k, sc, sc.Type)
		rwalker.r = *sc
		rwalker.attr = k
		rwalker.borders = []int{}
		rwalker.ladder = 0
		rwalker.keysTemplate = []any{}
		rwalker.attrsTemplate = strings.Builder{}
		// if k == "network_interface"{
		rwalker.WalkOverTopAttrs(k, sc)
		/* log.Println(rwalker.attrs[k])
		// }
		// block := fmt.Sprint(rwalker.attrsTemplate.String())
		// os.Stdout.Write([]byte(fmt.Sprint(block, "\n")))
		// allattrs[r.PickResourceId()] = append(allattrs[r.PickResourceId()], atts)
		// if  err == nil {
		// 	// data, _ := json.MarshalIndent(rwalker.attrs, "", "\t")
		// 	if _, err := file.Write([]byte(block)); err != nil {
		// 		log.Fatal("ERROR ON UPDATE TF FILE:", err)
		// 	}

		// } else {
		// 	log.Println(err)
		 } */
	}

	/* log.Println(allattrs)
	// for k, d := range  rwalker.data {
	// 	// log.Printf("%v---%T --- %v", k, sc, sc.Type)
	// 	// rwalker.r = *sc
	// 	rwalker.attr = k
	// 	// if k == "network_interface"{
	// 	rwalker.WalkOverDataTopKeys("",k,d)
	// 	log.Println(rwalker.attrs[k])
	// 	// }
	// 	// allattrs[r.PickResourceId()] = append(allattrs[r.PickResourceId()], atts)
	}*/

}

func DeconstructTFResourceAttributes(root string) {
	rootpath := func() string {
		if root == "" {
			return path.Join(globals.OrgRootPath, ".kitabus")
		}
		return root
	}
	if globals.ListOfAllTFResourcesPerPath == nil {
		file, _ := os.ReadFile(rootpath())
		var resources map[string][]*globals.Asset
		// var attributes []any
		// var attrType terraform.InstanceState
		if err := json.Unmarshal(file, &resources); err == nil {
		} else {
			log.Println("ERROR:", err)
		}
	} else {
		globals.MapOfAllTFResourcesWithData.Range(func(key, value any) bool {
			// statepath := func() string {
			// 	if strings.Contains(key.(string), ".ACTIVE"){
			// 		return strings.Split(key.(string), ".ACTIVE")[0]
			// 	} else {
			// 		return strings.Split(key.(string), ".META")[0]
			// 	}
			// }
			return true
		})
	}
}

func (lts *LocalTFStateContent) ReadTFStatesFile() {
	// lts.(LocalTFStateContent)["path"] = "/usr/local/google/home/magido/_reverse/cloudwork.joonix.net/Sandbox/cloudlabs/.ACTIVE/terraform.tfstate"
	statepath := lts.GetPath()
	// lts.SetTFStatesFilePath(statepath)
	var states LocalTFStateContent = *lts
	if file, err := os.ReadFile(statepath); err == nil {
		if err := json.Unmarshal(file, &states); err != nil {
			log.Println(err)
		} else {
			lts = &states
		}
	} else {
		log.Fatal(err)
	}
}

func (lts *LocalTFStateContent) SetTFStatesFilePath(p string) {
	var states LocalTFStateContent = *lts
	states["path"] = p
	lts = &states
}

func (lts *LocalTFStateContent) GetPath() string {
	var states LocalTFStateContent = *lts
	return states["path"].(string)
}

func (lts *LocalTFStateContent) GetTFResourceStateByType(rtype string) (any, string) {
	var states (LocalTFStateContent) = *lts
	var resources = (states["resources"]).([]any)
	var name string
	for _, res := range resources {
		if res.(map[string]any)["type"] == rtype {
			// res, _ := json.MarshalIndent(res, "","\b")
			// log.Println(string(res))
			r := res.(map[string]any)
			name = r["name"].(string)
			res = map[string]any{
				"name":      name,
				"type":      r["type"],
				"instances": r["instances"],
			}
			resources = []any{res}
			// return string(res)

			break
		}
	}
	if resources == nil {
		return nil, ""
	}
	states = LocalTFStateContent{
		"path":      lts.GetPath(),
		"resources": resources,
	}
	// states["resources"] = resources
	return states, name
}

func (lts *LocalTFStateContent) UpdateBaseValuesFile(rtype string) {
	mpath := lts.GetPath()
	states, f := lts.GetTFResourceStateByType(rtype)
	baseFilePath := func() string {
		if strings.Contains(mpath, ".ACTIVE") {
			return strings.Split(mpath, ".ACTIVE")[0] + "/.BASE" + "/." + f
		}
		if strings.Contains(mpath, ".META") {
			return strings.Split(mpath, ".META")[0] + "/.BASE" + "/." + f
		}
		return ""
	}
	if bf := baseFilePath(); bf == "" {
		log.Println("ERROR:", "Wrong file path")
	} else {
		state, _ := json.MarshalIndent(states.(LocalTFStateContent), "", "\t")
		if err := os.WriteFile(bf, state, 0444); err != nil {
			log.Fatal(err)
		}
		// log.Println(states)
	}

}

func (lts *LocalTFStateContent) GetTFManagedResourceName() string {
	var states LocalTFStateContent = *lts
	return states["name"].(string)
}
