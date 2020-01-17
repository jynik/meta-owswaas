package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	// I'll never run this on any other OS, so :shrug:
	"golang.org/x/sys/unix"
)

type GPIOPin struct {
	pin       uint
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
	event		unix.EpollEvent
}

type audioSamples []string

type audioSampleSet struct {
	selected uint
	paths	[]string
	samples map[string]audioSamples
}

type RandomSamplePlayer struct {
	Program string
	sampleSet audioSampleSet
}

// Recurse audio root directory and gather audio sample sets on a
// per-diretory basis
func (r *RandomSamplePlayer) LoadSamples(rootDir string) {
	r.sampleSet.samples = make(map[string]audioSamples)

	err := filepath.Walk(rootDir, func (currPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only concerned with files containing a .wav suffix
		if info.IsDir() { return nil }
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".wav") { return nil }

		dirPath := path.Dir(currPath)
		r.sampleSet.samples[dirPath] = append(r.sampleSet.samples[dirPath], currPath)
		//log.Printf("Added %s to r.sampleSet for %s\n", currPath, dirPath)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	i := uint(0)
	for dir := range r.sampleSet.samples {
		r.sampleSet.paths = append(r.sampleSet.paths, dir)

		// Favor our namesake as the default
		if path.Base(dir) == "wowz" {
			r.sampleSet.selected = i
		}
		i++
	}

	// Seed PRNG based on current time to randomize sample selection
	rand.Seed(time.Now().UnixNano())
}

func (r *RandomSamplePlayer) NextSet() {
	r.sampleSet.nextSet()
	log.Printf("Selected sample set: %s\n", r.sampleSet.SelectedName())
}

func (r *RandomSamplePlayer) PlayRandomSample() {
	sample := r.sampleSet.randomSample()
	log.Println("Playing file: " + sample)

	cmd := exec.Command(r.Program, sample)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
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

func (p *GPIOPin) Setup(isInput bool) {
	p.path = fmt.Sprintf("/sys/class/gpio/gpio%d/", p.pin)

	if _, err := os.Stat(p.path); err != nil {
		err := ioutil.WriteFile("/sys/class/gpio/export", []byte(fmt.Sprintf("%d", p.pin)), 0200)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := os.Stat(p.path); err != nil {
			log.Fatal("Failed to export " + p.path + " -- " + err.Error())
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
	if b, err = ioutil.ReadFile(p.path+"direction"); err != nil {
		log.Fatal(err)
	}
	currDir := strings.TrimSpace(string(b))

	if currDir != dir {
		if err := ioutil.WriteFile(p.path+"direction", []byte(dir), 0644); err != nil {
			log.Fatal(err)
		}
	}

	if isInput {
		p.valueFile, err = os.Open(p.path + "value")
	} else {
		p.valueFile, err = os.Create(p.path + "value")
	}

	if err != nil {
		log.Fatal(err)
	}
}

func MakeOutputPin(pin uint) OutputPin {
	return OutputPin{ GPIOPin{pin, "", nil} }
}

func (p *OutputPin) Setup() {
	p.GPIOPin.Setup(false)
}

func (p *OutputPin) SetValue(val bool) {
	data := []byte{'0', '\n'}
	if val {
		data[0] = '1'
	}

	if _, err := p.valueFile.Write(data); err != nil {
		log.Fatal(err)
	}
}

func MakeInputPin(pin uint, risingEdge, fallingEdge bool) InputPin {
	return InputPin { GPIOPin{pin, "", nil}, risingEdge, fallingEdge, unix.EpollEvent{} }
}

func (p *InputPin) Setup(epollFd int) {
	p.GPIOPin.Setup(true)

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
	if b, err = ioutil.ReadFile(p.path+"edge"); err != nil {
		log.Fatal(err)
	}

	currMode := strings.TrimSpace(string(b))
	if currMode != mode {
		log.Printf("Got %s, expected %s\n", string(b), mode)
		if err := ioutil.WriteFile(p.path+"edge", []byte(mode), 0644); err != nil {
			log.Fatal(err)
		}
	}

	p.event.Events = unix.EPOLLPRI | unix.EPOLLERR
	p.event.Fd = int32(p.valueFile.Fd())

	if epollFd >= 0 {
		err := unix.EpollCtl(epollFd, unix.EPOLL_CTL_ADD, int(p.valueFile.Fd()), &p.event)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Perform initial read
	p.Value()
}

func (p *InputPin) Value() bool {
	var b []byte
	var err error

	if _, err := p.valueFile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	if b, err = ioutil.ReadAll(p.valueFile); err != nil {
		log.Fatal(err)
	}

	ret := (b[0] == '1')
	return ret
}

func (p *InputPin) ChangedState(events []unix.EpollEvent) bool {
	for _, e := range events {
		if e.Fd == p.event.Fd {
			return true
		}
	}

	return false
}

func main() {
	var player RandomSamplePlayer

	player.Program = "/usr/bin/aplay"
	audioDir := "/audio"
	led := MakeOutputPin(18)
	btn := MakeInputPin(24, false, true)
	sw  := MakeInputPin(23, true, true)

	flag.UintVar(&led.pin, "gpio-led", led.pin, "LED control GPIO number")
	flag.UintVar(&btn.pin, "gpio-btn", btn.pin, "Button GPIO pin number")
	flag.UintVar(&sw.pin, "gpio-sw", sw.pin, "Switch GPIO number")
	flag.StringVar(&audioDir, "audio-dir", audioDir, "Audio sample root directory")
	flag.StringVar(&player.Program, "playback", player.Program, "Audio playback program")

	flag.Parse()

	// Load audio
	player.LoadSamples(audioDir)

	// Initialize GPIOs
	epollFd, err := unix.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
	}

	led.Setup()
	btn.Setup(epollFd)
	sw.Setup(epollFd)

	// Play a sample to denote that we're up and running
	led.SetValue(true)
	player.PlayRandomSample()
	led.SetValue(false)

	for {
		events := []unix.EpollEvent{ unix.EpollEvent{}, unix.EpollEvent{} }
		playSample := false

		_, err := unix.EpollWait(epollFd, events, -1)
		if err != nil {
			log.Fatal(err)
		}

		if sw.ChangedState(events) {
			sw.Value()
			player.NextSet()
			playSample = true
		} else if btn.ChangedState(events) {
			playSample = !btn.Value()
		}

		if playSample {
			led.SetValue(true)
			player.PlayRandomSample()
			led.SetValue(false)
		}
	}
}
