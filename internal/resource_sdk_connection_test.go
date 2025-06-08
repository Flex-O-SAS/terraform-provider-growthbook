package internal_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccSDKConnectionConfig(name string) string {
	return `
resource "growthbook_project" "test" {
  name = "` + name + `-proj"
}
resource "growthbook_environment" "test" {
  name        = "` + name + `-env"
  projects    = [growthbook_project.test.id]
}
resource "growthbook_sdk_connection" "test" {
  name        = "` + name + `"
  language    = "go"
  environment = growthbook_environment.test.id
  projects    = [growthbook_project.test.id]
}
data "growthbook_sdk_connection" "hh" {
  name = growthbook_sdk_connection.test.name
}
`
}

func TestAccGrowthBookSDKConnection_basic(t *testing.T) {
	connName := acctest.RandomWithPrefix("tf-acc-sdkconn-")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSDKConnectionConfig(connName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("growthbook_sdk_connection.test", "name", connName),
					resource.TestCheckResourceAttr("data.growthbook_sdk_connection.hh", "name", connName),
					resource.TestCheckResourceAttr("data.growthbook_sdk_connection.hh", "language", "go"),
					resource.TestCheckResourceAttrWith("data.growthbook_sdk_connection.hh", "projects.0", func(v string) error {
						if !strings.HasPrefix(v, "prj") {
							return fmt.Errorf("expected projects to start with 'prj', got %s", v)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("data.growthbook_sdk_connection.hh", "encryption_key", func(v string) error {
						if len(v) == 0 {
							return fmt.Errorf("unexpected empty encryption_key attribute")
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("data.growthbook_sdk_connection.hh", "proxy_signing_key", func(v string) error {
						if len(v) == 0 {
							return fmt.Errorf("unexpected empty proxy_signing_key attribute")
						}
						return nil
					}),
					resource.TestCheckResourceAttr("data.growthbook_sdk_connection.hh", "environment", connName+"-env"),
					resource.TestCheckResourceAttrWith("data.growthbook_sdk_connection.hh", "sdk_version", func(v string) error {
						if len(v) == 0 {
							return fmt.Errorf("unexpected empty proxy_signing_key attribute")
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("data.growthbook_sdk_connection.hh", "id", func(v string) error {
						if !strings.HasPrefix(v, "sdk_") {
							return fmt.Errorf("expected id to start with 'sdk_', got %s", v)
						}
						return nil
					}),
				),
			},
		},
	})
}
