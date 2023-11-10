/*
Use this data source to query a HashiCorp Vault Secret Lookup credential in AWX/AT

# Example Usage

```hcl

	data "awx_credential_hashivault_secret" "my_hashi_secret_lookup" {
	  credential_id = <my_hashi_secret_lookup_id>
	}

	output "my_hashi_secret_lookup" {
	  value     = data.awx_credential_hashivault_secret.my_hashi_secret_lookup
	  sensitive = true
	}

```
*/
package awx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func dataSourceCredentialHashiVault() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialHashiVaultRead,
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
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"cacert": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCredentialHashiVaultRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("organization_id", cred.OrganizationID)
	d.Set("url", cred.Inputs["url"])
	d.Set("token", cred.Inputs["token"])
	d.Set("cacert", cred.Inputs["cacert"])
	d.Set("api_version", cred.Inputs["api_version"])
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
