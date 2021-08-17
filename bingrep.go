package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func dump_region(bin []byte, offset, data_len int) {
	begin := offset - 0x10
	if begin < 0 {
		begin = 0
	}
	end := offset + data_len + 0x10
	if end > len(bin) {
		end = len(bin)
	}
	fmt.Println(hex.Dump(bin[begin:end]))

}
func check_if_file_contains_bytes(path string, target []byte) {
	bin, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Println("error reading ", path, e)
		return
	}

	show_file_once := true

	for len(bin) > 0 {
		offset := bytes.Index(bin, target)
		if offset == -1 {
			break
		}

		if show_file_once {
			show_file_once = false
			fmt.Println("found match: ", path)
		}

		dump_region(bin, offset, len(target))

		bin = bin[offset+len(target):] // rest data
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`usage: bingrep 010203...`)
		os.Exit(1)
	}

	target, e := hex.DecodeString(os.Args[1])
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	e = filepath.Walk(`.`,
		func(path string, info os.FileInfo, e error) error {
			if e != nil {
				return e
			}

			if info.IsDir() {
				return nil
			}

			check_if_file_contains_bytes(path, target)

			return nil
		})
	if e != nil {
		fmt.Println(e)
	}
}
