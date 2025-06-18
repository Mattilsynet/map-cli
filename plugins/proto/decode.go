package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"google.golang.org/protobuf/encoding/protowire"
)

// TODO: make this cli compliant, so it can be used as a plugin for cobra
func main() {
	data, _ := io.ReadAll(os.Stdin)
	decodeMessage(data, 0)
}

func decodeMessage(data []byte, depth int) {
	indent := strings.Repeat("  ", depth)
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		data = data[n:]
		var vlen int
		switch typ {
		case protowire.VarintType:
			val, l := protowire.ConsumeVarint(data)
			fmt.Printf("%s[%d]: varint %d\n", indent, num, val)
			vlen = l
		case protowire.BytesType:
			val, l := protowire.ConsumeBytes(data)
			if isUTF8(val) {
				fmt.Printf("%s[%d]: string \"%s\"\n", indent, num, string(val))
			} else {
				fmt.Printf("%s[%d]: nested bytes (%d bytes)\n", indent, num, len(val))
				decodeMessage(val, depth+1)
			}
			vlen = l
		case protowire.Fixed32Type:
			val, l := protowire.ConsumeFixed32(data)
			fmt.Printf("%s[%d]: fixed32 %d\n", indent, num, val)
			vlen = l
		case protowire.Fixed64Type:
			val, l := protowire.ConsumeFixed64(data)
			fmt.Printf("%s[%d]: fixed64 %d\n", indent, num, val)
			vlen = l
		default:
			fmt.Printf("%s[%d]: unknown wire type %d\n", indent, num, typ)
			vlen = len(data)
		}
		data = data[vlen:]
	}
}

func isUTF8(data []byte) bool {
	return utf8.Valid(data) && bytes.IndexFunc(data, func(r rune) bool {
		return r != '\n' && r < 32
	}) == -1
}
