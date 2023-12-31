package main

import (
  "time"
)

type ScheduledPost struct {
  Id int
  ImageId int
  Text string
  Site string
  Date time.Time
}

