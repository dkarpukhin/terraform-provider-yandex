package yandex

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/resourcemanager/v1"
	"github.com/yandex-cloud/go-sdk/sdkresolvers"
)

func dataSourceYandexResourceManagerCloud() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceYandexResourceManagerCloudRead,
		Schema: map[string]*schema.Schema{
			"cloud_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceYandexResourceManagerCloudRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := config.Context()

	err := checkOneOf(d, "cloud_id", "name")
	if err != nil {
		return err
	}

	cloudID := d.Get("cloud_id").(string)
	cloudName, cloudNameOk := d.GetOk("name")

	if cloudNameOk {
		cloudID, err = resolveCloudIDByName(ctx, config, cloudName.(string))
		if err != nil {
			return fmt.Errorf("failed to resolve data source cloud by name: %v", err)
		}
	}

	cloud, err := config.sdk.ResourceManager().Cloud().Get(ctx, &resourcemanager.GetCloudRequest{
		CloudId: cloudID,
	})

	if err != nil {
		return fmt.Errorf("failed to resolve data source cloud by id: %v", err)
	}

	d.Set("cloud_id", cloud.Id)
	d.Set("name", cloud.Name)
	d.Set("description", cloud.Description)
	d.Set("created_at", getTimestamp(cloud.CreatedAt))
	d.SetId(cloud.Id)

	return nil
}

func resolveCloudIDByName(ctx context.Context, config *Config, name string) (string, error) {
	var objectID string
	resolver := sdkresolvers.CloudResolver(name, sdkresolvers.Out(&objectID))

	err := config.sdk.Resolve(ctx, resolver)
	if err != nil {
		return "", err
	}

	return objectID, nil
}
