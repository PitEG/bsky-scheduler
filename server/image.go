package main

import (
  "fmt"
  "time"
)

type Image struct {
  Id int
  Name string
  Path string
  Date time.Time
}

func (i Image) String() string {
  return fmt.Sprintf("id:%v,name:%v,path:%v,date:%v", i.Id, i.Name, i.Path, i.Date)
}

