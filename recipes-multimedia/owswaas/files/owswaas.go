/* Owen Wilson Saying "Wow" as a Service
 *
 * This is a quickly slopped together program fit
 * for no one and for no purpose.
 *
 * Abandon all hope, ye who enter here.
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	// I'll never run this on any other OS, so :shrug:
	"golang.org/x/sys/unix"
)

type GPIOPin struct {
	Pin       uint // Pin's numeric ID in software. Must be set prior to Initialize().
	path      string
	valueFile *os.File
}

type OutputPin struct {
	GPIOPin
}

type InputPin struct {
	GPIOPin
	risingEdge  bool
	fallingEdge bool
	event       unix.EpollEvent
}

type audioSamples []string

type audioSampleSet struct {
	selected uint
	paths    []string
	samples  map[string]audioSamples

	// The sample we play at boot and for volume adjustment
	defaultSample string
}

// Randomly plays an audio sample from one or more audio sets
// organized by directory.
type RandomSamplePlayer struct {
	Program     string
	Mixer       string
	ControlFile string
	sampleSet   audioSampleSet
	volume      int64
}

// Load settings file. Currently, this is just the volume.
// Honestly, this is a sloppy hack that could be done with something like
// asound.state, but it was quicker to throw this when running
// owswaas as a non-root member of the audio group.
func (r *RandomSamplePlayer) LoadSettings(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		entries := strings.Split(line, "=")
		if len(entries) != 2 {
			return errors.New("Invalid line: " + line)
		}

		key := strings.TrimSpace(entries[0])
		value := strings.TrimSpace(entries[1])

		if key == "volume" {
			volume, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return errors.New("Failed to parse volume [0, 100]: " + value)
			}
			r.SetVolume(volume)
		} else {
			log.Printf("Unknown key: " + key)
		}
	}

	return nil
}

func (r *RandomSamplePlayer) Volume() int64 {
	return r.volume
}

func (r *RandomSamplePlayer) SetVolume(value int64) error {
	if value > 100 {
		r.volume = 100
	} else if value < 0 {
		r.volume = 0
	} else {
		r.volume = value
	}

	volumeStr := fmt.Sprintf("%d%%", r.volume)
	log.Printf("Setting volume to %s\n", volumeStr)

	cmd := exec.Command(r.Mixer, "set", "PCM", volumeStr)
	return cmd.Run()
}

func (r *RandomSamplePlayer) IncrementVolume() error {
	if r.volume < 20 {
		r.volume += 2
	} else {
		r.volume += 5
	}

	return r.SetVolume(r.volume)
}

func (r *RandomSamplePlayer) DecrementVolume() error {
	if r.volume < 20 {
		r.volume -= 2
	} else {
		r.volume -= 5
	}
	return r.SetVolume(r.volume)
}

// Recurse audio root directory and gather audio sample sets on a
// per-diretory basis
func (r *RandomSamplePlayer) LoadSamples(rootDir string) error {
	r.sampleSet.samples = make(map[string]audioSamples)

	err := filepath.Walk(rootDir, func(currPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only concerned with files containing a .wav suffix
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".wav") {
			return nil
		}

		dirPath := path.Dir(currPath)
		r.sampleSet.samples[dirPath] = append(r.sampleSet.samples[dirPath], currPath)
		//log.Printf("Added %s to r.sampleSet for %s\n", currPath, dirPath)

		return nil
	})

	if err != nil {
		return err
	}

	i := uint(0)
	for dir := range r.sampleSet.samples {
		r.sampleSet.paths = append(r.sampleSet.paths, dir)

		// Favor our namesake as the default
		if path.Base(dir) == "wowz" {
			r.sampleSet.selected = i
			r.sampleSet.defaultSample = r.sampleSet.samples[dir][0]
		}
		i++
	}

	// Seed PRNG based on current time to randomize sample selection
	rand.Seed(time.Now().UnixNano())

	return nil
}

func (r *RandomSamplePlayer) NextSet() {
	r.sampleSet.nextSet()
	log.Printf("Selected sample set: %s\n", r.sampleSet.SelectedName())
}

func (r *RandomSamplePlayer) PlayRandomSample() error {
	sample := r.sampleSet.randomSample()
	log.Println("Playing file: " + sample)

	cmd := exec.Command(r.Program, sample)
	return cmd.Run()
}

func (r *RandomSamplePlayer) PlayDefaultSample() error {
	cmd := exec.Command(r.Program, r.sampleSet.defaultSample)
	return cmd.Run()
}

// Advance to the next set of samples
func (s *audioSampleSet) nextSet() {
	s.selected = (s.selected + 1) % uint(len(s.paths))
}

// Return (directory) name of selected sample set
func (s *audioSampleSet) SelectedName() string {
	return path.Base(s.paths[s.selected])
}

// Get a random sample from the currently selected set
func (s *audioSampleSet) randomSample() string {
	dir := s.paths[s.selected]
	samples := s.samples[dir]

	idx := rand.Intn(len(samples))
	return samples[idx]
}

// Configure the pin
func (p *GPIOPin) Initialize(isInput bool) error {
	p.path = fmt.Sprintf("/sys/class/gpio/gpio%d/", p.Pin)

	if _, err := os.Stat(p.path); err != nil {

		log.Printf("Attempting to export %d\n", p.Pin)
		err := ioutil.WriteFile("/sys/class/gpio/export", []byte(fmt.Sprintf("%d", p.Pin)), 0200)
		if err != nil {
			return err
		}

		if _, err := os.Stat(p.path); err != nil {
			return errors.New("Failed to export " + p.path + " -- " + err.Error())
		}
	}

	var dir string
	if isInput {
		dir = "in"
	} else {
		dir = "out"
	}
	var b []byte

	var err error
	if b, err = ioutil.ReadFile(p.path + "direction"); err != nil {
		return err
	}
	currDir := strings.TrimSpace(string(b))

	if currDir != dir {
		if err := ioutil.WriteFile(p.path+"direction", []byte(dir), 0644); err != nil {
			return err
		}
	}

	if isInput {
		p.valueFile, err = os.Open(p.path + "value")
	} else {
		p.valueFile, err = os.Create(p.path + "value")
	}

	return err
}

func MakeOutputPin(pin uint) OutputPin {
	return OutputPin{GPIOPin{pin, "", nil}}
}

func (p *OutputPin) Initialize() error {
	return p.GPIOPin.Initialize(false)
}

func (p *OutputPin) SetValue(val bool) error {
	data := []byte{'0', '\n'}
	if val {
		data[0] = '1'
	}

	_, err := p.valueFile.Write(data)
	return err
}

func MakeInputPin(pin uint, risingEdge, fallingEdge bool) InputPin {
	return InputPin{GPIOPin{pin, "", nil}, risingEdge, fallingEdge, unix.EpollEvent{}}
}

func (p *InputPin) Initialize(epollFd int) error {
	p.GPIOPin.Initialize(true)

	var mode string
	if p.risingEdge && p.fallingEdge {
		mode = "both"
	} else if p.risingEdge && !p.fallingEdge {
		mode = "rising"
	} else if !p.risingEdge && p.fallingEdge {
		mode = "falling"
	} else {
		mode = "none"
	}

	var b []byte
	var err error
	if b, err = ioutil.ReadFile(p.path + "edge"); err != nil {
		return err
	}

	currMode := strings.TrimSpace(string(b))
	if currMode != mode {
		log.Printf("Got %s, expected %s\n", string(b), mode)
		if err := ioutil.WriteFile(p.path+"edge", []byte(mode), 0644); err != nil {
			return err
		}
	}

	p.event.Events = unix.EPOLLPRI | unix.EPOLLERR
	p.event.Fd = int32(p.valueFile.Fd())

	if epollFd >= 0 {
		err := unix.EpollCtl(epollFd, unix.EPOLL_CTL_ADD, int(p.valueFile.Fd()), &p.event)
		if err != nil {
			return err
		}
	}

	// Perform initial read
	p.Value()
	return nil
}

func (p *InputPin) Value() bool {
	var b []byte
	var err error

	if _, err := p.valueFile.Seek(0, 0); err != nil {
		log.Println(err)
		return false
	}

	if b, err = ioutil.ReadAll(p.valueFile); err != nil {
		log.Println(err)
		return false
	}

	ret := (b[0] == '1')
	return ret
}

func (p *InputPin) ChangedState(events []unix.EpollEvent) (bool, bool) {
	for _, e := range events {
		if e.Fd == p.event.Fd {
			return true, p.Value()
		}
	}

	return false, false
}

const (
	normal_await_btn_press   = iota
	normal_await_btn_release = iota
	voldec_await_btn_release = iota
	volinc_await_btn_release = iota
)

const (
	button_depressed = false
	button_open      = true
)

type Owswaas struct {
	Player RandomSamplePlayer
	Led    OutputPin
	Btn    InputPin
	Sw     InputPin

	settingsFile string
	epollFd      int
	ledCtrl      chan uint

	state     uint
	btn_timer time.Time
}

// Return a new Owen Wilson Saying "Wow" as a Service context instance with
// default values.
func NewOwswaas() *Owswaas {
	ctx := new(Owswaas)

	ctx.ledCtrl = make(chan uint)

	ctx.Player.Program = "/usr/bin/aplay"
	ctx.Player.Mixer = "/usr/bin/amixer"

	ctx.Led = MakeOutputPin(18)
	ctx.Btn = MakeInputPin(24, true, true)
	ctx.Sw = MakeInputPin(23, true, true)

	ctx.state = normal_await_btn_press

	// Just for the sake of having something in here at startup
	ctx.btn_timer = time.Now()

	return ctx
}

// Controls LED according to received uint:
//	0 - Turn the LED off
//	1 - Turn the LED on
//	n - Blink the LED at a rate n HZ, with a 50% duty cycle
func ledCtrl(led *OutputPin, ctrl <-chan uint) {
	request := <-ctrl
	for {

		switch request {
		case 0, 1:
			led.SetValue(request == 1)
			request = <-ctrl
		default:
			halfPeriod := time.Duration(int64(1000000) / int64(request))
			led.SetValue(true)
			time.Sleep(halfPeriod * time.Microsecond)
			led.SetValue(false)

			select {
			case request = <-ctrl:
				break
			default:
				time.Sleep(halfPeriod * time.Microsecond)
			}
		}
	}
}

// Perform initialization routines required before calls to HandleEvents()
func (ctx *Owswaas) Initialize(audioDir, settingsFile string) error {
	var err error

	// Launch LED control task
	go ledCtrl(&ctx.Led, ctx.ledCtrl)

	ctx.epollFd, err = unix.EpollCreate1(0)
	if err != nil {
		return err
	}

	ctx.settingsFile = settingsFile
	if err = ctx.Player.LoadSettings(settingsFile); err != nil { return err }
	if err = ctx.Player.LoadSamples(audioDir); err != nil { return err }

	if err = ctx.Led.Initialize(); err != nil { return err }
	if err = ctx.Btn.Initialize(ctx.epollFd); err != nil { return err }
	if err = ctx.Sw.Initialize(ctx.epollFd); err != nil { return err }

	// Play a sample to denote that we're up and running
	ctx.ledCtrl <- 1
	err = ctx.Player.PlayDefaultSample()
	ctx.ledCtrl <- 0
	return err
}

// Update settings file to reflect runtime changes
func (ctx *Owswaas) updateSettings() error {
	contents := []byte(fmt.Sprintf("volume = %d\n", ctx.Player.Volume()))
	return ioutil.WriteFile(ctx.settingsFile, contents, 0644)
}

// Play a random sample in the current set, with the LED illuminated for the
// duration of the audio sample.
func (ctx *Owswaas) PlayRandomSample() error {
	ctx.ledCtrl <- 1
	err := ctx.Player.PlayRandomSample()
	ctx.ledCtrl <- 0
	return err
}

func (ctx *Owswaas) HandleEvents() error {
	var err error

	events := []unix.EpollEvent{unix.EpollEvent{}, unix.EpollEvent{}}
	_, err = unix.EpollWait(ctx.epollFd, events, -1)
	if err != nil {
		return err
	}

	sw_rotated, _ := ctx.Sw.ChangedState(events)
	btn_changed, btn_state := ctx.Btn.ChangedState(events)

	switch ctx.state {

	case normal_await_btn_press:
		if sw_rotated {
			// Handle switch event, but remain in the current state
			ctx.Player.NextSet()
			err = ctx.PlayRandomSample()
		} else if btn_changed {
			if btn_state == button_depressed {
				// Button pressed - start timer
				ctx.btn_timer = time.Now()
				ctx.state = normal_await_btn_release
			}
		}

	case normal_await_btn_release:
		if btn_changed {
			press_dur := time.Now().Sub(ctx.btn_timer).Seconds()

			if press_dur <= 2 {
				// This is a request to play an audio sample
				err = ctx.PlayRandomSample()
				ctx.state = normal_await_btn_press
			} else if press_dur >= 3 {
				// Slow blink
				ctx.ledCtrl <- 4

				// Request to enter volume decrement mode
				ctx.state = voldec_await_btn_release
			}
		}

	case voldec_await_btn_release:
		if sw_rotated {
			// Decrement volume and remain in this mode
			if err = ctx.Player.DecrementVolume(); err == nil {
				err = ctx.Player.PlayDefaultSample()
			}
		} else if btn_changed && btn_state == button_open {
			// Fast blink
			ctx.ledCtrl <- 8

			// Advance to volume increment mode
			ctx.state = volinc_await_btn_release
		}

	case volinc_await_btn_release:
		if sw_rotated {
			// Increment volume and remain in this mode
			if err = ctx.Player.IncrementVolume(); err == nil {
				err = ctx.Player.PlayDefaultSample()
			}
		} else if btn_changed && btn_state == button_open {
			ctx.ledCtrl <- 0

			// Persist volume change and return to normal operation
			err = ctx.updateSettings()
			ctx.state = normal_await_btn_press
		}

	default:
		err = fmt.Errorf("Invalid state: %d", ctx.state)
	}

	return err
}

func main() {
	ctx := NewOwswaas()

	audioDir := "/audio"
	settings := filepath.Join(audioDir, "owswaas.conf")

	flag.UintVar(&ctx.Led.Pin, "gpio-led", ctx.Led.Pin, "LED control GPIO number")
	flag.UintVar(&ctx.Btn.Pin, "gpio-btn", ctx.Btn.Pin, "Button GPIO pin number")
	flag.UintVar(&ctx.Sw.Pin, "gpio-sw", ctx.Sw.Pin, "Switch GPIO number")
	flag.StringVar(&audioDir, "audio-dir", audioDir, "Audio sample root directory")
	flag.StringVar(&settings, "settings", settings, "Runtime settings file")
	flag.StringVar(&ctx.Player.Program, "playback", ctx.Player.Program, "Audio playback program")
	flag.StringVar(&ctx.Player.Mixer, "mixer", ctx.Player.Mixer, "Audio mixer program")

	flag.Parse()

	if err := ctx.Initialize(audioDir, settings); err != nil {
		log.Fatal(err)
	}

	for {
		if err := ctx.HandleEvents(); err != nil {
			log.Println(err)
		}
	}
}
