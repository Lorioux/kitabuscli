package mappings

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/lorioux/kitabus/globals"
)


type ReverseMapType map[string]string

type Reversor interface{
	ResourceReversorInstance() 
}

type ResourceReversor struct{
	resourceMapReverse           map[string]any
	UnsupportedResourceMap       map[string]any
}


var (
	BASE_FILE                    = path.Join(CloudRTypeToTFRNameMappingFile)
	err                    error = nil
)


func ResourceReversorInstance() any{
	rr := new(ResourceReversor)
	rr.resourceMapReverse, err = rr.readMappingsFile()
	// globals.Reversor = rr
	return rr
}

// Pick a json file then reverse the key and value
// Save the result as json file
func (rr *ResourceReversor) DoReverse(dsfile ...string) {

	// if err != nil {
	// 	panic(err)
	// }

	// Create a new map to store the swapped keys and values.
	newMap := make(map[string]string)

	// Iterate over the map and swap the keys and values.
	for _, _map := range rr.resourceMapReverse["maps"].([]ReverseMapType)[0:] {

		for key, value := range _map {
			newMap[value] = key
		}
		// fmt.Println(key, value)
	}

	// Print the swapped map.
	// fmt.Println(newMap)
	rr.resourceMapReverse["reverse"] = newMap
	data, err := json.MarshalIndent(rr.resourceMapReverse, "", "\t")
	if err != nil {
		fmt.Print(err)
	}
	// dsfile = scfile
	if err := os.WriteFile(BASE_FILE, data, 0777); err != nil {
		fmt.Printf("There was an error: %v", err)
	}
	// fmt.Println("Done!")
}


func (rr *ResourceReversor) readMappingsFile() (map[string]interface{}, error) {
	// Read the JSON file.
	file, err := globals.MetaTFMappingsFile.ReadFile(BASE_FILE)
	
	if err != nil {
		log.Print(err)
		// return nil, err
		file, _ = json.Marshal(DATA)
	}
	var mapObj map[string]interface{}
	if err = json.Unmarshal(file, &mapObj); err != nil {
		return nil, err
	}
	rr.UnsupportedResourceMap = mapObj["unsupported"].(map[string]interface{})

	return mapObj, nil
}

func (rr *ResourceReversor) MatchTFRNameByCloudRSType(cloudRSType string) any {

	v := rr.resourceMapReverse["reverse"].(map[string]interface{})[cloudRSType]

	switch v.(type) {
	case nil:
		{
			if rr.UnsupportedResourceMap[cloudRSType] == nil {
				// log.Printf("REQUIRED UNSUPPORTED TF RESOURCE TYPE:-----> %v",cloudRSType )
				rr.UnsupportedResourceMap[cloudRSType] = "YES"
			}
			return ""
		}
	default:
		return v
	}
}

func (rr *ResourceReversor) UpdateUnsupportedResourceMap() {
	rr.resourceMapReverse["unsupported"] = rr.UnsupportedResourceMap
	if data, err := json.MarshalIndent(rr.resourceMapReverse, "", "\t"); err == nil {
		os.WriteFile(BASE_FILE, data, 0777)
	} else {
		log.Fatalf("FAILED TO UPDATE UNSUPPORTED MAP: %v", err)
	}
}
