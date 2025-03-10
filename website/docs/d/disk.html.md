---
layout: "ignition"
page_title: "Ignition: ignition_disk"
sidebar_current: "docs-ignition-datasource-disk"
description: |-
  Describes the desired state of a system’s disk.
---

# ignition\_disk

Describes the desired state of a system’s disk.

## Example Usage

```hcl
data "ignition_disk" "foo" {
	device = "/dev/sda"
	partition {
		startmib = 2048
		sizemib = 196037632
	}
}
```

## Argument Reference

The following arguments are supported:

* `device` - (Required) The absolute path to the device. Devices are typically referenced by the _/dev/disk/by-*_ symlinks.

* `wipe_table` - (Optional) Whether or not the partition tables shall be wiped. When true, the partition tables are erased before any further manipulation. Otherwise, the existing entries are left intact.

* `partition` - (Optional) The list of partitions and their configuration for this particular disk..


The `partition` block supports:
 
* `label` - (Optional) The PARTLABEL for the partition.

* `number` - (Optional) The partition number, which dictates it’s position in the partition table (one-indexed). If zero, use the next available partition slot.

* `sizemib` - (Optional) The size of the partition (in MiB). If zero, the partition will fill the remainder of the disk.

* `startmib` - (Optional) The start of the partition (in MiB). If zero, the partition will be positioned at the earliest available part of the disk.

* `type_guid` - (Optional) The GPT [partition type GUID](http://en.wikipedia.org/wiki/GUID_Partition_Table#Partition_type_GUIDs). If omitted, the default will be _0FC63DAF-8483-4772-8E79-3D69D8477DE4_ (Linux filesystem data).

* `guid` - (Optional) The GPT unique partition GUID.

* `wipe_partition_entry` - (Optional) If true, Ignition will clobber an existing partition if it does not match the config. If false (default), Ignition will fail instead.

* `should_exist` - (Optional) Whether or not the partition with the specified number should exist. If omitted, it defaults to true. If false Ignition will either delete the specified partition or fail, depending on wipePartitionEntry. If false number must be specified and non-zero and label, start, size, guid, and typeGuid must all be omitted.

* `resize` - (Optional) Whether or not the existing partition should be resized. If omitted, it defaults to false. If true, Ignition will resize an existing partition if it matches the config in all respects except the partition size.

## Attributes Reference

The following attributes are exported:

* `rendered` - The rendered template to reference this resource in _ignition_config_.
