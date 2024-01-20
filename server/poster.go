package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
  "time"
	"net/http"
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

  postEndpoint := "https://bsky.social/xrpc/com.atproto.repo.createRecord"

  now := time.Now().Format("2006-01-02T15:04:05Z07:00")
  post := fmt.Sprintf(`{"$type":"app.bsky.feed.post","text":"%s","createdAt":"%s"}`,caption,now)
  reqJson := []byte(fmt.Sprintf(`{"repo":"%s","collection":"app.bsky.feed.post","record":%s}`, did, post))
  println(string(reqJson))
  req, err := http.NewRequest("POST",postEndpoint,bytes.NewBuffer(reqJson));
  req.Header.Set("Content-Type","application/json")
  req.Header.Set("Authorization","Bearer " + token)
  if err != nil {
    println("request for session failed")
  }

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }

  body,err := io.ReadAll(resp.Body)
  if err != nil {
    println("could not read response body")
  }
  resp.Body.Close()

  println(string(body))
}
