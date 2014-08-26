package main

import (
  "fmt"
  "html/template"
  "net/http"
  "regexp"
  "encoding/json"
  "io/ioutil"
  "go/build"
)

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

var templates = template.Must(template.ParseFiles(build.Default.GOPATH + "/src/github.com/sanchezjjose/go-sportify/template/homePage.html"))

var validPath = regexp.MustCompile("^/(home|roster)/")

/*
 * Convert JSON Response to Game Object
 */
func loadPage(title string) (*Game, error) {
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

  game := &Game{}
  json.Unmarshal([]byte(string(body)), &game)

  return game, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Game) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func homepageHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/home/" + title, http.StatusFound)
    return
  }

  renderTemplate(w, "homePage", p)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
      http.NotFound(w, r)
      return
    }
    fn(w, r, m[0])
  }
}

func main() {
  fmt.Printf("starting server...")

  http.HandleFunc("/home/", makeHandler(homepageHandler))

  http.ListenAndServe(":8000", nil)
}
