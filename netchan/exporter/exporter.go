package main

import (
  "netchan"
  "log"
  "fmt"
  "exec"
  "bufio"
  "os"
)

func main() {
    exp, err := netchan.NewExporter("tcp", "localhost:0")
    if err != nil {
        log.Exit("Error creating exporter: ", err)
    }

    sendchan := make(chan string)
    chname := "sendchan"
    err = exp.Export(chname, sendchan, netchan.Send)
    if err != nil {
        log.Exit("Error exporting channel: ", err)
    }
 
    args := []string{"../importer/importer", "-a", exp.Addr().String()}
    imp, err := exec.Run("../importer/importer", args, []string{}, "", exec.Pipe, exec.Pipe, exec.Pipe)
    if err != nil {
        log.Exit("Couldn't spawn command: ", err)
    } else {
        log.Println("Started command", imp, args)
    }

    fmt.Fprintf(imp.Stdin, "%s\n", chname)
    bufou := bufio.NewReader(imp.Stderr)
    go func() {
        for line, err := bufou.ReadString('\n'); err == nil; line, err = bufou.ReadString('\n') {
            log.Print("Importer: ", line)
        }
    }()
    bufin := bufio.NewReader(os.Stdin)
    go func() {
        for line, err := bufin.ReadString('\n'); err == nil; line, err = bufin.ReadString('\n') {
            sendchan <- line
        }
        close(sendchan)
    }()
    imp.Wait(0)
}
