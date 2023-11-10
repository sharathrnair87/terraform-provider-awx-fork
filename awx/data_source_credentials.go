/*
Use this data source to all available Credentials in AWX/AT

# Example Usage

```hcl
// Query all available credentials
data "awx_credentials" "all_creds" {}

	output "all_creds" {
	  value = length(data.awx_credentials.all_creds.credentials)
	}

// Query a specific credential by type

	output "ssh" {
	  value = toset([for each in data.awx_credentials.creds.credentials : each if each.kind == "ssh"])
	}

// Query all credentials in a given Organization

	data "awx_organization" "org" {
	  name = "My Org"
	}

	output "creds_in_my_org" {
	  value = toset([for each in data.awx_credentials.creds.credentials :
	    each if each.organization_id == data.awx_organization.org.id
	  ])
	}

```
*/
package awx

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func dataSourceCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialsRead,
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCredentialsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)

	creds, err := client.CredentialsService.ListCredentials(map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credentials",
			Detail:   "Unable to fetch credentials from AWX API",
		})
		return diags
	}

	parsedCreds := make([]map[string]interface{}, 0)
	for _, c := range creds {
		parsedCreds = append(parsedCreds, map[string]interface{}{
			"id":              c.ID,
			"username":        c.Inputs["username"],
			"kind":            c.Kind,
			"name":            c.Name,
			"organization_id": c.OrganizationID,
		})
	}

	err = d.Set("credentials", parsedCreds)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
