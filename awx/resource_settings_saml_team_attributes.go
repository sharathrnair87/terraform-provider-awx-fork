/*
Use this resource to globally set the SOCIAL_AUTH_SAML_TEAM_ATTR setting in the SAML config.
NOTE: This resource controls this setting globally across your entire AWX/AT setup, and must be
managed centrally from a single state file to avoid conflicts

# Example Usage

```hcl

	data "awx_organization" "default" {
	  name = "Default"
	}

	resource "awx_settings_saml_team_attributes" "global" {
	  saml_attr = "groups"
	  remove    = true
	  team_org_map {
	    team         = var.saml_provider_team_id // The team ID as it is displayed in your SAML Auth Provider
	    organization = data.awx_organization.default.name
	    team_alias   = "Admin"
	  }
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

var samlTeamAttrMapAccessMutex sync.Mutex

func resourceSettingsSAMLTeamAttrMap() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to globally set the SOCIAL_AUTH_SAML_TEAM_ATTR setting in the SAML config.\nNOTE: This resource controls this setting globally across your entire AWX/AT setup, and must be managed from a single state",
		CreateContext: resourceSettingsSAMLTeamAttrMapCreate,
		ReadContext:   resourceSettingsSAMLTeamAttrMapRead,
		DeleteContext: resourceSettingsSAMLTeamAttrMapDelete,
		UpdateContext: resourceSettingsSAMLTeamAttrMapUpdate,

		Schema: map[string]*schema.Schema{
			"saml_attr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SAML attribute",
			},
			"remove": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Name of the organization",
			},
			"team_org_map": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"team": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Team Name",
						},
						"organization": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Organization Name",
						},
						"team_alias": {
							// Only supported in AT >= 3.8.0 and AWX >= 12.0.0
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Team Alias",
						},
					},
				},
				Optional:    true,
				Description: "When True, a user who is not a member of ",
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

type samlTeamAttrEntry struct {
	Team         string `json:"team"`
	Organization string `json:"organization"`
	TeamAlias    string `json:"team_alias"`
}

type samlTeamAttrLegacyEntry struct {
	Team         string `json:"team"`
	Organization string `json:"organization"`
}

type samlTeamAttrs struct {
	SamlAttr   string              `json:"saml_attr"`
	Remove     bool                `json:"remove"`
	TeamOrgMap []samlTeamAttrEntry `json:"team_org_map"`
}

type samlTeamAttrsLegacy struct {
	SamlAttr   string                    `json:"saml_attr"`
	Remove     bool                      `json:"remove"`
	TeamOrgMap []samlTeamAttrLegacyEntry `json:"team_org_map"`
}

func resourceSettingsSAMLTeamAttrMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	samlTeamAttrMapAccessMutex.Lock()
	defer samlTeamAttrMapAccessMutex.Unlock()

	var teamAliasSupport bool

	if checkTeamAliasSupport(m) {
		teamAliasSupport = true
	}

	client := m.(*awx.AWX)
	awxService := client.SettingService

	res, err := awxService.GetSettingsBySlug("saml", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: failed to fetch settings",
			"Failed to fetch any saml setting, got: %s", err.Error(),
		)
	}

	var tmaps samlTeamAttrs
	var ltmaps samlTeamAttrsLegacy

	getTeamOrgMap := make([]samlTeamAttrEntry, 0)
	getTeamOrgLMap := make([]samlTeamAttrLegacyEntry, 0)

	if teamAliasSupport {
		err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"], &tmaps)
		if err != nil {
			return buildDiagnosticsMessage(
				"Create: failed to parse SOCIAL_AUTH_SAML_TEAM_ATTR setting",
				"Failed to parse SOCIAL_AUTH_SAML_TEAM_ATTR setting, got: %s with input %s", err.Error(), (*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"],
			)
		}

		if v, ok := d.GetOk("team_org_map"); ok {
			for _, e := range v.(*schema.Set).List() {
				lv := e.(map[string]interface{})
				var en samlTeamAttrEntry
				en.Team = lv["team"].(string)
				en.TeamAlias = lv["team_alias"].(string)
				en.Organization = lv["organization"].(string)

				getTeamOrgMap = append(getTeamOrgMap, en)
			}
		}

		newtmap := samlTeamAttrs{
			SamlAttr:   d.Get("saml_attr").(string),
			Remove:     d.Get("remove").(bool),
			TeamOrgMap: getTeamOrgMap,
		}

		tmaps = newtmap
		payload := map[string]interface{}{
			"SOCIAL_AUTH_SAML_TEAM_ATTR": tmaps,
		}

		_, err = awxService.UpdateSettings("saml", payload, make(map[string]string))
		if err != nil {
			return buildDiagnosticsMessage(
				"Create: organization map not created",
				"failed to save organization map data, got: %s", err.Error(),
			)
		}
	} else {
		err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"], &ltmaps)
		if err != nil {
			return buildDiagnosticsMessage(
				"Create: failed to parse SOCIAL_AUTH_SAML_TEAM_ATTR setting",
				"Failed to parse SOCIAL_AUTH_SAML_TEAM_ATTR setting, got: %s with input %s", err.Error(), (*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"],
			)
		}
		if v, ok := d.GetOk("team_org_map"); ok {
			for _, e := range v.(*schema.Set).List() {
				lv := e.(map[string]interface{})
				var en samlTeamAttrLegacyEntry
				en.Team = lv["team"].(string)
				en.Organization = lv["organization"].(string)

				getTeamOrgLMap = append(getTeamOrgLMap, en)
			}
		}

		newltmap := samlTeamAttrsLegacy{
			SamlAttr:   d.Get("saml_attr").(string),
			Remove:     d.Get("remove").(bool),
			TeamOrgMap: getTeamOrgLMap,
		}

		ltmaps = newltmap
		lpayload := map[string]interface{}{
			"SOCIAL_AUTH_SAML_TEAM_ATTR": ltmaps,
		}

		_, err = awxService.UpdateSettings("saml", lpayload, make(map[string]string))
		if err != nil {
			return buildDiagnosticsMessage(
				"Create: organization map not created",
				"failed to save organization map data, got: %s", err.Error(),
			)
		}
	}

	d.SetId("SOCIAL_AUTH_SAML_TEAM_ATTR")
	return resourceSettingsSAMLTeamAttrMapRead(ctx, d, m)
}

func resourceSettingsSAMLTeamAttrMapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	samlTeamAttrMapAccessMutex.Lock()
	defer samlTeamAttrMapAccessMutex.Unlock()

	client := m.(*awx.AWX)
	awxService := client.SettingService

	res, err := awxService.GetSettingsBySlug("saml", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Update: Unable to fetch settings",
			"Unable to load settings with slug saml: got %s", err.Error(),
		)
	}

	var teamAliasSupport bool

	if checkTeamAliasSupport(m) {
		teamAliasSupport = true
	}

	var tmaps samlTeamAttrs
	var ltmaps samlTeamAttrsLegacy
	getTeamOrgMap := make([]samlTeamAttrEntry, 0)
	getTeamOrgLMap := make([]samlTeamAttrLegacyEntry, 0)

	id := d.Id()

	if teamAliasSupport {
		err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"], &tmaps)
		if err != nil {
			return buildDiagnosticsMessage(
				"Update: failed to parse AUTH_SAML_TEAM_ATTR setting",
				"Failed to parse AUTH_SAML_TEAM_ATTR setting, got: %s", err.Error(),
			)
		}

		if v, ok := d.GetOk("team_org_map"); ok {
			for _, e := range v.(*schema.Set).List() {
				lv := e.(map[string]interface{})
				var en samlTeamAttrEntry
				en.Team = lv["team"].(string)
				en.TeamAlias = lv["team_alias"].(string)
				en.Organization = lv["organization"].(string)

				getTeamOrgMap = append(getTeamOrgMap, en)
			}
		}

		samlAttr := d.Get("saml_attr").(string)
		remove := d.Get("remove").(bool)
		teamOrgMap := getTeamOrgMap

		tmaps.SamlAttr = samlAttr
		tmaps.Remove = remove
		tmaps.TeamOrgMap = teamOrgMap

		payload := map[string]interface{}{
			"SOCIAL_AUTH_SAML_TEAM_ATTR": tmaps,
		}

		_, err = awxService.UpdateSettings("saml", payload, make(map[string]string))
		if err != nil {
			return buildDiagnosticsMessage(
				"Update: organization map not created",
				"failed to save organization map data: %v, got: %s", payload, err.Error(),
			)
		}
	} else {
		err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"], &ltmaps)
		if err != nil {
			return buildDiagnosticsMessage(
				"Update: failed to parse AUTH_SAML_TEAM_ATTR setting",
				"Failed to parse AUTH_SAML_TEAM_ATTR setting, got: %s", err.Error(),
			)
		}

		if v, ok := d.GetOk("team_org_map"); ok {
			for _, e := range v.(*schema.Set).List() {
				lv := e.(map[string]interface{})
				var en samlTeamAttrLegacyEntry
				en.Team = lv["team"].(string)
				en.Organization = lv["organization"].(string)

				getTeamOrgLMap = append(getTeamOrgLMap, en)
			}
		}

		samlAttr := d.Get("saml_attr").(string)
		remove := d.Get("remove").(bool)
		teamOrgLMap := getTeamOrgLMap

		ltmaps.SamlAttr = samlAttr
		ltmaps.Remove = remove
		ltmaps.TeamOrgMap = teamOrgLMap

		lpayload := map[string]interface{}{
			"SOCIAL_AUTH_SAML_TEAM_ATTR": ltmaps,
		}

		_, err = awxService.UpdateSettings("saml", lpayload, make(map[string]string))
		if err != nil {
			return buildDiagnosticsMessage(
				"Update: organization map not created",
				"failed to save organization map data: %v, got: %s", lpayload, err.Error(),
			)
		}
	}

	d.SetId(id)
	return resourceSettingsSAMLTeamAttrMapRead(ctx, d, m)
}

func resourceSettingsSAMLTeamAttrMapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	var tmaps samlTeamAttrs
	var ltmaps samlTeamAttrsLegacy

	var teamAliasSupport bool

	if checkTeamAliasSupport(m) {
		teamAliasSupport = true
	}

	if teamAliasSupport {
		err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"], &tmaps)
		if err != nil {
			return buildDiagnosticsMessage(
				"Unable to parse SOCIAL_AUTH_SAML_TEAM_ATTR",
				"Unable to parse SOCIAL_AUTH_SAML_TEAM_ATTR, got: %s", err.Error(),
			)
		}

		setTeamOrgMap := make([]map[string]interface{}, 0)

		for _, teamAttr := range tmaps.TeamOrgMap {
			lv := map[string]interface{}{
				"team":         "",
				"organization": "",
				"team_alias":   "",
			}

			lv["team"] = teamAttr.Team
			lv["organization"] = teamAttr.Organization
			lv["team_alias"] = teamAttr.TeamAlias

			setTeamOrgMap = append(setTeamOrgMap, lv)
		}
		d.Set("team_org_map", setTeamOrgMap)

	} else {
		err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"], &ltmaps)
		if err != nil {
			return buildDiagnosticsMessage(
				"Unable to parse SOCIAL_AUTH_SAML_TEAM_ATTR",
				"Unable to parse SOCIAL_AUTH_SAML_TEAM_ATTR, got: %s", err.Error(),
			)
		}

		setTeamOrgLMap := make([]map[string]interface{}, 0)

		for _, lteamAttr := range ltmaps.TeamOrgMap {
			lv := map[string]interface{}{
				"team":         "",
				"organization": "",
			}

			lv["team"] = lteamAttr.Team
			lv["organization"] = lteamAttr.Organization

			setTeamOrgLMap = append(setTeamOrgLMap, lv)
		}
		d.Set("team_org_map", setTeamOrgLMap)
	}

	d.SetId(d.Id())
	d.Set("saml_attr", ltmaps.SamlAttr)
	d.Set("remove", ltmaps.Remove)
	return diags
}

func resourceSettingsSAMLTeamAttrMapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	samlTeamAttrMapAccessMutex.Lock()
	defer samlTeamAttrMapAccessMutex.Unlock()

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

	var tmaps samlTeamAttrs
	tmaps.TeamOrgMap = make([]samlTeamAttrEntry, 0)

	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Delete: failed to parse SOCIAL_AUTH_SAML_TEAM_ATTR setting",
			"Failed to parse SOCIAL_AUTH_SAML_TEAM_ATTR setting, got: %s", err.Error(),
		)
	}

	payload := map[string]interface{}{
		"SOCIAL_AUTH_SAML_TEAM_ATTR": tmaps,
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
