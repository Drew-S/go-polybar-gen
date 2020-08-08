package main

import (
    "os"
    "fmt"
    "strings"
    "text/template"
)

// Data for writing to template
type BarTemplate struct {
    Name string
    Size string
    Monitor string
    Modules string
    Fonts string
}

/*
    Parse each of the modules to attach depth numbers and fill in the separators.
    e.g.:
    "left": [ [ "a" ], [ "b", "c" ], [ "d" ] ] -> a1 sl1 b2 ssl2 c2 sl2 d3 sl3b
*/
func parseModuleList(b Bar) (string, int) {
    var out string
    var tray int = 2
    if len(b.Left) > 0 {
        out += "modules-left = "
        for i, set := range b.Left {
            for j, mod := range set {
                if mod == "tray" {
                    tray = -1
                    continue
                }
                out += fmt.Sprintf("%s%d ", mod, i + 1)
                if j < len(set)-1 {
                    out += fmt.Sprintf("ssl%d ", i + 1)
                }
            }
            if i == len(b.Left) - 1 {
                out += fmt.Sprintf("sl%db ", i + 1)
            } else {
                out += fmt.Sprintf("sl%d ", i + 1)
            }
        }
    }
    out = strings.Trim(out, " ")
    out += "\n"

    var rout string
    if len(b.Right) > 0 {
        for i, set := range b.Right {
            for j, mod := range set {
                if mod == "tray" {
                    tray = 1
                    continue
                }
                rout += fmt.Sprintf("%s%d ", mod, i + 1)
                if j < len(set)-1 {
                    rout += fmt.Sprintf("ssr%d ", i + 1)
                }
            }
            if i == len(b.Right)-1 {
                rout += fmt.Sprintf("sr%db ", i + 1)
            } else {
                rout += fmt.Sprintf("sr%d ", i + 1)
            }
        }
    }
    rout = strings.Trim(rout, " ")
    if rout != "" {
        routs := strings.Split(rout, " ")
        out += "modules-right = "
        for i := len(routs) - 1; i >= 0; i-- {
            out += routs[i] + " "
        }
        out = strings.Trim(out, " ")
    }
    
    return out, tray
}

// Parse a bar and generate the module pattern
func parseBar(b Bar, f string, file *os.File) {
    t := template.Must(template.ParseFiles("bar.ini"))

    modules, trayloc := parseModuleList(b)
    var tray string
    if trayloc == -1 {
        tray = "tray-position = left\ntray-background = ${colors.color2}\n\n"
    } else if trayloc == 1 {
        tray = "tray-position = right\ntray-background = ${colors.color2}\n\n"
    }

    var data BarTemplate = BarTemplate{
        Name: b.Name,
        Size: fmt.Sprintf(
            "width = %s\nheight = %s\noffset-y = %s\noffset-x = %s",
            b.Size.Width, b.Size.Height, b.Size.Y, b.Size.X),
        Monitor: fmt.Sprintf("monitor = %s", b.Monitor),
        Modules: fmt.Sprintf("%s\n\n%s", modules, tray),
        Fonts: f,
    }

    t.Execute(file, data)
}
