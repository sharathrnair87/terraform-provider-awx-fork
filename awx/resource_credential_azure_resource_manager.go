/*
*TBD*

# Example Usage

```hcl
*TBD*
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

func resourceCredentialAzureRM() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialAzureRMCreate,
		ReadContext:   resourceCredentialAzureRMRead,
		UpdateContext: resourceCredentialAzureRMUpdate,
		DeleteContext: CredentialsServiceDeleteByID,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subscription": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant": {
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
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"username", "client"},
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"username"},
			},
			"client": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"username", "client"},
			},
			"secret": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"client"},
			},
			"cloud_environment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceCredentialAzureRMCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	inputs := make(map[string]interface{})

	inputs["subscription"] = d.Get("subscription").(string)
	inputs["tenant"] = d.Get("tenant").(string)
	if username, userOk := d.GetOk("username"); userOk {
		inputs["username"] = username
		inputs["password"] = d.Get("password")
	}

	if client, clientOk := d.GetOk("client"); clientOk {
		inputs["client"] = client
		inputs["secret"] = d.Get("secret")
	}

	if cloud_environment, cloudOk := d.GetOk("cloud_environment"); cloudOk {
		inputs["cloud_environment"] = cloud_environment
	}

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": 11, // Azure Resource Manager
		"inputs":          inputs,
	}

	client := m.(*awx.AWX)
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
	resourceCredentialAzureRMRead(ctx, d, m)

	return diags
}

func resourceCredentialAzureRMRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("subscription", cred.Inputs["subscription"])
	d.Set("tenant", cred.Inputs["tenant"])
	d.Set("client", cred.Inputs["client"])
	d.Set("secret", d.Get("secret").(string))
	d.Set("username", cred.Inputs["username"])
	d.Set("password", d.Get("password").(string))

	return diags
}

func resourceCredentialAzureRMUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"subscription",
		"client",
		"username",
		"tenant",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": 11, // Azure Resource Manager
			"inputs": map[string]interface{}{
				"url":      d.Get("url").(string),
				"client":   d.Get("client").(string),
				"secret":   d.Get("secret").(string),
				"tenant":   d.Get("tenant").(string),
				"username": d.Get("username").(string),
				"password": d.Get("password").(string),
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

	return resourceCredentialAzureRMRead(ctx, d, m)
}
