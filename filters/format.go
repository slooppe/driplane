package filters

import (
	"fmt"
	"github.com/Matrix86/driplane/data"
	html "html/template"
	"io/ioutil"
	"path/filepath"
	text "text/template"
)

type Format struct {
	Base

	template     interface{}
	templateType string // html or text

	params map[string]string
}

func NewFormatFilter(p map[string]string) (Filter, error) {
	f := &Format{
		params:       p,
		templateType: "text",
	}
	f.cbFilter = f.DoFilter

	if v, ok := f.params["type"]; ok && v == "html" {
		f.templateType = "html"
	}

	if v, ok := f.params["template"]; ok {
		if f.templateType == "html" {
			t, err := html.New("formatFilterTemplate").Parse(v)
			if err != nil {
				return nil, err
			}
			f.template = t
		} else {
			t, err := text.New("formatFilterTemplate").Parse(v)
			if err != nil {
				return nil, err
			}
			f.template = t
		}
	}
	if v, ok := f.params["file"]; ok {
		fpath := v
		if v, ok := p["general.templates_path"]; !ok {
			if r, ok := p["general.rules_path"]; !ok {
				return nil, fmt.Errorf("NewJsFilter: rules_path or js_path configs not found")
			} else {
				fpath = filepath.Join(r, fpath)
			}
		} else {
			fpath = filepath.Join(v, fpath)
		}
		content, err := ioutil.ReadFile(fpath)
		if err != nil {
			return nil, err
		}
		if f.templateType == "html" {
			t, err := html.New("formatFilterTemplate").Parse(string(content))
			if err != nil {
				return nil, err
			}
			f.template = t
		} else {
			t, err := text.New("formatFilterTemplate").Parse(string(content))
			if err != nil {
				return nil, err
			}
			f.template = t
		}
	}

	return f, nil
}

func (f *Format) DoFilter(msg *data.Message) (bool, error) {
	txt, err := msg.ApplyPlaceholder(f.template)
	if err != nil {
		return false, err
	}
	msg.SetMessage(txt)
	return true, nil
}

// Set the name of the filter
func init() {
	register("format", NewFormatFilter)
}
