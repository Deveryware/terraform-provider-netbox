---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netbox_virtualization_cluster Resource - terraform-provider-netbox"
subcategory: ""
description: |-
  Manage a tag (extra module) within Netbox.
---

# netbox_virtualization_cluster (Resource)

Manage a tag (extra module) within Netbox.

## Example Usage

```terraform
resource "netbox_virtualization_cluster" "cluster_test" {
	name = "Test cluster"
	type_id = netbox_virtualization_cluster_type.test.id

	comments = <<-EOT
	Test cluster
	EOT

	group_id = netbox_virtualization_cluster_group.test.id
	site_id = netbox_dcim_site.site_test.id
	tenant_id = netbox_tenancy_tenant.tenant_test.id

  tag {
    name = "tag1"
    slug = "tag1"
  }

  custom_field {
    name = "cf_boolean"
    type = "boolean"
    value = "true"
  }

  custom_field {
    name = "cf_date"
    type = "date"
    value = "2020-12-25"
  }

  custom_field {
    name = "cf_text"
    type = "text"
    value = "some text"
  }

  custom_field {
    name = "cf_integer"
    type = "integer"
    value = "10"
  }

  custom_field {
    name = "cf_selection"
    type = "select"
    value = "1"
  }

  custom_field {
    name = "cf_url"
    type = "url"
    value = "https://github.com"
  }

  custom_field {
    name = "cf_multi_selection"
    type = "multiselect"
    value = jsonencode([
      "0",
      "1"
    ])
  }

  custom_field {
    name = "cf_json"
    type = "json"
    value = jsonencode({
      stringvalue = "string"
      boolvalue = false
      dictionary = {
        numbervalue = 5
      }
    })
  }

  custom_field {
    name = "cf_object"
    type = "object"
    value = 1
  }

  custom_field {
    name = "cf_multi_object"
    type = "multiobject"
    value = jsonencode([
      1,
      2
    ])
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of this cluster (virtualization module).
- `type_id` (Number) Type of this cluster.

### Optional

- `comments` (String) Comments for this cluster (virtualization module).
- `custom_field` (Block Set) Existing custom fields to associate to this ressource. (see [below for nested schema](#nestedblock--custom_field))
- `group_id` (Number) The cluster group of this cluster.
- `site_id` (Number) The site of this cluster.
- `tag` (Block Set) Existing tag to associate to this resource. (see [below for nested schema](#nestedblock--tag))
- `tenant_id` (Number) ID of the tenant where this cluster is attached.

### Read-Only

- `content_type` (String) The content type of this cluster (virtualization module).
- `created` (String) Date when this cluster was created.
- `device_count` (Number) Number of devices in this cluster.
- `id` (String) The ID of this resource.
- `last_updated` (String) Date when this cluster was last updated.
- `url` (String) The link to this cluster (virtualization module).
- `virtualmachine_count` (Number) Number of virtual machines in this cluster.

<a id="nestedblock--custom_field"></a>
### Nested Schema for `custom_field`

Required:

- `name` (String) Name of the existing custom field.
- `type` (String) Type of the existing custom field (text, longtext, integer, boolean, date, url, json, select, multiselect, object, multiobject, selection (deprecated), multiple(deprecated)).
- `value` (String) Value of the existing custom field.


<a id="nestedblock--tag"></a>
### Nested Schema for `tag`

Required:

- `name` (String) Name of the existing tag.
- `slug` (String) Slug of the existing tag.


