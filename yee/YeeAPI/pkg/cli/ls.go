package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	"github.com/thomasf/lg"
	"github.com/thomasf/yeelight/pkg/ssdp"
)

type ListCmd struct {
	Common CommonFlags
	IDfmt  bool
}

func (*ListCmd) Name() string     { return "ls" }
func (*ListCmd) Synopsis() string { return "List bulbs" }
func (*ListCmd) Usage() string {
	return ``
}

func (p *ListCmd) SetFlags(fs *flag.FlagSet) {
	p.Common.Register(fs)
	fs.BoolVar(&p.IDfmt, "q", false, "only print id's")
}

func (p *ListCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	devices, err := ssdp.GetDevices("enp0s31f6")
	if err != nil {
		lg.Fatal(err)
	}
	for _, v := range devices {
		if p.IDfmt {
			fmt.Println(v.ID)
		} else {
			fmt.Printf("%v [%v] %v\n", v.ID, v.Name, v.LocationAddr())
		}
		// if v.Name == "tf-desk" {
		// 	conn := &yeel.Conn{
		// 		Device: v,
		// 	}
		// 	go conn.Open()
		// 	time.Sleep(time.Second)

		// 	conn.ExecCommand(yeel.SetBrightnessNormCommand{Brightness: 0.2, Duration: time.Second})

		// 	_, err := conn.ExecCommand(yeel.SetNameCommand{Name: "tf-desk"})
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// }
	}

	return subcommands.ExitSuccess
}
