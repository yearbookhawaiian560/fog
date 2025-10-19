package repl

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/0div/fog/internal/runtime"
)

func RunRuntime() {
	r, err := runtime.NewRuntimeGoja()
	if err != nil {
		slog.Error("error creating runtime", "err", err)
		return
	}

	fmt.Println()
	fmt.Println("\t·:*༺ ♱ ✮ ♱ ༻*:·")
	fmt.Println("\t🏃  Runtime  🏃")
	fmt.Println("\t·:*༺ ♱ ✮ ♱ ༻*:·")
	fmt.Println()
	fmt.Println()
	fmt.Println("Enter expressions (Ctrl+C to exit):")
	scanner := bufio.NewScanner(os.Stdin)

	initInput := "function f(a) { return a + 1 }"
	err = r.RunString(initInput)
	if err != nil {
		slog.Error("error running string", "err", err)
		return
	}

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		err = r.Compile(input)
		if err != nil {
			slog.Error("error compiling", "err", err)
			return
		}

		if input == "save" {
			err = r.Snapshot()
			if err != nil {
				slog.Error("error snapshotting", "err", err)
				return
			}
			continue
		}

		err = r.RunString(input)
		if err != nil {
			slog.Error("error running string", "err", err)
			return
		}

		globals, err := r.SerializeGlobal()
		if err != nil {
			slog.Error("error serializing global", "err", err)
			return
		}
		fmt.Println(string(globals))
		continue
		/*
			globals, err := r.SerializeGlobal()
			if err != nil {
				slog.Error("error serializing global", "err", err)
				return
			}
			result, err := parser.Parse(context.Background(), input, parse.ParserOpts{
				Globals: string(globals),
				Rules:   parse.JSRules,
			})
			if err != nil {
				slog.Error("parser error", "err", err)
				return
				continue
			}
			if !result.IsSolvableWithCode {
				fmt.Println("Not solvable with code")
				continue
			}

			for _, stmt := range result.Statements {
				err = r.RunString(stmt)
				if err != nil {
					slog.Error("error running string", "err", err)
					return
				}
			}
		*/
	}

	if err := scanner.Err(); err != nil {
		slog.Error("scanner error", "err", err)
		return
	}
}
