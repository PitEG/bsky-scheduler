package main 

import (
  "fmt"
  "time"
)

func main() {
  fmt.Println("hello world")
  conn := Connection{filepath:"schedule.db"}
  conn.Init()
  conn.AddImage("name","path",time.Now())
  image := conn.GetImage(1)
  images := conn.GetAllImages()
  fmt.Println(image) 
  fmt.Println(images)
}
