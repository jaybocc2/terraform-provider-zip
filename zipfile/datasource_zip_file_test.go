package zipfile

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testProviders = map[string]terraform.ResourceProvider{
	"zip": Provider(),
}

func TestZipGeneration(t *testing.T) {
	cases := []struct {
		files string
		want  string
	}{
		{
			`{derp.py = "print 'hello world'"}`,
			"PK\x03\x04\x14\x00\b\x00\b\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\a\x00\x00\x00derp.py*(\xca\xcc+QP\xcfH\xcd\xc9\xc9W(\xcf/\xcaIQ\a\x04\x00\x00\xff\xffPK\a\b\xe0\x14\xa1{\x19\x00\x00\x00\x13\x00\x00\x00PK\x01\x02\x14\x00\x14\x00\b\x00\b\x00\x00\x00\x00\x00\xe0\x14\xa1{\x19\x00\x00\x00\x13\x00\x00\x00\a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00derp.pyPK\x05\x06\x00\x00\x00\x00\x01\x00\x01\x005\x00\x00\x00N\x00\x00\x00\x00\x00",
		},
		{
			`
			{
					derp.py = "print 'Hello world'"
					herp.py = "print 'GoodBye world'"
			}`,
			"PK\x03\x04\x14\x00\b\x00\b\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\a\x00\x00\x00derp.py*(\xca\xcc+QP\xf7H\xcd\xc9\xc9W(\xcf/\xcaIQ\a\x04\x00\x00\xff\xffPK\a\b\x18ϐc\x19\x00\x00\x00\x13\x00\x00\x00PK\x03\x04\x14\x00\b\x00\b\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\a\x00\x00\x00herp.py*(\xca\xcc+QPw\xcf\xcfOq\xaaLU(\xcf/\xcaIQ\a\x04\x00\x00\xff\xffPK\a\b0\x01K;\x1b\x00\x00\x00\x15\x00\x00\x00PK\x01\x02\x14\x00\x14\x00\b\x00\b\x00\x00\x00\x00\x00\x18ϐc\x19\x00\x00\x00\x13\x00\x00\x00\a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00derp.pyPK\x01\x02\x14\x00\x14\x00\b\x00\b\x00\x00\x00\x00\x000\x01K;\x1b\x00\x00\x00\x15\x00\x00\x00\a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00N\x00\x00\x00herp.pyPK\x05\x06\x00\x00\x00\x00\x02\x00\x02\x00j\x00\x00\x00\x9e\x00\x00\x00\x00\x00",
		},
	}

	for _, tt := range cases {
		resource.UnitTest(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config: testZipFileConfig(tt.files),
					Check: func(s *terraform.State) error {
						got := s.RootModule().Outputs["generated"]
						if got.Value != tt.want {
							return fmt.Errorf("zip input:\n%s\ngot:\n%s\nwant:\n%s\n", tt.files, got, tt.want)
						}
						return nil
					},
				},
			},
		})
	}
}

func TestValidateFilesAttribute(t *testing.T) {
	cases := map[string]struct {
		Files     map[string]interface{}
		ExpectErr string
	}{
		"lists are invalid": {
			map[string]interface{}{
				"list": []interface{}{},
			},
			`files: cannot contain non-strings`,
		},
		"maps are invalid": {
			map[string]interface{}{
				"map": map[string]interface{}{},
			},
			`files: cannot contain non-strings`,
		},
		"ints are invalid": {
			map[string]interface{}{
				"int": 1,
			},
			`files: cannot contain non-strings`,
		},
		"bools are invalid": {
			map[string]interface{}{
				"bool": true,
			},
			`files: cannot contain non-strings`,
		},
		"floats are invalid": {
			map[string]interface{}{
				"float": float64(1.0),
			},
			`files: cannot contain non-strings`,
		},
		"strings are OK": {
			map[string]interface{}{
				"string": "foo",
			},
			``,
		},
	}

	for tn, tc := range cases {
		_, es := validateFilesAttribute(tc.Files, "files")

		if len(es) > 0 {
			if tc.ExpectErr == "" {
				t.Fatalf("%s: expected no err got %#v", tn, es)
			}
			if !strings.Contains(es[0].Error(), tc.ExpectErr) {
				t.Fatalf("%s: expected\n%s\nto contain\n%s", tn, es[0], tc.ExpectErr)
			}
		} else if tc.ExpectErr != "" {
			t.Fatalf("%s: expected err containing %q, got none!", tn, tc.ExpectErr)
		}
	}
}

func testZipFileConfig(files string) string {
	return fmt.Sprintf(`
				data "zip_file" "t0" {
						files %s
				}
				output "generated" {
						value = "${data.zip_file.t0.generated}"
				}`, files)
}
