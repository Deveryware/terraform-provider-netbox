package netbox

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/ipam"
)

func dataNetboxJSONIpamIPRangesList() *schema.Resource {
	return &schema.Resource{
		Read: dataNetboxJSONIpamIPRangesListRead,

		Schema: map[string]*schema.Schema{
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataNetboxJSONIpamIPRangesListRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	params := ipam.NewIpamIPRangesListParams()
	limit := int64(d.Get("limit").(int))
	params.Limit = &limit

	list, err := client.Ipam.IpamIPRangesList(params, nil)
	if err != nil {
		return err
	}

	j, _ := json.Marshal(list.Payload.Results)

	d.Set("json", string(j))
	d.SetId("NetboxJSONIpamIPRangesList")

	return nil
}
