/*
Use this data source to query a Project in AWX/AT

# Example Usage

```hcl

	    // By Name
		data "awx_project" "default" {
		  name = "Default"
		}

	    // By ID
	    data "awx_project" "sharedServices" {
	        id = var.shared_services_prj_id
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

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query a Project in AWX/AT",
		ReadContext: dataSourceProjectsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"local_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scm_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scm_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scm_credential_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scm_branch": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scm_clean": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"scm_delete_on_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"scm_update_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"scm_update_cache_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceProjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or project_id)",
		)
	}
	projects, _, err := client.ProjectService.ListProjects(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch Project",
			"Failed to find the group got: %s",
			err.Error(),
		)
	}
	if len(projects) > 1 {
		return buildDiagnosticsMessage(
			"Get: found more than one Element",
			"The Query Returns more than one Project, %d",
			len(projects),
		)
	}
	if len(projects) == 0 {
		return buildDiagnosticsMessage(
			"Get: No Project found",
			"The Query Returns no Project matching filter, %v",
			params,
		)
	}

	Project := projects[0]
	d = setProjectResourceData(d, Project)
	return diags
}
