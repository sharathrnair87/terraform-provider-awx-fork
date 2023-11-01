/*
*TBD*

# Example Usage

```hcl

	resource "awx_job_template" "my_job_template" {
	  name = "My JobTemplate"
	}

	data "awx_job_template_role" "job_template_admin_role" {
	  role_name       = "Admin"
	  job_template_id = data.awx_job_template.my_job_template.id
	}

```
*/
package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func dataSourceJobTemplateRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJobTemplateRoleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			"job_template_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func dataSourceJobTemplateRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	jobTmpl_id := d.Get("job_template_id").(int)
	jobTemplate, err := client.JobTemplateService.GetJobTemplateByID(jobTmpl_id, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch JobTemplate",
			"Fail to find the job template, got: %s",
			err.Error(),
		)
	}

	roleslist := []*awx.ApplyRole{
		jobTemplate.SummaryFields.ObjectRoles.AdminRole,
		jobTemplate.SummaryFields.ObjectRoles.ReadRole,
		jobTemplate.SummaryFields.ObjectRoles.ExecuteRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range roleslist {
			if v != nil && id == v.ID {
				d = setJobTemplateRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range roleslist {
			if v != nil && name == v.Name {
				d = setJobTemplateRoleData(d, v)
				return diags
			}
		}
	}

	return buildDiagnosticsMessage(
		"Failed to fetch job template role - Not Found",
		"The project role was not found",
	)
}

func setJobTemplateRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	d.Set("name", r.Name)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
