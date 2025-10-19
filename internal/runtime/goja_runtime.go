package runtime

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
)

type RuntimeGoja struct {
	vm *goja.Runtime
}

func NewRuntimeGoja() (*RuntimeGoja, error) {
	vm := goja.New()
	vm.SetParserOptions(parser.WithDisableSourceMaps)

	vm.Set("__secret__", func(call goja.FunctionCall) goja.Value {
		slog.Debug("__secret__ called", "call", call)
		return vm.ToValue("123xyz")
	})

	vm.Set("fetch", func(call goja.FunctionCall) goja.Value {
		url := call.Argument(0).String()
		slog.Debug("__fetch__ called with URL", "url", url)
		return vm.ToValue("100F")
	})

	return &RuntimeGoja{vm: vm}, nil
}

func (r *RuntimeGoja) Compile(code string) error {
	prog, err := goja.Compile(code, "math.js", true)
	if err != nil {
		return err
	}

	fmt.Printf("[Compile] prog: %v\n", prog)
	return nil
}

func (r *RuntimeGoja) RunString(code string) error {
	val, err := r.vm.RunString(code)
	if val != nil {
		fmt.Printf("[RunString] val: %s\n", val.String())
	}
	return err
}

func (r *RuntimeGoja) RunScript(name, script string) error {
	_, err := r.vm.RunScript(name, script)
	return err
}

func (r *RuntimeGoja) SerializeGlobal() ([]byte, error) {
	globals := r.vm.GlobalObject()
	return json.Marshal(globals)
}

func (r *RuntimeGoja) Snapshot() error {
	// Get the current VM context and symbols
	ctx := r.vm
	f1, err := os.Create("/tmp/ctx-" + time.Now().Format("20060102150405") + ".gob")
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(f1)
	err = enc.Encode(ctx)
	if err != nil {
		panic(err)
	}
	f1.Close()

	return nil
}
