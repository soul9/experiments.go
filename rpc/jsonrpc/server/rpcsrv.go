package main

import (
  "rpc"
  "rpc/jsonrpc"
  "net"
  "log"
  "fmt"
  "os"
)

const listenport = ":6666"

type testApp struct {
    name string
}

func (t *testApp) PrintStuff (a, c *string) os.Error {
    *c = fmt.Sprintf("%s was given by testApp %s", *a, t.name)
    return nil
}

func main () {
    ta := &testApp{"fooname"}
    rpc.Register(ta)
    l, err := net.Listen("tcp", listenport)
    if err != nil {
        log.Exit("Listen error: ", err)
    } else {
        log.Println("Listening: ", l)
    }
    a := "foo"
    var c string
    ta.PrintStuff(&a, &c)
    log.Println(c)
    for conn, err := l.Accept(); err == nil; conn, err = l.Accept() {
        log.Println(conn)
        go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
    log.Exit("Error: ", err)
}
