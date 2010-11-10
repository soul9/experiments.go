package main

import (
  "netchan"
  "log"
  "os"
  "bufio"
  "flag"
  "strings"
)

func importChan(imp *netchan.Importer, name string, done chan os.Error) {
    ch := make(chan string)
    err := imp.Import(name, ch, netchan.Recv)
    if err != nil {
        log.Println("Error during channel import: ", err)
        done <- err
    }
    for line := <- ch; ! closed(ch) ; line = <- ch {
        log.Print(line)
    }
    done <- os.NewError("No moar inputs")
}

func main() {
    var straddr string
    flag.StringVar(&straddr, "a", "", "Address of the exporter")
    flag.Parse()
    if straddr == "" {
        flag.Usage()
        log.Exit("No address given")
    }
    bstdin := bufio.NewReader(os.Stdin)
    var err os.Error
    done := make(chan os.Error)
    imp, err := netchan.NewImporter("tcp", straddr)
    if err != nil {
        log.Exit("Error creating importer: ", err)
    }
    go func() {
        for line, err := bstdin.ReadString('\n'); err == nil; line, err = bstdin.ReadString('\n') {
            go importChan(imp, strings.TrimSpace(line), done)
            go func() {
                err :=<- done
                log.Println(err)
            }()
        }
        log.Exit("Error reading stdin: ", err)
    }()
    <- done
}
