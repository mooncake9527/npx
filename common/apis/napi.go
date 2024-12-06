package apis

import (
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mooncake9527/x/xerrors/xerror"
	"reflect"
	"strings"
	"sync"
)

const (
	_ uint8 = iota
	json
	xml
	yaml
	form
	query
)

var constructor = &bindConstructor{}

type bindConstructor struct {
	cache map[string][]uint8
	mux   sync.Mutex
}

func (e *bindConstructor) GetBindingForGin(d interface{}) []binding.Binding {
	bs := e.getBinding(reflect.TypeOf(d).String())
	if bs == nil {
		//重新构建
		bs = e.resolve(d)
	}
	gbs := make([]binding.Binding, 0)
	mp := make(map[uint8]binding.Binding, 0)
	for _, b := range bs {
		switch b {
		case json:
			mp[json] = binding.JSON
		case xml:
			mp[xml] = binding.XML
		case yaml:
			mp[yaml] = binding.YAML
		case form:
			mp[form] = binding.Form
		case query:
			mp[query] = binding.Query
		default:
			mp[0] = nil
		}
	}
	for e := range mp {
		gbs = append(gbs, mp[e])
	}
	return gbs
}

func (e *bindConstructor) resolve(d interface{}) []uint8 {
	bs := make([]uint8, 0)
	qType := reflect.TypeOf(d).Elem()
	var tag reflect.StructTag
	var ok bool

	for i := 0; i < qType.NumField(); i++ {
		tag = qType.Field(i).Tag
		if _, ok = tag.Lookup("json"); ok {
			bs = append(bs, json)
		}
		if _, ok = tag.Lookup("xml"); ok {
			bs = append(bs, xml)
		}
		if _, ok = tag.Lookup("yaml"); ok {
			bs = append(bs, yaml)
		}
		if _, ok = tag.Lookup("form"); ok {
			bs = append(bs, form)
		}
		if _, ok = tag.Lookup("query"); ok {
			bs = append(bs, query)
		}
		if _, ok = tag.Lookup("uri"); ok {
			bs = append(bs, 0)
		}
		if t, ok := tag.Lookup("binding"); ok && strings.Index(t, "dive") > -1 {
			qValue := reflect.ValueOf(d)
			bs = append(bs, e.resolve(qValue.Field(i))...)
			continue
		}
		if t, ok := tag.Lookup("validate"); ok && strings.Index(t, "dive") > -1 {
			qValue := reflect.ValueOf(d)
			bs = append(bs, e.resolve(qValue.Field(i))...)
		}
	}
	return bs
}

func (e *bindConstructor) getBinding(name string) []uint8 {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.cache[name]
}

func (e *bindConstructor) setBinding(name string, bs []uint8) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e.cache == nil {
		e.cache = make(map[string][]uint8)
	}
	e.cache[name] = bs
}

type NApi struct {
	Api
}

// Bind 参数校验
func (e *Api) Bind(ctx *gin.Context, d interface{}, bindings ...binding.Binding) error {
	var err error
	if len(bindings) == 0 {
		bindings = constructor.GetBindingForGin(d)
	}
	for i := range bindings {
		if bindings[i] == nil {
			err = ctx.ShouldBindUri(d)
		} else {
			err = ctx.ShouldBindWith(d, bindings[i])
		}
		if err != nil && err.Error() == "EOF" {
			//slog.Warn("request body is not present anymore. ")
			err = nil
			continue
		}
		if err != nil {
			return err
		}
	}
	if err := vd.Validate(d); err != nil {
		return xerror.NewC(400, err.Error())
	}
	if dv, ok := d.(ValidI); ok {
		if err := dv.Valid(); err != nil {
			return xerror.NewC(400, err.Error())
		}
	}
	return nil
}

type ValidI interface {
	Valid() error
}
