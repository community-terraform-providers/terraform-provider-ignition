package ignition

import (
	"encoding/json"

	"github.com/coreos/ignition/v2/config/v3_6/types"
	"github.com/coreos/ignition/v2/config/validate"
	"github.com/coreos/vcontext/path"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLink() *schema.Resource {
	return &schema.Resource{
		Exists: resourceLinkExists,
		Read:   resourceLinkRead,
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
			"target": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hard": {
				Type:     schema.TypeBool,
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

func resourceLinkRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildLink(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceLinkExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildLink(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildLink(d *schema.ResourceData) (string, error) {
	link := &types.Link{}
	link.Path = d.Get("path").(string)
	if err := handleReport(link.Validate(path.ContextPath{})); err != nil {
		return "", err
	}

	overwrite := d.Get("overwrite").(bool)
	link.Overwrite = &overwrite

	target := d.Get("target").(string)
	link.Target = &target

	hard, hasHard := d.GetOk("hard")
	if hasHard {
		bhard := hard.(bool)
		link.Hard = &bhard
	}

	user := d.Get("user").(string)
	if user != "" {
		link.User = types.NodeUser{Name: &user}
	}

	uid := d.Get("uid").(int)
	if uid != 0 {
		link.User = types.NodeUser{ID: &uid}
	}

	group := d.Get("group").(string)
	if group != "" {
		link.Group = types.NodeGroup{Name: &group}
	}

	gid := d.Get("gid").(int)
	if gid != 0 {
		link.Group = types.NodeGroup{ID: &gid}
	}

	b, err := json.Marshal(link)
	if err != nil {
		return "", err
	}
	err = d.Set("rendered", string(b))
	if err != nil {
		return "", err
	}

	return hash(string(b)), handleReport(validate.ValidateWithContext(new(*types.Link), b))
}
