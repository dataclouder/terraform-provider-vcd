//go:build ALL || functional

package vcd

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"testing"
)

type listDef struct {
	name         string
	resourceType string
	parent       string
	knownItem    string
	vdc          string
	listMode     string
	importFile   bool
}

func TestAccVcdDatasourceResourceList(t *testing.T) {
	preTestChecks(t)

	var lists = []listDef{
		{name: "resources", resourceType: "resources", knownItem: "vcd_org"},
		{name: "global_role", resourceType: "vcd_global_role", knownItem: "vApp Author"},
		{name: "rights_bundle", resourceType: "vcd_rights_bundle", knownItem: "Default Rights Bundle"},
		{name: "right", resourceType: "vcd_right", knownItem: "Catalog: Change Owner"},

		// entities belonging to an Org don't require an explicit parent, as it is given from the Org passed in the provider
		// For each resource, we test with and without and explicit parent
		{name: "user", resourceType: "vcd_org_user"},

		// test for VM requires a VApp as parent, which may not be guaranteed, as there is none in the config file
		//{"vapp_vm", "vcd_vapp_vm", "TestVapp", ""},
	}

	if testConfig.VCD.Org != "" {
		lists = append(lists, listDef{name: "orgs", resourceType: "vcd_org", knownItem: testConfig.VCD.Org})

		// entities belonging to an Org don't require an explicit parent, as it is given from the Org passed in the provider
		// For each resource, we test with and without and explicit parent
		lists = append(lists, listDef{name: "user-parent", resourceType: "vcd_org_user", parent: testConfig.VCD.Org})
		lists = append(lists, listDef{name: "role-parent", resourceType: "vcd_role", parent: testConfig.VCD.Org, knownItem: "vApp Author"})
		if testConfig.Networking.ExternalNetwork != "" {
			lists = append(lists, listDef{name: "extent-parent", resourceType: "vcd_external_network", parent: testConfig.VCD.Org, knownItem: testConfig.Networking.ExternalNetwork})
		} else {
			fmt.Print("`Networking.ExternalNetwork` value isn't configured, datasource test will be skipped\n")
		}
		if testConfig.VCD.Catalog.Name != "" {
			// entities belonging to an Org don't require an explicit parent, as it is given from the Org passed in the provider
			// For each resource, we test with and without and explicit parent
			lists = append(lists, listDef{name: "catalog-parent", resourceType: "vcd_catalog", parent: testConfig.VCD.Org, knownItem: testConfig.VCD.Catalog.Name})
		} else {
			fmt.Print("`VCD.Catalog.Name` value isn't configured, datasource test using this will be skipped\n")
		}
		if testConfig.Nsxt.Vdc != "" {
			// entities belonging to a VDC don't require an explicit parent, as it is given from the VDC passed in the provider
			// For each resource, we test with and without and explicit parent
			lists = append(lists, listDef{name: "VDC-parent", resourceType: "vcd_org_vdc", parent: testConfig.VCD.Org, knownItem: testConfig.Nsxt.Vdc})
		} else {
			fmt.Print("`Nsxt.Vdc` value isn't configured, datasource test using this will be skipped\n")
		}
	} else {
		fmt.Print("`VCD.Org` value isn't configured, datasource test will be skipped\n")
	}

	if testConfig.Networking.ExternalNetwork != "" {
		// entities belonging to an Org don't require an explicit parent, as it is given from the Org passed in the provider
		// For each resource, we test with and without and explicit parent
		lists = append(lists, listDef{name: "extnet", resourceType: "vcd_external_network", knownItem: testConfig.Networking.ExternalNetwork})
	} else {
		fmt.Print("`Networking.ExternalNetwork` value isn't configured, datasource test using this will be skipped\n")
	}

	if testConfig.VCD.Catalog.Name != "" {
		// entities belonging to an Org don't require an explicit parent, as it is given from the Org passed in the provider
		// For each resource, we test with and without and explicit parent
		lists = append(lists, listDef{name: "catalog", resourceType: "vcd_catalog", knownItem: testConfig.VCD.Catalog.Name})

		if testConfig.VCD.Catalog.CatalogItem != "" {
			// tests in this last group always require an explicit parent
			lists = append(lists, listDef{name: "catalog_item", resourceType: "vcd_catalog_item", parent: testConfig.VCD.Catalog.Name, knownItem: testConfig.VCD.Catalog.CatalogItem})
		} else {
			fmt.Print("`VCD.CatalogItem` value isn't configured, datasource test using this will be skipped\n")
		}
		if testConfig.Media.MediaName != "" {
			// tests in this last group always require an explicit parent
			lists = append(lists, listDef{name: "catalog_media", resourceType: "vcd_catalog_media", parent: testConfig.VCD.Catalog.Name, knownItem: testConfig.Media.MediaName})
		} else {
			fmt.Print("`Media.MediaName` value isn't configured, datasource test using this will be skipped\n")
		}
	} else {
		fmt.Print("`VCD.Catalog.Name` value isn't configured, datasource test using this will be skipped\n")
	}

	if testConfig.VCD.Vdc != "" {
		// entities belonging to a VDC don't require an explicit parent, as it is given from the VDC passed in the provider
		// For each resource, we test with and without and explicit parent
		lists = append(lists, listDef{name: "network-parent", resourceType: "network", parent: testConfig.VCD.Vdc, vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "network_isolated-parent", resourceType: "vcd_network_isolated", parent: testConfig.VCD.Vdc, vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "network_routed-parent", resourceType: "vcd_network_routed", parent: testConfig.VCD.Vdc, vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "network_direct-parent", resourceType: "vcd_network_direct", parent: testConfig.VCD.Vdc, vdc: testConfig.VCD.Vdc})

		lists = append(lists, listDef{name: "network", resourceType: "network", vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "network_isolated", resourceType: "vcd_network_isolated", vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "network_routed", resourceType: "vcd_network_routed", vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "network_direct", resourceType: "vcd_network_direct", vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "ipset", resourceType: "vcd_ipset", vdc: testConfig.VCD.Vdc})

		if testConfig.Networking.EdgeGateway != "" {
			// entities belonging to a VDC don't require an explicit parent, as it is given from the VDC passed in the provider
			// For each resource, we test with and without and explicit parent
			lists = append(lists, listDef{name: "edge_gateway-parent", resourceType: "vcd_edgegateway", parent: testConfig.VCD.Vdc, knownItem: testConfig.Networking.EdgeGateway, vdc: testConfig.VCD.Vdc})
		} else {
			fmt.Print("`Networking.EdgeGateway` value isn't configured, datasource test using this will be skipped\n")
		}
	} else {
		fmt.Print("`" +
			"VCD.Vdc` value isn't configured, datasource test using this will be skipped\n")
	}

	if testConfig.Nsxt.Vdc != "" {
		// entities belonging to an Org don't require an explicit parent, as it is given from the Org passed in the provider
		// For each resource, we test with and without and explicit parent
		lists = append(lists, listDef{name: "VDC", resourceType: "vcd_org_vdc", knownItem: testConfig.Nsxt.Vdc})
		// entities belonging to a VDC don't require an explicit parent, as it is given from the VDC passed in the provider
		// For each resource, we test with and without and explicit parent
		lists = append(lists, listDef{name: "vcd_network_routed_v2", resourceType: "vcd_network_routed_v2", parent: testConfig.Nsxt.Vdc})
		lists = append(lists, listDef{name: "vcd_network_isolated_v2", resourceType: "vcd_network_isolated_v2", parent: testConfig.Nsxt.Vdc})
		lists = append(lists, listDef{name: "vcd_nsxt_network_imported", resourceType: "vcd_nsxt_network_imported", parent: testConfig.Nsxt.Vdc})
		lists = append(lists, listDef{name: "vapp-parent", resourceType: "vcd_vapp", parent: testConfig.Nsxt.Vdc})

		lists = append(lists, listDef{name: "vapp", resourceType: "vcd_vapp"})
	} else {
		fmt.Print("`Nsxt.Vdc` value isn't configured, datasource test using this will be skipped\n")
	}

	if testConfig.Networking.EdgeGateway != "" {
		// entities belonging to an Org don't require an explicit parent, as it is given from the Org passed in the provider
		// For each resource, we test with and without and explicit parent
		lists = append(lists, listDef{name: "edge_gateway", resourceType: "vcd_edgegateway", knownItem: testConfig.Networking.EdgeGateway, vdc: testConfig.VCD.Vdc})

		// tests in this last group always require an explicit parent
		lists = append(lists, listDef{name: "nsxv_data", resourceType: "vcd_nsxv_dnat", parent: testConfig.Networking.EdgeGateway, vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "nsxv_sat", resourceType: "vcd_nsxv_snat", parent: testConfig.Networking.EdgeGateway, vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "nsxv_firewall_rule", resourceType: "vcd_nsxv_firewall_rule", parent: testConfig.Networking.EdgeGateway, vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "lb_server_pool", resourceType: "vcd_lb_server_pool", parent: testConfig.Networking.EdgeGateway, vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "lb_service_monitor", resourceType: "vcd_lb_service_monitor", parent: testConfig.Networking.EdgeGateway, vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "lb_virtual_server", resourceType: "vcd_lb_virtual_server", parent: testConfig.Networking.EdgeGateway, vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "lb_app_profile", resourceType: "vcd_lb_app_profile", parent: testConfig.Networking.EdgeGateway, vdc: testConfig.VCD.Vdc})
		lists = append(lists, listDef{name: "lb_app_rule", resourceType: "vcd_lb_app_rule", parent: testConfig.Networking.EdgeGateway, vdc: testConfig.VCD.Vdc})

	} else {
		fmt.Print("`testConfig.Networking.EdgeGateway` value isn't configured, datasource test using this will be skipped\n")
	}

	lists = append(lists, listDef{name: "library_certificate", resourceType: "vcd_library_certificate"})

	lists = append(lists,
		listDef{
			name:         "testVm",
			resourceType: "vcd_vapp_vm",
			parent:       "TestVapp",
			knownItem:    "TestVm",
			vdc:          testConfig.VCD.Vdc,
			listMode:     "import",
			importFile:   true,
		},
		listDef{
			name:         "ldap-server",
			resourceType: "vcd_vm",
			knownItem:    "ldap-server",
			vdc:          testConfig.Nsxt.Vdc,
			listMode:     "import",
			importFile:   true,
		},
	)
	for _, def := range lists {
		t.Run(def.name+"-"+def.resourceType, func(t *testing.T) { runResourceInfoTest(def, t) })
	}
	postTestChecks(t)
}

