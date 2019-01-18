```
  ___ ____  ___/ /__ ____  _______    
 / _ '/ _ \/ _  / _ '/ _ \/ __/ -_)   
 \_, /\___/\_,_/\_,_/_//_/\__/\__/    
/___/
```


# godance - A password spraying SMB bruteforcer

SMB password sprayer

```
$ godance -h 192.168.75.173 -u users.txt -w passwords.txt -d WORKGROUP -t 200   
 
  ___ ____  ___/ /__ ____  _______    
 / _ '/ _ \/ _  / _ '/ _ \/ __/ -_)   
 \_, /\___/\_,_/\_,_/_//_/\__/\__/    
/___/

-----------------------------------------------------
 [*] Number of usernames: 4242
 [*] Number of passwords: 4
 [*] Test cases: 16968
 [*] Number of threads: 200
-----------------------------------------------------
 [*] In hacker voice *I'm in* // Username: pystyy // Password: vetaa

```

## Usage


```
Usage of godance:
  -d string
    	Domain
  -h string
    	Target host
  -p int
    	Target port (default 445)
  -t int
    	Number of threads (default 20)
  -u string
    	User wordlist
  -v	Debug
  -w string
    	Password list
```

## Installation

 - If you have go compiler installed: `go get github.com/joohoi/godance`

