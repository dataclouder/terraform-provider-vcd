//go:build vapp || vm || vmv2 || ALL || functional
// +build vapp vm vmv2 ALL functional

package vcd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdVAppVmV2_Basic(t *testing.T) {
	preTestChecks(t)
	var vappName = t.Name() + "-vapp"
	var vmName = t.Name() + "-vm"

	var params = StringMap{
		"Org":          testConfig.VCD.Org,
		"Vdc":          testConfig.VCD.Vdc,
		"EdgeGateway":  testConfig.Networking.EdgeGateway,
		"NetworkName":  t.Name() + "-net",
		"Catalog":      testConfig.VCD.Catalog.Name,
		"CatalogItem":  testConfig.VCD.Catalog.CatalogItem,
		"MediaItem":    testConfig.Media.MediaName,
		"VappName":     vappName,
		"VmName":       vmName,
		"ComputerName": vmName + "-unique",
		"Tags":         "vapp vm vmv2",
	}

	configText := templateFill(testAccCheckVcdVAppVmV2_basic, params)
	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	debugPrintf("#[DEBUG] CONFIGURATION: %s\n", configText)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckVcdVAppVmDestroy(vappName),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVcdVAppVmExists(vappName, vmName, "vcd_vapp_vm_v2."+vmName, nil, nil),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm_v2."+vmName+"1", "name", vmName+"1"),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm_v2."+vmName+"2", "name", vmName+"2"),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm_v2."+vmName+"3", "name", vmName+"3"),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm_v2."+vmName+"1", "computer_name", vmName+"-unique1"),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm_v2."+vmName+"2", "computer_name", vmName+"-unique2"),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm_v2."+vmName+"3", "computer_name", vmName+"-unique3"),
				),
			},
			resource.TestStep{
				ResourceName:      "vcd_vapp_vm." + vmName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdVappObject(testConfig, vappName2, vmName),
				// These fields can't be retrieved from user data
				ImportStateVerifyIgnore: []string{"template_name", "catalog_name",
					"accept_all_eulas", "power_on", "computer_name", "prevent_update_power_off"},
			},
		},
	})
	postTestChecks(t)
}

const testAccCheckVcdVAppVmV2_basic = `

resource "vcd_vapp" "{{.VappName}}" {
  name = "{{.VappName}}"
  org  = "{{.Org}}"
  vdc  = "{{.Vdc}}"
}


data "vcd_catalog_item"" "main_item" {
	catalog = "{{.Catalog}}"
    name    = "{{.CatalogItem}}"
}

data "vcd_media_item" "media_item" {
	catalog = "{{.Catalog}}"
    name    = "{{.MediaItem}}"
}

resource "vcd_vapp_vm_v2" "{{.VmName}}1" {
  org             = "{{.Org}}"
  vdc             = "{{.Vdc}}"
  vapp_id         = vcd_vapp.{{.VappName}}.id
  name            = "{{.VmName}}1"
  computer_name   = "{{.ComputerName}}1"
  catalog_item_id = data.vcd_catalog_item.main_item.id
  concurrent_vms  = 3
}

resource "vcd_vapp_vm_v2" "{{.VMName}}2" {
  org = "{{.Org}}"
  vdc = "{{.Vdc}}"

  power_on = false

  vapp_id        = vcd_vapp.{{.VAppName}}.id
  description    = "test empty VM"
  concurrent_vms = 3
  name           = "{{.VMName}}2"
  memory         = 512
  cpus           = 2
  cpu_cores      = 1 
  
  os_type                        = "sles11_64Guest"
  hardware_version               = "vmx-13"
  boot_image_id                  = data.vcd_media_item.media_item.id
  expose_hardware_virtualization = true
  computer_name                  = "{{.ComputerName}}2"

  cpu_hot_add_enabled    = true
  memory_hot_add_enabled = true

}

resource "vcd_vapp_vm_v2" "{{.VmName}}3" {
  org             = "{{.Org}}"
  vdc             = "{{.Vdc}}"
  vapp_id         = vcd_vapp.{{.VappName}}.id
  name            = "{{.VmName}}3"
  computer_name   = "{{.ComputerName}}3"
  catalog_item_id = data.vcd_catalog_item.main_item.id
  concurrent_vms  = 3
}

`
