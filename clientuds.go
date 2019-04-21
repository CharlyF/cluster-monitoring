package main

import (
 "io"
 "log"
 "os"

 "github.com/orisano/uds"
)

const (
 sockPath = "/Users/charly.fontaine/var/run/net.sock"
)

func main() {
 client := uds.NewClient(sockPath)
 resp, err := client.Get("http://unix/connections")
 if err != nil {
  log.Fatal(err)
 }
 defer resp.Body.Close()

 io.Copy(os.Stdout, resp.Body)
}
