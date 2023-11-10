/*
Use this data source to query an SCM Credential in AWX/AT

# Example Usage

```hcl

	data "awx_credential_scm" "my_scm_cred" {
	  credential_id = <my_scm_cred_id>
	}

	output "my_scm_cred" {
	  value     = data.awx_credential_scm.my_scm_cred
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

func dataSourceCredentialSCM() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialSCMRead,
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
			"ssh_key_unlock": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceCredentialSCMRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("ssh_key_unlock", cred.Inputs["ssh_key_unlock"])
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
