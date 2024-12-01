package ireland

import (
	"fmt"
	securitygroup "github.com/Yashprime1/genzinfra-cloudformation/lib/components/core/security-group"
	utils "github.com/Yashprime1/genzinfra-cloudformation/lib/utilities"
)

func GenerateSecurityGroupStack() {
	// Create the stacks
	fmt.Println("Generating Ireland Security Group stack template")

	// Generate the stack template
	SecurityGroupTemplate := securitygroup.GenerateCoreSecurityGroupStack("Eu-Network")	
	//Write the Json Templates
	SecurityGroupJsonTemplate, err := SecurityGroupTemplate.JSON()
	if err != nil {
		fmt.Println(err)
	} else {
		utils.WriteTemplatesToFile("./templates/ireland/core","ireland-securitygroup.json", SecurityGroupJsonTemplate)
	}
}
