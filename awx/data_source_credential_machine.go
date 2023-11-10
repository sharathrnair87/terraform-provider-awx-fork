/*
Use this data source to query a Machine Credential in AWX/AT

# Example Usage

```hcl

	data "awx_credential_machine" "my_machine_creds" {
	  credential_id = <my_machine_creds>
	}

	output "my_machine_creds" {
	  value     = data.awx_credential_machine.my_machine_creds
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

func dataSourceCredentialMachine() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialMachineRead,
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
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"ssh_key_data": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"ssh_public_key_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_key_unlock": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"become_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"become_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"become_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceCredentialMachineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("username", cred.Inputs["username"])
	d.Set("password", cred.Inputs["password"])
	d.Set("ssh_key_data", cred.Inputs["ssh_key_data"])
	d.Set("ssh_public_key_data", cred.Inputs["ssh_public_key_data"])
	d.Set("ssh_key_unlock", cred.Inputs["ssh_key_unlock"])
	d.Set("become_method", cred.Inputs["become_method"])
	d.Set("become_username", cred.Inputs["become_username"])
	d.Set("become_password", cred.Inputs["become_password"])
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
