package internal_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"terraform-provider-growthbook/internal"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"growthbook": func() (*schema.Provider, error) {
		return internal.Provider(), nil
	},
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GROWTHBOOK_API_KEY"); v == "" {
		t.Fatal("GROWTHBOOK_API_KEY must be set for acceptance tests")
	}
}

func testCheckResourceAttrPrefix(resourceName, attr, prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}
		if !strings.HasPrefix(rs.Primary.Attributes[attr], prefix) {
			return fmt.Errorf("expected %s to start with %s, got %s", attr, prefix, rs.Primary.Attributes[attr])
		}
		return nil
	}
}
