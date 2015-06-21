package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetSupervisedProcesses(t *testing.T) {
	expectedCommand := "supervisorctl"
	expectedArgs := []string{"status"}
	run := func(command string, args ...string) ([]byte, error) {
		assert.Equal(t, expectedCommand, command)
		assert.Equal(t, expectedArgs, args)
		return []byte("One STOPPED \n Two RUNNING"), nil
	}
	out, err := GetSupervisedProcesses(run)
	assert.Nil(t, err)
	assert.Equal(t, []string{"One", "Two"}, out)
}

func TestParseStatusReturn(t *testing.T) {
	output := `long_script                      STOPPED    Jun 04 02:52 PM
another            STOPPED

	`

	out := ParseSupervisorOutput(output)

	// TODO parse return into slice
	assert.Equal(t, []string{"long_script", "another"}, out)
}

func TestRun(t *testing.T) {
	output, err := Run("echo", "test")
	assert.Nil(t, err)
	assert.Equal(t, "test\n", string(output))
}

func TestTick(t *testing.T) {
	ticker := time.NewTicker(time.Millisecond * 1)
	processes := []string{"one", "two"}
	exit := make(chan int)
	shouldDo := func(chance int) bool {
		return true
	}

	doCalled := 0
	do := func(processes []string, run func(command string, args ...string) ([]byte, error)) {
		doCalled++
	}

	go Tick(ticker, shouldDo, do, processes, exit)
	time.Sleep(time.Millisecond * 2)
	ticker.Stop()
	assert.Equal(t, 2, doCalled)
}

func TestShouldDo(t *testing.T) {
	chance = 50
	seed = 42
	assert.Equal(t, true, ShouldDo(100))
	assert.Equal(t, false, ShouldDo(0))
	// With a seed of 42 the 3rd result is 47
	assert.Equal(t, true, ShouldDo(50))
}

func TestDoCallsSupervisor(t *testing.T) {
	expectedCommand := "supervisorctl"
	expectedArgs := []string{"restart", "TESTPROC"}
	runCalled := false

	run := func(command string, args ...string) ([]byte, error) {
		runCalled = true
		assert.Equal(t, expectedCommand, command)
		assert.Equal(t, expectedArgs, args)
		return []byte("ok"), nil
	}

	Do([]string{"TESTPROC"}, run)

	assert.Equal(t, true, runCalled)
}

func TestDoPicksRandomProcess(t *testing.T) {
	seed = 42
	expectedCommand := "supervisorctl"
	expectedArgs := []string{"restart", "TWO"}
	runCalled := false

	run := func(command string, args ...string) ([]byte, error) {
		runCalled = true
		assert.Equal(t, expectedCommand, command)
		assert.Equal(t, expectedArgs, args)
		return []byte("ok"), nil
	}

	Do([]string{"ONE", "TWO", "THREE"}, run)

	assert.Equal(t, true, runCalled)
}
