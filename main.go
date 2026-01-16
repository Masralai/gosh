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
	cmd:= &cli.Command{
		Name: "boom",
		Usage: "make an explosive entrance",
		Action: func(context.Context, *cli.Command)error{
			fmt.Println("boom! I say!")
			return nil
		},
	}
	
	if err:= cmd.Run(context.Background(),os.Args);err!=nil{
		log.Fatal(err)
	}
}

