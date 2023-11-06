/*
*TBD*

# Example Usage

```hcl

	data "awx_organization" "default" {
	  name = "Default"
	}

	resource "awx_team" "admin_team" {
	  name = "Admins"
	  organization_id = data.awx_organization.default.id
	}

	resource "awx_settings_saml_team_map" "admin_team_map" {
	  name         = resource.awx_team.admin_team.name
	  users        = ["CN=MyTeam,OU=Groups,DC=example,DC=com"]
	  organization = data.awx_organization.default.name
	  remove       = true
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

var samlTeamMapAccessMutex sync.Mutex

func resourceSettingsSAMLTeamMap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSettingsSAMLTeamMapCreate,
		ReadContext:   resourceSettingsSAMLTeamMapRead,
		DeleteContext: resourceSettingsSAMLTeamMapDelete,
		UpdateContext: resourceSettingsSAMLTeamMapUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this Team",
			},
			"users": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "SAML Group to map to this team",
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the team organization",
			},
			"remove": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When True, a user who is not a member of the given saml groups will be removed from the team",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

type samlTeamMapEntry struct {
	SamlGroups   interface{} `json:"users"`
	Organization string      `json:"organization"`
	Remove       bool        `json:"remove"`
}

type samlTeamMap map[string]samlTeamMapEntry

func resourceSettingsSAMLTeamMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	samlTeamMapAccessMutex.Lock()
	defer samlTeamMapAccessMutex.Unlock()

	client := m.(*awx.AWX)
	awxService := client.SettingService

	res, err := awxService.GetSettingsBySlug("saml", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: failed to fetch settings",
			"Failed to fetch any saml setting, got: %s", err.Error(),
		)
	}

	tmaps := make(samlTeamMap)
	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_MAP"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: failed to parse SOCIAL_AUTH_SAML_TEAM_MAP setting",
			"Failed to parse SOCIAL_AUTH_SAML_TEAM_MAP setting, got: %s with input %s", err.Error(), (*res)["SOCIAL_AUTH_SAML_TEAM_MAP"],
		)
	}

	name := d.Get("name").(string)

	_, ok := tmaps[name]
	if ok {
		return buildDiagnosticsMessage(
			"Create: team map already exists",
			"Map for saml to team map %v already exists", d.Id(),
		)
	}

	newtmap := samlTeamMapEntry{
		SamlGroups:   d.Get("users").([]interface{}),
		Organization: d.Get("organization").(string),
		Remove:       d.Get("remove").(bool),
	}

	tmaps[name] = newtmap

	payload := map[string]interface{}{
		"SOCIAL_AUTH_SAML_TEAM_MAP": tmaps,
	}

	_, err = awxService.UpdateSettings("saml", payload, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: team map not created",
			"failed to save team map data, got: %s", err.Error(),
		)
	}

	d.SetId(name)
	return resourceSettingsSAMLTeamMapRead(ctx, d, m)
}

func resourceSettingsSAMLTeamMapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	samlTeamMapAccessMutex.Lock()
	defer samlTeamMapAccessMutex.Unlock()

	client := m.(*awx.AWX)
	awxService := client.SettingService

	res, err := awxService.GetSettingsBySlug("saml", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Update: Unable to fetch settings",
			"Unable to load settings with slug saml: got %s", err.Error(),
		)
	}

	tmaps := make(samlTeamMap)
	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_MAP"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Update: failed to parse AUTH_SAML_TEAM_MAP setting",
			"Failed to parse AUTH_SAML_TEAM_MAP setting, got: %s", err.Error(),
		)
	}

	id := d.Id()
	name := d.Get("name").(string)
	organization := d.Get("organization").(string)
	users := d.Get("users").([]interface{})
	remove := d.Get("remove").(bool)

	if name != id {
		tmaps[name] = tmaps[id]
		delete(tmaps, id)
	}

	utmap := tmaps[name]
	utmap.SamlGroups = users
	utmap.Organization = organization
	utmap.Remove = remove
	tmaps[name] = utmap

	payload := map[string]interface{}{
		"SOCIAL_AUTH_SAML_TEAM_MAP": tmaps,
	}

	_, err = awxService.UpdateSettings("saml", payload, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Update: team map not created",
			"failed to save team map data, got: %s", err.Error(),
		)
	}

	d.SetId(name)
	return resourceSettingsSAMLTeamMapRead(ctx, d, m)
}

func resourceSettingsSAMLTeamMapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	tmaps := make(samlTeamMap)
	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_MAP"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Unable to parse SOCIAL_AUTH_SAML_TEAM_MAP",
			"Unable to parse SOCIAL_AUTH_SAML_TEAM_MAP, got: %s", err.Error(),
		)
	}
	mapdef, ok := tmaps[d.Id()]
	if !ok {
		return buildDiagnosticsMessage(
			"Unable to fetch saml team map",
			"Unable to load saml team map %v: not found", d.Id(),
		)
	}

	var users []string
	switch tt := mapdef.SamlGroups.(type) {
	case string:
		users = []string{tt}
	case []string:
		users = tt
	case []interface{}:
		for _, v := range tt {
			if dn, ok := v.(string); ok {
				users = append(users, dn)
			}
		}
	}

	d.Set("name", d.Id())
	d.Set("users", users)
	d.Set("organization", mapdef.Organization)
	d.Set("remove", mapdef.Remove)
	return diags
}

func resourceSettingsSAMLTeamMapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	samlTeamMapAccessMutex.Lock()
	defer samlTeamMapAccessMutex.Unlock()

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

	tmaps := make(samlTeamMap)
	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_MAP"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Delete: failed to parse SOCIAL_AUTH_SAML_TEAM_MAP setting",
			"Failed to parse SOCIAL_AUTH_SAML_TEAM_MAP setting, got: %s", err.Error(),
		)
	}

	id := d.Id()
	delete(tmaps, id)

	payload := map[string]interface{}{
		"SOCIAL_AUTH_SAML_TEAM_MAP": tmaps,
	}

	_, err = awxService.UpdateSettings("saml", payload, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Delete: team map not created",
			"failed to save team map data, got: %s", err.Error(),
		)
	}
	d.SetId("")
	return diags
}
