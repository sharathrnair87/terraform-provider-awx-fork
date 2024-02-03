/*
Use this resource to create an AWX/AT Organization

# Example Usage

```hcl

	resource "awx_organization" "default" {
	  name            = "acc-test"
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

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create an AWX/AT Organization",
		CreateContext: resourceOrganizationsCreate,
		ReadContext:   resourceOrganizationsRead,
		UpdateContext: resourceOrganizationsUpdate,
		DeleteContext: resourceOrganizationsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"max_hosts": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum number of hosts allowed to be managed by this organization",
			},
			"custom_virtualenv": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Local absolute file path containing a custom Python virtualenv to use",
			},
			"default_environment": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The default execution environment for jobs run by this organization.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceOrganizationsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.OrganizationsService

	result, err := awxService.CreateOrganization(map[string]interface{}{
		"name":                d.Get("name").(string),
		"description":         d.Get("description").(string),
		"max_hosts":           d.Get("max_hosts").(int),
		"custom_virtualenv":   d.Get("custom_virtualenv").(string),
		"default_environment": d.Get("default_environment").(int),
	}, map[string]string{})
	if err != nil {
		log.Printf("Failed to Create Organization %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Organizations",
			Detail:   fmt.Sprintf("Failed to create Organization with name %s: %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceOrganizationsRead(ctx, d, m)
}

func resourceOrganizationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.OrganizationsService
	id, diags := convertStateIDToNumeric("Update Organizations", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)

	_, err := awxService.GetOrganizationsByID(id, params)
	if err != nil {
		return buildDiagNotFoundFail("Organizations", id, err)
	}

	_, err = awxService.UpdateOrganization(id, map[string]interface{}{
		"name":                d.Get("name").(string),
		"description":         d.Get("description").(string),
		"max_hosts":           d.Get("max_hosts").(int),
		"custom_virtualenv":   d.Get("custom_virtualenv").(string),
		"default_environment": d.Get("default_environment").(int),
	}, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Organizations",
			Detail:   fmt.Sprintf("Organizations with name %s failed to update %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	return resourceOrganizationsRead(ctx, d, m)
}

func resourceOrganizationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.OrganizationsService
	id, diags := convertStateIDToNumeric("Read Organizations", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetOrganizationsByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("Organization", id, err)

	}
	d = setOrganizationsResourceData(d, res)
	return nil
}

func resourceOrganizationsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	digMessagePart := "Organization"
	client := m.(*awx.AWX)
	awxService := client.OrganizationsService
	id, diags := convertStateIDToNumeric("Delete Organization", d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteOrganization(id); err != nil {
		return buildDiagDeleteFail(digMessagePart, fmt.Sprintf("OrganizationID %v, got %s ", id, err.Error()))
	}
	d.SetId("")
	return diags
}

func setOrganizationsResourceData(d *schema.ResourceData, r *awx.Organization) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("max_hosts", r.MaxHosts)
	d.Set("custom_virtualenv", r.CustomVirtualenv)
	d.Set("default_environment", r.DefaultEnvironment)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
