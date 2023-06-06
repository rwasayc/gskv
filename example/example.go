package main

import (
	"fmt"

	"github.com/rwasayc/gskv"
)

func main() {
	db := gskv.Open()
	v1 := db.Get([]byte("k_1"))
	if v1 == nil {
		fmt.Println("k_1 is nil")
		db.Set([]byte("k_1"), []byte("v_1"))
		fmt.Printf("set k_1\n")
	} else {
		fmt.Printf("k_1 is %v\n", string(v1))
		db.Delete([]byte("k_1"))
		fmt.Printf("delete k_1\n")
	}
}
