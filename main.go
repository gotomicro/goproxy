package main

import (
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"goproxy/invoker"
)

func main() {
	err := ego.New(ego.WithHang(true)).Invoker(invoker.Init).Run()
	if err != nil {
		elog.Panic("app panic", elog.FieldErr(err))
	}
}
