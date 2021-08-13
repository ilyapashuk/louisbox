package gendis

import (
"brbox"
"os"
"strings"
"github.com/ilyapashuk/go-braille/translation"
)


func Gendis(args []string) {
table,err := os.ReadFile(args[0])
if err != nil {
panic(err)
}
ts := string(table)
ts = strings.ReplaceAll(ts, "\r", "")
td := strings.Split(ts, "\n")
var ll bool
disfile,err := os.Create(args[1])
if err != nil {
panic(err)
}
defer disfile.Close()
for _,line := range td {
if line == "" {
continue
}
if line == "#langlets" {
ll = true
continue
}
if line == "#endlanglets" {
ll = false
continue
}
if strings.HasPrefix(line,"#") {
disfile.WriteString(line + "\n")
continue
}
rl,err := translation.ParseRule(line)
if err != nil {
panic(err)
}
switch rl.Code {
case "sign":
if rl.Runes[0] == '\\' {
disfile.WriteString("display " + "\\\\" + " " + rl.Dots.String() + "\n")
} else {
disfile.WriteString("display " + string(rl.Runes[0]) + " " + rl.Dots.String() + "\n")
}
case "digit":
disfile.WriteString("display " + string(rl.Runes[0]) + " " + rl.Dots.String() + "b\n")
case "uplow":
if ll {
disfile.WriteString("display " + string(rl.Runes[0]) + " " + rl.Dots.String() + "9a\n")
} else {
disfile.WriteString("display " + string(rl.Runes[0]) + " " + rl.Dots.String() + "a\n")
}
if ll {
disfile.WriteString("display " + string(rl.Runes[1]) + " " + rl.Dots.String() + "9\n")
} else {
disfile.WriteString("display " + string(rl.Runes[1]) + " " + rl.Dots.String() + "\n")
}
}
}
}
func init() {
brbox.Subcommands["gendis"] = Gendis
}