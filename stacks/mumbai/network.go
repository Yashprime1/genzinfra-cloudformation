package mumbai

import (
	"fmt"
	networkcore "github.com/Yashprime1/genzinfra-cloudformation/lib/components/core/network"
	utils "github.com/Yashprime1/genzinfra-cloudformation/lib/utilities"
)

func GenerateNetworkStack() {
	// Create the stacks
	fmt.Println("Generating Mumbai Network stack template")

	// Generate the stack template
	var defaults networkcore.CoreNetworkDefaults 
	NetworkTemplate := networkcore.GenerateCoreNetworkStack(defaults)	
	//Write the Json Templates
	NetworkJsonTemplate, err := NetworkTemplate.JSON()
	if err != nil {
		fmt.Println(err)
	} else {
		utils.WriteTemplatesToFile("./templates/mumbai/core","mumbai-network.json", NetworkJsonTemplate)
	}
}
