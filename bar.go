package main

import (
    "os"
    "fmt"
    "strings"
    "text/template"
    "math"
)

// Data for writing to template
type BarTemplate struct {
    Name string
    Size string
    Monitor string
    Modules string
    Fonts string
}

func midleft(mid int, i int, set []string) string {
    var out string
    for j, mod := range set {
        out += fmt.Sprintf("%s%d ", mod, mid - i)
        if j < len(set) - 1 {
            out += fmt.Sprintf("ssr%d ", mid - i)
        }
    }
    if i < mid - 1 {
        out += fmt.Sprintf("sr%d ", mid - i - 1)
    }
    return out
}

func midright(mid int, i int, set []string) string {
    var out string
    for j, mod := range set {
        out += fmt.Sprintf("%s%d ", mod, i + 1)
        if j < len(set) - 1 {
            out += fmt.Sprintf("ssl%d ", i + 1)
        }
    }
    if i < mid - 1 {
        out += fmt.Sprintf("sl%d ", i + 1)
    }
    return out
}

/*
    Parse each of the modules to attach depth numbers and fill in the separators.
    e.g.:
    "left": [ [ "a" ], [ "b", "c" ], [ "d" ] ] -> a1 sl1 b2 ssl2 c2 sl2 d3 sl3b
*/
func parseModuleList(b Bar) (string, int) {
    var out string
    var tray int = 2

    // Parse the left modules
    if len(b.Left) > 0 {
        out += "modules-left = "
        for i, set := range b.Left {
            for j, mod := range set {
                if mod == "tray" {
                    tray = -1
                    continue
                }
                out += fmt.Sprintf("%s%d ", mod, i + 1)
                if j < len(set) - 1 {
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

    // Parse the center modules
    var cout string
    if len(b.Center) > 0 {
        cout = "modules-center = "
        var mid int = int(math.Floor(float64(len(b.Center)) / 2.0))

        // Handle even length major modules in center
        if len(b.Center) % 2 == 0 {
            cout += fmt.Sprintf("sr%db ", mid)
            for i := 0; i < mid; i++ {
                cout += midleft(mid, i, b.Center[i])
            }
            cout += "ssm1 "
            for i := 0; i < mid; i++ {
                cout += midright(mid, i, b.Center[mid + i])
            }
            cout += fmt.Sprintf("sl%db ", mid)

        // Handle odd length major modules in center
        } else {
            cout += fmt.Sprintf("sr%db ", mid + 1)
            for i := 0; i < mid; i++ {
                cout += midleft(mid + 1, i, b.Center[i])
            }
            var cmid int = int(math.Floor(float64(len(b.Center[mid])) / 2.0))
            
            // Handle even length minor modules in center
            if len(b.Center[mid]) % 2 == 0 {
                for i := 0; i < cmid; i++ {
                    cout += fmt.Sprintf("%s1 ", b.Center[mid][i])
                    if i < cmid - 1 {
                        cout += "ssr1 "
                    }
                }
                cout += "ssm1 "
                for i := 0; i < cmid; i++ {
                    cout += fmt.Sprintf("%s1 ", b.Center[mid][cmid + i])
                    if i < cmid - 1 {
                        cout += "ssl1 "
                    }
                }

            // Handle odd length minor modules in center
            } else {
                for i := 0; i < cmid; i++ {
                    cout += fmt.Sprintf("%s1 ", b.Center[mid][i])
                    if i < cmid {
                        cout += "ssr1 "
                    }
                }
                for i := 0; i < cmid + 1; i++ {
                    cout += fmt.Sprintf("%s1 ", b.Center[mid][cmid + i])
                    if i < cmid {
                        cout += "ssl1 "
                    }
                }
            }
            cout += "sl1 "
            for i := 0; i < mid; i++ {
                cout += midright(mid + 1, i + 1, b.Center[mid + i + 1])
            }
            cout += fmt.Sprintf("sl%db ", mid + 1)
        }
    }
    cout = strings.Trim(cout, " ")
    out += cout + "\n"

    // Parse the right modules
    var rout string
    if len(b.Right) > 0 {
        for i, set := range b.Right {
            for j, mod := range set {
                if mod == "tray" {
                    tray = 1
                    continue
                }
                rout += fmt.Sprintf("%s%d ", mod, i + 1)
                if j < len(set) - 1 {
                    rout += fmt.Sprintf("ssr%d ", i + 1)
                }
            }
            if i == len(b.Right) - 1 {
                rout += fmt.Sprintf("sr%db ", i + 1)
            } else {
                rout += fmt.Sprintf("sr%d ", i + 1)
            }
        }
    }
    rout = strings.Trim(rout, " ")

    // Reverse right
    if rout != "" {
        routs := strings.Split(rout, " ")
        out += "modules-right = "
        for i := len(routs) - 1; i >= 0; i-- {
            out += routs[i] + " "
        }
        out = strings.Trim(out, " ")
    }
    
    // Return modules
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
