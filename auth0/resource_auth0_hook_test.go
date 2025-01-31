package auth0

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHook(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHookCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "name", "pre-user-reg-hook"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "script", "function (user, context, callback) { callback(null, { user }); }"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "trigger_id", "pre-user-registration"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "enabled", "true"),
				),
			},
		},
	})
}

const testAccHookCreate = `
resource "auth0_hook" "my_hook" {
  name = "pre-user-reg-hook"
  script = "function (user, context, callback) { callback(null, { user }); }"
  trigger_id = "pre-user-registration"
  enabled = true
}
`

func TestAccHookSecrets(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHookSecrets("alpha"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "name", "pre-user-reg-hook"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "dependencies.auth0", "2.30.0"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "script", "function (user, context, callback) { callback(null, { user }); }"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "trigger_id", "pre-user-registration"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "enabled", "true"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "secrets.foo", "alpha"),
					resource.TestCheckNoResourceAttr("auth0_hook.my_hook", "secrets.bar"),
				),
			},
			{
				Config: testAccHookSecrets2("gamma", "kappa"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "name", "pre-user-reg-hook"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "dependencies.auth0", "2.30.0"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "script", "function (user, context, callback) { callback(null, { user }); }"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "trigger_id", "pre-user-registration"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "enabled", "true"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "secrets.foo", "gamma"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "secrets.bar", "kappa"),
				),
			},
			{
				Config: testAccHookSecrets("delta"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "name", "pre-user-reg-hook"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "script", "function (user, context, callback) { callback(null, { user }); }"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "trigger_id", "pre-user-registration"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "enabled", "true"),
					resource.TestCheckResourceAttr("auth0_hook.my_hook", "secrets.foo", "delta"),
					resource.TestCheckNoResourceAttr("auth0_hook.my_hook", "secrets.bar"),
				),
			},
		},
	})
}

func testAccHookSecrets(fooValue string) string {
	return fmt.Sprintf(`
resource "auth0_hook" "my_hook" {
  name = "pre-user-reg-hook"
  script = "function (user, context, callback) { callback(null, { user }); }"
  trigger_id = "pre-user-registration"
  enabled = true
  dependencies = {
    auth0 = "2.30.0"
  }
	secrets = {
    foo = "%s"
  }
}
`, fooValue)
}

func testAccHookSecrets2(fooValue string, barValue string) string {
	return fmt.Sprintf(`
resource "auth0_hook" "my_hook" {
  name = "pre-user-reg-hook"
  script = "function (user, context, callback) { callback(null, { user }); }"
  trigger_id = "pre-user-registration"
  dependencies = {
    auth0 = "2.30.0"
  }
  enabled = true
  secrets = {
    foo = "%s"
    bar = "%s"
  }
}
`, fooValue, barValue)
}

func TestHookNameRegexp(t *testing.T) {
	for name, valid := range map[string]bool{
		"my-hook-1":                 true,
		"hook 2 name with spaces":   true,
		" hook with a space prefix": false,
		"hook with a space suffix ": false,
		" ":                         false,
		"   ":                       false,
	} {
		fn := validateHookName()

		validationResult := fn(name, cty.Path{cty.GetAttrStep{Name: "name"}})
		if validationResult.HasError() && valid {
			t.Fatalf("Expected %q to be valid, but got validation errors %v", name, validationResult)
		}
		if !validationResult.HasError() && !valid {
			t.Fatalf("Expected %q to be invalid, but got no validation errors.", name)
		}
	}
}
