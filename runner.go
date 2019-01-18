package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/stacktitan/smb/smb"
)

type Runner struct {
	counter int
	running bool
	conf    *Config
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
			}()
		}
		// Reset the pwd inputlist position
		r.conf.users.position = -1
	}
	wg.Wait()
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
	defer session.Close()
	if err != nil {
		errstr := fmt.Sprintf("%s", err)
		if !strings.Contains(errstr, "Logon failed") {
			fmt.Printf(" [!] Error: %s\n", err)
		}
	}
	defer session.Close()

	if session.IsAuthenticated {
		fmt.Printf(" [*] In hacker voice *I'm in* // Username: %s // Password: %s\n", username, password)
	}
}

func (r *Runner) Stop() {
	r.running = false
}
