package provider

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Pair struct {
	canalKey, tfKey string
}

func setValue(d *schema.ResourceData, configs map[string]string, resourceKey string, iniKey string) {
	d.Set(resourceKey, configs[iniKey])
	log.Printf("[set value] %v => %v: %v, after: %v", iniKey, resourceKey, configs[iniKey], d.Get(resourceKey))
}
