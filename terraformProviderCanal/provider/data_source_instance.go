package provider

import (
	"context"
	"strings"

	"daas-terraform-provider-canal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for canal instance",
		ReadContext: dataSourceInstanceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"canal_cluster": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_server": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"tables": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:     schema.TypeString,
										Required: true,
									},
									"table": {
										Type:     schema.TypeString,
										Required: true,
									},
									"pk": {
										Type:     schema.TypeString,
										Required: true,
									},
									"topic": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_server_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"running_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	var diags diag.Diagnostics

	instanceID := d.Get("id").(string)

	instance, err := c.GetInstance(instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if instance == nil {
		d.SetId("")
		return diags
	}

	d.Set("name", instance.Name)
	d.Set("cluster_server_id", instance.ClusterServerId)
	d.Set("cluster_id", instance.ClusterID)
	d.Set("modified_time", instance.ModifiedTime)
	d.Set("node_server", instance.NodeServer)
	d.Set("running_status", instance.RunningStatus)
	d.Set("status", instance.Status)
	d.Set("server_id", instance.ServerId)
	d.Set("canal_cluster", instance.Cluster)

	ctt := strings.Split(instance.Content, `\n`)
	configs := make(map[string]string)

	for _, s := range ctt {
		cs := strings.Split(s, "=")
		if len(cs) == 2 {
			configs[cs[0]] = cs[1]
		}
	}

	if value, ok := configs["canal.instance.master.journal.name"]; ok {
		if err := d.Set("binlog_name", value); err != nil {
			return diag.FromErr(err)
		}
	}
	if value, ok := configs["canal.instance.master.position"]; ok {
		if err := d.Set("master_address", value); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("db_credential_id", configs["canal.instance.credential"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("partitions_num", configs["canal.mq.partitionsNum"]); err != nil {
		return diag.FromErr(err)
	}

	// partitionHash, linepay_finnet_db1\\.confirm_refund_request:id,linepay_finnet_db1\\.order_request:id,
	// dynamic topic, linepay_finnet_db1_confirm_refund_request:linepay_finnet_db1\\.confirm_refund_request,linepay_finnet_db1_order_request:linepay_finnet_db1\\.order_request,
	partitionHash := strings.Split(configs["canal.mq.partitionHash"], ",")
	dynamicTopic := strings.Split(configs["canal.mq.dynamicTopic"], ",")

	dbTablePk := make(map[string]string)
	dbTableTopic := make(map[string]string)

	for _, dp := range partitionHash {
		dps := strings.Split(dp, ":")
		dbTablePk[dps[0]] = dps[1]
	}
	for _, dt := range dynamicTopic {
		dts := strings.Split(dt, ",")
		dbTableTopic[dts[1]] = dts[0]
	}

	tables := make([]interface{}, 1)
	for dbTable, pk := range dbTablePk {
		t := make(map[string]interface{})
		topic := dbTableTopic[dbTable]
		dts := strings.Split(dbTable, `\\.`)
		t["database"] = dts[0]
		t["table"] = dts[1]
		t["pk"] = pk
		t["topic"] = topic
		tables = append(tables, t)
	}

	if err := d.Set("tables", tables); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceID)

	return diags
}
