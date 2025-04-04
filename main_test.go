package main

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

func TestFoo(t *testing.T) {
	name = t.Name()
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	cmd := exec.Command("./windows-process-signal"+suffix, "foo", "bar")
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = getSysProcAttr()

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// We need to wait for the process to start
	for i := range 5 {
		time.Sleep(time.Second)
		t.Logf("=================== %d", i)
	}
	t.Log("==================== Sending stop signal")
	if err := stopCmd(cmd); err != nil {
		t.Fatal("Stop error:", err)
	}
	t.Log("============================== Stop signal sent")

	time.Sleep(2 * time.Second)
	if err := cmd.Wait(); err != nil {
		t.Log("============================== Error", err)
		t.Fatal("Wait failed", err)
	}
	t.Log("============================== Process has stopped")
}
