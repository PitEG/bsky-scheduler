package main

import (
  "fmt"
  "time"
)

type Image struct {
  id int
  name string
  path string
  date time.Time
}

func (i Image) String() string {
  return fmt.Sprintf("id:%v,name:%v,path:%v,date:%v", i.id, i.name, i.path, i.date)
}

