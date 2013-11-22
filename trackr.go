package main

import (
        "fmt"
        "os"
        "net/http"
        "net/url"
        "encoding/json"
        "io/ioutil"
        "time"
       )

var trackrd string = "http://localhost:8080/" // Trailing / is important

type Message struct {
    State string
    Msg string
}

type Project struct {
    Name string
    Time time.Duration
}


func main(){
    if len(os.Args) > 1 {
        switch os.Args[1]{
            case "create":
                if len(os.Args) > 2 {
                    project := os.Args[2]
                    v := url.Values{}
                    v.Set("name", project)

                    if len(os.Args) == 4 {
                        duration_string := os.Args[3]
                        v.Set("duration", duration_string)
                    }

                    resp, err := http.Get(trackrd + "create.json?" + v.Encode())

                    if err != nil {
                        panic(err)
                    }

                    jsonString, _ := ioutil.ReadAll(resp.Body)
                    var msg Message
                    json.Unmarshal(jsonString, &msg)
                    fmt.Printf("[%s] %s\n", msg.State, msg.Msg)
                }
            case "log":
                if len(os.Args) == 4 {
                }
            case "list":
                var projects []Project
                resp, err := http.Get(trackrd + "list.json")

                if err != nil {
                    panic(err)
                }

                jsonString, _ := ioutil.ReadAll(resp.Body)
                json.Unmarshal(jsonString, &projects)

                for _, curr_project := range projects{
                    fmt.Printf("%s: [%s]\n", curr_project.Name, curr_project.Time)
                }

            default:
                fmt.Println("Invalid option")
        }
    } else {
        fmt.Println("Invalid usage")
    }
}
