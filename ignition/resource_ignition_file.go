package ignition

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/coreos/ignition/v2/config/v3_5/types"
	"github.com/coreos/vcontext/path"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		Exists: resourceFileExists,
		Read:   resourceFileRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"overwrite": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"contents": {
				ConflictsWith: []string{"content", "source"},
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				Elem:          configReferenceResource,
			},
			"content": {
				ConflictsWith: []string{"contents"},
				Deprecated:    "Use contents.source instead",
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mime": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "text/plain",
						},

						"content": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"source": {
				ConflictsWith: []string{"contents"},
				Deprecated:    "Use contents instead",
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				Elem:          configReferenceResource,
			},
			"mode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"user": {
				ConflictsWith: []string{"uid"},
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
			},
			"uid": {
				ConflictsWith: []string{"user"},
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
			},
			"group": {
				ConflictsWith: []string{"gid"},
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
			},
			"gid": {
				ConflictsWith: []string{"group"},
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
			},
			"rendered": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFileRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildFile(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceFileExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildFile(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildFile(d *schema.ResourceData) (string, error) {
	_, hasContent := d.GetOk("content")
	_, hasContents := d.GetOk("contents")
	_, hasSource := d.GetOk("source")
	if (hasContent || hasContents) && hasSource {
		return "", fmt.Errorf("contents and source options are incompatible")
	}

	if !(hasContent || hasContents) && !hasSource {
		return "", fmt.Errorf("contents or source options must be present")
	}

	var contents types.Resource
	if hasContent {
		s := encodeDataURL(
			d.Get("content.0.mime").(string),
			d.Get("content.0.content").(string),
		)
		contents.Source = &s
	}

	if hasContents {
		if err := fillResource(d, &contents, "contents"); err != nil {
			return "", err
		}
	}

	if hasSource {
		if err := fillResource(d, &contents, "source"); err != nil {
			return "", err
		}
	}

	file := &types.File{}

	file.Path = d.Get("path").(string)

	overwrite := d.Get("overwrite").(bool)
	file.Overwrite = &overwrite

	file.Contents = contents

	mode, hasMode := d.GetOk("mode")
	if hasMode {
		imode := mode.(int)
		file.Mode = &imode
	}

	user := d.Get("user").(string)
	if user != "" {
		file.User = types.NodeUser{Name: &user}
	}

	uid := d.Get("uid").(int)
	if uid != 0 {
		file.User = types.NodeUser{ID: &uid}
	}

	group := d.Get("group").(string)
	if group != "" {
		file.Group = types.NodeGroup{Name: &group}
	}

	gid := d.Get("gid").(int)
	if gid != 0 {
		file.Group = types.NodeGroup{ID: &gid}
	}

	if err := handleReport(file.Validate(path.ContextPath{})); err != nil {
		return "", err
	}

	b, err := json.Marshal(file)
	if err != nil {
		return "", err
	}
	err = d.Set("rendered", string(b))
	if err != nil {
		return "", err
	}

	return hash(string(b)), nil
}

func fillResource(d *schema.ResourceData, contents *types.Resource, name string) error {
	src := d.Get(name + ".0.source").(string)
	if src != "" {
		contents.Source = &src
	}
	compression := d.Get(name + ".0.compression").(string)
	if compression != "" {
		contents.Compression = &compression
	}
	v := d.Get(name + ".0.verification").(string)
	if v != "" {
		contents.Verification.Hash = &v
	}
	for _, hh := range d.Get(name + ".0.http_headers").([]interface{}) {
		h, err := buildConfigHTTPHeaderReference(hh.(map[string]interface{}))
		if err != nil {
			return err
		}
		contents.HTTPHeaders = append(contents.HTTPHeaders, h)
	}
	return nil
}

func encodeDataURL(mime, content string) string {
	base64 := base64.StdEncoding.EncodeToString([]byte(content))
	return fmt.Sprintf("data:%s;charset=utf-8;base64,%s", mime, base64)
}
