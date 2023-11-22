/*
Use this data source to lookup an Azure Resource Manager Credential in AWX/AT

# Example Usage

```hcl

	data "awx_credential_azure_resource_manager" "my_azrm_creds" {
	  credential_id = var.my_azrm_cred_id
	}

	output "my_azrm_creds" {
	  value     = data.awx_credential_azure_resource_manager.my_azrm_creds
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

func dataSourceCredentialAzureRM() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to lookup an Azure Resource Manager Credential in AWX/AT",
		ReadContext: dataSourceCredentialAzureRMRead,
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
			"subscription": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
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
			"tenant": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCredentialAzureRMRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("subscription", cred.Inputs["subscription"])
	d.Set("client", cred.Inputs["client"])
	d.Set("secret", d.Get("secret").(string))
	d.Set("username", cred.Inputs["username"])
	d.Set("password", d.Get("password").(string))
	d.Set("tenant", cred.Inputs["tenant"])
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
