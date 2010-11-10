package main

import (
  "log"
  "rpc"
  "rpc/jsonrpc"
)

type TestRPCSrv struct {
    *rpc.Client
}

func main () {
    r, err := jsonrpc.Dial("tcp", "localhost:6666")
    if err != nil {
        log.Exit("Error dialing host: ", err)
    }
    remote := &TestRPCSrv{r}
    a := "foo"
    var b string
    err = remote.Call("testApp.PrintStuff", &a, &b)
    if err != nil {
        log.Exit("Error calling function: ", err)
    }
    log.Exit(b)
}
