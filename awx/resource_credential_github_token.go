/*
Use this resource to manage a Github Token Credential in AWX/AT.
For more details see [Github Token](https://docs.ansible.com/automation-controller/latest/html/userguide/credentials.html#github-personal-access-token)

# Example Usage

```hcl

	data "awx_organization" "devops" {
	  name = "DevOpsOrg"
	}

	resource "awx_credential_github_token" "gh_pat" {
	  name            = "devops_gh_pat"
	  token           = "..."
	  organization_id = data.awx_organization.devops.id
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

func resourceCredentialGithubPAT() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialGithubPATCreate,
		ReadContext:   resourceCredentialGithubPATRead,
		UpdateContext: resourceCredentialGithubPATUpdate,
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
			"token": {
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

func resourceCredentialGithubPATCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error
	var credentialTypeID int

	client := m.(*awx.AWX)

	credType, err := client.CredentialTypeService.GetCredentialTypeByName(map[string]string{
		"name": "GitHub Personal Access Token",
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
		//"credential_type": 12, // Github PAT
		"credential_type": credentialTypeID,
		"inputs": map[string]interface{}{
			"token": d.Get("token").(string),
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
	resourceCredentialGithubPATRead(ctx, d, m)

	return diags
}

func resourceCredentialGithubPATRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("token", d.Get("token").(string))
	d.Set("organization_id", cred.OrganizationID)

	return diags
}

func resourceCredentialGithubPATUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var credentialTypeID int

	client := m.(*awx.AWX)

	credType, err := client.CredentialTypeService.GetCredentialTypeByName(map[string]string{
		"name": "GitHub Personal Access Token",
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
		"organization_id",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":         d.Get("name").(string),
			"description":  d.Get("description").(string),
			"organization": d.Get("organization_id").(int),
			//"credential_type": 12, // Github PAT
			"credential_type": credentialTypeID, // Github PAT
			"inputs": map[string]interface{}{
				"token": d.Get("token").(string),
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

	return resourceCredentialGithubPATRead(ctx, d, m)
}
