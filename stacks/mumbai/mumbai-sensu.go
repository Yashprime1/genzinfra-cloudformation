package mumbai

import (
	"fmt"

	sensubase "github.com/Yashprime1/genzinfra-cloudformation/lib/components/sensu/base"
	sensuservice "github.com/Yashprime1/genzinfra-cloudformation/lib/components/sensu/service"
	utils "github.com/Yashprime1/genzinfra-cloudformation/lib/utilities"
)

func GenerateSensuStacks() {
	// Create the stacks
	fmt.Println("Generating Mumbai-Sensu stack templates")

	// Assign the values to the DS stack infra requires

	// Generate the stack templates
	var defaults sensubase.SensuBaseDefaults 
	defaults.NetworkStack = "Mu-Network"
	defaults.SecurityGroupStack = "Mu-SecurityGroup"

	BaseTemplate := sensubase.GenerateSensuBaseStack(defaults)
	ServiceTemplate := sensuservice.GenerateSensuServiceStack()


	//Write the Json Templates
	BaseJsonTemplate, err := BaseTemplate.JSON()
	if err != nil {
		fmt.Println(err)
	} else {
		utils.WriteTemplatesToFile("./templates/mumbai/ds","mumbai-sensu.json", BaseJsonTemplate)
	}
	ServiceJsonTemplate, err := ServiceTemplate.JSON()
	if err != nil {
		fmt.Println(err)
	} else {
		utils.WriteTemplatesToFile("./templates/mumbai/ds","mumbai-sensu-service.json", ServiceJsonTemplate)
	}
}
