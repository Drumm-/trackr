package main

import (
        "fmt"
        "time"
        "flag"
        "encoding/json"
        "io/ioutil"
       )

type Project struct {
    Name string
    Time time.Duration
}

func main(){
    var create = flag.Bool("create", false, "Create a new project")
    var log = flag.Bool("log", false, "Log time to project")
    var list = flag.Bool("list", false, "List projects")
    var project_name = flag.String("name", "", "Project name")
    var duration = flag.Duration("time", time.Duration(0), "Time to log")
    flag.Parse()

    projects := make([]Project, 0)

    b,_ := ioutil.ReadFile("projects.json")

    json.Unmarshal(b, &projects)


    if ( *create == true ){
        if ( len(*project_name) == 0){
            fmt.Printf("Missing -name parameter\n")
            return
        }
        fmt.Printf("Creating project %s\n", *project_name)
        p := Project{*project_name, *duration}
        projects = append(projects, p)
    } else if ( *log == true ){
        fmt.Printf("Logging time to %s\n", *project_name)
        if ( len(*project_name) == 0){
            fmt.Printf("Missing -name parameter\n")
            return
        }

        var project Project
        var key int = -1

        for i, curr_project := range projects{
            if (*project_name == curr_project.Name){
                key = i
                project = curr_project
                break
            }
        }

        if (key < 0){
            fmt.Printf("Invalid project name")
        }

        projects[key] = Project{project.Name, project.Time+*duration}

    } else if ( *list == true ) {
        for _, project := range projects{
            fmt.Printf("%s - %s\n", project.Name, project.Time)
        }
    } else {
        fmt.Printf("No action specified\n")
    }

    b, _ = json.Marshal(projects)
    ioutil.WriteFile("projects.json", b, 0644)

}
