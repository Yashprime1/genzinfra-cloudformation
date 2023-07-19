package ireland

import "fmt"

func GenerateStacks(){
	// Create the stacks
	fmt.Println("Generating all Ireland stack templates")
	GenerateNetworkStack()	
	GenerateDsStacks() 

}