---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netbox_json_extras_job_results_list Data Source - terraform-provider-netbox"
subcategory: ""
description: |-
  Get json output from the extrasjobresults_list Netbox endpoint.
---

# netbox_json_extras_job_results_list (Data Source)

Get json output from the extras_job_results_list Netbox endpoint.

## Example Usage

```terraform
data "netbox_json_extras_job_results_list" "test" {
  limit = 0
}

output "example" {
  value = jsondecode(data.netbox_json_extras_job_results_list.test.json)
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `limit` (Number) The max number of returned results. If 0 is specified, all records will be returned.

### Read-Only

- `id` (String) The ID of this resource.
- `json` (String) JSON output of the list of objects for this Netbox endpoint.


