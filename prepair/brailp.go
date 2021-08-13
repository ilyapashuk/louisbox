// Copyright 2021 ilyapashuk<ilusha.paschuk@gmail.com>

//    This file is part of Louisbox.

//    Louisbox is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.

//    Louisbox is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.

//    You should have received a copy of the GNU General Public License
//    along with Louisbox.  If not, see <https://www.gnu.org/licenses/>.

// main file for this module

package prepair
import "strings"
import "flag"
import "fmt"
import "brbox"
import "path/filepath"



var BomSequence = []byte("\xef\xbb\xbf")
// такой формат сохранения множиств символов введен для компактности кода
type SymbolGroupe string
func (c SymbolGroupe) Contains(r rune) bool {
for _,i := range c {
if r == i {
return true
}
}
return false
}
var Digits SymbolGroupe = SymbolGroupe("0123456789")

type Handler interface {
Handle(string, []string) *string
}
type HandlerFunc func(string, []string) *string
func (c HandlerFunc) Handle(st string, opts []string) *string {
return c(st, opts)
}
var Handlers map[string]Handler = make(map[string]Handler)

func CallHandler(t string, hn string, args []string) *string {
res := Handlers[hn].Handle(t,args)
return res
}
func ScriptCallHandler(t,hn string, args []string) interface{} {
res := CallHandler(t,hn,args)
if res == nil {
return false
} else {
return *res
}
}
func ListHandlers() []string {
res := make([]string,0,len(Handlers))
for k,_ := range Handlers {
res = append(res,k)
}
return res
}

// this function will check weather provided rune is a digit
func isdig(s rune) bool {
return Digits.Contains(s)
}

type HandlerChain [][]string
func (c HandlerChain) Handle(t string,_ []string) *string {
res := t
for _,v := range c {
var r *string
if len(v) == 1 {
r = Handlers[v[0]].Handle(res,nil)
} else {
r = Handlers[v[0]].Handle(res,v[1:])
}
if r == nil {
return nil
}
res = *r
}
return &res
}

func Prepair(args []string) {
cmdline := flag.NewFlagSet("prepair", flag.ExitOnError)
cmdline.Usage = func() {
fmt.Println("usage: prepair [options] <infile>")
cmdline.PrintDefaults()
}
outext := cmdline.String("outext","","extension for new file")
cmdline.Parse(args)
fn := cmdline.Arg(0)
t,err := brbox.ReadInputFile(fn)
if err != nil {
panic(err)
}
tp := Handlers["louis"].Handle(t,nil)
t = *tp
t = strings.ReplaceAll(t,"\r","")
lines := strings.Split(t,"\n")
res := make([]*string, len(lines))
for i,l := range lines {
res[i] = Handlers["linewrap"].Handle(l, []string{"1"})
}
sb := new(strings.Builder)
sb.Grow(len(t))
for _,l := range res {
if l != nil {
sb.WriteString(*l)
sb.WriteString("\n")
}
}
t = sb.String()
var rfn string
if *outext != "" {
rfn = strings.TrimSuffix(fn, filepath.Ext(fn)) + "." + *outext
} else {
rfn = cmdline.Arg(1)
}
err = brbox.WriteOutputFile(rfn, t, true)
if err != nil {
panic(err)
}
}
func init() {
brbox.Subcommands["prepair"] = Prepair
}