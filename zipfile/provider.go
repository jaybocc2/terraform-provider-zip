package zipfile

import "github.com/hashicorp/terraform/helper/schema"
import "github.com/hashicorp/terraform/terraform"

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"zip_file": dataSourceZip(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"zip_file": schema.DataSourceResourceShim(
				"zip_file",
				dataSourceZip(),
			),
		},
	}
}
