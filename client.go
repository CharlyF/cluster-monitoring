package main

import "net"
import "time"

func main() {
    c,err := net.Dial("unix","/Users/charly.fontaine/var/run/net.sock")
    if err != nil {
        panic(err)
    }
    for {
        _,err := c.Write([]byte("hi\n"))
        if err != nil {
            panic(err)
        }
        time.Sleep(3e9)
    }
}
