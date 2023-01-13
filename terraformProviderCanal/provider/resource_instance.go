package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"daas-terraform-provider-canal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var keyPairs = []Pair{
	{"canal.instance.gtidon", "instance_gtid_on"},
	{"canal.instance.master.journal.name", "instance_master_journal_name"},
	{"canal.instance.master.position", "instance_master_position"},
	{"canal.instance.master.timestamp", "instance_master_timestamp"},
	{"canal.instance.master.gtid", "instance_master_gtid"},
	{"canal.instance.rds.accesskey", "instance_rds_access_key"},
	{"canal.instance.rds.secretkey", "instance_rds_secret_key"},
	{"canal.instance.rds.instanceId", "instance_rds_instance_id"},
	{"canal.instance.tsdb.enable", "instance_tsdb_enable"},
	{"canal.instance.connectionCharset", "instance_connection_charset"},
	{"canal.instance.enableDruid", "instance_enable_druid"},
	{"canal.mq.partitionsNum", "mq_partitions_num"},
	{"canal.instance.credential", "db_credential_id"},
}

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Resource data for canal instance",

		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"instance_gtid_on": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_master_journal_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"instance_master_position": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"instance_master_timestamp": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"instance_master_gtid": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"instance_rds_access_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"instance_rds_secret_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"instance_rds_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"instance_tsdb_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"instance_connection_charset": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UTF-8",
			},
			"instance_enable_druid": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"db_credential_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mq_partitions_num": {
				Type:     schema.TypeInt,
				Required: true,
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
									"database_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"table_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"primary_key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"mq_topic": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("resourceInstanceCreate")
	c := m.(*client.Client)

	instanceId, err := c.CreateInstance(assembleInstanceConfig(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceId)

	return diag.Diagnostics{}
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("resourceInstanceRead")
	c := m.(*client.Client)

	instance, err := c.GetInstance(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	disperseInstanceConfig(d, instance)

	return diag.Diagnostics{}
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	instanceConfig := assembleInstanceConfig(d)
	instanceConfig.InstanceName += "/instance.propertios" // Need to add this only when update
	err := c.UpdateInstance(instanceConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	err := c.DeleteInstance(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func assembleTableConfig(d *schema.ResourceData) (string, string, string) {
	tablesBlock := d.Get("tables").([]interface{})[0].(map[string]interface{})
	tableResources := tablesBlock["table"].([]interface{})
	instanceFilterRegex := make([]string, len(tableResources))
	mqDynamicTopic := make([]string, len(tableResources))
	mqPartitionHash := make([]string, len(tableResources))
	for i, t := range tableResources {
		table := t.(map[string]interface{})
		databaseName := table["database_name"].(string)
		tableName := table["table_name"].(string)
		primaryKey := table["primary_key"].(string)
		mqTopic := table["mq_topic"].(string)
		instanceFilterRegex[i] = fmt.Sprintf(`%s\\.%s`, databaseName, tableName)
		mqDynamicTopic[i] = fmt.Sprintf(`%s:%s\\.%s`, mqTopic, databaseName, tableName)
		mqPartitionHash[i] = fmt.Sprintf(`%s\\.%s:%s`, databaseName, tableName, primaryKey)
	}
	return strings.Join(instanceFilterRegex, ","), strings.Join(mqDynamicTopic, ","), strings.Join(mqPartitionHash, ",")
}

func disperseTableConfig(instanceFilterRegex, mqDynamicTopic, mqPartitionHash string) []interface{} {
	dbTables := strings.Split(instanceFilterRegex, ",")
	tableResources := make([]map[string]string, len(dbTables))
	dbTable2TableResource := make(map[string]map[string]string)
	for _, dbTable := range dbTables {
		elements := strings.Split(dbTable, `\\.`)
		databaseName := elements[0]
		tableName := elements[1]
		dbTable2TableResource[dbTable] = map[string]string{
			"database_name": databaseName,
			"table_name":    tableName,
		}
	}
	for _, topicDbTable := range strings.Split(mqDynamicTopic, ",") {
		elements := strings.Split(topicDbTable, `:`)
		mqTopic := elements[0]
		dbTable := elements[1]
		dbTable2TableResource[dbTable]["mq_topic"] = mqTopic
	}
	for _, dbTablePk := range strings.Split(mqPartitionHash, ",") {
		elements := strings.Split(dbTablePk, `:`)
		dbTable := elements[0]
		primaryKey := elements[1]
		dbTable2TableResource[dbTable]["primary_key"] = primaryKey
	}

	for i, dbTable := range dbTables {
		tableResources[i] = dbTable2TableResource[dbTable]
	}

	tablesBlock := map[string]interface{}{
		"table": tableResources,
	}
	return []interface{}{tablesBlock}
}

// From resources to canal config
func assembleInstanceConfig(d *schema.ResourceData) client.InstanceConfig {
	log.Println("assembleInstanceConfig")
	instanceFilterRegex, mqDynamicTopic, mqPartitionHash := assembleTableConfig(d)

	var configs []string
	for _, pair := range keyPairs {
		configs = append(configs, fmt.Sprintf("%s=%v", pair.canalKey, d.Get(pair.tfKey)))
	}

	configs = append(configs, fmt.Sprintf("canal.instance.filter.regex=%v", instanceFilterRegex))
	configs = append(configs, "canal.mq.topic=example")
	configs = append(configs, fmt.Sprintf("canal.mq.dynamicTopic=%v", mqDynamicTopic))
	configs = append(configs, fmt.Sprintf("canal.mq.partitionHash=%v", mqPartitionHash))
	configs = append(configs, "canal.instance.filter.black.regex=")

	instanceName := d.Get("name").(string)
	clusterId := fmt.Sprintf("cluster:%d", d.Get("cluster_id").(int))

	status := "0"
	if d.Get("enabled").(bool) {
		status = "1"
	}
	instanceConfig := client.InstanceConfig{
		InstanceName: instanceName,
		ClusterId:    clusterId,
		Content:      strings.Join(configs, "\n"),
		Status:       status,
	}
	if d.Id() != "" {
		log.Printf("assembleInstanceConfig, setId: %+v", d.Id())
		instanceConfig.ID = json.Number(d.Id())
	}

	return instanceConfig
}

// From canal config to resources
func disperseInstanceConfig(d *schema.ResourceData, instance *client.Instance) {
	if instance.ID == "" {
		d.SetId("")
		return
	}

	log.Println("disperseInstanceConfig")
	configs := make(map[string]string)
	lines := strings.Split(instance.Content, "\n")
	for _, line := range lines {
		elements := strings.Split(line, "=")
		configs[elements[0]] = strings.Join(elements[1:], "=")
	}
	log.Printf("%+v", configs)

	for _, pair := range keyPairs {
		setValue(d, configs, pair.tfKey, pair.canalKey)
	}
	d.Set("tables", disperseTableConfig(configs["canal.instance.filter.regex"], configs["canal.mq.dynamicTopic"], configs["canal.mq.partitionHash"]))
	d.Set("name", instance.Name)
	d.Set("cluster_id", instance.ClusterID)
	if instance.Status == "1" {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}
}
