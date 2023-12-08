package awx

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWX_HOSTNAME", "http://localhost"),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable SSL verification of API calls",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWX_USERNAME", "admin"),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AWX_PASSWORD", "password"),
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AWX_TOKEN", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"awx_credential_azure_key_vault":                          resourceCredentialAzureKeyVault(),
			"awx_credential_azure_resource_manager":                   resourceCredentialAzureRM(),
			"awx_credential_google_compute_engine":                    resourceCredentialGoogleComputeEngine(),
			"awx_credential_input_source":                             resourceCredentialInputSource(),
			"awx_credential":                                          resourceCredential(),
			"awx_credential_type":                                     resourceCredentialType(),
			"awx_credential_machine":                                  resourceCredentialMachine(),
			"awx_credential_scm":                                      resourceCredentialSCM(),
			"awx_credential_galaxy":                                   resourceCredentialGalaxy(),
			"awx_credential_github_token":                             resourceCredentialGithubPAT(),
			"awx_credential_hashivault_secret":                        resourceCredentialHashiVaultSecret(),
			"awx_credential_hashivault_signed_ssh":                    resourceCredentialHashiVaultSSH(),
			"awx_credential_vault":                                    resourceCredentialVault(),
			"awx_execution_environment":                               resourceExecutionEnvironment(),
			"awx_host":                                                resourceHost(),
			"awx_instance_group":                                      resourceInstanceGroup(),
			"awx_inventory_group":                                     resourceInventoryGroup(),
			"awx_inventory_source":                                    resourceInventorySource(),
			"awx_inventory":                                           resourceInventory(),
			"awx_job_template_credential":                             resourceJobTemplateCredentials(),
			"awx_job_template":                                        resourceJobTemplate(),
			"awx_job_template_launch":                                 resourceJobTemplateLaunch(),
			"awx_job_template_notification_template_error":            resourceJobTemplateNotificationTemplateError(),
			"awx_job_template_notification_template_started":          resourceJobTemplateNotificationTemplateStarted(),
			"awx_job_template_notification_template_success":          resourceJobTemplateNotificationTemplateSuccess(),
			"awx_notification_template":                               resourceNotificationTemplate(),
			"awx_organization":                                        resourceOrganization(),
			"awx_organization_galaxy_credential":                      resourceOrganizationsGalaxyCredentials(),
			"awx_project":                                             resourceProject(),
			"awx_schedule":                                            resourceSchedule(),
			"awx_settings_ldap_team_map":                              resourceSettingsLDAPTeamMap(),
			"awx_settings_saml_team_map":                              resourceSettingsSAMLTeamMap(),
			"awx_settings_saml_organization_map":                      resourceSettingsSAMLOrganizationMap(),
			"awx_settings_saml_team_attributes":                       resourceSettingsSAMLTeamAttrMap(),
			"awx_setting":                                             resourceSetting(),
			"awx_team":                                                resourceTeam(),
			"awx_user":                                                resourceUser(),
			"awx_workflow_job_template_node_always":                   resourceWorkflowJobTemplateNodeAlways(),
			"awx_workflow_job_template_node_failure":                  resourceWorkflowJobTemplateNodeFailure(),
			"awx_workflow_job_template_node_success":                  resourceWorkflowJobTemplateNodeSuccess(),
			"awx_workflow_job_template_node":                          resourceWorkflowJobTemplateNode(),
			"awx_workflow_job_template":                               resourceWorkflowJobTemplate(),
			"awx_workflow_job_template_schedule":                      resourceWorkflowJobTemplateSchedule(),
			"awx_workflow_job_template_notification_template_error":   resourceWorkflowJobTemplateNotificationTemplateError(),
			"awx_workflow_job_template_notification_template_started": resourceWorkflowJobTemplateNotificationTemplateStarted(),
			"awx_workflow_job_template_notification_template_success": resourceWorkflowJobTemplateNotificationTemplateSuccess(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"awx_credential_azure_key_vault":        dataSourceCredentialAzure(),
			"awx_credential_azure_resource_manager": dataSourceCredentialAzureRM(),
			"awx_credential_machine":                dataSourceCredentialMachine(),
			"awx_credential_scm":                    dataSourceCredentialSCM(),
			"awx_credential_github_token":           dataSourceCredentialGithubPAT(),
			"awx_credential_hashivault_secret":      dataSourceCredentialHashiVault(),
			"awx_credential_hashivault_signed_ssh":  dataSourceCredentialHashiVaultSSH(),
			"awx_credential_vault":                  dataSourceCredentialVault(),
			"awx_credential":                        dataSourceCredentialByID(),
			"awx_credential_role":                   dataSourceCredentialRole(),
			"awx_credential_type":                   dataSourceCredentialTypeByID(),
			"awx_credentials":                       dataSourceCredentials(),
			"awx_execution_environment":             dataSourceExecutionEnvironment(),
			"awx_host":                              dataSourceHost(),
			"awx_inventory_group":                   dataSourceInventoryGroup(),
			"awx_inventory":                         dataSourceInventory(),
			"awx_inventory_role":                    dataSourceInventoryRole(),
			"awx_inventory_source":                  dataSourceInventorySource(),
			"awx_job_template":                      dataSourceJobTemplate(),
			"awx_job_template_role":                 dataSourceJobTemplateRole(),
			"awx_notification_template":             dataSourceNotificationTemplate(),
			"awx_organization":                      dataSourceOrganization(),
			"awx_organization_role":                 dataSourceOrganizationRole(),
			"awx_organizations":                     dataSourceOrganizations(),
			"awx_project":                           dataSourceProject(),
			"awx_project_role":                      dataSourceProjectRole(),
			"awx_schedule":                          dataSourceSchedule(),
			"awx_workflow_job_template":             dataSourceWorkflowJobTemplate(),
			"awx_team":                              dataSourceTeam(),
			"awx_user":                              dataSourceUser(),
			//TODO
			// "awx_setting": dataSourceSetting(),
			// "awx_credential_galaxy": dataSourceGalaxy(),
			// "awx_organization_galaxy_credential": dataSourceOrganizationGalaxyCredentials(),

		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	hostname := d.Get("hostname").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	token := d.Get("token").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := http.DefaultClient
	if d.Get("insecure").(bool) {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client.Transport = customTransport
	}

	var c *awx.AWX
	var err error
	if token != "" {
		c, err = awx.NewAWXToken(hostname, token, client)
	} else {
		c, err = awx.NewAWX(hostname, username, password, client)
	}
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create AWX client",
			Detail:   "Unable to auth user against AWX API: check the hostname, username and password",
		})
		return nil, diags
	}

	return c, diags
}
