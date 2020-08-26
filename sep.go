package main

import (
    "fmt"
	"os"
	"text/template"
)

const separators string =
`;================================================SEPARATORS
;======================================================BASE

[module/sepright]
type = custom/text
content = {{.Sep.Right}}

[module/sepleft]
type = custom/text
content = {{.Sep.Left}}

[module/subsepright]
type = custom/text
content = {{.Subsep.Right}}

[module/subsepleft]
type = custom/text
content = {{.Subsep.Left}}

;=======================================================MID

[module/ssm1]
type = custom/text
content = {{.Sep.Center}}
content-background = ${colors.color2}
content-foreground = ${colors.color1}

;======================================================LEFT
{{range $i, $el := colors .Colors.Colors}}
[module/sl{{add $i 1}}]
inherit = module/sepleft
content-background = {{color $i 3}}
content-foreground = {{color $i 2}}
{{end}}
; BACKGROUND
{{range $i, $el := colors .Colors.Colors}}
[module/sl{{add $i 1}}b]
inherit = module/sepleft
content-background = ${colors.background}
content-foreground = {{color $i 2}}
{{end}}
;===================================================SUBLEFT
{{range $i, $el := colors .Colors.Colors}}
[module/ssl{{add $i 1}}]
inherit = module/subsepleft
content-background = {{color $i 2}}
content-foreground = {{color $i 1}}
{{end}}
; BACKGROUND
{{range $i, $el := colors .Colors.Colors}}
[module/ssl{{add $i 1}}b]
inherit = module/subsepleft
content-background = ${colors.background}
content-foreground = {{color $i 2}}
{{end}}
;=====================================================RIGHT
{{range $i, $el := colors .Colors.Colors}}
[module/sr{{add $i 1}}]
inherit = module/sepright
content-background = {{color $i 3}}
content-foreground = {{color $i 2}}
{{end}}
; BACKGROUND
{{range $i, $el := colors .Colors.Colors}}
[module/sr{{add $i 1}}b]
inherit = module/sepright
content-background = ${colors.background}
content-foreground = {{color $i 2}}
{{end}}
;==================================================SUBRIGHT
{{range $i, $el := colors .Colors.Colors}}
[module/ssr{{add $i 1}}]
inherit = module/subsepright
content-background = {{color $i 2}}
content-foreground = {{color $i 1}}
{{end}}
; BACKGROUND
{{range $i, $el := colors .Colors.Colors}}
[module/ssr{{add $i 1}}b]
inherit = module/subsepright
content-background = ${colors.background}
content-foreground = {{color $i 2}}
{{end}}
;===================================================MODULES

`

// Generates all the different separators for the config, (major left/right,
// minor left/right, and straight to background varients)
func generateSeparators(config Config, file *os.File) {
    t := template.Must(template.New("separators").Funcs(template.FuncMap{
        "add": func(i int, a int) int {
            return i + a
        },
        "colors": func(colors []string) []int {
            return make([]int, len(colors) - 2)
        },
        "color": func(i int, a int) string {
            if i + a == len(config.Colors.Colors) {
                return "${colors.background}"
            } else {
                return fmt.Sprintf("${colors.color%d}", i + a)
            }
        },
    }).Parse(separators))

    t.Execute(file, config)
}
