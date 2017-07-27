package main

import (
	"crypto/md5"
	"fmt"
	"crypto/sha1"
)

func main() {
	TestString := "Hi, pandaman!"

	Md5Inst := md5.New()
	Md5Inst.Write([]byte(TestString))
	Result := Md5Inst.Sum([]byte(""))
	fmt.Printf("%x\n\n", Result)

	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(TestString))
	Result=Sha1Inst.Sum([]byte(""))
	fmt.Printf("%x\n\n",Result)
}
