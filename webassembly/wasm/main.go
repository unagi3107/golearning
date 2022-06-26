package main

import (
	"bytes"
	"io"
	"syscall/js"
)

func main() {
	js.Global().Set("compress", compressFunc)
	select {}
}

func newUnit8Array(size int) js.Value {
	ua := js.Global().Get("Unit8Array")
	return ua.New(size)
}

var compressFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	jsSrc := args[0]
	srcLen := jsSrc.Get("length").Int()
	srcBytes := make([]byte, srcLen)
	js.CopyBytesToGo(srcBytes, jsSrc)

	src := bytes.NewReader(srcBytes)

	r, err := compressor.Compress(src)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic(err)
	}

	ua := newUnit8Array(buf.Len())
	js.CopyBytesToJS(ua, buf.Bytes())
	return ua
})
