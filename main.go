package main
import (
	"fmt"
	"log"
	// "io"
	"os"
    "context"
    "github.com/urfave/cli/v3"
)

func main(){
	(&cli.Command{}).Run(context.Background(), os.Args)
}

