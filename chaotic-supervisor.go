package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

var chance = 10

var seed = time.Now().UnixNano()

func main() {
	log.Println("Chaotic Supervisor")
	rand.Seed(seed)
	exit := make(chan int)
	out, err := GetSupervisedProcesses(Run)
	if err != nil {
		log.Fatal("Error getting supervised processes: ", err)
	}

	fmt.Println(out)
	go Tick(time.NewTicker(time.Second), ShouldDo, Do, out, exit)
	<-exit
	close(exit)
}

func Run(command string, args ...string) ([]byte, error) {
	out, err := exec.Command(command, args...).CombinedOutput()
	return out, err
}

func ShouldDo(percentageChance int) bool {
	if rand.Intn(100) <= percentageChance {
		return true
	} else {
		return false
	}
}

// Pick an action to perform on the processes
// First pass will support restart only
func Do(processes []string, run func(command string, args ...string) ([]byte, error)) {
	process := processes[rand.Intn(len(processes))]
	fmt.Printf("%s restarting %s\n", time.Now(), process)
	out, err := run("supervisorctl", "restart", process)
	if err != nil {
		fmt.Println("ERROR ", err)
	}
	fmt.Println(string(out))
}

func Tick(
	ticker *time.Ticker,
	shouldDo func(chance int) bool,
	do func(processes []string, run func(command string, args ...string) ([]byte, error)),
	processes []string,
	exit chan<- int) {

	for _ = range ticker.C {
		if shouldDo(chance) {
			do(processes, Run)
		}
	}

	exit <- 0
}

func GetSupervisedProcesses(run func(command string, args ...string) ([]byte, error)) ([]string, error) {
	out, err := run("supervisorctl", []string{"status"}...)
	return ParseSupervisorOutput(string(out)), err
}

func ParseSupervisorOutput(output string) []string {
	lines := strings.Split(output, "\n")
	processes := []string{}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			name := fields[0]
			processes = append(processes, name)
		}
	}

	return processes
}
