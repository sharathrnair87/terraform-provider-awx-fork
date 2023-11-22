/*
Use this data source to query a HashiCorp Vault Signed SSH Credential in AWX/AT

# Example Usage

```hcl

	data "awx_credential_hashivault_signed_ssh" "my_hashi_signed_ssh" {
	  credential_id = var.my_hashi_signed_ssh_id
	}

	output "my_hashi_signed_ssh" {
	  value     = data.awx_credential_hashivault_signed_ssh.my_hashi_signed_ssh
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

func dataSourceCredentialHashiVaultSSH() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query a HashiCorp Vault Signed SSH Credential in AWX/AT",
		ReadContext: dataSourceCredentialHashiVaultSSHRead,
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
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kubernetes_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_auth_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCredentialHashiVaultSSHRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("role_id", cred.Inputs["role_id"])
	d.Set("secret_id", cred.Inputs["secret_id"])
	d.Set("namespace", cred.Inputs["namespace"])
	d.Set("kubernetes_role", cred.Inputs["kubernetes_role"])
	d.Set("default_auth_path", cred.Inputs["default_auth_path"])
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
