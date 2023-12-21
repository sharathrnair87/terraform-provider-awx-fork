/*
This resource configures generic AWX settings.
Please note that resource deletion only deletes the object from terraform state and does not reset the setting to its initial value.

See available settings list here: https://docs.ansible.com/ansible-tower/latest/html/towerapi/api_ref.html#/Settings/Settings_settings_update

# Example Usage

```hcl

	resource "awx_setting" "social_auth_saml_technical_contact" {
	  name  = "SOCIAL_AUTH_SAML_TECHNICAL_CONTACT"
	  value = <<EOF
	  {
	    "givenName": "Myorg",
	    "emailAddress": "test@foo.com"
	  }
	  EOF
	}

	resource "awx_setting" "social_auth_saml_sp_entity_id" {
	  name  = "SOCIAL_AUTH_SAML_SP_ENTITY_ID"
	  value = "test"
	}

	resource "awx_setting" "schedule_max_jobs" {
	  name  = "SCHEDULE_MAX_JOBS"
	  value = 15
	}

	resource "awx_setting" "remote_host_headers" {
	  name  = "REMOTE_HOST_HEADERS"
	  value = <<EOF
	  [
	    "HTTP_X_FORWARDED_FOR",
	    "REMOTE_ADDR",
	    "REMOTE_HOST"
	  ]
	  EOF
	}

```
*/
package awx

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func resourceSetting() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to manage the AWX/AT Generic Controller Configuration.\nNOTE: Settings managed using this resource must be in a single state to avoid unexpected behaviour",
		CreateContext: resourceSettingUpdate,
		ReadContext:   resourceSettingRead,
		DeleteContext: resourceSettingDelete,
		UpdateContext: resourceSettingUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of setting to modify",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value to be modified for given setting.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Supress Diffs when the Stringified JSON is logically equivalent
					var map_decoded map[string]interface{}
					var slice_decoded []interface{}
					err := json.Unmarshal([]byte(old), &map_decoded)
					if err != nil {
						err = json.Unmarshal([]byte(old), &slice_decoded)

						if err != nil {
							// string
							if strings.TrimSpace(old) == strings.TrimSpace(new) {
								return true
							}
						} else {
							// stringified List
							var new_slice_decoded []interface{}
							err = json.Unmarshal([]byte(new), &new_slice_decoded)
							if err != nil {
								return false
							} else {
								return reflect.DeepEqual(slice_decoded, new_slice_decoded)
							}
						}
					} else {
						// stringified JSON
						var new_map_decoded map[string]interface{}
						err = json.Unmarshal([]byte(new), &new_map_decoded)
						if err != nil {
							return false
						} else {
							return reflect.DeepEqual(map_decoded, new_map_decoded)
						}
					}

					return false
				},
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

type setting map[string]string

func resourceSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.SettingService

	_, err := awxService.GetSettingsBySlug("all", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: failed to fetch settings",
			"Failed to fetch setting, got: %s", err.Error(),
		)
	}

	var map_decoded map[string]interface{}
	var array_decoded []interface{}
	var formatted_value interface{}

	name := d.Get("name").(string)
	value := d.Get("value").(string)

	// Attempt to unmarshall string into a map
	err = json.Unmarshal([]byte(value), &map_decoded)

	if err != nil {
		// Attempt to unmarshall string into an array
		err = json.Unmarshal([]byte(value), &array_decoded)

		if err != nil {
			formatted_value = value
		} else {
			formatted_value = array_decoded
		}
	} else {
		formatted_value = map_decoded
	}

	payload := map[string]interface{}{
		name: formatted_value,
	}

	_, err = awxService.UpdateSettings("all", payload, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: setting not created",
			"failed to save setting data, got: %s, %s", err.Error(), value,
		)
	}

	d.SetId(name)
	return resourceSettingRead(ctx, d, m)
}

func resourceSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.SettingService

	name := d.Id()

	setting, err := awxService.GetSettingsBySlug("all", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Unable to fetch settings",
			"Unable to load settings with slug all: got %s", err.Error(),
		)
	}

	d.Set("name", name)
	d.Set("value", string((*setting)[name]))
	return diags
}

func resourceSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}
