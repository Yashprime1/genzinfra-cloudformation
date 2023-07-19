package ireland

import (
	"fmt"

	dsbase "github.com/Yashprime1/genzinfra-cloudformation/lib/components/ds/base"
	dsservice "github.com/Yashprime1/genzinfra-cloudformation/lib/components/ds/service"
	utils "github.com/Yashprime1/genzinfra-cloudformation/lib/utilities"
)

func GenerateDsStacks() {
	// Create the stacks
	fmt.Println("Generating Ireland-DS stack templates")

	// Assign the values to the DS stack infra requires

	// Generate the stack templates
	var defaults dsbase.DsBaseDefaults 

	BaseTemplate := dsbase.GenerateDsBaseStack(defaults)
	ServiceTemplate := dsservice.GenerateDsServiceStack()

	//Write the Json Templates
	BaseJsonTemplate, err := BaseTemplate.JSON()
	if err != nil {
		fmt.Println(err)
	} else {
		utils.WriteTemplatesToFile("./templates/ireland/ds","ireland-ds.json", BaseJsonTemplate)
	}
	ServiceJsonTemplate, err := ServiceTemplate.JSON()
	if err != nil {
		fmt.Println(err)
	} else {
		utils.WriteTemplatesToFile("./templates/ireland/ds","ireland-ds-service.json", ServiceJsonTemplate)
	}
}
