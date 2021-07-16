package vcd

import (
	//"bufio"
	"context"
	//"encoding/json"
	"fmt"
	//"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVcdVappVmV2() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceVcdVappVmV2Create,
		ReadContext:   resourceVcdVappVmV2Read,
		UpdateContext: resourceVcdVappVmV2Update,
		DeleteContext: resourceVcdVappVmV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceVcdVappVmImport,
		},
		Schema: map[string]*schema.Schema{
			"vapp_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The vApp ID where this VM belongs to",
			},
			"vm_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: fmt.Sprintf("Type of VM: one of ['%s' '%s' '%s']", vappVmType, standaloneVmType, vappVmV2Type),
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "A name for the VM, unique within the vApp",
			},
			"computer_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Computer name to assign to this virtual machine",
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
			"template_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the vApp Template to use",
			},
			"vm_name_in_template": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the VM in vApp Template to use. In cases when vApp template has more than one VM",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The VM description",
			},
			"memory": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "The amount of RAM (in MB) to allocate to the VM",
				ValidateFunc: validateMultipleOf4(),
			},
			"cpus": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of virtual CPUs to allocate to the VM",
			},
			"cpu_cores": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of cores per socket",
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				// For now underlying go-vcloud-director repo only supports
				// a value of type String in this map.
				Description: "Key value map of metadata to assign to this VM",
			},
			"accept_all_eulas": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Automatically accept EULA if OVA has it",
			},
			"power_on": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "A boolean value stating if this VM should be powered on",
			},
			"storage_profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Storage profile to override the default one",
			},
			"os_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Operating System type. Possible values can be found in documentation.",
			},
			"hardware_version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Virtual Hardware Version (e.g.`vmx-14`, `vmx-13`, `vmx-12`, etc.)",
			},
			"boot_image": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Media name to add as boot image.",
			},
			"network_dhcp_wait_seconds": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
				Description: "Optional number of seconds to try and wait for DHCP IP (valid for " +
					"'network' block only)",
			},
			"network": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: " A block to define network interface. Multiple can be used.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Required:     true,
							Type:         schema.TypeString,
							ValidateFunc: vmNetworkTypeValidator(vappVmV2Type),
							Description:  "Network type to use: 'vapp', 'org' or 'none'. Use 'vapp' for vApp network, 'org' to attach Org VDC network. 'none' for empty NIC.",
						},
						"ip_allocation_mode": {
							Optional:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"POOL", "DHCP", "MANUAL", "NONE"}, false),
							Description:  "IP address allocation mode. One of POOL, DHCP, MANUAL, NONE",
						},
						"name": {
							ForceNew:    false,
							Optional:    true, // In case of type = none it is not required
							Type:        schema.TypeString,
							Description: "Name of the network this VM should connect to. Always required except for `type` `NONE`",
						},
						"ip": {
							Computed:     true,
							Optional:     true,
							Type:         schema.TypeString,
							ValidateFunc: checkEmptyOrSingleIP(), // Must accept empty string to ease using HCL interpolation
							Description:  "IP of the VM. Settings depend on `ip_allocation_mode`. Omitted or empty for DHCP, POOL, NONE. Required for MANUAL",
						},
						"is_primary": {
							Optional: true,
							Computed: true,
							// By default if the value is omitted it will report schema change
							// on every terraform operation. The below function
							// suppresses such cases "" => "false" when applying.
							DiffSuppressFunc: falseBoolSuppress(),
							Type:             schema.TypeBool,
							Description:      "Set to true if network interface should be primary. First network card in the list will be primary by default",
						},
						"mac": {
							Computed:    true,
							Optional:    true,
							Type:        schema.TypeString,
							Description: "Mac address of network interface",
						},
						"adapter_type": {
							Type:             schema.TypeString,
							Computed:         true,
							Optional:         true,
							DiffSuppressFunc: suppressCase,
							Description:      "Network card adapter type. (e.g. 'E1000', 'E1000E', 'SRIOVETHERNETCARD', 'VMXNET3', 'PCNet32')",
						},
						"connected": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "It defines if NIC is connected or not.",
						},
					},
				},
			},
			"disk": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Independent disk name",
					},
					"bus_number": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Bus number on which to place the disk controller",
					},
					"unit_number": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Unit number (slot) on the bus specified by BusNumber",
					},
					"size_in_mb": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "The size of the disk in MB.",
					},
				}},
				Optional: true,
				Set:      resourceVcdVmIndependentDiskHash,
			},
			"override_template_disk": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "A block to match internal_disk interface in template. Multiple can be used. Disk will be matched by bus_type, bus_number and unit_number.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"bus_type": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice([]string{"ide", "parallel", "sas", "paravirtual", "sata"}, false),
						Description:  "The type of disk controller. Possible values: ide, parallel( LSI Logic Parallel SCSI), sas(LSI Logic SAS (SCSI)), paravirtual(Paravirtual (SCSI)), sata",
					},
					"size_in_mb": {
						Type:        schema.TypeInt,
						ForceNew:    true,
						Required:    true,
						Description: "The size of the disk in MB.",
					},
					"bus_number": {
						Type:        schema.TypeInt,
						ForceNew:    true,
						Required:    true,
						Description: "The number of the SCSI or IDE controller itself.",
					},
					"unit_number": {
						Type:        schema.TypeInt,
						ForceNew:    true,
						Required:    true,
						Description: "The device number on the SCSI or IDE controller of the disk.",
					},
					"iops": {
						Type:        schema.TypeInt,
						ForceNew:    true,
						Optional:    true,
						Description: "Specifies the IOPS for the disk. Default is 0.",
					},
					"storage_profile": &schema.Schema{
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "Storage profile to override the VM default one",
					},
				}},
			},
			"internal_disk": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A block will show internal disk details",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"disk_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The disk ID.",
					},
					"bus_type": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The type of disk controller. Possible values: ide, parallel( LSI Logic Parallel SCSI), sas(LSI Logic SAS (SCSI)), paravirtual(Paravirtual (SCSI)), sata",
					},
					"size_in_mb": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "The size of the disk in MB.",
					},
					"bus_number": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "The number of the SCSI or IDE controller itself.",
					},
					"unit_number": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "The device number on the SCSI or IDE controller of the disk.",
					},
					"thin_provisioned": {
						Type:        schema.TypeBool,
						Computed:    true,
						Description: "Specifies whether the disk storage is pre-allocated or allocated on demand.",
					},
					"iops": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Specifies the IOPS for the disk. Default is 0.",
					},
					"storage_profile": &schema.Schema{
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Storage profile to override the VM default one",
					},
				}},
			},
			"expose_hardware_virtualization": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Expose hardware-assisted CPU virtualization to guest OS.",
			},
			"guest_properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Key/value settings for guest properties",
			},
			"customization": &schema.Schema{
				Optional:    true,
				Computed:    true,
				MinItems:    1,
				MaxItems:    1,
				Type:        schema.TypeList,
				Description: "Guest customization block",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"force": {
							ValidateFunc: noopValueWarningValidator(true,
								"Using 'true' value for field 'vcd_vapp_vm.customization.force' will reboot VM on every 'terraform apply' operation"),
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							// This settings is used as a 'flag' and it does not matter what is set in the
							// state. If it is 'true' - then it means that 'update' procedure must set the
							// VM for customization at next boot and reboot it.
							DiffSuppressFunc: suppressFalse(),
							Description:      "'true' value will cause the VM to reboot on every 'apply' operation",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "'true' value will enable guest customization. It may occur on first boot or when 'force' is used",
						},
						"change_sid": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "'true' value will change SID. Applicable only for Windows VMs",
						},
						"allow_local_admin_password": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Allow local administrator password",
						},
						"must_change_password_on_first_login": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Require Administrator to change password on first login",
						},
						"auto_generate_password": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Auto generate password",
						},
						"admin_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Sensitive:   true,
							Description: "Manually specify admin password",
						},
						"number_of_auto_logons": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							Description:  "Number of times to log on automatically. '0' - disabled.",
							ValidateFunc: validation.IntAtLeast(0),
						},
						"join_domain": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable this VM to join a domain",
						},
						"join_org_domain": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Use organization's domain for joining",
						},
						"join_domain_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Custom domain name for join",
						},
						"join_domain_user": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Username for custom domain name join",
						},
						"join_domain_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Sensitive:   true,
							Description: "Password for custom domain name join",
						},
						"join_domain_account_ou": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Account organizational unit for domain name join",
						},
						"initscript": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Script to run on initial boot or with customization.force=true set",
						},
					},
				},
			},
			"cpu_hot_add_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "True if the virtual machine supports addition of virtual CPUs while powered on.",
			},
			"memory_hot_add_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "True if the virtual machine supports addition of memory while powered on.",
			},
			"prevent_update_power_off": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "True if the update of resource should fail when virtual machine power off needed.",
			},
			"sizing_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "VM sizing policy ID. Has to be assigned to Org VDC.",
			},
		},
		Description: "Vapp VM V2",
	}
}

