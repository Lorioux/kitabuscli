package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleFolderOrganizationPolicy() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(ResourceGoogleFolderOrganizationPolicy().Schema)

	addRequiredFieldsToSchema(dsSchema, "folder")
	addRequiredFieldsToSchema(dsSchema, "constraint")

	return &schema.Resource{
		Read:   datasourceGoogleFolderOrganizationPolicyRead,
		Schema: dsSchema,
	}
}

func datasourceGoogleFolderOrganizationPolicyRead(d *schema.ResourceData, meta interface{}) error {

	d.SetId(fmt.Sprintf("%s/%s", d.Get("folder"), d.Get("constraint")))

	return resourceGoogleFolderOrganizationPolicyRead(d, meta)
}
