package invoker

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
)

func tcp(wg *sync.WaitGroup, info Proxy) {
	defer wg.Done()
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4zero, Port: info.DstPort})
	if err != nil {
		elog.Error("tcp listen err", elog.FieldErr(err))
		return
	}
	elog.Info("listen tcp", zap.String("srcAddr", info.SrcAddr), zap.Int("dstPort", info.DstPort))
	defer listen.Close()
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			elog.Error("tcp listen AcceptTCP err", elog.FieldErr(err))
			break
		}
		go handle(conn, info)
	}
}

func handle(conn *net.TCPConn, info Proxy) {
	defer conn.Close()
	realAddr := ""
	if info.ProxyAddr == "" {
		realAddr = info.SrcAddr
	} else {
		realAddr = info.ProxyAddr
	}

	tcpAddrTemp, err := net.ResolveTCPAddr("tcp", realAddr)
	if err != nil {
		return
	}
	dialConn, err := net.DialTCP("tcp", nil, tcpAddrTemp)
	if err != nil {
		elog.Error("handle DialTCP err", elog.FieldErr(err))
		return
	}
	defer dialConn.Close()
	dialConnReader := bufio.NewReader(dialConn)
	if info.ProxyAddr != "" {
		proxyText := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nProxy-Authorization: Basic %s\r\nProxy-Connection: Keep-Alive\r\n\r\n", realAddr, base64.StdEncoding.EncodeToString([]byte(info.ProxyUser)))
		_, err = dialConn.Write([]byte(proxyText))
		if err != nil {
			return
		}
		success := false
		for {
			command, _, err := dialConnReader.ReadLine()
			if err != nil {
				break
			}
			if string(command) == "" {
				success = true
				break
			}
		}
		if !success {
			return
		}
	}
	if info.Protocol {
		remoteTmp, err := net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())
		if err != nil {
			return
		}
		localTmp, err := net.ResolveTCPAddr("tcp", conn.LocalAddr().String())
		if err != nil {
			return
		}
		proxyProtocolText := fmt.Sprintf("PROXY TCP4 %s %s %d %d\r\n", remoteTmp.IP.String(), localTmp.IP.String(), remoteTmp.Port, localTmp.Port)
		_, err = dialConn.Write([]byte(proxyProtocolText))
		if err != nil {
			return
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go transport(wg, conn, dialConn, conn, dialConn)
	go transport(wg, dialConn, conn, dialConnReader, conn)
	wg.Wait()
}

func transport(wg *sync.WaitGroup, left, right *net.TCPConn, src io.Reader, dst io.Writer) {
	defer wg.Done()
	defer left.CloseRead()
	defer right.CloseWrite()
	_, _ = io.Copy(dst, src)
}
