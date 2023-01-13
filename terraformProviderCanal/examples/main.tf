terraform {
  required_providers {
    canal = {
      version = "0.0.0-sample"
      source  = "paypay/daas/canal"
    }
  }
}

provider "canal" {
  username = "admin"
  password = "haveaniceday"
  host     = "http://localhost:8089"
}

resource "canal_cluster" "test_cluster" {
  name = "test_cluster"
  zk_hosts = "localhost:2181"
  mq_servers = "test"
}

resource "canal_instance" "test_instance" {
  name              = "test_instance"
  cluster_id        = 1
  db_credential_id  = ""
  mq_partitions_num = 3
  enabled = true

  tables {
    table {
      database_name = "testdb"
      table_name    = "test"
      primary_key   = "id"
      mq_topic      = "testdb_test" // is created by steaming ingester resource
    }
  }
}
