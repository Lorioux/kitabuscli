package playbooks


type MetaMapFactory interface {
	/*Check if terraform resource import meta mappings file's exists
	 For instance, a resource catalog file correspond to resource (e.g. compute.googlepais.com/[Instance|Network]) part: "Instance or Network"
	 Example 1:
		/Projects/
		|--CloudLabs/       # project_name
		|  |--meta_mappings # TF importing meta mappings file 
	*/
	CheckTFImportMetaFileExistsOrCreate () bool
}


type MetaMapType struct {
	import_signature map[string]string // e.g. module.compute.google_compute_instance.worker : "projects/[project_id]/zones/europe-central2-a/instances/worker"
}


