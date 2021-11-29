// Code generated by pluginator on PrefixTransformer; DO NOT EDIT.
// pluginator {unknown  1970-01-01T00:00:00Z  }

package builtins

import (
	"errors"

	"sigs.k8s.io/kustomize/api/filters/prefix"

	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Add the given prefix to the field
type PrefixTransformerPlugin struct {
	Prefix     string        `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	FieldSpecs types.FsSlice `json:"fieldSpecs,omitempty" yaml:"fieldSpecs,omitempty"`
}

var prefixFieldSpecsToSkip = types.FsSlice{
	{Gvk: resid.Gvk{Kind: "CustomResourceDefinition"}},
	{Gvk: resid.Gvk{Group: "apiregistration.k8s.io", Kind: "APIService"}},
	{Gvk: resid.Gvk{Kind: "Namespace"}},
}

func (p *PrefixTransformerPlugin) Config(
	_ *resmap.PluginHelpers, c []byte) (err error) {
	p.Prefix = ""
	p.FieldSpecs = nil
	err = yaml.Unmarshal(c, p)
	if err != nil {
		return
	}
	if p.FieldSpecs == nil {
		return errors.New("fieldSpecs is not expected to be nil")
	}
	return
}

func (p *PrefixTransformerPlugin) Transform(m resmap.ResMap) error {
	// Even if the Prefix is empty we want to proceed with the
	// transformation. This allows to add contextual information
	// to the resources (AddNamePrefix).
	for _, r := range m.Resources() {
		// TODO: move this test into the filter (i.e. make a better filter)
		if p.shouldSkip(r.OrgId()) {
			continue
		}
		id := r.OrgId()
		// current default configuration contains
		// only one entry: "metadata/name" with no GVK
		for _, fs := range p.FieldSpecs {
			// TODO: this is redundant to filter (but needed for now)
			if !id.IsSelected(&fs.Gvk) {
				continue
			}
			// TODO: move this test into the filter.
			if fs.Path == "metadata/name" {
				// "metadata/name" is the only field.
				// this will add a prefix to the resource
				// even if it is empty

				r.AddNamePrefix(p.Prefix)
				if p.Prefix != "" {
					r.StorePreviousId()
				}
			}
			if err := r.ApplyFilter(prefix.Filter{
				Prefix:    p.Prefix,
				FieldSpec: fs,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *PrefixTransformerPlugin) shouldSkip(id resid.ResId) bool {
	for _, path := range prefixFieldSpecsToSkip {
		if id.IsSelected(&path.Gvk) {
			return true
		}
	}
	return false
}

func NewPrefixTransformerPlugin() resmap.TransformerPlugin {
	return &PrefixTransformerPlugin{}
}
