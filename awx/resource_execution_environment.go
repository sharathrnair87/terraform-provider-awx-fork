/*
Use this resource to create an Execution Environment in AWX/AT

# Example Usage

```hcl

	data "awx_organization" "devops" {
	  name = "DevOps"
	}

	resource "awx_execution_environment" "default" {
	  name         = "acc-test"
	  image        = "us-docker.pkg.dev/cloudrun/container/hello"
	  pull         = "never"
	  organization = data.awx_organization.devops.id
	}

```
*/
package awx

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func resourceExecutionEnvironment() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create an Execution Environment in AWX/AT",
		CreateContext: resourceExecutionEnvironmentsCreate,
		ReadContext:   resourceExecutionEnvironmentsRead,
		UpdateContext: resourceExecutionEnvironmentsUpdate,
		DeleteContext: resourceExecutionEnvironmentsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"credential": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"pull": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "always" && v != "missing" && v != "never" && v != "" {
						errs = append(errs, fmt.Errorf("%q must be one of 'always', 'missing', 'never' or '' (blank), got %s", key, v))
					}
					return
				},
			},
		},
	}
}

func resourceExecutionEnvironmentsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ExecutionEnvironmentsService

	result, err := awxService.CreateExecutionEnvironment(map[string]interface{}{
		"name":         d.Get("name").(string),
		"image":        d.Get("image").(string),
		"description":  d.Get("description").(string),
		"organization": d.Get("organization").(string),
		"credential":   AtoipOr(d.Get("credential").(string), nil),
		"pull":         d.Get("pull").(string),
	}, map[string]string{})
	if err != nil {
		log.Printf("Failed to Create ExecutionEnvironment %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create ExecutionEnvironments",
			Detail:   fmt.Sprintf("ExecutionEnvironments with name, failed to create %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceExecutionEnvironmentsRead(ctx, d, m)
}

func resourceExecutionEnvironmentsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ExecutionEnvironmentsService
	id, diags := convertStateIDToNumeric("Update ExecutionEnvironments", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)

	_, err := awxService.GetExecutionEnvironmentByID(id, params)
	if err != nil {
		return buildDiagNotFoundFail("ExecutionEnvironments", id, err)
	}

	_, err = awxService.UpdateExecutionEnvironment(id, map[string]interface{}{
		"name":         d.Get("name").(string),
		"image":        d.Get("image").(string),
		"description":  d.Get("description").(string),
		"organization": d.Get("organization").(string),
		"credential":   AtoipOr(d.Get("credential").(string), nil),
		"pull":         d.Get("pull").(string),
	}, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update ExecutionEnvironments",
			Detail:   fmt.Sprintf("ExecutionEnvironments with name %s failed to update %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	return resourceExecutionEnvironmentsRead(ctx, d, m)
}

func resourceExecutionEnvironmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ExecutionEnvironmentsService
	id, diags := convertStateIDToNumeric("Read ExecutionEnvironments", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetExecutionEnvironmentByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("ExecutionEnvironment", id, err)

	}
	d = setExecutionEnvironmentsResourceData(d, res)
	return nil
}

func resourceExecutionEnvironmentsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	digMessagePart := "ExecutionEnvironment"
	client := m.(*awx.AWX)
	awxService := client.ExecutionEnvironmentsService
	id, diags := convertStateIDToNumeric("Delete ExecutionEnvironment", d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteExecutionEnvironment(id); err != nil {
		return buildDiagDeleteFail(digMessagePart, fmt.Sprintf("ExecutionEnvironmentID %v, got %s ", id, err.Error()))
	}
	d.SetId("")
	return diags
}

func setExecutionEnvironmentsResourceData(d *schema.ResourceData, r *awx.ExecutionEnvironment) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("image", r.Image)
	d.Set("description", r.Description)
	d.Set("organization", r.Organization)
	d.Set("credential", r.Credential)
	d.Set("pull", r.Pull)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
