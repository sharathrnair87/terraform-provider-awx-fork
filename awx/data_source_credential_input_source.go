/*
Use this data source to query an Input Source mapping for an AWX/AT Credential

# Example Usage

```hcl
data "awx_credential_machine" "ssh_key" {
  name = "iaas-ssh-key"
}

data "awx_credential_hashivault_signed_ssh" "hcv_signer" {
  name = "hcv-sig-ssh"
}

data "awx_credential_input_source" "hcv-iaas-sig-map" {
    input_source_id = <id_of_credential_input_source>
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

func dataSourceCredentialInputSource() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceCredentialInputSourceRead,
		Schema: map[string]*schema.Schema{
            "input_source_id": {
                Type: schema.TypeInt,
                Required: true,
            },
            "description": {
                Type: schema.TypeString,
                Computed: true,
            },
			"input_field_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"metadata": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}


func dataSourceCredentialInputSourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := d.Get("input_source_id").(int)
	inputSource, err := client.CredentialInputSourceService.GetCredentialInputSourceByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credentials",
			Detail:   fmt.Sprintf("Unable to credentials with id %d: %s", id, err.Error()),
		})
		return diags
	}

	d.Set("description", inputSource.Description)
	d.Set("input_field_name", inputSource.InputFieldName)
	d.Set("target", inputSource.TargetCredential)
	d.Set("source", inputSource.SourceCredential)
	d.Set("metadata", inputSource.Metadata)
    d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
