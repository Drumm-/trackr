package main

import (
        "fmt"
        "time"
        "encoding/json"
        "io/ioutil"
        "net/http"
       )

type Project struct {
    Name string
    Time time.Duration
}
var projects []Project

func saveProjects(){
    b, _ := json.Marshal(projects)
    ioutil.WriteFile("projects.json", b, 0644)
}

type Message struct {
    State string
    Msg string
}

func jsonMessage(state string, msg string) ([]byte){
    message := Message{state, msg}
    b, _ := json.Marshal(message)
    return b
}

func listHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    b, err := json.Marshal(projects)
    if err != nil {
        fmt.Fprintf(w, "{}")
    }
    fmt.Fprintf(w, "%s", b)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    project_name := r.FormValue("name")
    duration_string := r.FormValue("duration")
    var duration time.Duration
    duration, err := time.ParseDuration(duration_string)

    if err != nil {
        w.Write(jsonMessage("fail", "Invalid duration"))
        return
    }

    if len(project_name) == 0 {
        w.Write(jsonMessage("fail", "Missing name parameter"))
        return
    }

    var project Project
    var key int = -1

    for i, curr_project := range projects{
        if project_name == curr_project.Name {
            key = i
            project = curr_project
            break
        }
    }

    if (key < 0){
        w.Write(jsonMessage("fail", "Invalid project name"))
    }

    projects[key] = Project{project.Name, project.Time+duration}
    saveProjects()
    w.Write(jsonMessage("success", "Successfully logged to project"))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    project_name := r.FormValue("name")
    duration_string := r.FormValue("duration")
    var duration time.Duration
    duration, err := time.ParseDuration(duration_string)

    if err != nil {
        duration = time.Duration(0)
    }

    if len(project_name) == 0 {
        w.Write(jsonMessage("fail", "Missing name parameter"))
        return
    }

    p := Project{project_name, duration}
    projects = append(projects, p)
    saveProjects()
    w.Write(jsonMessage("success", "Created new project"))
}

func main(){
    projects = make([]Project, 0)
    b, err := ioutil.ReadFile("projects.json")

    if err != nil {
        b, _ = json.Marshal(projects)
    }
    err = json.Unmarshal(b, &projects)

    if err != nil {
        panic(err)
    }

    http.HandleFunc("/log.json", logHandler)
    http.HandleFunc("/list.json", listHandler)
    http.HandleFunc("/create.json", createHandler)
    http.ListenAndServe(":8080", nil)
}
