package invoker

import (
	"sync"

	"github.com/gotomicro/ego/core/econf"
)

type Proxy struct {
	DstPort   int
	SrcAddr   string
	ProxyAddr string
	ProxyUser string
	Protocol  bool
	Test      bool
}

func Init() (err error) {
	var list []Proxy
	err = econf.UnmarshalKey("proxy", &list)
	if err != nil {
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(list))
	for _, value := range list {
		go tcp(wg, value)
	}
	wg.Wait()
	return nil
}
