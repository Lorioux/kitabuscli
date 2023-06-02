// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const workloadIdentityPoolProviderIdRegexp = `^[0-9a-z-]+$`

func validateWorkloadIdentityPoolProviderId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if strings.HasPrefix(value, "gcp-") {
		errors = append(errors, fmt.Errorf(
			"%q (%q) can not start with \"gcp-\"", k, value))
	}

	if !regexp.MustCompile(workloadIdentityPoolProviderIdRegexp).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must contain only lowercase letters (a-z), numbers (0-9), or dashes (-)", k))
	}

	if len(value) < 4 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be smaller than 4 characters", k))
	}

	if len(value) > 32 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be greater than 32 characters", k))
	}

	return
}

func ResourceIAMBetaWorkloadIdentityPoolProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceIAMBetaWorkloadIdentityPoolProviderCreate,
		Read:   resourceIAMBetaWorkloadIdentityPoolProviderRead,
		Update: resourceIAMBetaWorkloadIdentityPoolProviderUpdate,
		Delete: resourceIAMBetaWorkloadIdentityPoolProviderDelete,

		Importer: &schema.ResourceImporter{
			State: resourceIAMBetaWorkloadIdentityPoolProviderImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"workload_identity_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The ID used for the pool, which is the final component of the pool resource name. This
value should be 4-32 characters, and may contain the characters [a-z0-9-]. The prefix
'gcp-' is reserved for use by Google, and may not be specified.`,
			},
			"workload_identity_pool_provider_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateWorkloadIdentityPoolProviderId,
				Description: `The ID for the provider, which becomes the final component of the resource name. This
value must be 4-32 characters, and may contain the characters [a-z0-9-]. The prefix
'gcp-' is reserved for use by Google, and may not be specified.`,
			},
			"attribute_condition": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `[A Common Expression Language](https://opensource.google/projects/cel) expression, in
plain text, to restrict what otherwise valid authentication credentials issued by the
provider should not be accepted.

The expression must output a boolean representing whether to allow the federation.

The following keywords may be referenced in the expressions:
  * 'assertion': JSON representing the authentication credential issued by the provider.
  * 'google': The Google attributes mapped from the assertion in the 'attribute_mappings'.
  * 'attribute': The custom attributes mapped from the assertion in the 'attribute_mappings'.

The maximum length of the attribute condition expression is 4096 characters. If
unspecified, all valid authentication credential are accepted.

The following example shows how to only allow credentials with a mapped 'google.groups'
value of 'admins':
'''
"'admins' in google.groups"
'''`,
			},
			"attribute_mapping": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `Maps attributes from authentication credentials issued by an external identity provider
to Google Cloud attributes, such as 'subject' and 'segment'.

Each key must be a string specifying the Google Cloud IAM attribute to map to.

The following keys are supported:
  * 'google.subject': The principal IAM is authenticating. You can reference this value
    in IAM bindings. This is also the subject that appears in Cloud Logging logs.
    Cannot exceed 127 characters.
  * 'google.groups': Groups the external identity belongs to. You can grant groups
    access to resources using an IAM 'principalSet' binding; access applies to all
    members of the group.

You can also provide custom attributes by specifying 'attribute.{custom_attribute}',
where '{custom_attribute}' is the name of the custom attribute to be mapped. You can
define a maximum of 50 custom attributes. The maximum length of a mapped attribute key
is 100 characters, and the key may only contain the characters [a-z0-9_].

You can reference these attributes in IAM policies to define fine-grained access for a
workload to Google Cloud resources. For example:
  * 'google.subject':
    'principal://iam.googleapis.com/projects/{project}/locations/{location}/workloadIdentityPools/{pool}/subject/{value}'
  * 'google.groups':
    'principalSet://iam.googleapis.com/projects/{project}/locations/{location}/workloadIdentityPools/{pool}/group/{value}'
  * 'attribute.{custom_attribute}':
    'principalSet://iam.googleapis.com/projects/{project}/locations/{location}/workloadIdentityPools/{pool}/attribute.{custom_attribute}/{value}'

Each value must be a [Common Expression Language](https://opensource.google/projects/cel)
function that maps an identity provider credential to the normalized attribute specified
by the corresponding map key.

You can use the 'assertion' keyword in the expression to access a JSON representation of
the authentication credential issued by the provider.

The maximum length of an attribute mapping expression is 2048 characters. When evaluated,
the total size of all mapped attributes must not exceed 8KB.

For AWS providers, the following rules apply:
  - If no attribute mapping is defined, the following default mapping applies:
    '''
    {
      "google.subject":"assertion.arn",
      "attribute.aws_role":
        "assertion.arn.contains('assumed-role')"
        " ? assertion.arn.extract('{account_arn}assumed-role/')"
        "   + 'assumed-role/'"
        "   + assertion.arn.extract('assumed-role/{role_name}/')"
        " : assertion.arn",
    }
    '''
  - If any custom attribute mappings are defined, they must include a mapping to the
    'google.subject' attribute.

For OIDC providers, the following rules apply:
  - Custom attribute mappings must be defined, and must include a mapping to the
    'google.subject' attribute. For example, the following maps the 'sub' claim of the
    incoming credential to the 'subject' attribute on a Google token.
    '''
    {"google.subject": "assertion.sub"}
    '''`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"aws": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `An Amazon Web Services identity provider. Not compatible with the property oidc.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The AWS account ID.`,
						},
					},
				},
				ExactlyOneOf: []string{"aws", "oidc"},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A description for the provider. Cannot exceed 256 characters.`,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `Whether the provider is disabled. You cannot use a disabled provider to exchange tokens.
However, existing tokens still grant access.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A display name for the provider. Cannot exceed 32 characters.`,
			},
			"oidc": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `An OpenId Connect 1.0 identity provider. Not compatible with the property aws.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issuer_uri": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The OIDC issuer URL.`,
						},
						"allowed_audiences": {
							Type:     schema.TypeList,
							Optional: true,
							Description: `Acceptable values for the 'aud' field (audience) in the OIDC token. Token exchange
requests are rejected if the token audience does not match one of the configured
values. Each audience may be at most 256 characters. A maximum of 10 audiences may
be configured.

If this list is empty, the OIDC token audience must be equal to the full canonical
resource name of the WorkloadIdentityPoolProvider, with or without the HTTPS prefix.
For example:
'''
//iam.googleapis.com/projects/<project-number>/locations/<location>/workloadIdentityPools/<pool-id>/providers/<provider-id>
https://iam.googleapis.com/projects/<project-number>/locations/<location>/workloadIdentityPools/<pool-id>/providers/<provider-id>
'''`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				ExactlyOneOf: []string{"aws", "oidc"},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The resource name of the provider as
'projects/{project_number}/locations/global/workloadIdentityPools/{workload_identity_pool_id}/providers/{workload_identity_pool_provider_id}'.`,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The state of the provider.
* STATE_UNSPECIFIED: State unspecified.
* ACTIVE: The provider is active, and may be used to validate authentication credentials.
* DELETED: The provider is soft-deleted. Soft-deleted providers are permanently deleted
  after approximately 30 days. You can restore a soft-deleted provider using
  UndeleteWorkloadIdentityPoolProvider. You cannot reuse the ID of a soft-deleted provider
  until it is permanently deleted.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceIAMBetaWorkloadIdentityPoolProviderCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	displayNameProp, err := expandIAMBetaWorkloadIdentityPoolProviderDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandIAMBetaWorkloadIdentityPoolProviderDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	disabledProp, err := expandIAMBetaWorkloadIdentityPoolProviderDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disabled"); !isEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}
	attributeMappingProp, err := expandIAMBetaWorkloadIdentityPoolProviderAttributeMapping(d.Get("attribute_mapping"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("attribute_mapping"); !isEmptyValue(reflect.ValueOf(attributeMappingProp)) && (ok || !reflect.DeepEqual(v, attributeMappingProp)) {
		obj["attributeMapping"] = attributeMappingProp
	}
	attributeConditionProp, err := expandIAMBetaWorkloadIdentityPoolProviderAttributeCondition(d.Get("attribute_condition"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("attribute_condition"); !isEmptyValue(reflect.ValueOf(attributeConditionProp)) && (ok || !reflect.DeepEqual(v, attributeConditionProp)) {
		obj["attributeCondition"] = attributeConditionProp
	}
	awsProp, err := expandIAMBetaWorkloadIdentityPoolProviderAws(d.Get("aws"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("aws"); !isEmptyValue(reflect.ValueOf(awsProp)) && (ok || !reflect.DeepEqual(v, awsProp)) {
		obj["aws"] = awsProp
	}
	oidcProp, err := expandIAMBetaWorkloadIdentityPoolProviderOidc(d.Get("oidc"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("oidc"); !isEmptyValue(reflect.ValueOf(oidcProp)) && (ok || !reflect.DeepEqual(v, oidcProp)) {
		obj["oidc"] = oidcProp
	}

	url, err := replaceVars(d, config, "{{IAMBetaBasePath}}projects/{{project}}/locations/global/workloadIdentityPools/{{workload_identity_pool_id}}/providers?workloadIdentityPoolProviderId={{workload_identity_pool_provider_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new WorkloadIdentityPoolProvider: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for WorkloadIdentityPoolProvider: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating WorkloadIdentityPoolProvider: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/locations/global/workloadIdentityPools/{{workload_identity_pool_id}}/providers/{{workload_identity_pool_provider_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = IAMBetaOperationWaitTime(
		config, res, project, "Creating WorkloadIdentityPoolProvider", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create WorkloadIdentityPoolProvider: %s", err)
	}

	log.Printf("[DEBUG] Finished creating WorkloadIdentityPoolProvider %q: %#v", d.Id(), res)

	return resourceIAMBetaWorkloadIdentityPoolProviderRead(d, meta)
}

func resourceIAMBetaWorkloadIdentityPoolProviderRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{IAMBetaBasePath}}projects/{{project}}/locations/global/workloadIdentityPools/{{workload_identity_pool_id}}/providers/{{workload_identity_pool_provider_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for WorkloadIdentityPoolProvider: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("IAMBetaWorkloadIdentityPoolProvider %q", d.Id()))
	}

	res, err = resourceIAMBetaWorkloadIdentityPoolProviderDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing IAMBetaWorkloadIdentityPoolProvider because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading WorkloadIdentityPoolProvider: %s", err)
	}

	if err := d.Set("state", flattenIAMBetaWorkloadIdentityPoolProviderState(res["state"], d, config)); err != nil {
		return fmt.Errorf("Error reading WorkloadIdentityPoolProvider: %s", err)
	}
	if err := d.Set("display_name", flattenIAMBetaWorkloadIdentityPoolProviderDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading WorkloadIdentityPoolProvider: %s", err)
	}
	if err := d.Set("description", flattenIAMBetaWorkloadIdentityPoolProviderDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading WorkloadIdentityPoolProvider: %s", err)
	}
	if err := d.Set("name", flattenIAMBetaWorkloadIdentityPoolProviderName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading WorkloadIdentityPoolProvider: %s", err)
	}
	if err := d.Set("disabled", flattenIAMBetaWorkloadIdentityPoolProviderDisabled(res["disabled"], d, config)); err != nil {
		return fmt.Errorf("Error reading WorkloadIdentityPoolProvider: %s", err)
	}
	if err := d.Set("attribute_mapping", flattenIAMBetaWorkloadIdentityPoolProviderAttributeMapping(res["attributeMapping"], d, config)); err != nil {
		return fmt.Errorf("Error reading WorkloadIdentityPoolProvider: %s", err)
	}
	if err := d.Set("attribute_condition", flattenIAMBetaWorkloadIdentityPoolProviderAttributeCondition(res["attributeCondition"], d, config)); err != nil {
		return fmt.Errorf("Error reading WorkloadIdentityPoolProvider: %s", err)
	}
	if err := d.Set("aws", flattenIAMBetaWorkloadIdentityPoolProviderAws(res["aws"], d, config)); err != nil {
		return fmt.Errorf("Error reading WorkloadIdentityPoolProvider: %s", err)
	}
	if err := d.Set("oidc", flattenIAMBetaWorkloadIdentityPoolProviderOidc(res["oidc"], d, config)); err != nil {
		return fmt.Errorf("Error reading WorkloadIdentityPoolProvider: %s", err)
	}

	return nil
}

func resourceIAMBetaWorkloadIdentityPoolProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for WorkloadIdentityPoolProvider: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	displayNameProp, err := expandIAMBetaWorkloadIdentityPoolProviderDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandIAMBetaWorkloadIdentityPoolProviderDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	disabledProp, err := expandIAMBetaWorkloadIdentityPoolProviderDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disabled"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}
	attributeMappingProp, err := expandIAMBetaWorkloadIdentityPoolProviderAttributeMapping(d.Get("attribute_mapping"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("attribute_mapping"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, attributeMappingProp)) {
		obj["attributeMapping"] = attributeMappingProp
	}
	attributeConditionProp, err := expandIAMBetaWorkloadIdentityPoolProviderAttributeCondition(d.Get("attribute_condition"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("attribute_condition"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, attributeConditionProp)) {
		obj["attributeCondition"] = attributeConditionProp
	}
	awsProp, err := expandIAMBetaWorkloadIdentityPoolProviderAws(d.Get("aws"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("aws"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, awsProp)) {
		obj["aws"] = awsProp
	}
	oidcProp, err := expandIAMBetaWorkloadIdentityPoolProviderOidc(d.Get("oidc"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("oidc"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, oidcProp)) {
		obj["oidc"] = oidcProp
	}

	url, err := replaceVars(d, config, "{{IAMBetaBasePath}}projects/{{project}}/locations/global/workloadIdentityPools/{{workload_identity_pool_id}}/providers/{{workload_identity_pool_provider_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating WorkloadIdentityPoolProvider %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("display_name") {
		updateMask = append(updateMask, "displayName")
	}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("disabled") {
		updateMask = append(updateMask, "disabled")
	}

	if d.HasChange("attribute_mapping") {
		updateMask = append(updateMask, "attributeMapping")
	}

	if d.HasChange("attribute_condition") {
		updateMask = append(updateMask, "attributeCondition")
	}

	if d.HasChange("aws") {
		updateMask = append(updateMask, "aws")
	}

	if d.HasChange("oidc") {
		updateMask = append(updateMask, "oidc.allowed_audiences",
			"oidc.issuer_uri")
	}
	// updateMask is a URL parameter but not present in the schema, so replaceVars
	// won't set it
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "PATCH", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating WorkloadIdentityPoolProvider %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating WorkloadIdentityPoolProvider %q: %#v", d.Id(), res)
	}

	err = IAMBetaOperationWaitTime(
		config, res, project, "Updating WorkloadIdentityPoolProvider", userAgent,
		d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return err
	}

	return resourceIAMBetaWorkloadIdentityPoolProviderRead(d, meta)
}

func resourceIAMBetaWorkloadIdentityPoolProviderDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for WorkloadIdentityPoolProvider: %s", err)
	}
	billingProject = project

	url, err := replaceVars(d, config, "{{IAMBetaBasePath}}projects/{{project}}/locations/global/workloadIdentityPools/{{workload_identity_pool_id}}/providers/{{workload_identity_pool_provider_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting WorkloadIdentityPoolProvider %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "WorkloadIdentityPoolProvider")
	}

	err = IAMBetaOperationWaitTime(
		config, res, project, "Deleting WorkloadIdentityPoolProvider", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting WorkloadIdentityPoolProvider %q: %#v", d.Id(), res)
	return nil
}

func resourceIAMBetaWorkloadIdentityPoolProviderImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/global/workloadIdentityPools/(?P<workload_identity_pool_id>[^/]+)/providers/(?P<workload_identity_pool_provider_id>[^/]+)",
		"(?P<project>[^/]+)/(?P<workload_identity_pool_id>[^/]+)/(?P<workload_identity_pool_provider_id>[^/]+)",
		"(?P<workload_identity_pool_id>[^/]+)/(?P<workload_identity_pool_provider_id>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/locations/global/workloadIdentityPools/{{workload_identity_pool_id}}/providers/{{workload_identity_pool_provider_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenIAMBetaWorkloadIdentityPoolProviderState(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIAMBetaWorkloadIdentityPoolProviderDisplayName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIAMBetaWorkloadIdentityPoolProviderDescription(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIAMBetaWorkloadIdentityPoolProviderName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIAMBetaWorkloadIdentityPoolProviderDisabled(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIAMBetaWorkloadIdentityPoolProviderAttributeMapping(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIAMBetaWorkloadIdentityPoolProviderAttributeCondition(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIAMBetaWorkloadIdentityPoolProviderAws(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["account_id"] =
		flattenIAMBetaWorkloadIdentityPoolProviderAwsAccountId(original["accountId"], d, config)
	return []interface{}{transformed}
}
func flattenIAMBetaWorkloadIdentityPoolProviderAwsAccountId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIAMBetaWorkloadIdentityPoolProviderOidc(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["allowed_audiences"] =
		flattenIAMBetaWorkloadIdentityPoolProviderOidcAllowedAudiences(original["allowedAudiences"], d, config)
	transformed["issuer_uri"] =
		flattenIAMBetaWorkloadIdentityPoolProviderOidcIssuerUri(original["issuerUri"], d, config)
	return []interface{}{transformed}
}
func flattenIAMBetaWorkloadIdentityPoolProviderOidcAllowedAudiences(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIAMBetaWorkloadIdentityPoolProviderOidcIssuerUri(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandIAMBetaWorkloadIdentityPoolProviderDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderDisabled(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderAttributeMapping(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderAttributeCondition(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderAws(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAccountId, err := expandIAMBetaWorkloadIdentityPoolProviderAwsAccountId(original["account_id"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAccountId); val.IsValid() && !isEmptyValue(val) {
		transformed["accountId"] = transformedAccountId
	}

	return transformed, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderAwsAccountId(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderOidc(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAllowedAudiences, err := expandIAMBetaWorkloadIdentityPoolProviderOidcAllowedAudiences(original["allowed_audiences"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowedAudiences); val.IsValid() && !isEmptyValue(val) {
		transformed["allowedAudiences"] = transformedAllowedAudiences
	}

	transformedIssuerUri, err := expandIAMBetaWorkloadIdentityPoolProviderOidcIssuerUri(original["issuer_uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIssuerUri); val.IsValid() && !isEmptyValue(val) {
		transformed["issuerUri"] = transformedIssuerUri
	}

	return transformed, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderOidcAllowedAudiences(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIAMBetaWorkloadIdentityPoolProviderOidcIssuerUri(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func resourceIAMBetaWorkloadIdentityPoolProviderDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	if v := res["state"]; v == "DELETED" {
		return nil, nil
	}

	return res, nil
}