package internal_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceGrowthBookAttribute_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "growthbook_attribute" "test" {
  property    = "test"
  description = "simple test"
  datatype    = "enum"
  enum_values = "test1,test2,test3"
}

data "growthbook_attribute" "by_property" {
  property = growthbook_attribute.test.property
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.growthbook_attribute.by_property", "property", "test"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.by_property", "datatype", "enum"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.by_property", "description", "simple test"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.by_property", "enum_values", "test1,test2,test3"),
				),
			},
		},
	})
}

// Génère N ressources growthbook_attribute en HCL et renvoie la string.
func GenerateManyAttributes(n int) string {
	var b strings.Builder
	b.Grow(n * 120)

	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, `
resource "growthbook_attribute" "attribute_%d" {
  property    = "tf_attr_%d"
  description = "simple test"
  datatype    = "enum"
  enum_values = "test1,test2,test3"
}
`, i, i)
	}

	return b.String()
}

func TestAccDataSourceGrowthbookAttribute_manyAttributes(t *testing.T) {
	t.Parallel()

	config := GenerateManyAttributes(10) + `
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_1", "property", "tf_attr_1"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_1", "datatype", "enum"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_1", "description", "simple test"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_1", "enum_values", "test1,test2,test3"),

					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_5", "property", "tf_attr_5"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_5", "datatype", "enum"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_5", "description", "simple test"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_5", "enum_values", "test1,test2,test3"),

					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_10", "property", "tf_attr_10"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_10", "datatype", "enum"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_10", "description", "simple test"),
					resource.TestCheckResourceAttr("data.growthbook_attribute.attr_10", "enum_values", "test1,test2,test3"),
				),
			},
		},
	})
}
