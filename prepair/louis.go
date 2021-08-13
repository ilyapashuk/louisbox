package prepair
import "strings"
import "os/exec"
import "path/filepath"
import "brbox"
import "os"

type LouisHandler struct {
}
func (c *LouisHandler) Handle(t string, _ []string) *string {
tb := filepath.Join(brbox.ExeDir,"louisdata")
os.Setenv("LOUIS_TABLEPATH",tb)
lt := filepath.Join(brbox.ExeDir,"lou_translate.exe")
cm := exec.Command(lt,"disfile.dis,ru-litbrl.ctb")
cm.Stdin = strings.NewReader(t)
cm.Stderr = os.Stderr
bres,err := cm.Output()
if err != nil {
panic(err)
}
res := string(bres)
return &res
}

func init() {
Handlers["louis"] = new(LouisHandler)
}