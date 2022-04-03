package provider

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceIdentity(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIdentityConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceIdentityCheck("ockam_identity.foo"),
					testAccResourceIdentityCheck("ockam_identity.bar"),
					testAccResourceIdentityCheckDup("ockam_identity.foo", "ockam_identity.bar"),
				),
			},
		},
	})
}

func testAccResourceIdentityCheck(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id := rs.Primary.Attributes["id"]
		if len(id) != 64 {
			return fmt.Errorf("id length incorrect, is %d, expected 64", len(id))
		}

		identity := rs.Primary.Attributes["identity"]
		if !isJSON(identity) {
			return fmt.Errorf("identity is not valid JSON")
		}
		vault := rs.Primary.Attributes["vault"]
		if !isJSON(vault) {
			return fmt.Errorf("vault is not valid JSON")
		}

		return nil
	}
}

func testAccResourceIdentityCheckDup(id1, id2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs1, ok1 := s.RootModule().Resources[id1]
		rs2, ok2 := s.RootModule().Resources[id2]
		if !ok1 || !ok2 {
			return fmt.Errorf("Not found")
		}

		id1 := rs1.Primary.Attributes["id"]
		id2 := rs2.Primary.Attributes["id"]
		if id1 == id2 {
			return fmt.Errorf("ids are identical")
		}

		identity1 := rs1.Primary.Attributes["identity"]
		identity2 := rs2.Primary.Attributes["identity"]
		if identity1 == identity2 {
			return fmt.Errorf("identities are identical")
		}

		vault1 := rs1.Primary.Attributes["vault"]
		vault2 := rs2.Primary.Attributes["vault"]
		if vault1 == vault2 {
			return fmt.Errorf("vaults are identical")
		}

		return nil
	}
}

const (
	testAccResourceIdentityConfig = `
resource "ockam_identity" "foo" { }

resource "ockam_identity" "bar" { }
`
)

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
