package internal_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceGrowthBookAttribute_basic(t *testing.T) {
	t.Parallel()

	property := acctest.RandomWithPrefix("tf-acc-attr-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "growthbook_attribute" "test" {
  property    = "` + property + `"
  description = "simple test"
  datatype    = "enum"
  enum_values = "test1,test2,test3"
}

data "growthbook_attribute" "by_property" {
  property = growthbook_attribute.test.property
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.growthbook_attribute.by_property", "property", property),
					resource.TestCheckResourceAttr("data.growthbook_attribute.by_property", "datatype", "enum"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.by_property", "description", "simple test"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.by_property", "enum_values", "test1,test2,test3"),
				),
			},
		},
	})
}

// generateManyAttributes generates N growthbook_attribute HCL resources with the given property prefix.
func generateManyAttributes(n int, prefix string) string {
	var b strings.Builder
	b.Grow(n * 120)

	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, `
resource "growthbook_attribute" "attribute_%d" {
  property    = "%s%d"
  description = "simple test"
  datatype    = "enum"
  enum_values = "test1,test2,test3"
}
`, i, prefix, i)
	}

	return b.String()
}

func TestAccDataSourceGrowthbookAttribute_manyAttributes(t *testing.T) {
	t.Parallel()

	prefix := acctest.RandomWithPrefix("tf-acc-attr-")

	config := generateManyAttributes(10, prefix) + `
data "growthbook_attribute" "attr_1" {
  property = growthbook_attribute.attribute_1.property
}
data "growthbook_attribute" "attr_5" {
  property = growthbook_attribute.attribute_5.property
}
data "growthbook_attribute" "attr_10" {
  property = growthbook_attribute.attribute_10.property
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrPrefix("data.growthbook_attribute.attr_1", "property", prefix),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_1", "datatype", "enum"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_1", "description", "simple test"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_1", "enum_values", "test1,test2,test3"),

					testCheckResourceAttrPrefix("data.growthbook_attribute.attr_5", "property", prefix),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_5", "datatype", "enum"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_5", "description", "simple test"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_5", "enum_values", "test1,test2,test3"),

					testCheckResourceAttrPrefix("data.growthbook_attribute.attr_10", "property", prefix),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_10", "datatype", "enum"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_10", "description", "simple test"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_10", "enum_values", "test1,test2,test3"),
				),
			},
		},
	})
}
