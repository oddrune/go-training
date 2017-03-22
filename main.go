package main

import (  
  "os"
  "strconv"
  "net/http"
  "encoding/json"
  "strings"
  "github.com/spf13/viper"
)

func hello(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Howdy!"))
}

type weatherData struct {
  Name string `json:"name"`
  Weather []struct {
    Main string `json:"main"`
    Description string `json:"description"`
  }
  Main struct {
    Kelvin float64 `json:"temp"`
  } `json:"main"`
}

func query(city string) (weatherData, error) {
  var APIKEY string = viper.GetString("weather.apikey")
  resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID="+APIKEY+"&q=" + city)
  if err != nil {
    return weatherData{}, err
  }

  defer resp.Body.Close()

  var d weatherData

  if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
    return weatherData{}, err
  }
  return d, nil
}



func KelvinToString(kelvin float64) string {
    return strconv.FormatFloat( 273.16 - kelvin, 'f', 1, 64)
}

func main() {
  println("Cry 'Havoc!', and let slip the dogs of war")

  viper.SetConfigName("config")
  viper.AddConfigPath(".")
  err := viper.ReadInConfig()
  if err != nil {
    println("Could not load config file.")
    os.Exit(1)
  } 

  http.HandleFunc("/", hello)

  http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
    println(">> " + r.URL.Path)
    city := strings.SplitN(r.URL.Path, "/", 3)[2]

    data, err := query(city)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    println("<< " + data.Name + ": " + data.Weather[0].Description + " / " + KelvinToString(data.Main.Kelvin) + "°C")
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    json.NewEncoder(w).Encode(data)
  })

  http.ListenAndServe(":8080", nil)
}
