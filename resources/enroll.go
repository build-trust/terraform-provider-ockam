package resources

import (
	"context"
	"encoding/json"

	"github.com/build-trust/terraform-provider-ockam/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Enroll struct {
	ID    string
	name  string
	users []string
}

func ResourceEnroll() *schema.Resource {
	e := Enroll{}

	return &schema.Resource{
		Description:   "Enroll with Ockam Orchestrator",
		CreateContext: e.createOckamEnroll,
		ReadContext:   e.readOckamEnroll,
		DeleteContext: e.deleteOckamEnroll,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Space identifier",
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
			},
			"name": {
				Description: "Space name",
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
			},
			"users": {
				Description: "Space user",
				Type:        schema.TypeSet,
				Computed:    true,
				Required:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func (e Enroll) createOckamEnroll(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	command := []string{"enroll"}

	if n, ok := d.Get("node").(string); ok {
		command = append(command, n)
	}

	_, err := client.Run(command...)
	if err != nil {
		return diag.FromErr(err)
	}

	return e.ReadSpaceProjectHelper(client, d)
}

func (e Enroll) readOckamEnroll(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	return e.ReadSpaceProjectHelper(client, d)
}

func (e Enroll) deleteOckamEnroll(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	return e.DeleteSpaceProjectHelper(client, d)
}

func (e Enroll) DeleteSpaceProjectHelper(client *client.Client, d *schema.ResourceData) diag.Diagnostics {
	ID := d.Id()

	_, err := client.Run("space", "delete", ID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func (e Enroll) ReadSpaceProjectHelper(client *client.Client, d *schema.ResourceData) diag.Diagnostics {
	space, err := client.Run("space", "list", "--output json")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := json.Unmarshal([]byte(space), &e); err != nil {
		return diag.Errorf("Error retrieving created space")
	}

	d.SetId(e.name)
	d.Set("id", e.ID)
	d.Set("name", e.name)
	d.Set("users", e.users)

	return nil
}