/*
type transferVmData struct {
	Org string `json:"org,omitempty"`
	Vdc string `json:"vdc,omitempty"`
	VappId string `json:"vapp_id"`
	Name string `json:"name"`
	ComputerName string `json:"computer_name"`
	TemplateId string `json:"template_id"`
	VmNameInTemplate string `json:"vm_name_in_template,omitempty"`
	Description string `json:"description,omitempty"`
	OsType string `json:"os_type"`
	Memory int `json:"memory"`
	Cpus int `json:"cpus"`
	CpuCores int `json:"cpu_cores"`
	AcceptAllEulas bool `json:"accept_all_eulas"`
	PowerOn bool `json:"power_on"`
	HardwareVersion string `json:"hardware_version"`
	BootImage string `json:"boot_image,omitempty"`

}
*/

func resourceVcdVappVmV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	/*
		// When initializing the resource, we
		if d.Get("vapp_id").(string) == "" {
			err := saveVmResourceToFile(d)
			if err != nil {
				return diag.Errorf("error saving VM resource contents to JSON file: %s", err)
			}
			return nil
		}

	*/

	return nil
}

func resourceVcdVappVmV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceVcdVappVmV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceVcdVappVmV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

/*
func saveVmResourceToFile(d *schema.ResourceData) error {

	var outRec = transferVmData{
		Org:              d.Get("org").(string),
		Vdc:              d.Get("vdc").(string),
		VappId:           "",
		Name:             d.Get("name").(string),
		ComputerName:     d.Get("computer_name").(string),
		TemplateId:       d.Get("template_id").(string),
		VmNameInTemplate: d.Get("vm_name_in_template").(string),
		Description:      d.Get("description").(string),
		OsType:           d.Get("os_type").(string),
		Memory:           d.Get("memory").(int),
		Cpus:             d.Get("cpus").(int),
		CpuCores:         d.Get("cpu_cores").(int),
		AcceptAllEulas:   d.Get("accept_all_eulas").(bool),
		PowerOn:          d.Get("template_id").(bool),
		HardwareVersion:  d.Get("hardware_version").(string),
		BootImage:        d.Get("boot_image").(string),
	}

	jsonStr, err := json.MarshalIndent(outRec, " ", " ")
	if err != nil {
		return fmt.Errorf("error encoding data for %s: %s",outRec.Name, err)
	}
	recFile := outRec.Name + ".json"
	file, err := os.Create(recFile)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %s", recFile,  err)
	}
	writer := bufio.NewWriter(file)
	count, err := writer.Write(jsonStr)
	if err != nil || count == 0 {
		return fmt.Errorf("error writing to file %s. Reported %d bytes written. %s", recFile, count, err)
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing file %s. %s", recFile, err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("error closing file %s. %s", recFile, err)
	}
	return nil
}

*/
