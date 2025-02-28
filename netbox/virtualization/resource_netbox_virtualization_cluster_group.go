package virtualization

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/virtualization"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/util"
)

func ResourceNetboxVirtualizationClusterGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a cluster group (virtualization module) within Netbox.",
		CreateContext: resourceNetboxVirtualizationClusterGroupCreate,
		ReadContext:   resourceNetboxVirtualizationClusterGroupRead,
		UpdateContext: resourceNetboxVirtualizationClusterGroupUpdate,
		DeleteContext: resourceNetboxVirtualizationClusterGroupDelete,
		Exists:        resourceNetboxVirtualizationClusterGroupExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this cluster group (virtualization module).",
			},
			"cluster_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of clusters in this cluster group (virtualization module).",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this cluster group was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 100),
				Description:  "The description of this cluster group (virtualization module).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this cluster group was last updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this cluster group (virtualization module).",
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The slug of this cluster group (virtualization module).",
			},
			"tag": &tag.TagSchema,
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this cluster group (virtualization module).",
			},
		},
	}
}

var clusterGroupRequiredFields = []string{
	"created",
	"last_updated",
	"name",
	"slug",
	"tags",
}

func resourceNetboxVirtualizationClusterGroupCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.ClusterGroup{
		CustomFields: customFields,
		Description:  d.Get("description").(string),
		Name:         &name,
		Slug:         &slug,
		Tags:         tag.ConvertTagsToNestedTags(tags),
	}

	resource := virtualization.NewVirtualizationClusterGroupsCreateParams().WithData(newResource)

	resourceCreated, err := client.Virtualization.VirtualizationClusterGroupsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxVirtualizationClusterGroupRead(ctx, d, m)
}

func resourceNetboxVirtualizationClusterGroupRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := virtualization.NewVirtualizationClusterGroupsListParams().WithID(&resourceID)
	resources, err := client.Virtualization.VirtualizationClusterGroupsList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(resources.Payload.Results) != 1 {
		d.SetId("")
		return nil
	}

	resource := resources.Payload.Results[0]

	if err = d.Set("cluster_count", resource.ClusterCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("content_type", util.ConvertURIContentType(resource.URL)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("created", resource.Created.String()); err != nil {
		return diag.FromErr(err)
	}
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

	if err = d.Set("custom_field", customFields); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("description", resource.Description); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", resource.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("slug", resource.Slug); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("url", resource.URL); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxVirtualizationClusterGroupUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	modifiedFields := make(map[string]interface{})

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	params := &models.ClusterGroup{}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		params.Description = description
		modifiedFields["description"] = description
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}
	if d.HasChange("slug") {
		slug := d.Get("slug").(string)
		params.Slug = &slug
	}
	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		params.Tags = tag.ConvertTagsToNestedTags(tags)
	}

	resource := virtualization.NewVirtualizationClusterGroupsPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Virtualization.VirtualizationClusterGroupsPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modifiedFields, clusterGroupRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxVirtualizationClusterGroupRead(ctx, d, m)
}

func resourceNetboxVirtualizationClusterGroupDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxVirtualizationClusterGroupExists(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource := virtualization.NewVirtualizationClusterGroupsDeleteParams().WithID(id)
	if _, err := client.Virtualization.VirtualizationClusterGroupsDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxVirtualizationClusterGroupExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := virtualization.NewVirtualizationClusterGroupsListParams().WithID(&resourceID)
	resources, err := client.Virtualization.VirtualizationClusterGroupsList(params, nil)
	if err != nil {
		return resourceExist, err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			resourceExist = true
		}
	}

	return resourceExist, nil
}
