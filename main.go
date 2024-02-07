package main
import (
	"fmt"
	"github.com/Yashprime1/genzinfra-cloudformation/stacks/ireland"
	"github.com/Yashprime1/genzinfra-cloudformation/stacks/mumbai"
)


func main() {
	fmt.Println("Generating all stackstemplates")
	ireland.GenerateStacks()
	mumbai.GenerateStacks()
}
