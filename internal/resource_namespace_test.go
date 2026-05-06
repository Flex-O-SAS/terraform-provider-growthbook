package internal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccNamespaceConfig(name string) string {
	return `
resource "growthbook_namespace" "test" {
  display_name   = "` + name + `"
  description    = "Acceptance test namespace"
  status         = "active"
  format         = "legacy"
  hash_attribute = "id"
}
data "growthbook_namespace" "by_display_name" {
  display_name = growthbook_namespace.test.display_name
}
`
}

func testAccNamespaceConfigUpdate(name string) string {
	return `
resource "growthbook_namespace" "test" {
  display_name   = "` + name + `"
  description    = "Updated namespace"
  status         = "inactive"
  format         = "legacy"
  hash_attribute = "user_id"
}
data "growthbook_namespace" "by_display_name" {
  display_name = growthbook_namespace.test.display_name
}
`
}

func TestAccNamespace_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-acc-namespace-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_namespace.test", "display_name", name),
					resource.TestCheckResourceAttr("data.growthbook_namespace.by_display_name", "status", "active"),
				),
			},
			{
				Config: testAccNamespaceConfigUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_namespace.test", "description", "Updated namespace"),
					resource.TestCheckResourceAttr("data.growthbook_namespace.by_display_name", "hash_attribute", "user_id"),
				),
			},
		},
	})
}
