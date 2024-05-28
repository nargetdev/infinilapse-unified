package oscmanager

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/hypebeast/go-osc/osc"
	"strconv"
	"time"
)

const (
	LX_ADDR = "127.0.0.1"
	//LX_ADDR = "100.101.1.7"
	LX_PORT = 3031

	//BLENDER_ADDR = "100.106.164.86"
	//BLENDER_ADDR = "100.101.1.7"
	BLENDER_ADDR = "127.0.0.1"
	BLENDER_PORT = 9001
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

/**
`CubeSchedule.blend` is a blender project with 691200 frames and set to play at 8 frames per second for 24 hours of frames.
The X position of the flaming monkey head sends an OSC message as /lx/output/brightness which will set our grow lights
for the correct time of day.
*/

func StartBlenderDailyCycle() {
	resetMsg := "/blender/reset"
	sendOscTrue(resetMsg)

	//client := osc.NewClient(BLENDER_ADDR, BLENDER_PORT)
	//msg := osc.NewMessage(resetMsg)
	//msg.Append(true)
	//_ = client.Send(msg)

	playMsg := "/blender/play"
	sendOscTrue(playMsg)

	//msg = osc.NewMessage(playMsg)
	//msg.Append(true)
	//_ = client.Send(msg)
}

func SynchronizeBlenderPlayhead() {
	SetBlenderFrame()
	playMsg := "/blender/play"
	sendOscTrue(playMsg)
}

func SetBlenderFrame() {
	//oscPath := "/scenes/Scene/frame_current"
	//oscPath := "/blender/out/frame"
	//oscPath := "/blender/frame/jump"
	currentFrameOfDay := getCurrentFrameOfDay()
	fmt.Printf("current frame is: %d\n", currentFrameOfDay)
	//sendOscIntValue(oscPath, currentFrameOfDay)

	sendoscCmd := "sendosc " + BLENDER_ADDR + " 9001 /scenes/Scene/frame_current i " + strconv.Itoa(currentFrameOfDay)

	stdout, err := script.Exec(sendoscCmd).Stdout()
	if err != nil {
		fmt.Printf("error running "+sendoscCmd+" ::: %s", err)
	}
	fmt.Println(stdout)
}

func frameJumpString(path string, deltaString string) {
	sendOscString(path, deltaString)
}

func sendOscString(path string, s string) {
	client := osc.NewClient(BLENDER_ADDR, BLENDER_PORT)
	msg := osc.NewMessage(path)
	//msg.Append(s)
	msg.Append(27)
	_ = client.Send(msg)
}

func sendOscByPathAndValue(path string, value interface{}) {
	client := osc.NewClient(BLENDER_ADDR, BLENDER_PORT)
	msg := osc.NewMessage(path)
	msg.Append(value)
	_ = client.Send(msg)
}

func sendOscIntValue(path string, value int) {
	client := osc.NewClient(BLENDER_ADDR, BLENDER_PORT)
	msg := osc.NewMessage(path)
	msg.Append(value)
	_ = client.Send(msg)
}

func sendOscTrue(path string) {
	client := osc.NewClient(BLENDER_ADDR, BLENDER_PORT)
	msg := osc.NewMessage(path)
	msg.Append(true)
	_ = client.Send(msg)
}

func sendOscBool(path string, tf bool) {
	client := osc.NewClient(BLENDER_ADDR, BLENDER_PORT)
	msg := osc.NewMessage(path)
	msg.Append(tf)
	msg.Append("delta=1000")
	_ = client.Send(msg)
}

func daySeconds(t time.Time) int {
	year, month, day := t.Date()
	t2 := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return int(t.Sub(t2).Seconds())
}

func dayFrames(seconds int) int {
	return seconds * 8
}

func getCurrentFrameOfDay() int {
	return dayFrames(daySeconds(time.Now()))
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
