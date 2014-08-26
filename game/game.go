package main

import "fmt"
import "encoding/json"
import "net/http"
import "io/ioutil"


type Game struct {
  Id int `json:"game_id"`
  Seq int `json:"game_seq"`
  StartTime string `json:"startTime"`
  Address string `json:"address"`
  Gym string `json:"gym"`
  LocationDetails string `json:"locationDetails"`
  Opponent string `json:"opponent"`
  Result string `json:"result"`
  PlayersIn []string `json:"playersIn"`
  PlayersOut []string `json:"playersOut"`
  PlayoffGame bool `json:"is_playoff_game"`
  Season string `json:"season"`
}

func main() {

  fmt.Println()


  fmt.Println("\n===================")
  fmt.Println("USING CUSTOM LIBRARY")
  fmt.Println("===================\n")
  fmt.Println(helper.SayHello("Special Ops"))


  fmt.Println("\n===================")
  fmt.Println("OBJECT TO JSON")
  fmt.Println("===================\n")

  gameJson := &Game {
    Id: 59,
    Seq: 6,
    StartTime: "Thu, 7:00 pm" }
  gameObj, _ := json.Marshal(gameJson)
  fmt.Println(gameJson)
  fmt.Println()
  fmt.Println(string(gameObj))


  fmt.Println("\n===================")
  fmt.Println("JSON TO OBJECT")
  fmt.Println("===================\n")

  str := `{"id": 1, "seq": 3, "start_time": "Thu, 7:00 PM", "Address": "2 Park Ave.", "Gym": "JR2",
           "location_details": "Show up.", "Opponent": "NY Knicks", "Result": "W 32-30",
           "PlayersIn": ["1176064196"], "PlayersOut": [],
           "PlayoffGame": false, "Season" : "Summer 2014"}`
  obj := &Game{}
  json.Unmarshal([]byte(str), &obj)
  fmt.Println(str)
  fmt.Println()
  fmt.Println(obj)


  fmt.Println("\n===================")
  fmt.Println("HTTP URL RESPONSE")
  fmt.Println("===================\n")

  url := "http://sportify-gilt.heroku.com/json/nextgame"
  res, err := http.Get(url)
  if err != nil {
    fmt.Printf("%s", err)
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Printf("%s", err)
  }
  fmt.Println(url)
  fmt.Println()
  fmt.Printf("%s\n", string(body))


  fmt.Println("\n===================")
  fmt.Println("HTTP URL RESPONSE TO OBJECT")
  fmt.Println("===================\n")

  gameObj2 := &Game{}
  json.Unmarshal([]byte(string(body)), &gameObj2)
  fmt.Println(gameObj2)
  fmt.Println()
}
