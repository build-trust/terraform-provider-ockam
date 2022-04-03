package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentity() *schema.Resource {
	return &schema.Resource{
		Description: "Create an Ockam identity.",

		Create: createIdentityFunc,
		Read:   schema.Noop,
		Delete: schema.RemoveFromState,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID of the Ockam identity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"identity": {
				Description: "The JSON representation of the Ockam identity.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"vault": {
				Description: "The JSON representation of the Ockam vault.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func createIdentityFunc(d *schema.ResourceData, c interface{}) error {
	client := c.(*Client)

	dir, err := os.MkdirTemp("", "ockam")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	_, err = client.Run("create-identity", dir)
	if err != nil {
		return err
	}

	id, err := client.Run("print-identity", dir)
	if err != nil {
		return err
	}

	identity, err := client.Read("identity.json", dir)
	if err != nil {
		return err
	}
	vault, err := client.Read("vault.json", dir)
	if err != nil {
		return err
	}

	d.SetId(id)
	if err := d.Set("id", id); err != nil {
		return err
	}
	if err := d.Set("identity", identity); err != nil {
		return err
	}
	if err := d.Set("vault", vault); err != nil {
		return err
	}

	return nil
}
