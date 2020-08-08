package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "encoding/json"
    "text/template"
)

/*
    A list of modules to load for a bar.
    Modules are loaded in a similar manner to that of
    powerline.
*/
type Module [][]string

// The size parameters of a bar, see polybar wiki
type Size struct {
    Width string `json:"w"`
    Height string `json:"h"`
    X string `json:"x"`
    Y string `json:"y"`
}

// Bar settings
type Bar struct {
    Name string `json:"name"`
    Size Size `json:"size"`
    Monitor string `json:"monitor"`
    Left Module `json:"left"`
    Right Module `json:"right"`
    Center Module `json:"center"`
}

// Information on a separator, automatically genorates output
type Separator struct {
    Left string `json:"left"`
    Right string `json:"right"`
}

type Colors struct {
    Foreground string `json:"foreground"`
    Background string `json:"background"`
    Text string `json:"text"`
    Colors []string `json:"colors"`
}

// Config setup
type Config struct {
    Bars []Bar `json:"bars"`
    Fonts []string `json:"fonts"`
    Sep Separator `json:"sep"`
    Subsep Separator `json:"subsep"`
    Data map[string]interface{} `json:"data"`
    Colors Colors `json:"colors"`
}

const SEPARATOR string = ";==========================================================\n\n"
const MODSEPARATOR string =
`;==========================================================
;                                                   MODULES
;==========================================================

`
const BARSEPARATOR string =
`;==========================================================
;                                                      BARS
;==========================================================

`

const preambleTemplate string =
`;==========================================================
;
;
;   ██████╗  ██████╗ ██╗  ██╗   ██╗██████╗  █████╗ ██████╗
;   ██╔══██╗██╔═══██╗██║  ╚██╗ ██╔╝██╔══██╗██╔══██╗██╔══██╗
;   ██████╔╝██║   ██║██║   ╚████╔╝ ██████╔╝███████║██████╔╝
;   ██╔═══╝ ██║   ██║██║    ╚██╔╝  ██╔══██╗██╔══██║██╔══██╗
;   ██║     ╚██████╔╝███████╗██║   ██████╔╝██║  ██║██║  ██║
;   ╚═╝      ╚═════╝ ╚══════╝╚═╝   ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝
;
;
;   Created By Aditya Shakya.
;
;   Config generated by Go-Polybar-Gen
;
;==========================================================

{{.Global}}

;==========================================================

[colors]
background = {{.Colors.Background}} 
foreground = {{.Colors.Foreground}}
{{range $i, $col := .Colors.Colors}}color{{inc $i}} = {{$col}}
{{end}}
`

func main() {
    var config Config

    fmt.Println("Generating Polybar config")

    filedata, err := ioutil.ReadFile("config.json")
    if err != nil {
        fmt.Printf("An error occurred reading json: %s\n", err)
        return
    }
    
    json.Unmarshal(filedata, &config)

    file, err := os.Create("config")
    if err != nil {
        fmt.Printf("An error occurred: %s\n", err)
    }
    defer file.Close()

    t := template.Must(template.New("preamble").Funcs(template.FuncMap{
        "inc": func(i int) int {
            return i + 1
        },
    }).Parse(preambleTemplate))

    t.Execute(file, struct{
        Global string
        Colors Colors
    }{
        "[global/wm]\nmargin-bottom = -12",
        config.Colors,
    })

    file.WriteString(BARSEPARATOR)

    var fonts string

    for i, f := range config.Fonts {
        fonts += fmt.Sprintf("font-%d = \"%s\"\n", i, f)
    }

    for i, b := range config.Bars {
        parseBar(b, fonts, file)
        if i != len(config.Bars) - 1 {
            file.WriteString(SEPARATOR)
        }
    }

    file.WriteString(MODSEPARATOR)

    generateSeparators(config, file)

    for _, b := range config.Bars {
        for i, set := range b.Left {
            for _, mod := range set {
                if !parseMod(mod, i+1, config, file) {
                    continue
                }
                file.WriteString("\n" + SEPARATOR)
            }
        }
        for i, set := range b.Right {
            for _, mod := range set {
                if !parseMod(mod, i+1, config, file) {
                    continue
                }
                file.WriteString("\n" + SEPARATOR)
            }
        }
    }
}