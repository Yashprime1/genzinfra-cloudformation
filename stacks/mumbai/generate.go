package mumbai
import(
	"fmt"
)

func GenerateStacks(){
	// Create the stacks
	fmt.Println("Generating all Mumbai stack templates")
	GenerateDsStacks()
	GenerateSensuStacks()
	GenerateNetworkStack()	
	GenerateSecurityGroupStack()
}