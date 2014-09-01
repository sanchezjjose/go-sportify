package main

import (
  "fmt"
  "html/template"
  "net/http"
  "regexp"
  "encoding/json"
  "io/ioutil"
  "os"
  "path"
  "log"
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

var templates = template.Must(template.ParseFiles("templates/layout.html", "templates/home.html", "templates/roster.html"))

var validPath = regexp.MustCompile("^/(home|roster)/")

/*
 * Convert JSON Response to Game Object
 */
func getNextGame() (*Game, error) {
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

func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := getNextGame()
  if err != nil {
    http.Redirect(w, r, "/home/" + title, http.StatusFound)
    return
  }

  renderTemplate(w, "home", p)
}

func rosterHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := getNextGame()
  if err != nil {
    http.Redirect(w, r, "/roster/" + title, http.StatusFound)
    return
  }

  renderTemplate(w, "roster", p)
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

  fs := http.FileServer(http.Dir("static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))

  // Method A
  http.HandleFunc("/", serveTemplate)

  // Method B
  http.HandleFunc("/home/", makeHandler(homeHandler))
  http.HandleFunc("/roster/", makeHandler(rosterHandler))

  http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
  lp := path.Join("templates", "layout.html")
  fp := path.Join("templates", r.URL.Path) // e.g, home.html

  //Return a 404 if the template doesn't exist
  info, err := os.Stat(fp)
  if err != nil {
    if os.IsNotExist(err) {
      http.NotFound(w, r)
      return
    }
  }

  //Return a 404 if the request is for a directory
  if info.IsDir() {
    http.NotFound(w, r)
    return
  }

  templates, err := template.ParseFiles(lp, fp)
  if err != nil {
    // Log the detailed error
    log.Println(err.Error())
    // Return a generic "Internal Server Error" message
    http.Error(w, http.StatusText(500), 500)
    return
  }

  if err := templates.ExecuteTemplate(w, "layout", nil); err != nil {
    log.Println(err.Error())
    http.Error(w, http.StatusText(500), 500)
  }

}
