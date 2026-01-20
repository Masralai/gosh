package main

import (
	"bufio"
	// "context"
	"fmt"
	// "log"
	"os"
    // "strings"

	// "github.com/urfave/cli/v3"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan(){
        if(scanner.Text()=="exit"){
            os.Exit(0)
        }else{
            //fmt.Println(scanner.Text())
            continue
        }
        
    }
    if err := scanner.Err(); err!= nil{
        fmt.Fprintln(os.Stderr, "shouldnt  see an error scanning a string")
    }
}