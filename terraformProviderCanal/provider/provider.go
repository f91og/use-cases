package provider

import (
	"context"
	"fmt"

	"daas-terraform-provider-canal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CANAL_HOST", ""),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CANAL_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CANAL_PASSWORD", ""),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"canal_instance": dataSourceInstance(),
			"canal_cluster":  dataSourceCluster(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"canal_instance": resourceInstance(),
			"canal_cluster":  resourceCluster(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	host := d.Get("host").(string)
	if host == "" {
		return nil, diag.FromErr(fmt.Errorf("host cannot be empty"))
	}
	username := d.Get("username").(string)
	if username == "" {
		return nil, diag.FromErr(fmt.Errorf("username cannot be empty"))
	}
	password := d.Get("password").(string)
	if password == "" {
		return nil, diag.FromErr(fmt.Errorf("password cannot be empty"))
	}

	canalAPIClient, err := client.NewClient(host, username, password)
	if err != nil {
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create canal api client with username and password",
				Detail:   err.Error(),
			},
		}
	}
	return canalAPIClient, diag.Diagnostics{}
}
