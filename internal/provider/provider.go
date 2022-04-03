package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:               map[string]*schema.Schema{},
		ConfigureContextFunc: configureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"ockam_identity": resourceIdentity(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	c, err := NewClient()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, nil
}
