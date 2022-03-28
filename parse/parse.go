package parse

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func UnmarshalKV(kv map[string]string, v interface{}, tagName string) error {
	p := &props{kv: kv, tagName: tagName}
	return p.unmarshal(v)
}

type props struct {
	kv      map[string]string
	tagName string
}

func (p *props) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Elem().Type().Kind() != reflect.Struct {
		return InvalidUnmarshalError
	}

	return p.value("", rv)
}

func (p *props) value(key string, v reflect.Value) (err error) {
	switch v.Kind() {
	default:
		err = p.valueBasicType(key, v)
	case reflect.Ptr:
		err = p.value(key, v.Elem())
	case reflect.Struct:
		err = p.valueStruct(key, v)
	case reflect.Map:
		err = p.valueMap(key, v)
	case reflect.Slice:
		err = p.valueSlice(key, v)
	}

	return err
}

func (p *props) valueStruct(key string, v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		vf, tf := v.Field(i), v.Type().Field(i)

		if !vf.CanInterface() {
			continue
		}

		if vf.Kind() == reflect.Ptr {
			vf.Set(reflect.New(tf.Type.Elem()))
		}

		tag := parseTag(tf.Tag.Get(p.tagName))

		if tag == singleHorizontalBar {
			continue
		}

		if key != "" {
			tag = fmt.Sprintf("%s.%s", key, tag)
		}

		if err := p.value(tag, vf); err != nil {
			return nil
		}
	}
	return nil
}

// valueBasicType deal with int, float, bool, string
func (p *props) valueBasicType(key string, v reflect.Value) error {
	s, ok := p.get(key)
	// NOTE: if key not found, just skip over it.
	if !ok {
		return nil
	}

	switch v.Kind() {
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		spu, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(spu)
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		spi, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(spi)
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		spf, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(spf)
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		spb, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(spb)
	default:
		return UnsupportedTypeError
	}

	return nil
}

func (p *props) valueMap(key string, v reflect.Value) (err error) {
	m := reflect.MakeMap(v.Type())
	pp := p.subPropsWithDot(key)
	for kk := range pp.kv {
		mv := reflect.New(v.Type().Elem())

		vv := mv
		// allocate a new value for the pointer
		valueIsPtr := mv.Elem().Type().Kind() == reflect.Ptr
		if valueIsPtr {
			vv = reflect.New(v.Type().Elem().Elem())
		}

		mk := strings.Split(kk, ".")[0]
		err = pp.value(mk, vv)
		if err != nil {
			return
		}

		if valueIsPtr {
			mv.Elem().Set(vv)
		}

		m.SetMapIndex(reflect.ValueOf(mk), mv.Elem())
	}
	v.Set(m)
	return
}

func (p *props) valueSlice(key string, v reflect.Value) (err error) {
	var spp = map[string]*props{}
	var sepp = map[string]*props{}

	i := 0
	for {
		sk := fmt.Sprintf("%s[%d]", key, i)
		if !p.hasKeyPrefix(sk) {
			break
		}

		if pp := p.subPropsWithDot(sk); !pp.isEmpty() {
			spp[sk] = pp
		}

		if epp := p.exactSubProps(sk); !epp.isEmpty() {
			sepp[sk] = epp
		}

		i += 1
	}

	slice := reflect.MakeSlice(v.Type(), 0, len(spp))

	for ii := 0; ii < len(spp); ii++ {
		sk := fmt.Sprintf("%s[%d]", key, ii)
		pp := spp[sk]

		var ev reflect.Value
		if v.Type().Elem().Kind() == reflect.Ptr {
			ev = reflect.New(v.Type().Elem().Elem())
			if err := pp.value("", ev); err != nil {
				return err
			}
			slice = reflect.Append(slice, ev)
		} else {
			ev = reflect.New(v.Type().Elem())
			if err := pp.value("", ev); err != nil {
				return err
			}
			slice = reflect.Append(slice, ev.Elem())
		}
	}

	for ii := 0; ii < len(sepp); ii++ {
		sk := fmt.Sprintf("%s[%d]", key, ii)
		epp := sepp[sk]

		ev := reflect.New(v.Type().Elem())
		err := epp.value(sk, ev)
		if err != nil {
			return err
		}
		slice = reflect.Append(slice, ev.Elem())
	}

	v.Set(slice)
	return nil
}

func (p *props) subPropsWithDot(prefix string) *props {
	var kv = map[string]string{}

	for k, v := range p.kv {
		if strings.HasPrefix(k, prefix+dot) {
			kv[k[len(prefix)+1:]] = v
		}
	}

	return &props{kv, p.tagName}
}

func (p *props) exactSubProps(name string) *props {
	var kv = map[string]string{}

	for k, v := range p.kv {
		if k == name {
			kv[k] = v
		}
	}

	return &props{kv, p.tagName}
}

func (p *props) isEmpty() bool {
	return len(p.kv) == 0
}

func (p *props) get(k string) (string, bool) {
	v, ok := p.kv[k]
	return v, ok
}

func (p *props) hasKeyPrefix(prefix string) bool {
	for k := range p.kv {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}
