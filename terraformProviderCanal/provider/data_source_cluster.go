package provider

import (
	"context"

	"daas-terraform-provider-canal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for canal cluster",
		ReadContext: dataSourceClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceClusterRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	canalAPIClient := meta.(*client.Client)

	clusterName := resourceData.Get("name").(string)
	cluster, err := canalAPIClient.GetClusterByName(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceData.SetId(string(cluster.ID))
	resourceData.Set("name", cluster.Name)

	return diag.Diagnostics{}
}
