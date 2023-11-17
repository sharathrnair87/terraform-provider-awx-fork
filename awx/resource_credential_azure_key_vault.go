/*
Use this resource to manage an Azure Key Vault Credential in AWX/AT.
for more details see [Microsoft Azure Key Vault](https://docs.ansible.com/automation-controller/latest/html/userguide/credential_plugins.html#microsoft-azure-key-vault)

# Example Usage

```hcl

	data "awx_organization" "infra" {
	  name = "Infrastructure"
	}

	resource "awx_credential_azure_key_vault" "kv" {
	  name            = "Infra KV"
	  description     = "Azure KV for Infra project"
	  organization_id = data.awx_organization.infra.id
	  url             = "https://infra-vault-example.vault.azure.net"
	  client          = var.azrm_client_id
	  secret          = var.azrm_client_secret
	  tenant          = var.azrm_tenant_id
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

func resourceCredentialAzureKeyVault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialAzureKeyVaultCreate,
		ReadContext:   resourceCredentialAzureKeyVaultRead,
		UpdateContext: resourceCredentialAzureKeyVaultUpdate,
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
			"client": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"tenant": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceCredentialAzureKeyVaultCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error
	var credentialTypeID int

	client := m.(*awx.AWX)

	credType, err := client.CredentialTypeService.GetCredentialTypeByName(map[string]string{
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

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		//"credential_type": 19, // Azure Key Vault
		"credential_type": credentialTypeID,
		"inputs": map[string]interface{}{
			"url":    d.Get("url").(string),
			"client": d.Get("client").(string),
			"secret": d.Get("secret").(string),
			"tenant": d.Get("tenant").(string),
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
	resourceCredentialAzureKeyVaultRead(ctx, d, m)

	return diags
}

func resourceCredentialAzureKeyVaultRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("client", cred.Inputs["client"])
	d.Set("secret", d.Get("secret").(string))
	d.Set("tenant", cred.Inputs["tenant"])

	return diags
}

func resourceCredentialAzureKeyVaultUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var credentialTypeID int

	client := m.(*awx.AWX)

	credType, err := client.CredentialTypeService.GetCredentialTypeByName(map[string]string{
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
		"client",
		//"secret",
		"tenant",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			//"credential_type": 19, // Azure Key Vault
			"credential_type": credentialTypeID,
			"inputs": map[string]interface{}{
				"url":    d.Get("url").(string),
				"client": d.Get("client").(string),
				"secret": d.Get("secret").(string),
				"tenant": d.Get("tenant").(string),
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

	return resourceCredentialAzureKeyVaultRead(ctx, d, m)
}
