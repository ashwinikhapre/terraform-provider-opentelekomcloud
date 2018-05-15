package cluster

import "github.com/huaweicloud/golangsdk"

const (
	rootpath = "clusters"
	certificatePath = "certificates"
)

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootpath)
}

func certificateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootpath, id, certificatePath)
}