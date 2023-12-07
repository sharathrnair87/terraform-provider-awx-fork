/*
Use this resource to map a SAML derived Group to an AWX/AT Organization.
NOTE: This resource manages the global SAML group to AWX/AT Organization mapping i.e. SOCIAL_AUTH_SAML_ORGANIZATION_MAP,
ensure that this resource is managed by a single terraform state to avoid unexpected behaviour

# Example Usage

```hcl

	data "awx_organization" "default" {
	  name = "Default"
	}

	resource "awx_organization" "admin_organization" {
	  name            = "Admins"
	  organization_id = data.awx_organization.default.id
	}

	resource "awx_settings_saml_organization_map" "admin_organization_map" {
	  name          = resource.awx_organization.admin_organization.name
	  users         = ["myorg-infra-users"]
	  admins        = ["myorg-infra-admins"]
	  organization  = data.awx_organization.default.name
	  remove_users  = true
	  remove_admins = true
	}

```
*/
package awx

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

var samlOrganizationMapAccessMutex sync.Mutex

func resourceSettingsSAMLOrganizationMap() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to map a SAML derived Group to an AWX/AT Organization.",
		CreateContext: resourceSettingsSAMLOrganizationMapCreate,
		ReadContext:   resourceSettingsSAMLOrganizationMapRead,
		DeleteContext: resourceSettingsSAMLOrganizationMapDelete,
		UpdateContext: resourceSettingsSAMLOrganizationMapUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this Organization",
			},
			"users": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "SAML Group to map users to this organization",
			},
			"admins": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "SAML Group to map admins to this organization",
			},
			"remove_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When True, a user who is not a member of the given saml groups will be removed from the organization",
			},
			"remove_admins": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When True, an admin who is not a member of the given saml groups will be removed from the organization",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

type samlOrganizationMapEntry struct {
	SamlUserGroups  interface{} `json:"users"`
	SamlAdminGroups interface{} `json:"admins"`
	RemoveUsers     bool        `json:"remove_users"`
	RemoveAdmins    bool        `json:"remove_admins"`
}

type samlOrganizationMap map[string]samlOrganizationMapEntry

func resourceSettingsSAMLOrganizationMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	samlOrganizationMapAccessMutex.Lock()
	defer samlOrganizationMapAccessMutex.Unlock()

	client := m.(*awx.AWX)
	awxService := client.SettingService

	res, err := awxService.GetSettingsBySlug("saml", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: failed to fetch settings",
			"Failed to fetch any saml setting, got: %s", err.Error(),
		)
	}

	tmaps := make(samlOrganizationMap)
	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_ORGANIZATION_MAP"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: failed to parse SOCIAL_AUTH_SAML_ORGANIZATION_MAP setting",
			"Failed to parse SOCIAL_AUTH_SAML_ORGANIZATION_MAP setting, got: %s with input %s", err.Error(), (*res)["SOCIAL_AUTH_SAML_ORGANIZATION_MAP"],
		)
	}

	name := d.Get("name").(string)

	_, ok := tmaps[name]
	if ok {
		return buildDiagnosticsMessage(
			"Create: organization map already exists",
			"Map for saml to organization map %v already exists", d.Id(),
		)
	}

	newtmap := samlOrganizationMapEntry{
		SamlUserGroups:  d.Get("users").([]interface{}),
		SamlAdminGroups: d.Get("admins").([]interface{}),
		RemoveUsers:     d.Get("remove_users").(bool),
		RemoveAdmins:    d.Get("remove_admins").(bool),
	}

	tmaps[name] = newtmap

	payload := map[string]interface{}{
		"SOCIAL_AUTH_SAML_ORGANIZATION_MAP": tmaps,
	}

	_, err = awxService.UpdateSettings("saml", payload, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: organization map not created",
			"failed to save organization map data, got: %s", err.Error(),
		)
	}

	d.SetId(name)
	return resourceSettingsSAMLOrganizationMapRead(ctx, d, m)
}

func resourceSettingsSAMLOrganizationMapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	samlOrganizationMapAccessMutex.Lock()
	defer samlOrganizationMapAccessMutex.Unlock()

	client := m.(*awx.AWX)
	awxService := client.SettingService

	res, err := awxService.GetSettingsBySlug("saml", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Update: Unable to fetch settings",
			"Unable to load settings with slug saml: got %s", err.Error(),
		)
	}

	tmaps := make(samlOrganizationMap)
	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_ORGANIZATION_MAP"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Update: failed to parse AUTH_SAML_ORGANIZATION_MAP setting",
			"Failed to parse AUTH_SAML_ORGANIZATION_MAP setting, got: %s", err.Error(),
		)
	}

	id := d.Id()
	name := d.Get("name").(string)
	users := d.Get("users").([]interface{})
	admins := d.Get("admins").([]interface{})
	remove_users := d.Get("remove_users").(bool)
	remove_admins := d.Get("remove_admins").(bool)

	if name != id {
		tmaps[name] = tmaps[id]
		delete(tmaps, id)
	}

	utmap := tmaps[name]
	utmap.SamlUserGroups = users
	utmap.SamlAdminGroups = admins
	utmap.RemoveUsers = remove_users
	utmap.RemoveAdmins = remove_admins
	tmaps[name] = utmap

	payload := map[string]interface{}{
		"SOCIAL_AUTH_SAML_ORGANIZATION_MAP": tmaps,
	}

	_, err = awxService.UpdateSettings("saml", payload, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Update: organization map not created",
			"failed to save organization map data, got: %s", err.Error(),
		)
	}

	d.SetId(name)
	return resourceSettingsSAMLOrganizationMapRead(ctx, d, m)
}

func resourceSettingsSAMLOrganizationMapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.SettingService

	res, err := awxService.GetSettingsBySlug("saml", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Unable to fetch settings",
			"Unable to load settings with slug saml: got %s",
			err.Error(),
		)
	}
	tmaps := make(samlOrganizationMap)
	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_ORGANIZATION_MAP"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Unable to parse SOCIAL_AUTH_SAML_ORGANIZATION_MAP",
			"Unable to parse SOCIAL_AUTH_SAML_ORGANIZATION_MAP, got: %s", err.Error(),
		)
	}
	mapdef, ok := tmaps[d.Id()]
	if !ok {
		return buildDiagnosticsMessage(
			"Unable to fetch saml organization map",
			"Unable to load saml organization map %v: not found", d.Id(),
		)
	}

	var users []string
	switch tt := mapdef.SamlUserGroups.(type) {
	case string:
		users = []string{tt}
	case []string:
		users = tt
	case []interface{}:
		for _, v := range tt {
			if t, ok := v.(string); ok {
				users = append(users, t)
			}
		}
	}

	var admins []string
	switch tt := mapdef.SamlAdminGroups.(type) {
	case string:
		admins = []string{tt}
	case []string:
		admins = tt
	case []interface{}:
		for _, v := range tt {
			if t, ok := v.(string); ok {
				admins = append(admins, t)
			}
		}
	}

	d.Set("name", d.Id())
	d.Set("users", users)
	d.Set("admins", admins)
	d.Set("remove_users", mapdef.RemoveUsers)
	d.Set("remove_admins", mapdef.RemoveAdmins)
	return diags
}

func resourceSettingsSAMLOrganizationMapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	samlOrganizationMapAccessMutex.Lock()
	defer samlOrganizationMapAccessMutex.Unlock()

	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.SettingService

	res, err := awxService.GetSettingsBySlug("saml", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Delete: Unable to fetch settings",
			"Unable to load settings with slug saml: got %s", err.Error(),
		)
	}

	tmaps := make(samlOrganizationMap)
	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_ORGANIZATION_MAP"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Delete: failed to parse SOCIAL_AUTH_SAML_ORGANIZATION_MAP setting",
			"Failed to parse SOCIAL_AUTH_SAML_ORGANIZATION_MAP setting, got: %s", err.Error(),
		)
	}

	id := d.Id()
	delete(tmaps, id)

	payload := map[string]interface{}{
		"SOCIAL_AUTH_SAML_ORGANIZATION_MAP": tmaps,
	}

	_, err = awxService.UpdateSettings("saml", payload, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Delete: organization map not created",
			"failed to save organization map data, got: %s", err.Error(),
		)
	}
	d.SetId("")
	return diags
}
