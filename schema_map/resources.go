package schema_map

import (
	"github.com/build-trust/terraform-provider-ockam/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"enroll": resources.ResourceEnroll(),
	}
}
