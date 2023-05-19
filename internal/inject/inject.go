package inject

// TODO(black_desk): remove this package and switch to google/wire

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	. "github.com/black-desk/deepin-network-proxy-manager/internal/log"
	"github.com/black-desk/deepin-network-proxy-manager/pkg/location"
)

type Container struct {
	store sync.Map
}

var defaultContainer = &Container{}

func Default() *Container {
	return defaultContainer
}

func New() *Container {
	return &Container{}
}

func (c *Container) Register(v any) (err error) {
	rtype := reflect.TypeOf(v)
	if _, loaded := c.store.LoadOrStore(rtype, reflect.ValueOf(v)); loaded {
		err = fmt.Errorf(location.Capture()+
			`Type "%s" had been registered.`, rtype.String())
		return
	} else {
		Log.Debugw("Register new type", "type", rtype.String())
	}

	return
}

func (c *Container) RegisterI(ptrToI any) (err error) {
	rtype := reflect.TypeOf(ptrToI)
	if rtype.Kind() != reflect.Pointer {
		err = fmt.Errorf(location.Capture()+
			`Wrong type: %s`, rtype.String())
		return
	}

	elem := rtype.Elem()
	if elem.Kind() != reflect.Interface {
		err = fmt.Errorf(location.Capture()+
			`Wrong type: %s`, rtype.String())
		return
	}

	if _, loaded := c.store.LoadOrStore(elem, reflect.ValueOf(ptrToI).Elem()); loaded {
		err = fmt.Errorf(location.Capture()+
			`Interface "%s" had been registered.`, elem.String())
		return
	} else {
		Log.Debugw("Register new interface", "interface", elem.String())
	}

	return
}

func (c *Container) Fill(v any) (err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf(location.Capture()+
			"Failed to fill %#v:\n%w", v, err)
	}()

	if v == nil {
		err = fmt.Errorf(location.Capture() +
			"Fill should not take a nil.")
		return
	}

	rvalue := reflect.ValueOf(v)
	if rvalue.Kind() != reflect.Pointer {
		err = fmt.Errorf(location.Capture() +
			`Fill should always take a pointer as argument.`)
		return
	}

	elem := rvalue.Elem()
	if value, loaded := c.store.Load(elem.Type()); loaded {
		rvalue := reflect.ValueOf(v).Elem()
		rvalue.Set(value.(reflect.Value))
		return
	}

	if elem.Kind() != reflect.Struct {
		err = fmt.Errorf(location.Capture()+
			`Type %s not found in this container.`, elem.Type().String())
		return
	}

	for i := 0; i < elem.NumField(); i++ {
		if _, ok := elem.Type().Field(i).Tag.Lookup("inject"); !ok {
			continue
		}
		if err = c.Fill(
			reflect.NewAt(
				elem.Field(i).Type(),
				unsafe.Pointer(elem.Field(i).Addr().Pointer()),
			).Interface(),
		); err != nil {
			err = fmt.Errorf(location.Capture()+
				"Failed on field %d:\n%w", i, err)
			return
		}
	}

	return nil
}
