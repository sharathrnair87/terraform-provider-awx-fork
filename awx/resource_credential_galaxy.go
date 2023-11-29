/*
Use this resource to create an Ansible Galaxy Credential in AWX/AT

# Example Usage
```hcl

	data "awx_organization" "devops" {
	  name = "DevOps"
	}

	resource "awx_credential_galaxy" "devops_galaxy_cred_01" {
	  name            = "DevOps_Cred_01"
	  description     = "Galaxy Cred for DevOps Org"
	  organization_id = data.awx_organization.devops
	  url             = "https://galaxy.ansible.com"
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

func resourceCredentialGalaxy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialGalaxyCreate,
		ReadContext:   resourceCredentialGalaxyRead,
		UpdateContext: resourceCredentialGalaxyUpdate,
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
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCredentialGalaxyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error
	var credentialTypeID int

	client := m.(*awx.AWX)

	credType, err := awx.GetAllPages[awx.CredentialType](client, awx.CredentialTypesAPIEndpoint, map[string]string{
		"name": "Ansible Galaxy/Automation Hub API Token",
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
			"url":      d.Get("url").(string),
			"auth_url": d.Get("auth_url").(string),
			"token":    d.Get("token").(string),
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
	resourceCredentialGalaxyRead(ctx, d, m)

	return diags
}

func resourceCredentialGalaxyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("url", cred.Inputs["url"])
	d.Set("auth_url", cred.Inputs["auth_url"])
	d.Set("token", d.Get("token").(string))
	d.Set("organization_id", cred.OrganizationID)

	return diags
}

func resourceCredentialGalaxyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var credentialTypeID int

	client := m.(*awx.AWX)

	credType, err := awx.GetAllPages[awx.CredentialType](client, awx.CredentialTypesAPIEndpoint, map[string]string{
		"name": "Microsoft Azure Key Vault",
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
		"url",
		"auth_url",
		"token",
		"organization_id",
		"team_id",
		"owner_id",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": credentialTypeID,
			"inputs": map[string]interface{}{
				"url":      d.Get("url").(string),
				"auth_url": d.Get("auth_url").(string),
				"token":    d.Get("token").(string),
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

	return resourceCredentialGalaxyRead(ctx, d, m)
}
