package looker

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

func resourcePermissionSet() *schema.Resource {
	return &schema.Resource{
		Create: resourcePermissionSetCreate,
		Read:   resourcePermissionSetRead,
		Update: resourcePermissionSetUpdate,
		Delete: resourcePermissionSetDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePermissionSetImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePermissionSetCreate(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	permissionSetName := d.Get("name").(string)

	var permissions []string
	for _, permission := range d.Get("permissions").(*schema.Set).List() {
		permissions = append(permissions, permission.(string))
	}

	writePermissionSet := apiclient.WritePermissionSet{
		Name:        &permissionSetName,
		Permissions: &permissions,
	}

	permissionSet, err := client.CreatePermissionSet(writePermissionSet, nil)
	if err != nil {
		return err
	}

	permissionSetID := *permissionSet.Id
	d.SetId(strconv.Itoa(int(permissionSetID)))

	return resourcePermissionSetRead(d, m)
}

func resourcePermissionSetRead(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	permissionSetID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	permissionSet, err := client.PermissionSet(permissionSetID, "", nil)
	if err != nil {
		return err
	}

	if err = d.Set("name", permissionSet.Name); err != nil {
		return err
	}
	if err = d.Set("permissions", permissionSet.Permissions); err != nil {
		return err
	}
	return nil
}

func resourcePermissionSetUpdate(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	permissionSetID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	permissionSetName := d.Get("name").(string)
	var permissions []string
	for _, permission := range d.Get("permissions").(*schema.Set).List() {
		permissions = append(permissions, permission.(string))
	}
	writePermissionSet := apiclient.WritePermissionSet{
		Name:        &permissionSetName,
		Permissions: &permissions,
	}
	_, err = client.UpdatePermissionSet(permissionSetID, writePermissionSet, nil)
	if err != nil {
		return err
	}

	return resourcePermissionSetRead(d, m)
}

func resourcePermissionSetDelete(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	permissionSetID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	_, err = client.DeletePermissionSet(permissionSetID, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourcePermissionSetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourcePermissionSetRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
