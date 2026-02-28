package ignition

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/coreos/ignition/v2/config/v3_6/types"
)

func TestIgnitionDirectory(t *testing.T) {
	testIgnition(t, `
		data "ignition_directory" "foo" {
			path = "/foo"
			mode = 420
			uid = 42
			gid = 84
		}

		data "ignition_directory" "bar" {
			path = "/bar"
			overwrite = true
		}

		data "ignition_directory" "named_foo" {
			path = "/foo"
			mode = 420
			user = "foo"
			group = "foo"
		}

		data "ignition_config" "test" {
			directories = [
				data.ignition_directory.foo.rendered,
				data.ignition_directory.bar.rendered,
				data.ignition_directory.named_foo.rendered,
			]
		}
	`, func(c *types.Config) error {
		if len(c.Storage.Directories) != 3 {
			return fmt.Errorf("arrays, found %d", len(c.Storage.Directories))
		}

		f := c.Storage.Directories[0]
		if f.Path != "/foo" {
			return fmt.Errorf("path, found %q", f.Path)
		}

		if *f.Overwrite != false {
			return fmt.Errorf("overwrite, found %t", *f.Overwrite)
		}

		if int(*f.Mode) != 420 {
			return fmt.Errorf("mode, found %q", *f.Mode)
		}

		if *f.User.ID != 42 {
			return fmt.Errorf("uid, found %q", *f.User.ID)
		}

		if *f.Group.ID != 84 {
			return fmt.Errorf("gid, found %q", *f.Group.ID)
		}

		f = c.Storage.Directories[1]
		if f.Path != "/bar" {
			return fmt.Errorf("path, found %q", f.Path)
		}

		if *f.Overwrite != true {
			return fmt.Errorf("overwrite, found %t", *f.Overwrite)
		}

		f = c.Storage.Directories[2]
		if f.Path != "/foo" {
			return fmt.Errorf("path, found %q", f.Path)
		}

		if *f.Overwrite != false {
			return fmt.Errorf("overwrite, found %t", *f.Overwrite)
		}

		if int(*f.Mode) != 420 {
			return fmt.Errorf("mode, found %q", *f.Mode)
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

func TestIgnitionDirectoryInvalidMode(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_directory" "foo" {
			path = "/foo"
			mode = 999999
		}

		data "ignition_config" "test" {
			directories = [data.ignition_directory.foo.rendered]
		}
	`, regexp.MustCompile("illegal file mode"))
}

func TestIgnitionDirectoryInvalidPath(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_directory" "foo" {
			path = "foo"
			mode = 999999
		}

		data "ignition_config" "test" {
			directories = [data.ignition_directory.foo.rendered]
		}
	`, regexp.MustCompile("path not absolute"))
}

func TestIgnitionDirectoryUIDUserConflict(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_directory" "foo" {
			path = "foo"
			uid = 1000
			user = "foo"
		}

		data "ignition_config" "test" {
			directories = [data.ignition_directory.foo.rendered]
		}
	`, regexp.MustCompile("Conflicting configuration arguments"))
}

func TestIgnitionDirectoryGIDGroupConflict(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_directory" "foo" {
			path = "foo"
			gid = 1000
			group = "foo"
		}

		data "ignition_config" "test" {
			directories = [data.ignition_directory.foo.rendered]
		}
	`, regexp.MustCompile("Conflicting configuration arguments"))
}
