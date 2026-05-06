package internal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccSavedGroupConfig(name string) string {
	return `
resource "growthbook_saved_group" "test" {
  name          = "` + name + `"
  type          = "list"
  attribute_key = "id"
  values        = ["one", "two"]
  description   = "Acceptance test saved group"
}
data "growthbook_saved_group" "by_name" {
  name = growthbook_saved_group.test.name
}
`
}

func testAccSavedGroupConfigUpdate(name string) string {
	return `
resource "growthbook_saved_group" "test" {
  name          = "` + name + `"
  type          = "list"
  attribute_key = "user_id"
  values        = ["three", "four"]
  description   = "Updated saved group"
}
data "growthbook_saved_group" "by_name" {
  name = growthbook_saved_group.test.name
}
`
}

func TestAccSavedGroup_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("tf-acc-saved-group-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSavedGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_saved_group.test", "name", name),
					resource.TestCheckResourceAttr("growthbook_saved_group.test", "attribute_key", "id"),
					resource.TestCheckResourceAttr("growthbook_saved_group.test", "description", "Acceptance test saved group"),
					resource.TestCheckResourceAttr("data.growthbook_saved_group.by_name", "name", name),
				),
			},
			{
				Config: testAccSavedGroupConfigUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_saved_group.test", "attribute_key", "user_id"),
					resource.TestCheckResourceAttr("growthbook_saved_group.test", "description", "Updated saved group"),
					resource.TestCheckResourceAttr("data.growthbook_saved_group.by_name", "description", "Updated saved group"),
				),
			},
		},
	})
}
