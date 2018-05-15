package cluster

import (
	"github.com/huaweicloud/golangsdk"
	"reflect"
)

// ListOpts allows the filtering  of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned.
type ListOpts struct {
	Name 	string     `json:"name"`
	ID 		string     `json:"uuid"`
	AZ 		string 	   `json:"az"`
	Type 	string     `json:"clustertype"`
	VPC		string     `json:"vpc"`
	VpcId 	string     `json:"vpcid"`
	Status  string   `json:"status"`
}

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json"},
}

// List returns collection of
// clusters. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those clusters that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(client *golangsdk.ServiceClient) (r ListResult) {

	_,r.Err = client.Get(rootURL(client), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})

	return
}

func FilterClusters(clusters []RetrievedCluster, opts ListOpts) ([]RetrievedCluster, error) {

	var refinedClusters []RetrievedCluster
	var  clusfield string
	var matched bool
	m := map[string]interface{}{}

	if opts.Name != "" {
		clusfield="Metadata"
		m["Name"] = opts.Name
	}
	if opts.ID != "" {
		clusfield="Metadata"
		m["ID"] = opts.ID
	}
	if opts.AZ != "" {
		clusfield="Clusterspec"
		m["AZ"] = opts.AZ
	}
	if opts.Type != "" {
		clusfield="Clusterspec"
		m["ClusterType"] = opts.Type
	}
	if opts.VPC != "" {
		clusfield="Clusterspec"
		m["VPC"] = opts.VPC
	}
	if opts.VpcId != "" {
		clusfield="Clusterspec"
		m["VpcId"] = opts.VpcId
	}
	if opts.Status != "" {
		clusfield="ClusterStatus"
		m["Status"] = opts.Status
	}

	if len(m) > 0 && len(clusters) > 0 {
		for _, cluster := range clusters {
			matched = true

			for key, value := range m {

				if sVal := getStructNestedField(&cluster, clusfield, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedClusters = append(refinedClusters, cluster)
			}

		}

	} else {
		refinedClusters = clusters
	}

	return refinedClusters, nil
}

func getStructNestedField(v *RetrievedCluster, clusfield string , field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(clusfield).Interface()
	r1 := reflect.ValueOf(f)
	f1 := reflect.Indirect(r1).FieldByName(field)
	return string(f1.String())
}

// GetCertificate retrieves a particular cluster certificates based on cluster ID.
func GetCertificate(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(certificateURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}
