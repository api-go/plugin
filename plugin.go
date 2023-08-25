package plugin

import (
	"github.com/ssgo/log"
	"github.com/ssgo/s"
	"reflect"
	"sync"
)

type Plugin struct {
	Name    string
	Objects map[string]interface{}
}

type Context struct {
	Args     map[string]interface{}
	Logger   *log.Logger
	Request  *s.Request
	Response *s.Response
	injects  map[string]interface{}
	data     map[string]interface{}
}

func NewContext(args map[string]interface{}, logger *log.Logger, request *s.Request, response *s.Response, injects map[string]interface{}) *Context {
	return &Context{
		Args:     args,
		Logger:   logger,
		Request:  request,
		Response: response,
		injects:  injects,
		data:     map[string]interface{}{},
	}
}

func (ctx *Context) SetInject(obj interface{}) {
	ctx.injects[reflect.ValueOf(obj).Type().String()] = obj
}

func (ctx *Context) GetInject(typ string) interface{} {
	return ctx.injects[typ]
}

func (ctx *Context) GetData(k string) interface{} {
	return ctx.data[k]
}

func (ctx *Context) SetData(k string, v interface{}) {
	ctx.data[k] = v
}

var pluginByName = map[string]*Plugin{}
var pluginsLock = sync.RWMutex{}

func Register(p Plugin) {
	pluginsLock.Lock()
	defer pluginsLock.Unlock()
	pluginByName[p.Name] = &p
}

func Update(plugName, objectName string, value interface{}) {
	pluginsLock.RLock()
	p := pluginByName[plugName]
	pluginsLock.RUnlock()

	if p != nil {
		pluginsLock.Lock()
		if p.Objects == nil {
			p.Objects = map[string]interface{}{}
		}
		p.Objects[objectName] = value
		pluginsLock.Unlock()
	}
}

func List() []Plugin {
	pluginsLock.RLock()
	defer pluginsLock.RUnlock()
	out := make([]Plugin, len(pluginByName))
	i := 0
	for _, p := range pluginByName {
		out[i] = *p
		i++
	}
	return out
}

func Get(name string) Plugin {
	pluginsLock.RLock()
	defer pluginsLock.RUnlock()
	return *pluginByName[name]
}
