package main 

import (
  "fmt"
  "time"
  "strconv"
  "net/http"
  "github.com/gin-gonic/gin"
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

func main() {
  fmt.Println("hello world")
  conn := Connection{filepath:"schedule.db"}
  conn.Init()
  _ = time.Now()
  // conn.AddImage("name","path",time.Now())
  image := conn.GetImage(1)
  images := conn.GetAllImages()
  fmt.Println(image) 
  fmt.Println(images)

  router := gin.Default()
  router.GET("/images",getAllImages)
  router.POST("/images",postImage)
  router.GET("/images/:id",getImage)

  router.Run("localhost:8080")
}
