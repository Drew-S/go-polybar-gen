# Go Polybar Generator

Generate a polybar configuration file using a simple `config.json` file along with templates for each module. This program is designed to automatically fill out module foreground and background colours in a similar fashion to that of powerline.

The config.json file is where the polybar setup is done:

```json
{
    "bars": [
        {
            "name": "main",
            "size": {
                "w": "100%",
                "h": "24",
                "y": "0",
                "x": "0"
            },
            "monitor": "DVI-D-0",
            "left": [
                [ "cpu" ],
                [ "cpu", "cpu" ]
            ],
            "right": [ ... ]
        }
    ],
    "fonts": [
        "DroidSansNerdFont Mono:size=12;3",
        "..."
    ],
    "sep": {
        "left": ">>",
        "right": "<<"
    },
    "subsep": {
        "left": ">",
        "right": "<"
    },
    "data": {
        "somedata": "something"
    }
}
```

in the above json example the `cpu` module (which must match the file `modules/cpu.ini`) is repeated three times, the same cpu template as defined in the ini is generated into two modules: cpu1 and cpu2, each with different backgrounds for the different levels.

The `cpu.ini` file has the following template:

```dosini
[module/cpu{{.Depth}}]
type = internal/cpu

label = CPU %percentage%%
label-padding = 1

format-background = $$Background$$
format-foreground = $$Foreground$$
```

The template is the exact same as polybar config, only with the depth variable, which is needed to differentiate between the different depths and the background and foreground variables. The variables `$$Foreground$$` and `$$Background$$` get replaced with color references according to the depth. The colors default to using Xresource color values.

With this you do not need to manually fill in all the different separators, for left, right or sub separators, nor need to worry about module background colors matching. Simply use the templates and fill in the config to generate a bar.

The `fonts` will automatically get filled in for each bar the index matching the same as polybar (index 0 = font-0, use with %{T1}).

The `sep` and `subsep` is where you define the separators between modules. `subsep` will appear between modules on the same level : `[ "mod", "mod2" ]`, while `sep` will appear between groups : `[ "mod" ], [ "mod2" ]`

The `data` object houses some user based data that can be reference in template. The above example is referenced like:

```dosini
{{.Data.somedata}}
```

Currently the program is limited to only producing left and right modules and not center ones.

To see a complete example, see the example folder.
