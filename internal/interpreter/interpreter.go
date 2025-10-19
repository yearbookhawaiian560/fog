package interpreter

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dop251/goja"

	"github.com/0div/fog/internal/ast"
)

type Interpreter struct {
	vm *goja.Runtime
}

func NewInterpreter() (*Interpreter, error) {
	return &Interpreter{goja.New()}, nil
}

func (i *Interpreter) Eval(ctx context.Context, node ast.Node) (goja.Value, error) {
	slog.Debug("entering node", "type", node.Type, "name", node.Name, "value", node.Value, "js", node.JS, "childrenLen", len(node.Children))

	// Recursively evaluate children first and accumulate their results.
	evaluatedChildren := make([]goja.Value, len(node.Children))
	for idx, child := range node.Children {
		slog.Debug("evaluating child", "idx", idx+1, "len", len(node.Children), "type", child.Type, "parentType", node.Type, "parentName", node.Name)
		val, err := i.Eval(ctx, *child)
		if err != nil {
			slog.Error("error evaluating child", "idx", idx, "type", child.Type, "err", err)
			return nil, err
		}
		evaluatedChildren[idx] = val
		slog.Debug("finished evaluating child", "idx", idx+1, "len", len(node.Children), "val", val)
	}

	// Now handle this node, assuming .JS field is a snippet of JS to execute.
	// If JS is empty, just return the last evaluated child (if any), or Value as goja string.

	if node.Type == "FunctionCall" {
		slog.Debug("handling function call", "name", node.Name, "argsLen", len(evaluatedChildren))
		// First, run the JS code to define the function in the VM.
		_, err := i.vm.RunString(node.JS)
		if err != nil {
			slog.Error("error running JS code for function call", "name", node.Name, "js", node.JS, "err", err)
			return nil, err
		}

		// Then, call the function with the evaluated children.
		f, ok := goja.AssertFunction(i.vm.Get(node.Name))
		if !ok {
			slog.Error("JS value is not a function", "name", node.Name)
			return nil, fmt.Errorf("JS function is not a function")
		}
		res, err := f(goja.Undefined(), evaluatedChildren...)
		if err != nil {
			slog.Error("error calling JS function", "name", node.Name, "err", err)
			return nil, err
		}
		slog.Debug("JS function returned", "name", node.Name, "type", res.ExportType(), "value", res.Export())
		return res, nil
	}

	if node.Value != "" {
		slog.Debug("node has a value, returning as goja.Value", "value", node.Value)
		return i.vm.ToValue(node.Value), nil
	}

	// If leaf node (no children), return Value as goja.Value if present
	if len(node.Children) == 0 && node.Value != "" {
		val := i.vm.ToValue(node.Value)
		slog.Debug("leaf node with value", "value", val)
		return val, nil
	}

	slog.Debug("no result from node", "type", node.Type, "name", node.Name)
	return nil, nil
}
