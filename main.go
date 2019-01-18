package main

import (
	"flag"
	"fmt"
	"os"
)

type CLIConfig struct {
	host     string
	port     int
	debug    bool
	domain   string
	userFile string
	pwdFile  string
	threads  int
}

type Config struct {
	host    string
	port    int
	debug   bool
	domain  string
	threads int
	users   *WordlistInput
	passwds *WordlistInput
}

const (
	BANNER = ` 
  ___ ____  ___/ /__ ____  _______    
 / _ '/ _ \/ _  / _ '/ _ \/ __/ -_)   
 \_, /\___/\_,_/\_,_/_//_/\__/\__/    
/___/
`
	SEP = "-----------------------------------------------------"
)

func createConfig(cliconf *CLIConfig) (Config, error) {
	var conf Config
	var err error
	if cliconf.host == "" {
		return conf, fmt.Errorf("Host not defined")
	}
	conf.host = cliconf.host
	if cliconf.domain == "" {
		return conf, fmt.Errorf("Domain not defined")
	}
	conf.domain = cliconf.domain
	if cliconf.userFile == "" {
		return conf, fmt.Errorf("Userfile not defined")
	}
	conf.users, err = NewWordlistInput(cliconf.userFile)
	if err != nil {
		return conf, fmt.Errorf("Could not read user file: %s", err)
	}
	if cliconf.pwdFile == "" {
		return conf, fmt.Errorf("Passwordfile not defined")
	}
	conf.passwds, err = NewWordlistInput(cliconf.pwdFile)
	if err != nil {
		return conf, fmt.Errorf("Could not read password file: %s", err)
	}
	conf.threads = cliconf.threads
	conf.port = cliconf.port
	conf.debug = cliconf.debug
	return conf, nil
}

func main() {
	var cliconf CLIConfig
	flag.StringVar(&cliconf.host, "h", "", "Target host")
	flag.IntVar(&cliconf.port, "p", 445, "Target port")
	flag.IntVar(&cliconf.threads, "t", 10, "Number of threads")
	flag.StringVar(&cliconf.userFile, "u", "", "User wordlist")
	flag.StringVar(&cliconf.pwdFile, "w", "", "Password list")
	flag.StringVar(&cliconf.domain, "d", "WORKGROUP", "Domain")
	flag.BoolVar(&cliconf.debug, "v", false, "Debug")
	flag.Parse()
	conf, err := createConfig(&cliconf)
	if err != nil {
		fmt.Printf("  [!] Error: %s\n\n", err)
		flag.Usage()
		os.Exit(1)
	}
	runner := NewRunner(&conf)
	runner.Start()
}
