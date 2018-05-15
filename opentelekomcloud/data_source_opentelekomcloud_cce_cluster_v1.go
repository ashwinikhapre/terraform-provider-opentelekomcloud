package opentelekomcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
	"log"
	"github.com/huaweicloud/golangsdk/openstack/cce/v1/cluster"
)

func dataSourceCCEClusterV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCCEClusterV1Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"k8s_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"az": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"hosts": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"az": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_size": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume_type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"sshkey": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cacrt": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_crt": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCCEClusterV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV1Client(GetRegion(d, config))

	listOpts := cluster.ListOpts{
		ID:     d.Get("id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
		Type:	d.Get("type").(string),
		AZ:     d.Get("az").(string),
		VPC:    d.Get("vpc_name").(string),
		VpcId:  d.Get("vpc_id").(string),
	}

	refinedClusters, err := cluster.List(cceClient).ExtractCluster(listOpts)
	log.Printf("[DEBUG] Value of allClusters: %#v", refinedClusters)
	if err != nil {
		return fmt.Errorf("Unable to retrieve clusters: %s", err)
	}

	if len(refinedClusters) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedClusters) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Cluster := refinedClusters[0]

	var v []map[string]interface{}
	for _, volume := range Cluster.Clusterspec.ClusterHostList.HostListSpec.HostList[0].Hostspec.Volume{
		mapping := map[string]interface{}{
			"disk_type":   volume.DiskType,
			"disk_size":   volume.DiskSize,
			"volume_type": volume.VolumeType,
		}
		v = append(v, mapping)
	}

	var h []map[string]interface{}
	for _, hosts := range Cluster.Clusterspec.ClusterHostList.HostListSpec.HostList{
		mapping := map[string]interface{}{
			"name":		  hosts.Metadata.Name,
			"id":         hosts.Metadata.ID,
			"private_ip": hosts.Hostspec.PrivateIp,
			"public_ip":  hosts.Hostspec.PublicIp,
			"flavor":     hosts.Hostspec.Flavor,
			"az":         hosts.Hostspec.AZ,
			"sshkey":     hosts.Hostspec.SshKey,
			"status":     hosts.Status,
			"volume":     v,
		}
		h = append(h, mapping)
	}

	log.Printf("[DEBUG] Retrieved Clusters using given filter %s: %+v", Cluster.Metadata.ID, Cluster)
	d.SetId(Cluster.Metadata.ID)

	d.Set("id", Cluster.Metadata.ID)
	d.Set("name", Cluster.Metadata.Name)
	d.Set("status", Cluster.ClusterStatus.Status)
	d.Set("k8s_version", Cluster.K8sVersion)
	d.Set("az", Cluster.Clusterspec.AZ)
	d.Set("cpu", Cluster.Clusterspec.CPU)
	d.Set("type", Cluster.Clusterspec.ClusterType)
	d.Set("vpc_name", Cluster.Clusterspec.VPC)
	d.Set("vpc_id", Cluster.Clusterspec.VpcId)
	d.Set("subnet", Cluster.Clusterspec.Subnet)
	d.Set("endpoint", Cluster.Clusterspec.Endpoint)
	d.Set("external_endpoint", Cluster.Clusterspec.ExternalEndpoint)
	d.Set("security_group_id", Cluster.Clusterspec.SecurityGroupId)
	d.Set("region", GetRegion(d, config))
	if err := d.Set("hosts", h); err != nil {
		return err
	}

	n, err := cluster.GetCertificate(cceClient,Cluster.Metadata.ID).Extract()
	log.Printf("[DEBUG] Retrieved n %+v", n)
	d.Set("cacrt", n.Cacrt)
	d.Set("client_crt", n.ClientCrt)
	d.Set("client_key", n.ClientKey)

	return nil
}