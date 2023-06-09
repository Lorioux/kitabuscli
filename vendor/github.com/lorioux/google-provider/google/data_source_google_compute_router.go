package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeRouter() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(ResourceComputeRouter().Schema)
	addRequiredFieldsToSchema(dsSchema, "name")
	addRequiredFieldsToSchema(dsSchema, "network")
	addOptionalFieldsToSchema(dsSchema, "region")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceComputeRouterRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeRouterRead(d *schema.ResourceData, meta interface{}) error {
	routerName := d.Get("name").(string)

	d.SetId(routerName)
	return resourceComputeRouterRead(d, meta)
}