func runResourceInfoTest(def listDef, t *testing.T) {

	var data = StringMap{
		"ResName":    def.name,
		"ResType":    def.resourceType,
		"ResParent":  def.parent,
		"ListMode":   "name",
		"ImportFile": "",
		"FuncName":   fmt.Sprintf("ResourceList-%s", def.name+"-"+def.resourceType),
	}
	if def.listMode != "" {
		data["ListMode"] = def.listMode
	}
	if def.importFile {
		data["ImportFile"] = fmt.Sprintf("import-%s.tf", def.resourceType)
	}
	if def.vdc == "" {
		data["Vdc"] = testConfig.Nsxt.Vdc
	} else {
		data["Vdc"] = def.vdc
	}
	var configText string
	if def.parent == "" {
		configText = templateFill(testAccCheckVcdDatasourceInfoSimple, data)
	} else {
		configText = templateFill(testAccCheckVcdDatasourceInfoWithParent, data)
	}

	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	if !usingSysAdmin() && (def.resourceType == "vcd_external_network" ||
		def.resourceType == "vcd_global_role" ||
		def.resourceType == "vcd_rights_bundle") {
		t.Skipf("test with %s requires system administrator privileges", def.resourceType)
	}
	debugPrintf("#[DEBUG] CONFIGURATION: %s", configText)

	// networks aren't in the configuration file, but we can easily search for existing ones
	if strings.HasPrefix(def.resourceType, "vcd_network") {
		err := getAvailableNetworks()
		if err == nil {
			network, ok := availableNetworks[def.resourceType]
			if ok {
				def.knownItem = network.network.Name
			}
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vcd_resource_list."+def.name, "name", def.name),
					checkListForKnownItem(def.name, def.knownItem, true, def.importFile),
				),
			},
		},
	})
}

