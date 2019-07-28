package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	rate     = flag.Uint64("rate", 100, "Measurments per second")
	duration = flag.Uint64("duration", 10, "Measurement duration in seconds")
	dev      = flag.String("device", "eth0", "Network device")
)

func main() {
	flag.Parse()

	m := newMeasurement(*rate, time.Duration(*duration)*time.Second)
	err := m.run()
	if err != nil {
		log.Errorf("Failed to run measurement: %v", err)
		os.Exit(1)
	}

	err = m.stats(*dev)
	if err != nil {
		log.Errorf("Unable to print stats: %v", err)
		os.Exit(1)
	}
}

type ifCounter struct {
	name      string
	time      uint64
	txBytes   uint64
	rxBytes   uint64
	txPackets uint64
	rxPackets uint64
}

type measurement struct {
	rate     uint64
	duration time.Duration
	results  map[string][]ifCounter
	min      map[string]ifCounter
	max      map[string]ifCounter
	avg      map[string]ifCounter
	median   map[string]ifCounter
}

func newMeasurement(rate uint64, duration time.Duration) *measurement {
	return &measurement{
		rate:     rate,
		duration: duration,
		results:  make(map[string][]ifCounter),
		min:      make(map[string]ifCounter),
		max:      make(map[string]ifCounter),
		avg:      make(map[string]ifCounter),
		median:   make(map[string]ifCounter),
	}
}

func (m *measurement) run() error {
	count := uint64(time.Duration(m.rate) * (m.duration / time.Second))
	t := time.NewTicker(time.Second / time.Duration(m.rate))
	for i := uint64(0); i < count; i++ {
		counters, err := m.getInterfaceCounters()
		if err != nil {
			return err
		}

		for i := range counters {
			if _, ok := m.results[counters[i].name]; !ok {
				m.results[counters[i].name] = make([]ifCounter, 0, count)
			}
			m.results[counters[i].name] = append(m.results[counters[i].name], counters[i])
		}

		<-t.C
	}

	return nil
}

func (m *measurement) stats(dev string) error {
	devResult, ok := m.results[dev]
	if !ok {
		return fmt.Errorf("Unable to print stats: Device %s not found", dev)
	}

	fmt.Printf("Time (ms)\tTX Mbits/s\tTX packets/s\tRX Mbits/s\tRX packets/s\n")

	prev := ifCounter{}
	for i, r := range devResult {
		if i != 0 {
			fmt.Printf("%d\t%d\t%d\t%d\t%d\n", (1000/m.rate)*uint64(i), (r.txBytes-prev.txBytes)*8*m.rate/1000000, (r.txPackets-prev.txPackets)*m.rate, (r.rxBytes-prev.rxBytes)*8*m.rate/1000000, (r.rxPackets-prev.rxPackets)*m.rate)
		}

		prev = r
	}

	return nil
}

func (m *measurement) getInterfaceCounters() ([]ifCounter, error) {
	c, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, err
	}

	space := regexp.MustCompile(`\s+`)

	ret := make([]ifCounter, 0, 20)
	for _, l := range strings.Split(string(c), "\n") {
		if !strings.Contains(l, ":") {
			continue
		}

		fields := strings.Split(strings.TrimSpace(space.ReplaceAllString(l, " ")), " ")
		fields[0] = strings.Replace(fields[0], ":", "", 1)

		x := ifCounter{
			name: fields[0],
		}

		rxBytes, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}

		rxPackets, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}

		txBytes, err := strconv.Atoi(fields[9])
		if err != nil {
			return nil, err
		}

		txPackets, err := strconv.Atoi(fields[10])
		if err != nil {
			return nil, err
		}

		x.rxBytes = uint64(rxBytes)
		x.txBytes = uint64(txBytes)
		x.rxPackets = uint64(rxPackets)
		x.txPackets = uint64(txPackets)

		ret = append(ret, x)
	}

	return ret, nil
}
