package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
  "bufio"
  "time"
	"net/http"
  "os"
)

func createSession(handle string, password string) (map[string]interface{},bool) {
  respEndpoint := "https://bsky.social/xrpc/com.atproto.server.createSession"

  reqJson := []byte(fmt.Sprintf(`{"identifier":"%s","password":"%s"}`, handle, password))
  req, err := http.NewRequest("POST",respEndpoint,bytes.NewBuffer(reqJson));
  req.Header.Set("Content-Type","application/json")
  if err != nil {
    println("request for session failed")
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  // defer resp.Body.Close()

  body,err := io.ReadAll(resp.Body)
  if err != nil {
    println("could not read response body")
  }
  resp.Body.Close()
  // println(string(body))

  var session map[string]interface{}
  err = json.Unmarshal([]byte(body),&session)
  if err != nil {
    println("could not unmarshal body")
  }
  _,ok := session["accessJwt"].(string)
  if ok != true {
    return session,false
  }
  _,ok = session["did"].(string)
  if ok != true {
    return session,false
  }
  return session,true
}

func post(imagePath string, caption string, user string, pass string) {
  session, ok := createSession(user,pass)
  if ok != true {
    println("could not create session")
    return
  }

  did := session["did"].(string)
  token := session["accessJwt"].(string)
  client := &http.Client{}

  // make image blob
  fmt.Println("posting blob")
  imageFile, err := os.Open(imagePath)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer imageFile.Close()
  stat, err := imageFile.Stat()
  if err != nil {
    fmt.Println(err)
    return
  }
  blob := make([]byte, stat.Size())
  _,err = bufio.NewReader(imageFile).Read(blob)
  if err != nil && err != io.EOF {
    fmt.Println(err)
    return
  } 

  blobEndpoint := "https://bsky.social/xrpc/com.atproto.repo.uploadBlob"
  req, err := http.NewRequest("POST",blobEndpoint,bytes.NewBuffer(blob));
  req.Header.Set("Content-Type","image/png")
  req.Header.Set("Authorization","Bearer " + token)
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  fmt.Println(resp)
  body,err := io.ReadAll(resp.Body)
  if err != nil {
    println("could not read response body")
  }
  resp.Body.Close()
  var blobResp map[string]interface{}
  err = json.Unmarshal([]byte(body),&blobResp)
  if err != nil {
    fmt.Println("can't unmarshal json")
  }
  // fmt.Println(blobResp["blob"]["ref"]["$link"].(string))
  blobLink := blobResp["blob"].(map[string]interface{})["ref"].(map[string]interface{})["$link"]
  println(blobLink)

  // make post
  postEndpoint := "https://bsky.social/xrpc/com.atproto.repo.createRecord"
  now := time.Now().Format("2006-01-02T15:04:05Z07:00")
  post := fmt.Sprintf(
    `{
      "$type":"app.bsky.feed.post",
      "text":"%s",
      "createdAt":"%s",
      "langs":["en-US"],
      "embed": {
        "$type": "app.bsky.embed.images",
        "images": [
          {
            "alt":"",
            "image": {
              "$type":"blob",
              "ref": {
                "$link":"%s"
              },
              "mimeType":"image/png",
              "size":%d
            }
          }
        ]
      }
    }`,caption,now,blobLink,stat.Size())
  reqJson := []byte(fmt.Sprintf(`{"repo":"%s","collection":"app.bsky.feed.post","record":%s}`, did, post))
  println(string(reqJson))
  req, err = http.NewRequest("POST",postEndpoint,bytes.NewBuffer(reqJson));
  req.Header.Set("Content-Type","application/json")
  req.Header.Set("Authorization","Bearer " + token)
  if err != nil {
    println("request for session failed")
  }
  resp, err = client.Do(req)
  if err != nil {
      panic(err)
  }

  body,err = io.ReadAll(resp.Body)
  if err != nil {
    println("could not read response body")
  }
  resp.Body.Close()

  println(string(body))
}
