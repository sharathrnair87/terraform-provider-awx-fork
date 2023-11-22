/*
Use this resource to manage a HashiCorp Vault Signed SSH Credential in AWX/AT
For more details see [HashiCorp Vault Signed SSH](https://docs.ansible.com/automation-controller/latest/html/userguide/credential_plugins.html#hashicorp-vault-signed-ssh)

# Example Usage

```hcl

	data "awx_organization" "cybersec" {
	  name = "CyberSec"
	}

	resource "awx_credential_hashivault_signed_ssh" "hv_cyber_signed_ssh" {
	  name            = "HV Cyber Sig SSH"
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

func resourceCredentialHashiVaultSSH() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to manage a HashiCorp Vault Signed SSH Credential in AWX/AT",
		CreateContext: resourceCredentialHashiVaultSSHCreate,
		ReadContext:   resourceCredentialHashiVaultSSHRead,
		UpdateContext: resourceCredentialHashiVaultSSHUpdate,
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
			"role_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kubernetes_role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_auth_path": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceCredentialHashiVaultSSHCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error
	var credentialTypeID int

	client := m.(*awx.AWX)

	//params := make(map[string]string)
	//params["name"] = "HashiCorp Vault Signed SSH"

	credType, err := awx.GetAllPages[awx.CredentialType](client, awx.CredentialTypesAPIEndpoint, map[string]string{
		"name": "HashiCorp Vault Signed SSH",
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
		//"credential_type": 22, // Hashicorp Vault Signed SSH
		"credential_type": credentialTypeID,
		"inputs": map[string]interface{}{
			"url":               d.Get("url").(string),
			"token":             d.Get("token").(string),
			"cacert":            d.Get("cacert").(string),
			"role_id":           d.Get("role_id").(string),
			"secret_id":         d.Get("secret_id").(string),
			"namespace":         d.Get("namespace").(string),
			"kubernetes_role":   d.Get("kubernetes_role").(string),
			"default_auth_path": d.Get("default_auth_path").(string),
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
	resourceCredentialHashiVaultSSHRead(ctx, d, m)

	return diags
}

func resourceCredentialHashiVaultSSHRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("role_id", cred.Inputs["role_id"])
	d.Set("secret_id", cred.Inputs["secret_id"])
	d.Set("namespace", cred.Inputs["namespace"])
	d.Set("kubernetes_role", cred.Inputs["kubernetes_role"])
	d.Set("default_auth_path", cred.Inputs["default_auth_path"])

	return diags
}

func resourceCredentialHashiVaultSSHUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var credentialTypeID int

	client := m.(*awx.AWX)

	credType, err := awx.GetAllPages[awx.CredentialType](client, awx.CredentialTypesAPIEndpoint, map[string]string{
		"name": "HashiCorp Vault Signed SSH",
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
		"role_id",
		"secret_id",
		"namespace",
		"kubernetes_role",
		"default_auth_path",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":         d.Get("name").(string),
			"description":  d.Get("description").(string),
			"organization": d.Get("organization_id").(int),
			//"credential_type": 22, // Hashicorp Vault Signed SSH
			"credential_type": credentialTypeID,
			"inputs": map[string]interface{}{
				"url":               d.Get("url").(string),
				"token":             d.Get("token").(string),
				"cacert":            d.Get("cacert").(string),
				"role_id":           d.Get("role_id").(string),
				"secret_id":         d.Get("secret_id").(string),
				"namespace":         d.Get("namespace").(string),
				"kubernetes_role":   d.Get("kubernetes_role").(string),
				"default_auth_path": d.Get("default_auth_path").(string),
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

	return resourceCredentialHashiVaultSSHRead(ctx, d, m)
}
