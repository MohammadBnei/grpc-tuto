package main

import (
  "fmt"
  "time"
)

func main() {
  for {
    fmt.Println("Hello Universe")
    time.Sleep(time.Second * 3)
  }
}