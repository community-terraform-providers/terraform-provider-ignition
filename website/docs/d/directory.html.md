---
layout: "ignition"
page_title: "Ignition: ignition_directory"
sidebar_current: "docs-ignition-datasource-directory"
description: |-
  Describes a directory to be created in a particular filesystem.
---

# ignition\_directory

Describes a directory to be created in a particular filesystem.

## Example Usage

```hcl
data "ignition_directory" "folder" {
	path = "/folder"
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required) The absolute path to the directory.

* `overwrite` - (Optional) Whether to delete preexisting nodes at the path. Defaults to false.

* `mode` - (Optional) The directory's permission mode. Note that the mode must be properly specified as a decimal value, not octal (i.e. 0755 -> 493).

* `user` - (Optional) The user name of the owner.

  __Note__: `user` and `uid` are mutually exclusive.

* `uid` - (Optional) The user ID of the owner.

  __Note__: `uid` and `user` are mutually exclusive.

* `group` - (Optional) The group name of the owner.

  __Note__: `group` and `gid` are mutually exclusive.

* `gid` - (Optional) The group ID of the owner.

  __Note__: `gid` and `group` are mutually exclusive.


## Attributes Reference

The following attributes are exported:

* `rendered` - The rendered template to reference this resource in _ignition_config_.
