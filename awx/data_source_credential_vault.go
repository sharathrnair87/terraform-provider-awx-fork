/*
Use this data source to query an Ansible Vault Credential in AWX/AT

# Example Usage

```hcl

	data "awx_credential_vault" "my_vault_cred" {
	    credential_id = <my_vault_cred_id>
	}

	output "my_vault_cred_vault_id" {
	    value = data.awx_credential_vault.my_vault_cred.vault_id
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

func dataSourceCredentialVault() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialVaultRead,
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
			"vault_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vault_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceCredentialVaultRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("vault_id", cred.Inputs["vault_id"])
	d.Set("vault_password", cred.Inputs["vault_password"])
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
