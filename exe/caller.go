package exe

import (
	"fmt"
	"os/exec"
)

func RunBin() {
	cmd := exec.Command("./exe/RunExternally.exe")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