func checkListForKnownItem(resName, target string, isWanted, importing bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if target == "" {
			return nil
		}

		resourcePath := "data.vcd_resource_list." + resName

		res, ok := s.RootModule().Resources[resourcePath]
		if !ok {
			return fmt.Errorf("resource %s not found", resName)
		}

		var list = make([]string, 0)

		for key, value := range res.Primary.Attributes {
			if strings.HasPrefix(key, "list.") {
				list = append(list, value)
			}
		}

		for _, item := range list {
			found := item == target
			if importing {
				found = strings.Contains(item, target)
			}
			if found {
				if isWanted {
					return nil
				} else {
					return fmt.Errorf("item '%s' found in '%s'", target, resName)
				}
			}
		}
		if isWanted {
			return fmt.Errorf("item '%s' not found in list %s", target, resourcePath)
		} else {
			return nil
		}
	}
}

const testAccCheckVcdDatasourceInfoSimple = `
data "vcd_resource_list" "{{.ResName}}" {
  vdc              = "{{.Vdc}}"
  name             = "{{.ResName}}"
  resource_type    = "{{.ResType}}"
  import_file_name = "{{.ImportFile}}"
  list_mode        = "{{.ListMode}}"
}
`
const testAccCheckVcdDatasourceInfoWithParent = `
data "vcd_resource_list" "{{.ResName}}" {
  vdc              = "{{.Vdc}}"
  name             = "{{.ResName}}"
  resource_type    = "{{.ResType}}"
  parent           = "{{.ResParent}}"
  import_file_name = "{{.ImportFile}}"
  list_mode        = "{{.ListMode}}"
}
`
