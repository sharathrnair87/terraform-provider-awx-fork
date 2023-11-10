/*
Use this data source to query an AWX/AT Job Template

# Example Usage

```hcl

	data "awx_job_template" "default" {
	  name = "Default"
	}

    output "def_job_templ_playbook" {
        value = data.awx_job_template.default.playbook
    }

```
*/
package awx

import (
	"context"
	"strconv"

	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func dataSourceJobTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJobTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
                ConflictsWith: []string{"name"},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
                ConflictsWith: []string{"id"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
			},
			"inventory_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"playbook": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"forks": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"limit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"verbosity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed: true,
			},
			"extra_vars": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_tags": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"force_handlers": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"skip_tags": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"start_at_task": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"use_fact_cache": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"host_config_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ask_diff_mode_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ask_limit_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ask_tags_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ask_verbosity_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ask_inventory_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ask_variables_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ask_credential_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"survey_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"become_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"diff_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ask_skip_tags_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"allow_simultaneous": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"custom_virtualenv": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ask_job_type_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"execution_environment": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceJobTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if jtName, okName := d.GetOk("name"); okName {
		params["name"] = jtName.(string)
	}

	if jtID, okJtID := d.GetOk("id"); okJtID {
		params["id"] = strconv.Itoa(jtID.(int))
	}

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or jt_id)",
		)
	}

	jobTemplate, _, err := client.JobTemplateService.ListJobTemplates(params)

	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Job Template",
			"Fail to find the jt got: %s",
			err.Error(),
		)
	}

	for _, template := range jobTemplate {
		log.Printf("loop %v", template.Name)
		if template.Name == params["name"] {
			d = setJobTemplateResourceData(d, template)
			return diags
		}
	}

	if _, okGroupID := d.GetOk("id"); okGroupID {
		log.Printf("byid %v", len(jobTemplate))
		if len(jobTemplate) != 1 {
			return buildDiagnosticsMessage(
				"Get: found more than one Element",
				"The Query Returns more than one Job Template, %d",
				len(jobTemplate),
			)
		}
		d = setJobTemplateResourceData(d, jobTemplate[0])
		return diags
	}
	return buildDiagnosticsMessage(
		"Get: found more than one Element",
		"The Query Returns more than one Job Template, %d",
		len(jobTemplate),
	)
}
