package main

import (
	"english/myerror"
	"errors"
	"fmt"
)

func main() {
	err1 := myerror.ErrRecordNotFound
	err2 := myerror.ErrRecordNotFound
	fmt.Println(errors.Is(err1, err2))
}
