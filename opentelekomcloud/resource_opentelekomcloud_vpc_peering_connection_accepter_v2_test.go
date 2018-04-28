package opentelekomcloud

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/peerings"
)

func TestAccOTCVpcPeeringConnectionAccepterV2_basic(t *testing.T) {
	var peering peerings.Peering

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCVpcPeeringConnectionAccepterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCVpcPeeringConnectionAccepterV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCVpcPeeringConnectionV2Exists("opentelekomcloud_vpc_peering_connection_accepter_v2.peer", &peering),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_peering_connection_accepter_v2.peer", "name", "opentelekomcloud_acc_peering"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_vpc_peering_connection_accepter_v2.peer", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckOTCVpcPeeringConnectionAccepterDestroy(s *terraform.State) error {
	// We don't destroy the underlying VPC Peering Connection.
	return nil
}

var testAccOTCVpcPeeringConnectionAccepterV2_basic = fmt.Sprintf(`

provider "opentelekomcloud"  {
    alias = "main"
}

provider "opentelekomcloud"  {
    alias = "peer"
    tenant_id   = "%s"
}

resource "opentelekomcloud_vpc_v1" "vpc_1" {
	provider = "opentelekomcloud.main"
	name = "otc_vpc_1"
  	cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_v1" "vpc_2" {
	provider = "opentelekomcloud.peer"
	name = "otc_vpc_2"
	cidr = "192.168.0.0/16"
}

resource "opentelekomcloud_vpc_peering_connection_v2" "peering_1" {
    provider = "opentelekomcloud.main"
	name = "opentelekomcloud_acc_peering"
    vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
    peer_vpc_id = "${opentelekomcloud_vpc_v1.vpc_2.id}"
	peer_tenant_id = "%s"
  }

resource "opentelekomcloud_vpc_peering_connection_accepter_v2" "peer" {
	provider = "opentelekomcloud.peer"
  	vpc_peering_connection_id = "${opentelekomcloud_vpc_peering_connection_v2.peering_1.id}"
  	accept = true

}
`, OS_PEER_TENANT_ID, OS_PEER_TENANT_ID)
