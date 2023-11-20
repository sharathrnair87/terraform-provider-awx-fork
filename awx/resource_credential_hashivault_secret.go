/*
Use this resource to manage a HashiCorp Vault Secret Credential in AWX/AT.
For more details see [HashiCorp Vault Secret Lookup](https://docs.ansible.com/automation-controller/latest/html/userguide/credential_plugins.html#ug-credentials-hashivault)

# Example Usage

```hcl

	data "awx_organization" "cybersec" {
	  name = "CyberSec"
	}

	resource "awx_credential_hashivault_secret" "hv_cyber" {
	  name            = "HV Cyber"
	  organization_id = data.awx_organization.cybersec.id
	  url             = var.hashicorp_vault_url
	  token           = var.hashicorp_vault_token
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

func resourceCredentialHashiVaultSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialHashiVaultSecretCreate,
		ReadContext:   resourceCredentialHashiVaultSecretRead,
		UpdateContext: resourceCredentialHashiVaultSecretUpdate,
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
			"token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"cacert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "v1",
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

func resourceCredentialHashiVaultSecretCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error
	var credentialTypeID int

	client := m.(*awx.AWX)

	//credType, err := client.CredentialTypeService.GetCredentialTypeByName(map[string]string{
	//	"name": "HashiCorp Vault Secret Lookup",
	//})

	credType, err := awx.GetAllPages[awx.CredentialType](client, awx.CredentialTypesAPIEndpoint, map[string]string{
		"name": "HashiCorp Vault Secret Lookup",
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
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"organization": d.Get("organization_id").(int),
		//"credential_type": 21, // Hashicorp Vault Secret Lookup
		"credential_type": credentialTypeID, // Hashicorp Vault Secret Lookup
		"inputs": map[string]interface{}{
			"url":         d.Get("url").(string),
			"token":       d.Get("token").(string),
			"cacert":      d.Get("cacert").(string),
			"api_version": d.Get("api_version").(string),
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
	resourceCredentialHashiVaultSecretRead(ctx, d, m)

	return diags
}

func resourceCredentialHashiVaultSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("url", cred.Inputs["url"])
	d.Set("token", d.Get("token").(string))
	d.Set("cacert", cred.Inputs["cacert"])
	d.Set("api_version", cred.Inputs["api_version"])

	return diags
}

func resourceCredentialHashiVaultSecretUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var credentialTypeID int

	client := m.(*awx.AWX)

	credType, err := awx.GetAllPages[awx.CredentialType](client, awx.CredentialTypesAPIEndpoint, map[string]string{
		"name": "HashiCorp Vault Secret Lookup",
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
		"token",
		"cacert",
		"api_version",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":         d.Get("name").(string),
			"description":  d.Get("description").(string),
			"organization": d.Get("organization_id").(int),
			//"credential_type": 21, // Hashicorp Vault Secret Lookup
			"credential_type": credentialTypeID, // Hashicorp Vault Secret Lookup
			"inputs": map[string]interface{}{
				"url":         d.Get("url").(string),
				"token":       d.Get("token").(string),
				"cacert":      d.Get("cacert").(string),
				"api_version": d.Get("api_version").(string),
			},
		}

		client := m.(*awx.AWX)
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

	return resourceCredentialHashiVaultSecretRead(ctx, d, m)
}
