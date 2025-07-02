package ignition

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/coreos/ignition/v2/config/v3_5/types"
)

func TestIgnitionLink(t *testing.T) {
	testIgnition(t, `
		data "ignition_link" "foo" {
			path = "/foo"
			target = "/bar"
			hard = true
			uid = 42
			gid = 84
		}

		data "ignition_link" "baz" {
			path = "/baz"
			target = "/qux"
			overwrite = true
		}

		data "ignition_link" "named_foo" {
			path = "/foo"
			target = "/bar"
			user = "foo"
			group = "foo"
		}

		data "ignition_config" "test" {
			links = [
				data.ignition_link.foo.rendered,
				data.ignition_link.baz.rendered,
				data.ignition_link.named_foo.rendered,
			]
		}
	`, func(c *types.Config) error {
		if len(c.Storage.Links) != 3 {
			return fmt.Errorf("arrays, found %d", len(c.Storage.Links))
		}

		f := c.Storage.Links[0]
		if f.Path != "/foo" {
			return fmt.Errorf("path, found %q", f.Path)
		}

		if *f.Overwrite != false {
			return fmt.Errorf("overwrite, found %v", *f.Overwrite)
		}

		if *f.Target != "/bar" {
			return fmt.Errorf("target, found %q", *f.Target)
		}

		if *f.Hard != true {
			return fmt.Errorf("hard, found %v", *f.Hard)
		}

		if *f.User.ID != 42 {
			return fmt.Errorf("uid, found %q", *f.User.ID)
		}

		if *f.Group.ID != 84 {
			return fmt.Errorf("gid, found %q", *f.Group.ID)
		}

		f = c.Storage.Links[1]
		if f.Path != "/baz" {
			return fmt.Errorf("path, found %q", f.Path)
		}

		if *f.Target != "/qux" {
			return fmt.Errorf("target, found %q", *f.Target)
		}

		if *f.Overwrite != true {
			return fmt.Errorf("overwrite, found %v", *f.Overwrite)
		}

		f = c.Storage.Links[2]
		if f.Path != "/foo" {
			return fmt.Errorf("path, found %q", f.Path)
		}

		if *f.Target != "/bar" {
			return fmt.Errorf("target, found %q", *f.Target)
		}

		if *f.User.Name != "foo" {
			return fmt.Errorf("user, found %q", *f.User.Name)
		}

		if *f.Group.Name != "foo" {
			return fmt.Errorf("group, found %q", *f.Group.Name)
		}

		return nil
	})
}

func TestIgnitionLinkInvalidPath(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_link" "foo" {
			path = "foo"
			target = "bar"
		}

		data "ignition_config" "test" {
			links = [data.ignition_link.foo.rendered]
		}
	`, regexp.MustCompile("absolute"))
}

func TestIgnitionLinkUIDUserConflict(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_link" "foo" {
			path = "/foo"
			target = "/bar"
			uid = 1000
			user = "foo"
		}

		data "ignition_config" "test" {
			links = [data.ignition_link.foo.rendered]
		}
	`, regexp.MustCompile("Conflicting configuration arguments"))
}

func TestIgnitionLinkGIDGroupConflict(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_link" "foo" {
			path = "/foo"
			target = "/bar"
			gid = 1000
			group = "foo"
		}

		data "ignition_config" "test" {
			links = [data.ignition_link.foo.rendered]
		}
	`, regexp.MustCompile("Conflicting configuration arguments"))
}
