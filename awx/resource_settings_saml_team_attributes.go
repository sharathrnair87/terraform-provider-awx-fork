/*
*TBD*

# Example Usage

```hcl

	data "awx_organization" "default" {
	  name = "Default"
	}

	resource "awx_organization" "admin_organization" {
	  name = "Admins"
	  organization_id = data.awx_organization.default.id
	}

	resource "awx_settings_saml_team_attr" "admin_team_attr" {
	  name         = resource.awx_organization.admin_organization.name
	  users        = ["CN=MyTeamAttr,OU=Groups,DC=example,DC=com"]
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

var samlTeamAttrMapAccessMutex sync.Mutex

func resourceSettingsSAMLTeamAttrMap() *schema.Resource {
	return &schema.Resource{
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
							Required:    true,
							Description: "Team Name",
						},
						"organization": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Organization Name",
						},
						"team_alias": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Team Alias",
						},
					},
				},
				Optional:    true,
				Description: "When True, a user who is not a member of ",
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

type samlTeamAttrEntry struct {
	Team         string `json:"team"`
	Organization string `json:"organization"`
	TeamAlias    string `json:"team_alias"`
}

type samlTeamAttrs struct {
	SamlAttr   string              `json:"saml_attr"`
	Remove     bool                `json:"remove"`
	TeamOrgMap []samlTeamAttrEntry `json:"team_org_map"`
}

func resourceSettingsSAMLTeamAttrMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	samlTeamAttrMapAccessMutex.Lock()
	defer samlTeamAttrMapAccessMutex.Unlock()

	client := m.(*awx.AWX)
	awxService := client.SettingService

	res, err := awxService.GetSettingsBySlug("saml", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: failed to fetch settings",
			"Failed to fetch any saml setting, got: %s", err.Error(),
		)
	}

	//tmaps := map[string]interface{} {
	//    []samlTeamAttrMap{},
	//}

	//tmaps := []samlTeamAttrMap{}

	var tmaps samlTeamAttrs

	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: failed to parse SOCIAL_AUTH_SAML_TEAM_ATTR setting",
			"Failed to parse SOCIAL_AUTH_SAML_TEAM_ATTR setting, got: %s with input %s", err.Error(), (*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"],
		)
	}

	//team := d.Get("team").(string)
	//organization := d.Get("organization").(string)
	//teamAlias := d.Get("team_alias").(string)

	var getTeamOrgMap []samlTeamAttrEntry

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

	//if d == tmaps {
	//    return buildDiagnosticsMessage(
	//        "Create: organization map already exists",
	//        "Map for saml to organization map %v already exists", d.Id(),
	//    )
	//}

	//for _, t := range tmaps.TeamOrgMap {
	//    // Check for duplicates
	//    //if (team == t.Team) && (organization == t.Organization) && (teamAlias == t.TeamAlias) {
	//    if newtmap == t {
	//    }
	//}

	//tmaps.TeamOrgMap = append(tmaps.TeamOrgMap, newtmap)
	//id := len(tmaps.TeamOrgMap)

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

	//d.SetId(genHash(fmt.Sprintf("%s%s%s", newtmap.Team, newtmap.Organization, newtmap.TeamAlias)))
	//id := uuid.New()
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

	var tmaps samlTeamAttrs
	//var utmaps samlTeamAttrs

	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Update: failed to parse AUTH_SAML_TEAM_ATTR setting",
			"Failed to parse AUTH_SAML_TEAM_ATTR setting, got: %s", err.Error(),
		)
	}

	id := d.Id()

	var getTeamOrgMap []samlTeamAttrEntry

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

	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Unable to parse SOCIAL_AUTH_SAML_TEAM_ATTR",
			"Unable to parse SOCIAL_AUTH_SAML_TEAM_ATTR, got: %s", err.Error(),
		)
	}

	//mapdef, ok := tmaps[d.Id()]
	//if !ok {
	//	return buildDiagnosticsMessage(
	//		"Unable to fetch saml organization map",
	//		"Unable to load saml organization map %v: not found", d.Id(),
	//	)
	//}

	//var users []string
	//switch tt := mapdef.SamlUserGroups.(type) {
	//case string:
	//	users = []string{tt}
	//case []string:
	//	users = tt
	//case []interface{}:
	//	for _, v := range tt {
	//		if t, ok := v.(string); ok {
	//			users = append(users, t)
	//		}
	//	}
	//}

	//var admins []string
	//switch tt := mapdef.SamlAdminGroups.(type) {
	//case string:
	//	admins = []string{tt}
	//case []string:
	//	admins = tt
	//case []interface{}:
	//	for _, v := range tt {
	//		if t, ok := v.(string); ok {
	//			admins = append(admins, t)
	//		}
	//	}
	//}

	var setTeamOrgMap []map[string]interface{}

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

	d.SetId(d.Id())
	d.Set("saml_attr", tmaps.SamlAttr)
	d.Set("remove", tmaps.Remove)
	d.Set("team_org_map", setTeamOrgMap)
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

	err = json.Unmarshal((*res)["SOCIAL_AUTH_SAML_TEAM_ATTR"], &tmaps)
	if err != nil {
		return buildDiagnosticsMessage(
			"Delete: failed to parse SOCIAL_AUTH_SAML_TEAM_ATTR setting",
			"Failed to parse SOCIAL_AUTH_SAML_TEAM_ATTR setting, got: %s", err.Error(),
		)
	}

	//id := d.Id()
	//delete(tmaps, id)

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
