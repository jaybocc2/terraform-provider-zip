package zipfile

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/pathorcontents"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceZip() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZipRead,
		Schema: map[string]*schema.Schema{
			"files": &schema.Schema{
				Type:         schema.TypeMap,
				Required:     true,
				Description:  "files to zip",
				ValidateFunc: validateFilesAttribute,
			},
			"generated": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "generate the zipfile",
			},
			"sha256": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "sha256 of the generated zipfile",
			},
		},
	}
}

func dataSourceZipRead(d *schema.ResourceData, meta interface{}) error {
	generated, err := zipFiles(d)
	if err != nil {
		return err
	}

	d.Set("generated", generated)
	d.Set("sha256", hash256(generated))
	d.SetId(hash(generated))
	return nil
}

func zipFiles(d *schema.ResourceData) (string, error) {
	buf := new(bytes.Buffer)

	w := zip.NewWriter(buf)

	vars := d.Get("files").(map[string]interface{})

	for k, v := range vars {
		contents, _, err := pathorcontents.Read(v.(string))
		if err != nil {
			return "", err
		}

		f, err := w.Create(k)
		if err != nil {
			return "", err
		}

		_, err = f.Write([]byte(contents))
		if err != nil {
			return "", err
		}
	}
	err := w.Close()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func hash(s string) string {
	sha := sha512.Sum512([]byte(s))
	return hex.EncodeToString(sha[:])
}

func hash256(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}

func validateFilesAttribute(i interface{}, key string) (ws []string, es []error) {
	var badFiles []string
	for k, v := range i.(map[string]interface{}) {
		switch v.(type) {
		case []interface{}:
			badFiles = append(badFiles, fmt.Sprintf("%s (list)", k))
		case map[string]interface{}:
			badFiles = append(badFiles, fmt.Sprintf("%s (map)", k))
		case int:
			badFiles = append(badFiles, fmt.Sprintf("%s (int)", k))
		case bool:
			badFiles = append(badFiles, fmt.Sprintf("%s (bool)", k))
		case float64:
			badFiles = append(badFiles, fmt.Sprintf("%s (float)", k))
		}
	}
	if len(badFiles) > 0 {
		es = append(es, fmt.Errorf(
			"%s: cannot contain non-strings; bad keys: %s",
			key, strings.Join(badFiles, ", ")))
	}
	return
}
