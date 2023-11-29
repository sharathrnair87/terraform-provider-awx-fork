/*
Use this resource to create an Inventory Source for an existing Inventory in AWX/AT

# Example Usage

```hcl

	data "awx_inventory" "db_inventory" {
	  name            = "DB_Inventory"
	  organization_id = data.awx_organization.db.id
	}

	data "awx_project" "db_project" {
	  name = "DB_Infra"
	}

	resource "awx_inventory_source" "db_inventory_source" {
	  name              = "DB Inventory Src"
	  inventory_id      = data.awx_inventory.db_inventory.id
	  source_project_id = data.awx_project.db_project.id
	  source_path       = "inventory/db-hosts.yml"
	}

```
*/
package awx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func resourceInventorySource() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create an Inventory Source for an existing Inventory in AWX/AT",
		CreateContext: resourceInventorySourceCreate,
		ReadContext:   resourceInventorySourceRead,
		UpdateContext: resourceInventorySourceUpdate,
		DeleteContext: resourceInventorySourceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled_var": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"overwrite": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"overwrite_vars": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"update_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"inventory_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"credential_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Default:  "scm",
				Optional: true,
			},
			"source_vars": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_cache_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"verbosity": {
				Type:     schema.TypeInt,
				Default:  1,
				Optional: true,
			},
			"instance_filters": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_project_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"source_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// obsolete schema added so terraform doesn't break
			// these don't do anything in later versions of AWX! Update your code.
			"source_regions": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceInventorySourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventorySourcesService

	inventorySourceMap := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"description":          d.Get("description").(string),
		"enabled_var":          d.Get("enabled_var").(string),
		"enabled_value":        d.Get("enabled_value").(string),
		"overwrite":            d.Get("overwrite").(bool),
		"overwrite_vars":       d.Get("overwrite_vars").(bool),
		"update_on_launch":     d.Get("update_on_launch").(bool),
		"inventory":            d.Get("inventory_id").(int),
		"source":               d.Get("source").(string),
		"source_vars":          d.Get("source_vars").(string),
		"host_filter":          d.Get("host_filter").(string),
		"update_cache_timeout": d.Get("update_cache_timeout").(int),
		"verbosity":            d.Get("verbosity").(int),
		"instance_filters":     d.Get("instance_filters").(string),
		"group_by":             d.Get("group_by").(string),
		"source_path":          d.Get("source_path").(string),
		// obsolete schema added so terraform doesn't break
		// these don't do anything in later versions of AWX! Update your code.
		"source_regions": d.Get("source_regions").(string),
	}

	if credential, ok := d.GetOk("credential_id"); ok {
		inventorySourceMap["credential"] = credential.(int)
	}

	if sourceProjectID, ok := d.GetOk("source_project_id"); ok {
		inventorySourceMap["credential"] = sourceProjectID.(int)
	}

	result, err := awxService.CreateInventorySource(inventorySourceMap, map[string]string{})
	if err != nil {
		return buildDiagCreateFail(diagElementInventorySourceTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventorySourceRead(ctx, d, m)

}

func resourceInventorySourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventorySourcesService
	id, diags := convertStateIDToNumeric(diagElementInventorySourceTitle, d)
	if diags.HasError() {
		return diags
	}

	inventorySourceMap := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"description":          d.Get("description").(string),
		"enabled_var":          d.Get("enabled_var").(string),
		"enabled_value":        d.Get("enabled_value").(string),
		"overwrite":            d.Get("overwrite").(bool),
		"overwrite_vars":       d.Get("overwrite_vars").(bool),
		"update_on_launch":     d.Get("update_on_launch").(bool),
		"inventory":            d.Get("inventory_id").(int),
		"source":               d.Get("source").(string),
		"source_vars":          d.Get("source_vars").(string),
		"host_filter":          d.Get("host_filter").(string),
		"update_cache_timeout": d.Get("update_cache_timeout").(int),
		"verbosity":            d.Get("verbosity").(int),
		"instance_filters":     d.Get("instance_filters").(string),
		"group_by":             d.Get("group_by").(string),
		"source_path":          d.Get("source_path").(string),
		// obsolete schema added so terraform doesn't break
		// these don't do anything in later versions of AWX! Update your code.
		"source_regions": d.Get("source_regions").(string),
	}

	if credential, ok := d.GetOk("credential_id"); ok {
		inventorySourceMap["credential"] = credential.(int)
	}

	if sourceProjectID, ok := d.GetOk("source_project_id"); ok {
		inventorySourceMap["credential"] = sourceProjectID.(int)
	}

	_, err := awxService.UpdateInventorySource(id, inventorySourceMap, nil)
	if err != nil {
		return buildDiagUpdateFail(diagElementInventorySourceTitle, id, err)
	}

	return resourceInventorySourceRead(ctx, d, m)
}

func resourceInventorySourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventorySourcesService
	id, diags := convertStateIDToNumeric(diagElementInventorySourceTitle, d)
	if diags.HasError() {
		return diags
	}
	if _, err := awxService.DeleteInventorySource(id); err != nil {
		return buildDiagDeleteFail(
			"inventory source",
			fmt.Sprintf("inventory source %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func resourceInventorySourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventorySourcesService
	id, diags := convertStateIDToNumeric(diagElementInventorySourceTitle, d)
	if diags.HasError() {
		return diags
	}
	res, err := awxService.GetInventorySourceByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail(diagElementInventorySourceTitle, id, err)
	}
	d = setInventorySourceResourceData(d, res)
	return nil
}

func setInventorySourceResourceData(d *schema.ResourceData, r *awx.InventorySource) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("enabled_var", r.EnabledVar)
	d.Set("enabled_value", r.EnabledValue)
	d.Set("overwrite", r.Overwrite)
	d.Set("overwrite_vars", r.OverwriteVars)
	d.Set("update_on_launch", r.UpdateOnLaunch)
	d.Set("inventory_id", r.Inventory)
	d.Set("credential_id", r.Credential)
	d.Set("source", r.Source)
	d.Set("source_vars", normalizeJsonYaml(r.SourceVars))
	d.Set("host_filter", r.HostFilter)
	d.Set("update_cache_timeout", r.UpdateCacheTimeout)
	d.Set("verbosity", r.Verbosity)
	d.Set("instance_filters", r.InstanceFilters)
	d.Set("group_by", r.GroupBy)
	d.Set("source_project_id", r.SourceProject)
	d.Set("source_path", r.SourcePath)
	// obsolete schema added so terraform doesn't break
	// these don't do anything in later versions of AWX! Update your code.
	d.Set("source_regions", r.SourceRegions)

	return d
}
