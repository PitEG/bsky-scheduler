package main

import (
  "time"
)

type ScheduledPost struct {
  id int
  imageId int
  text string
  site string
  date time.Time
}

