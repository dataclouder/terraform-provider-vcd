package vcd

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var elementVMCreation = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "VM name",
		},
		"description": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "VM description",
		},
		"computer_name": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "VM computer name",
		},
		"os_type": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "OS type for the VM",
		},
		"memory": &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Memory in MB",
		},
		"boot_image": &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Optional boot image",
		},
		"cpus": &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "How many CPUs",
		},
		"cpu_cores": &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "How many CPU cores",
		},
		"power_on": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "A boolean value stating if this VM should be powered on",
		},
		"hardware_version": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Hardware version for the VM",
		},
	},
}

var elementVMTemplate = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "VM name",
		},
		"template_id": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ID of the vApp template",
		},
		"vm_name_in_template": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the VM within the template",
		},
		"power_on": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "A boolean value stating if this VM should be powered on",
		},
	},
}

func resourceVcdVAppV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcdVAppV2Create,
		UpdateContext: resourceVcdVAppV2Update,
		ReadContext:   resourceVcdVAppV2Read,
		DeleteContext: resourceVcdVAppV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceVcdVappImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A name for the vApp, unique withing the VDC",
			},
			"org": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "The name of organization to use, optional if defined at provider " +
					"level. Useful when connected as sysadmin working across different organizations",
			},
			"vdc": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of VDC to use, optional if defined at provider level",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of the vApp",
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				// For now underlying go-vcloud-director repo only supports
				// a value of type String in this map.
				Description: "Key value map of metadata to assign to this vApp. Key and value can be any string.",
			},
			"power_on": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "A boolean value stating if this vApp should be powered on",
			},
			"guest_properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Key/value settings for guest properties. Will be picked up by new VMs when created.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Shows the status code of the vApp",
			},
			"status_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Shows the status of the vApp",
			},
			"vm_creation_list": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "VM creation data",
				Elem:        elementVMCreation,
			},
			"vm_from_template_list": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "VM from template data",
				Elem:        elementVMTemplate,
			},
		},
	}
}

func resourceVcdVAppV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceVcdVAppV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceVcdVAppV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceVcdVAppV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
