package provider

import (
	"context"
	"daas-terraform-provider-canal/provider/client"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Resource data for canal cluster",

		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zk_hosts": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  11111,
			},
			"metrics_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  11112,
			},
			"admin_host": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "canal-admin-service",
			},
			"admin_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  11110,
			},
			"admin_user": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "admin",
			},
			"admin_password": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "4ACFE3202A5FF5CF467898FC58AAB1D615029441",
			},
			"admin_register_auto": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"zk_flush_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},
			"instance_without_netty": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_file_data_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/home/deployer/data/canal-deployer",
			},
			"instance_file_flush_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10000,
			},
			"instance_memory_buffer_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  16384,
			},
			"instance_memory_buffer_memunit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1024,
			},
			"instance_memory_batch_mode": { // MEMSIZE, ITEMSIZE
				Type:     schema.TypeString,
				Optional: true,
				Default:  "MEMSIZE",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !contains([]string{"MEMSIZE", "ITEMSIZE"}, val.(string)) {
						errs = append(errs, fmt.Errorf("%q must be eiather: MEMSIZE, ITEMSIZE ", key))
					}
					return warns, errs
				},
			},
			"instance_memory_raw_entry": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"instance_detecting_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"instance_detecting_sql": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "select 1",
			},
			"instance_detecting_interval_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"instance_detecting_retry_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"instance_detecting_heartbeat_ha_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_transaction_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1024,
			},
			"instance_fallback_interval_in_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"instance_network_receive_buffer_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  16384,
			},
			"instance_network_send_buffer_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  16384,
			},
			"instance_network_so_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"instance_filter_druid_ddl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"instance_filter_query_dcl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_filter_query_dml": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_filter_query_ddl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_filter_table_error": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_filter_rows": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_filter_transaction_entry": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_binlog_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ROW,STATEMENT,MIXED",
			},
			"instance_binlog_image": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "FULL,MINIMAL,NOBLOB",
			},
			"instance_get_ddl_isolation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_parser_parallel": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"instance_parser_parallel_buffer_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  256,
			},
			"instance_tsdb_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"instance_tsdb_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "${canal.file.data.dir:../conf}/${canal.instance.destination:}",
			},
			"instance_tsdb_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "jdbc:h2:${canal.instance.tsdb.dir}/h2;CACHE_SIZE=1000;MODE=MYSQL;",
			},
			"instance_tsdb_db_username": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "canal",
			},
			"instance_tsdb_db_password": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "canal",
			},
			"instance_tsdb_snapshot_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  24,
			},
			"instance_tsdb_snapshot_expire": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  360,
			},
			"destination_destinations": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "example",
			},
			"destination_conf_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "../conf",
			},
			"destination_auto_scan": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"destination_auto_scan_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"instance_tsdb_spring_xml": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "classpath:spring/tsdb/h2-tsdb.xml",
			},
			"instance_global_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "spring",
			},
			"instance_global_lazy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_global_manager_address": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "${canal.admin.manager}",
			},
			"instance_global_spring_xml": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "classpath:spring/file-instance.xml",
			},
			"mq_servers": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "kafka.bootstrap.servers",
			},
			"mq_acks": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "all",
				Description: "kafka.acks",
			},
			"mq_compression_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "lz4",
				Description: "kafka.compression.type",
			},
			"mq_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     15,
				Description: "kafka.retries",
			},
			"mq_batch_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5248800,
				Description: "kafka.batch.size",
			},
			"mq_linger_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     50,
				Description: "kafka.linger.ms",
			},
			"mq_max_request_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5242880,
				Description: "kafka.max.request.size, default 5120KB",
			},
			"mq_buffer_memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     52488000,
				Description: "kafka.buffer.memory 128MB, (max_batch_size * max_request_size)",
			},
			"mq_kerberos_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "kafka.kerberos.enable",
			},
			"mq_kerberos_krb5_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "\"../conf/kerberos/krb5.conf\"",
				Description: "kafka.kerberos.krb5.file, java.security.krb5.conf",
			},
			"mq_kerberos_jaas_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "\"../conf/kerberos/jaas.conf\"",
				Description: "kafka.kerberos.jaas.file, java.security.auth.loging.config",
			},
			"mq_canal_batch_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1050,
				Description: "batch size to get message from canal server",
			},
			"mq_canal_get_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "timeout threshold to get message from canal server",
			},
			"mq_flat_message": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "false: protobuf, true: json",
			},
			"mq_access_channel": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "local",
			},
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("resourceClusterCreate")
	c := m.(*client.Client)

	clusterId, err := c.CreateCluster(assembleClusterAndConfig(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterId)

	return resourceClusterRead(ctx, d, m)
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("resourceClusterRead")
	c := m.(*client.Client)

	var diags diag.Diagnostics

	clusterAndConfig, err := c.GetCluster(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = assembleResourceData(d, clusterAndConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	err := c.UpdateCluster(assembleClusterAndConfig(d))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceClusterRead(ctx, d, m)
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	err := c.DeleteCluster(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func assembleClusterAndConfig(d *schema.ResourceData) client.ClusterAndConfig {
	log.Println("assembleClusterAndConfig")

	cluster := client.Cluster{
		Name:    d.Get("name").(string),
		ZkHosts: d.Get("zk_hosts").(string),
	}

	if d.Id() != "" {
		log.Printf("assembleClusterAndConfig, setId: %+v", d.Id())
		cluster.ID = json.Number(d.Id())
	}

	clusterAndConfig := client.ClusterAndConfig{
		CanalCluster: cluster,
		CanalConfig:  assembleCanalConfig(d),
	}

	return clusterAndConfig
}

func assembleResourceData(d *schema.ResourceData, clusterAndConfig *client.ClusterAndConfig) error {
	log.Println("assembleResourceData")

	d.Set("name", clusterAndConfig.CanalCluster.Name)
	d.Set("zk_hosts", clusterAndConfig.CanalCluster.ZkHosts)

	assembleCanalConfigResourceData(d, &clusterAndConfig.CanalConfig)

	return nil
}

// CanalConfig
func assembleCanalConfig(d *schema.ResourceData) client.Config {
	log.Println("assembleCanalConfig")

	// build canal configs from tfstate
	var configs []string
	configs = append(configs, fmt.Sprintf("canal.port=%d", d.Get("port")))
	configs = append(configs, fmt.Sprintf("canal.metrics.pull.port=%d", d.Get("metrics_port")))
	configs = append(configs, fmt.Sprintf("canal.admin.manager=%s", d.Get("admin_host")))
	configs = append(configs, fmt.Sprintf("canal.admin.port=%d", d.Get("admin_port")))
	configs = append(configs, fmt.Sprintf("canal.admin.user=%s", d.Get("admin_user")))
	configs = append(configs, fmt.Sprintf("canal.admin.passwd=%s", d.Get("admin_password")))
	configs = append(configs, fmt.Sprintf("canal.admin.register.auto=%t", d.Get("admin_register_auto")))
	configs = append(configs, fmt.Sprintf("canal.admin.register.cluster=%s", d.Get("name")))

	configs = append(configs, fmt.Sprintf("canal.zkServers=%s", d.Get("zk_hosts")))
	configs = append(configs, fmt.Sprintf("canal.zookeeper.flush.period=%d", d.Get("zk_flush_period")))

	configs = append(configs, fmt.Sprintf("canal.withoutNetty=%t", d.Get("instance_without_netty")))
	configs = append(configs, fmt.Sprintf("canal.serverMode=%s", "kafka"))
	configs = append(configs, fmt.Sprintf("canal.file.data.dir=%s", d.Get("instance_file_data_dir")))
	configs = append(configs, fmt.Sprintf("canal.file.flush.period=%d", d.Get("instance_file_flush_period")))

	configs = append(configs, fmt.Sprintf("canal.instance.memory.buffer.size=%d", d.Get("instance_memory_buffer_size")))
	configs = append(configs, fmt.Sprintf("canal.instance.memory.buffer.memunit=%d", d.Get("instance_memory_buffer_memunit")))
	configs = append(configs, fmt.Sprintf("canal.instance.memory.batch.mode=%s", d.Get("instance_memory_batch_mode")))
	configs = append(configs, fmt.Sprintf("canal.instance.memory.rawEntry=%t", d.Get("instance_memory_raw_entry")))
	configs = append(configs, fmt.Sprintf("canal.instance.detecting.enable=%t", d.Get("instance_detecting_enable")))
	configs = append(configs, fmt.Sprintf("canal.instance.detecting.sql=%s", d.Get("instance_detecting_sql")))
	configs = append(configs, fmt.Sprintf("canal.instance.detecting.interval.time=%d", d.Get("instance_detecting_interval_time")))
	configs = append(configs, fmt.Sprintf("canal.instance.detecting.retry.threshold=%d", d.Get("instance_detecting_retry_threshold")))
	configs = append(configs, fmt.Sprintf("canal.instance.detecting.heartbeatHaEnable=%t", d.Get("instance_detecting_heartbeat_ha_enable")))
	configs = append(configs, fmt.Sprintf("canal.instance.transaction.size=%d", d.Get("instance_transaction_size")))
	configs = append(configs, fmt.Sprintf("canal.instance.fallbackIntervalInSeconds=%d", d.Get("instance_fallback_interval_in_seconds")))
	configs = append(configs, fmt.Sprintf("canal.instance.network.receiveBufferSize=%d", d.Get("instance_network_receive_buffer_size")))
	configs = append(configs, fmt.Sprintf("canal.instance.network.sendBufferSize=%d", d.Get("instance_network_send_buffer_size")))
	configs = append(configs, fmt.Sprintf("canal.instance.network.soTimeout=%d", d.Get("instance_network_so_timeout")))
	configs = append(configs, fmt.Sprintf("canal.instance.filter.druid.ddl=%t", d.Get("instance_filter_druid_ddl")))
	configs = append(configs, fmt.Sprintf("canal.instance.filter.query.dcl=%t", d.Get("instance_filter_query_dcl")))
	configs = append(configs, fmt.Sprintf("canal.instance.filter.query.dml=%t", d.Get("instance_filter_query_dml")))
	configs = append(configs, fmt.Sprintf("canal.instance.filter.query.ddl=%t", d.Get("instance_filter_query_ddl")))
	configs = append(configs, fmt.Sprintf("canal.instance.filter.table.error=%t", d.Get("instance_filter_table_error")))
	configs = append(configs, fmt.Sprintf("canal.instance.filter.rows=%t", d.Get("instance_filter_rows")))
	configs = append(configs, fmt.Sprintf("canal.instance.filter.transaction.entry=%t", d.Get("instance_filter_transaction_entry")))
	configs = append(configs, fmt.Sprintf("canal.instance.binlog.format=%s", d.Get("instance_binlog_format")))
	configs = append(configs, fmt.Sprintf("canal.instance.binlog.image=%s", d.Get("instance_binlog_image")))
	configs = append(configs, fmt.Sprintf("canal.instance.get.ddl.isolation=%t", d.Get("instance_get_ddl_isolation")))
	configs = append(configs, fmt.Sprintf("canal.instance.parser.parallel=%t", d.Get("instance_parser_parallel")))
	configs = append(configs, fmt.Sprintf("canal.instance.parser.parallelBufferSize=%d", d.Get("instance_parser_parallel_buffer_size")))
	configs = append(configs, fmt.Sprintf("canal.instance.tsdb.enable=%t", d.Get("instance_tsdb_enable")))
	configs = append(configs, fmt.Sprintf("canal.instance.tsdb.dir=%s", d.Get("instance_tsdb_dir")))
	configs = append(configs, fmt.Sprintf("canal.instance.tsdb.url=%s", d.Get("instance_tsdb_url")))
	configs = append(configs, fmt.Sprintf("canal.instance.tsdb.dbUsername=%s", d.Get("instance_tsdb_db_username")))
	configs = append(configs, fmt.Sprintf("canal.instance.tsdb.dbPassword=%s", d.Get("instance_tsdb_db_password")))
	configs = append(configs, fmt.Sprintf("canal.instance.tsdb.snapshot.interval=%d", d.Get("instance_tsdb_snapshot_interval")))
	configs = append(configs, fmt.Sprintf("canal.instance.tsdb.snapshot.expire=%d", d.Get("instance_tsdb_snapshot_expire")))

	configs = append(configs, fmt.Sprintf("canal.destinations=%s", d.Get("destination_destinations")))
	configs = append(configs, fmt.Sprintf("canal.conf.dir=%s", d.Get("destination_conf_dir")))
	configs = append(configs, fmt.Sprintf("canal.auto.scan=%t", d.Get("destination_auto_scan")))
	configs = append(configs, fmt.Sprintf("canal.auto.scan.interval=%d", d.Get("destination_auto_scan_interval")))
	configs = append(configs, fmt.Sprintf("canal.instance.tsdb.spring.xml=%s", d.Get("instance_tsdb_spring_xml")))
	configs = append(configs, fmt.Sprintf("canal.instance.global.mode=%s", d.Get("instance_global_mode")))
	configs = append(configs, fmt.Sprintf("canal.instance.global.lazy=%t", d.Get("instance_global_lazy")))
	configs = append(configs, fmt.Sprintf("canal.instance.global.manager.address=%s", d.Get("instance_global_manager_address")))
	configs = append(configs, fmt.Sprintf("canal.instance.global.spring.xml=%s", d.Get("instance_global_spring_xml")))

	// Configuration of org.apache.kafka.clients.producer.KafkaProducer
	configs = append(configs, fmt.Sprintf("canal.mq.servers=%s", d.Get("mq_servers")))
	configs = append(configs, fmt.Sprintf("canal.mq.acks=%s", d.Get("mq_acks")))
	configs = append(configs, fmt.Sprintf("canal.mq.compressionType=%s", d.Get("mq_compression_type")))
	configs = append(configs, fmt.Sprintf("canal.mq.retries=%d", d.Get("mq_retries")))
	configs = append(configs, fmt.Sprintf("canal.mq.batchSize=%d", d.Get("mq_batch_size")))
	configs = append(configs, fmt.Sprintf("canal.mq.lingerMs=%d", d.Get("mq_linger_ms")))
	configs = append(configs, fmt.Sprintf("canal.mq.maxRequestSize=%d", d.Get("mq_max_request_size")))
	configs = append(configs, fmt.Sprintf("canal.mq.bufferMemory=%d", d.Get("mq_buffer_memory")))
	configs = append(configs, fmt.Sprintf("canal.mq.kafka.kerberos.enable=%t", d.Get("mq_kerberos_enable")))
	configs = append(configs, fmt.Sprintf("canal.mq.kafka.kerberos.krb5FilePath=%s", d.Get("mq_kerberos_krb5_file_path")))
	configs = append(configs, fmt.Sprintf("canal.mq.kafka.kerberos.jaasFilePath=%s", d.Get("mq_kerberos_jaas_file_path")))

	// CanalMQStarter, CanalKafkaProducer
	configs = append(configs, fmt.Sprintf("canal.mq.canalBatchSize=%d", d.Get("mq_canal_batch_size")))
	configs = append(configs, fmt.Sprintf("canal.mq.canalGetTimeout=%d", d.Get("mq_canal_get_timeout")))
	configs = append(configs, fmt.Sprintf("canal.mq.flatMessage=%t", d.Get("mq_flat_message")))
	configs = append(configs, fmt.Sprintf("canal.mq.accessChannel=%s", d.Get("mq_access_channel")))

	canalConfig := client.Config{
		Name:    "canal.properties",
		Content: strings.Join(configs, "\n"),
	}

	if d.Id() != "" {
		canalConfig.ClusterId = json.Number(d.Id())
	}

	return canalConfig
}

func assembleCanalConfigResourceData(d *schema.ResourceData, config *client.Config) error {
	if config.ID == "" {
		d.SetId("")
		return nil
	}

	log.Println("assembleCanalConfigResourceData")
	configs := make(map[string]string)
	lines := strings.Split(config.Content, "\n")
	for _, line := range lines {
		elements := strings.Split(line, "=")
		configs[elements[0]] = strings.Join(elements[1:], "=")
	}
	log.Printf("%+v", configs)

	// admin
	setValue(d, configs, "port", "canal.port")
	setValue(d, configs, "metrics_port", "canal.metrics.pull.port")
	setValue(d, configs, "admin_host", "canal.admin.manager")
	setValue(d, configs, "admin_port", "canal.admin.port")
	setValue(d, configs, "admin_password", "canal.admin.passwd")
	setValue(d, configs, "admin_register_auto", "canal.admin.register.auto")

	// zk
	setValue(d, configs, "zk_flush_period", "canal.zookeeper.flush.period")

	// instance
	setValue(d, configs, "instance_without_netty", "canal.withoutNetty")
	setValue(d, configs, "instance_file_data_dir", "canal.file.data.dir")
	setValue(d, configs, "instance_file_flush_period", "canal.file.flush.period")
	setValue(d, configs, "instance_memory_buffer_size", "canal.instance.memory.buffer.size")
	setValue(d, configs, "instance_memory_buffer_memunit", "canal.instance.memory.buffer.memunit")
	setValue(d, configs, "instance_memory_batch_mode", "canal.instance.memory.batch.mode")
	setValue(d, configs, "instance_memory_raw_entry", "canal.instance.memory.rawEntry")
	setValue(d, configs, "instance_detecting_enable", "canal.instance.detecting.enable")
	setValue(d, configs, "instance_detecting_sql", "canal.instance.detecting.sql")
	setValue(d, configs, "instance_detecting_interval_time", "canal.instance.detecting.interval.time")
	setValue(d, configs, "instance_detecting_retry_threshold", "canal.instance.detecting.retry.threshold")
	setValue(d, configs, "instance_detecting_heartbeat_ha_enable", "canal.instance.detecting.heartbeatHaEnable")
	setValue(d, configs, "instance_transaction_size", "canal.instance.transaction.size")
	setValue(d, configs, "instance_fallback_interval_in_seconds", "canal.instance.fallbackIntervalInSeconds")
	setValue(d, configs, "instance_network_receive_buffer_size", "canal.instance.network.receiveBufferSize")
	setValue(d, configs, "instance_network_send_buffer_size", "canal.instance.network.sendBufferSize")
	setValue(d, configs, "instance_network_so_timeout", "canal.instance.network.soTimeout")
	setValue(d, configs, "instance_filter_druid_ddl", "canal.instance.filter.druid.ddl")
	setValue(d, configs, "instance_filter_query_dcl", "canal.instance.filter.query.dcl")
	setValue(d, configs, "instance_filter_query_dml", "canal.instance.filter.query.dml")
	setValue(d, configs, "instance_filter_query_ddl", "canal.instance.filter.query.ddl")
	setValue(d, configs, "instance_filter_table_error", "canal.instance.filter.table.error")
	setValue(d, configs, "instance_filter_rows", "canal.instance.filter.rows")
	setValue(d, configs, "instance_filter_transaction_entry", "canal.instance.filter.transaction.entry")
	setValue(d, configs, "instance_binlog_format", "canal.instance.binlog.format")
	setValue(d, configs, "instance_binlog_image", "canal.instance.binlog.image")
	setValue(d, configs, "instance_get_ddl_isolation", "canal.instance.get.ddl.isolation")
	setValue(d, configs, "instance_parser_parallel", "canal.instance.parser.parallel")
	setValue(d, configs, "instance_parser_parallel_buffer_size", "canal.instance.parser.parallelBufferSize")
	setValue(d, configs, "instance_tsdb_enable", "canal.instance.tsdb.enable")
	setValue(d, configs, "instance_tsdb_dir", "canal.instance.tsdb.dir")
	setValue(d, configs, "instance_tsdb_url", "canal.instance.tsdb.url")
	setValue(d, configs, "instance_tsdb_db_username", "canal.instance.tsdb.dbUsername")
	setValue(d, configs, "instance_tsdb_db_passowrd", "canal.instance.tsdb.dbPassword")
	setValue(d, configs, "instance_tsdb_snapshot_interval", "canal.instance.tsdb.snapshot.interval")
	setValue(d, configs, "instance_tsdb_snapshot_expire", "canal.instance.tsdb.snapshot.expire")

	// destination

	setValue(d, configs, "destination_destinations", "canal.destinations")
	setValue(d, configs, "destination_conf_dir", "canal.conf.dir")
	setValue(d, configs, "destination_auto_scan", "canal.auto.scan")
	setValue(d, configs, "destination_auto_scan_interval", "canal.auto.scan.interval")
	setValue(d, configs, "instance_tsdb_spring_xml", "canal.instance.tsdb.spring.xml")
	setValue(d, configs, "instance_global_mode", "canal.instance.global.mode")
	setValue(d, configs, "instance_global_lazy", "canal.instance.global.lazy")
	setValue(d, configs, "instance_global_manager_address", "canal.instance.global.manager.address")
	setValue(d, configs, "instance_global_spring_xml", "canal.instance.spring.xml")

	// mq
	setValue(d, configs, "mq_servers", "canal.mq.servers")
	setValue(d, configs, "mq_retries", "canal.mq.retries")
	setValue(d, configs, "mq_batch_size", "canal.mq.batchSize")
	setValue(d, configs, "mq_max_request_size", "canal.mq.maxRequestSize")
	setValue(d, configs, "mq_linger_ms", "canal.mq.lingerMs")
	setValue(d, configs, "mq_buffer_memory", "canal.mq.bufferMemory")
	setValue(d, configs, "mq_canal_batch_size", "canal.mq.canalBatchSize")
	setValue(d, configs, "mq_canal_get_timeout", "canal.mq.canalGetTimeout")
	setValue(d, configs, "mq_flat_message", "canal.mq.flatMessage")
	setValue(d, configs, "mq_compression_type", "canal.mq.compressionType")
	setValue(d, configs, "mq_acks", "canal.mq.acks")
	setValue(d, configs, "mq_kerberos_enable", "canal.mq.kafka.kerberos.enable")
	setValue(d, configs, "mq_kerberos_krb5_file_path", "canal.mq.kafka.kerberos.krb5FilePath")
	setValue(d, configs, "mq_kerberos_jaas_file_path", "canal.mq.kafka.kerberos.jaasFilePath")

	return nil
}
