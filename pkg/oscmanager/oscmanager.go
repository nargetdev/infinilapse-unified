package oscmanager

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
)

const (
	LX_ADDR = "100.101.1.7"
	LX_PORT = 3031
)

func FadeMaster(brightness float64) {
	if brightness > 1.0 {
		brightness = 1.0
	} else if brightness < 0.0 {
		brightness = 0.0
	}

	client := osc.NewClient(LX_ADDR, LX_PORT)
	msg := osc.NewMessage("/lx/output/brightness")
	msg.Append(brightness)
	_ = client.Send(msg)
}

func OscListen(connectionString string) {
	var err error

	addr := connectionString
	d := osc.NewStandardDispatcher()
	err = d.AddMsgHandler("*", func(msg *osc.Message) {
		osc.PrintMessage(msg)
	})
	check(err)

	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}
	err = server.ListenAndServe()
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Printf("ERROR ;; %s\n", err)
	}
}
