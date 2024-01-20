package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
  defer resp.Body.Close()

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
  if ok != false {
    return
  }

  did := session["did"].(string)
  token := session["accessJwt"].(string)

  endpoint := "https://bsky.social/xrpc/com.atproto.repo.createRecord"
}
