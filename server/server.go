package main 

import (
  "fmt"
  "time"
  "strconv"
  "net/http"
  "github.com/gin-gonic/gin"
  "flag"
)

type test struct {
  Name string `json:"name"`
  Age int `json:"age"`
}

var things = []test{
    {Name:"man",Age:21},
    {Name:"child",Age:12},
    {Name:"corn",Age:3},
  }

func getAllImages(context *gin.Context) {
  context.IndentedJSON(http.StatusOK,things)
}

func getImage(context *gin.Context) {
  age, err := strconv.Atoi(context.Param("id"))
  if err != nil {
    age = 0
  }
  var test = test{Name:"haha",Age:age}
  context.IndentedJSON(http.StatusOK,test)
}

func postImage(context *gin.Context) {
  var thing test
  
  if err := context.BindJSON(&thing); err != nil {
    return
  }

  things = append(things, thing)
  context.IndentedJSON(http.StatusOK, thing)
}

// this function just posts something when it's time
func socialMediaWorker(imgDir string, dbPath string) {
  for {
    fmt.Println("hi")

    conn := Connection{filepath:"schedule.db"}

    earliestPost := conn.GetEarliestPost()
    fmt.Println(earliestPost)

    // post if the scheduled time is before current time
    fmt.Println(time.Now())
    if time.Now().After(earliestPost.Date) {
      // time to post it on social media
    }

    time.Sleep(5 * time.Second)
    // time.Sleep(1 * time.Minute)
  }
}

func main() {
  fmt.Println("hello world")
  // parse flags
  username := flag.String("user","","bsky username")
  password := flag.String("pass","","bsky api key")
  flag.Parse()
  println(*username)
  println(*password)

  conn := Connection{filepath:"schedule.db"}
  conn.Init()
  _ = time.Now()
  /*
  conn.AddImage("name","path",time.Now())
  image := conn.GetImage(1)
  images := conn.GetAllImages()
  fmt.Println(image) 
  fmt.Println(images)
  conn.SchedulePost(conn.GetImage(1), "filler text","bsky",time.Now())
  */
  post("","hi",*username,*password)

  router := gin.Default()
  router.GET("/images",getAllImages)
  router.POST("/images",postImage)
  router.GET("/images/:id",getImage)

  go socialMediaWorker(".","schedule.db")

  router.Run("localhost:8080")
}
