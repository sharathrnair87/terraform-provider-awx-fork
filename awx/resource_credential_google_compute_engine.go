/*
Use this resource to manage a Google Compute Engine credential in AWX/AT.
For more details see [Google Compute Engine](https://docs.ansible.com/automation-controller/latest/html/userguide/credentials.html#google-compute-engine)

# Example Usage

```hcl

	data "awx_organization" "gcp_iaas" {
	  name = "Google Infra"
	}

	resource "awx_credential_google_compute_engine" "gce" {
	  name            = "GCE Credential"
	  organization_id = data.awx_organization.gcp_iaas.id
	  username        = "svc_acccount@gcp-prj.iam.gserviceaccount.com"
	  ssh_key_data    = <<-EOT
	    -----BEGIN RSA PRIVATE KEY-----
	    ...
	    EOT
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

func resourceCredentialGoogleComputeEngine() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to manage a Google Compute Engine credential in AWX/AT.",
		CreateContext: resourceCredentialGoogleComputeEngineCreate,
		ReadContext:   resourceCredentialGoogleComputeEngineRead,
		UpdateContext: resourceCredentialGoogleComputeEngineUpdate,
		DeleteContext: CredentialsServiceDeleteByID,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_key_data": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceCredentialGoogleComputeEngineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error
	var credentialTypeID int

	client := m.(*awx.AWX)

	credType, err := awx.GetAllPages[awx.CredentialType](client, awx.CredentialTypesAPIEndpoint, map[string]string{
		"name": "Google Compute Engine",
	})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to find Credential Type",
			Detail:   fmt.Sprintf("Unable to find Credential Type: %s", err.Error()),
		})
		return diags
	}

	credentialTypeID = credType[0].ID

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": credentialTypeID,
		"inputs": map[string]interface{}{
			"username":     d.Get("username").(string),
			"project":      d.Get("project").(string),
			"ssh_key_data": d.Get("ssh_key_data").(string),
		},
	}

	cred, err := client.CredentialsService.CreateCredentials(newCredential, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new credentials",
			Detail:   fmt.Sprintf("Unable to create new credentials: %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(cred.ID))
	resourceCredentialGoogleComputeEngineRead(ctx, d, m)

	return diags
}

func resourceCredentialGoogleComputeEngineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := strconv.Atoi(d.Id())
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
	d.Set("project", cred.Inputs["project"])

	return diags
}

func resourceCredentialGoogleComputeEngineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var credentialTypeID int

	client := m.(*awx.AWX)

	credType, err := awx.GetAllPages[awx.CredentialType](client, awx.CredentialTypesAPIEndpoint, map[string]string{
		"name": "Google Compute Engine",
	})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to find Credential Type",
			Detail:   fmt.Sprintf("Unable to find Credential Type: %s", err.Error()),
		})
		return diags
	}

	credentialTypeID = credType[0].ID

	keys := []string{
		"name",
		"description",
		"username",
		"project",

		"ssh_key_data",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": credentialTypeID, // Google Compute Engine
			"inputs": map[string]interface{}{
				"username":     d.Get("username").(string),
				"project":      d.Get("project").(string),
				"ssh_key_data": d.Get("ssh_key_data").(string),
			},
		}

		_, err = client.CredentialsService.UpdateCredentialsByID(id, updatedCredential, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update existing credentials",
				Detail:   fmt.Sprintf("Unable to update existing credentials with id %d: %s", id, err.Error()),
			})
			return diags
		}
	}

	return resourceCredentialGoogleComputeEngineRead(ctx, d, m)
}
