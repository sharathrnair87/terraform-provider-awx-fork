/*
Use this resource to create an Ansible Galaxy Credential in AWX/AT

# Example Usage
```hcl

	data "awx_organization" "devops" {
	  name = "DevOps"
	}

	resource "awx_credential_galaxy" "devops_galaxy_cred_01" {
	  name            = "DevOps_Cred_01"
	  description     = "Galaxy Cred for DevOps Org"
	  organization_id = data.awx_organization.devops
	  url             = "https://galaxy.ansible.com"
	}

```
*/
package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func dataSourceCredentialGalaxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to query an Ansible Galaxy Credential in AWX/AT",
		ReadContext: dataSourceCredentialGalaxyRead,
		Schema: map[string]*schema.Schema{
			"credential_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func dataSourceCredentialGalaxyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := d.Get("credential_id").(int)
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credentials",
			Detail:   fmt.Sprintf("Unable to credentials with id %d: %s", id, err.Error()),
		})
		return diags
	}

	d.Set("name", cred.Name)
	d.Set("description", cred.Description)
	d.Set("url", cred.Inputs["url"])
	d.Set("auth_url", cred.Inputs["auth_url"])
	d.Set("token", d.Get("token").(string))
	d.Set("organization_id", cred.OrganizationID)
	d.SetId(strconv.Itoa(id))

	return diags
}
