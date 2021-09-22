//go:build vapp || vm || ALL || functional
// +build vapp vm ALL functional

package vcd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdVAppVmV1_Basic(t *testing.T) {
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
		"ComputerName": "compname",
		"Tags":         "vapp vm",
	}

	configText := templateFill(testAccCheckVcdVAppVmV1_basic, params)
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
					testAccCheckVcdVAppVmExists(vappName, vmName+"1", "vcd_vapp_vm."+vmName+"1", nil, nil),
					testAccCheckVcdVAppVmExists(vappName, vmName+"2", "vcd_vapp_vm."+vmName+"2", nil, nil),
					testAccCheckVcdVAppVmExists(vappName, vmName+"3", "vcd_vapp_vm."+vmName+"3", nil, nil),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm."+vmName+"1", "name", vmName+"1"),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm."+vmName+"2", "name", vmName+"2"),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm."+vmName+"3", "name", vmName+"3"),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm."+vmName+"1", "computer_name", "compname1"),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm."+vmName+"2", "computer_name", "compname2"),
					resource.TestCheckResourceAttr(
						"vcd_vapp_vm."+vmName+"3", "computer_name", "compname3"),
				),
			},
			resource.TestStep{
				ResourceName:      "vcd_vapp_vm." + vmName + "1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdVappObject(testConfig, vappName, vmName+"1"),
				// These fields can't be retrieved from user data
				ImportStateVerifyIgnore: []string{"template_name", "catalog_name",
					"accept_all_eulas", "power_on", "computer_name", "prevent_update_power_off",
					"concurrent_vms", "vapp_id", "catalog_item_id"},
			},
		},
	})
	postTestChecks(t)
}

const testAccCheckVcdVAppVmV1_basic = `

resource "vcd_vapp" "{{.VappName}}" {
  name = "{{.VappName}}"
  org  = "{{.Org}}"
  vdc  = "{{.Vdc}}"
}

resource "vcd_vapp_vm" "{{.VmName}}1" {
  org             = "{{.Org}}"
  vdc             = "{{.Vdc}}"
  vapp_name       = vcd_vapp.{{.VappName}}.name
  name            = "{{.VmName}}1"
  computer_name   = "{{.ComputerName}}1"
  catalog_name    = "{{.Catalog}}"
  template_name  = "{{.CatalogItem}}"
  memory          = 1024
  cpus            = 2
  cpu_cores       = 1
  network {
     adapter_type       = "VMXNET3"
     connected          = false
     ip_allocation_mode = "NONE"
     is_primary         = true
     mac                = "00:50:56:29:00:de"
     type               = "none"
  }
}

resource "vcd_vapp_vm" "{{.VmName}}2" {
  org = "{{.Org}}"
  vdc = "{{.Vdc}}"

  power_on = false

  vapp_name       = vcd_vapp.{{.VappName}}.name
  description    = "test empty VM"
  name           = "{{.VmName}}2"
  memory         = 512
  cpus           = 2
  cpu_cores      = 1 
  
  os_type                        = "sles11_64Guest"
  hardware_version               = "vmx-13"
  catalog_name                   = "{{.Catalog}}"
  boot_image                     = "{{.MediaItem}}"
  expose_hardware_virtualization = true
  computer_name                  = "{{.ComputerName}}2"

  cpu_hot_add_enabled    = true
  memory_hot_add_enabled = true

}

resource "vcd_vapp_vm" "{{.VmName}}3" {
  org             = "{{.Org}}"
  vdc             = "{{.Vdc}}"
  vapp_name       = vcd_vapp.{{.VappName}}.name
  name            = "{{.VmName}}3"
  computer_name   = "{{.ComputerName}}3"
  catalog_name    = "{{.Catalog}}"
  template_name   = "{{.CatalogItem}}"
  memory          = 1024
  cpus            = 2
  cpu_cores       = 1
  network {
     adapter_type       = "VMXNET3"
     connected          = false
     ip_allocation_mode = "NONE"
     is_primary         = true
     mac                = "00:50:56:29:00:de"
     type               = "none"
  }

}

`
