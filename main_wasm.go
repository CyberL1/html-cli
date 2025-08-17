//go:build js && wasm

package main

import (
	"fmt"
	"html-cli/utils"
	"path/filepath"
	"syscall/js"

	"github.com/spf13/afero"
)

func main() {
	global := js.Global()

	global.Set("writeFile", js.FuncOf(fsWriteFile))
	global.Set("mkdir", js.FuncOf(fsMkdir))
	global.Set("mkdirAll", js.FuncOf(fsMkdirAll))
	global.Set("readFile", js.FuncOf(fsReadFile))
	global.Set("readDir", js.FuncOf(fsReadDir))
	global.Set("build", js.FuncOf(buildWrapper))

	select {}
}

var fs = afero.NewMemMapFs()

func fsWriteFile(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return "Too little arguments. Got 1, want 2"
	}

	err := afero.WriteFile(fs, args[0].String(), []byte(args[1].String()), 0644)
	if err != nil {
		return errorToJS(err)
	}
	return nil
}

func fsMkdir(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return "Too little arguments. Got 0, want 1"
	}

	err := fs.Mkdir(args[0].String(), 0755)
	if err != nil {
		return errorToJS(err)
	}
	return nil
}

func fsMkdirAll(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return "Too little arguments. Got 0, want 1"
	}

	err := fs.MkdirAll(args[0].String(), 0755)
	if err != nil {
		return errorToJS(err)
	}
	return nil
}

func fsReadFile(this js.Value, args []js.Value) any {
	fileContents, err := afero.ReadFile(fs, args[0].String())
	if err != nil {
		return errorToJS(err)
	}

	uint8Array := js.Global().Get("Uint8Array").New(len(fileContents))
	js.CopyBytesToJS(uint8Array, fileContents)

	return uint8Array
}

func fsReadDir(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return "Too little arguments. Got 0, want 1"
	}

	dirContents, err := afero.ReadDir(fs, args[0].String())
	if err != nil {
		return errorToJS(err)
	}

	array := js.Global().Get("Array").New()

	for _, fileInfo := range dirContents {
		obj := js.Global().Get("Object").New()
		obj.Set("name", fileInfo.Name())
		obj.Set("isDir", fileInfo.IsDir())
		obj.Set("size", fileInfo.Size())
		obj.Set("mode", uint32(fileInfo.Mode()))
		obj.Set("modTime", fileInfo.ModTime().UnixNano())

		array.Call("push", obj)
	}
	return array
}

func buildWrapper(this js.Value, args []js.Value) any {
	if err := build("."); err != nil {
		fmt.Println("Build failed:", err)
		return errorToJS(err)
	}

	fmt.Println("\nBuild complete")
	return nil
}

func build(directory string) error {
	buildDirectory := "build"

	dirContents, err := afero.ReadDir(fs, directory)
	if err != nil {
		return err
	}

	for _, file := range dirContents {
		if file.Name() == "html-cli.json" || (file.IsDir() && file.Name() == buildDirectory) {
			continue
		}

		if file.IsDir() {
			build(filepath.Join(directory, file.Name()))
		} else {
			fmt.Println("Building:", filepath.Join(directory, file.Name()))

			fileContents, err := afero.ReadFile(fs, filepath.Join(directory, file.Name()))
			if err != nil {
				return err
			}

			if filepath.Ext(file.Name()) == ".html" {
				fileContents = utils.ApplyBoilerplate(fileContents, false)
				fileContents = utils.Minify(fileContents)
			}

			fs.MkdirAll(filepath.Join(buildDirectory, directory), 0775)
			afero.WriteFile(fs, filepath.Join(buildDirectory, directory, file.Name()), fileContents, 0644)
		}
	}
	return nil
}

func errorToJS(err error) js.Value {
	return js.Global().Get("Error").New(err.Error())
}
