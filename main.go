package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var name string

func main() {
	if len(os.Args) == 1 {
		log.Println("Usage: ./windows-process-signal [name] [child name]")
		os.Exit(1)
	}

	name = os.Args[1]
	log.Printf("[%s] Starting, PID %d", name, os.Getpid())
	defer log.Printf("[%s] Done", name)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	var cmd *exec.Cmd
	if len(os.Args) >= 3 {
		suffix := ""
		if runtime.GOOS == "windows" {
			suffix = ".exe"
		}

		go func() {
			for s := range signals {
				log.Printf("[%s] Got signal: %s and will KEEP RUNNING the child", name, s)
			}
		}()

		childName := os.Args[2]
		cmd = exec.Command("./windows-process-signal"+suffix, childName)
		cmd.Stderr = os.Stdout
		cmd.Stdout = os.Stdout
		cmd.SysProcAttr = getSysProcAttr()

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		for i := range 10 {
			time.Sleep(time.Second)
			log.Printf("[%s] Waiting %d", name, i)
		}

		if err := stopCmd(cmd); err != nil {
			log.Fatalf("[%s] cannot stop child: %s", name, err)
		}

		log.Printf("[%s] Calling wait on children", name)

		if err := cmd.Wait(); err != nil {
			log.Fatalf("[%s] error waiting child to return: %s", name, err)
		}

		log.Printf("[%s] Done!", name)
		return
	}

	log.Printf("[%s] Wait for signal", name)
	s := <-signals
	log.Printf("[%s] Got signal: %s", name, s)

	// if cmd != nil {
	// 	log.Printf("[%s] Waiting 2s", name)
	// 	time.Sleep(2 * time.Second)

	// 	log.Printf("[%s] sending signal to children", name)
	// 	if err := stopCmd(cmd); err != nil {
	// 		log.Fatalf("[%s] cannot stop child: %s", name, err)
	// 	}
	// 	if err := cmd.Wait(); err != nil {
	// 		log.Fatalf("[%s] error waiting child to return: %s", name, err)
	// 	}
	// }
}
