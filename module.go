package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var mods map[string][]int = make(map[string][]int)

func parseMod(mod string, depth int, config Config, file *os.File) bool {
    if mod == "tray" {
        return false
    }
    if _, exists := mods[mod]; exists {
        for _, i := range mods[mod] {
            if i == depth {
                return false
            }
        }
    } else {
        mods[mod] = make([]int, 0)
    }
    mods[mod] = append(mods[mod], depth)

    var b string
    var f string
    if depth == len(config.Colors.Colors) {
        b = config.Colors.Background
        f = config.Colors.Foreground
    } else {
        b = config.Colors.Colors[depth]
        f = config.Colors.Text
    }

    var fileTemplate string = fmt.Sprintf("modules/%s.ini", mod)

    if _, err := os.Stat(fileTemplate); err != nil {
        if os.IsNotExist(err) {
            fmt.Printf("Error, module not found: %s\n", mod)
            return false
        }
    }

    filedata, err := ioutil.ReadFile(fileTemplate)
    if err != nil {
        fmt.Println("Error, unable to read file")
        return false
    }
    
    var fileString string = string(filedata)
    fileString = strings.ReplaceAll(fileString, "$$Foreground$$", f)
    fileString = strings.ReplaceAll(fileString, "$$Background$$", b)
    fileString = strings.ReplaceAll(fileString, "$$Text$$", config.Colors.Text)

    t := template.Must(template.New("mod").Parse(fileString))

    t.Execute(file, struct{
        Depth string
        Foreground string
        Background string
        Text string
        Data map[string]interface{}
    }{
        Depth: fmt.Sprintf("%d", depth),
        Foreground: f,
        Background: b,
        Text: config.Colors.Text,
        Data: config.Data,
    })

    return true
}
