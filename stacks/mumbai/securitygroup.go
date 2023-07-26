package mumbai

import (
	"fmt"
	securitygroup "github.com/Yashprime1/genzinfra-cloudformation/lib/components/core/security-group"
	utils "github.com/Yashprime1/genzinfra-cloudformation/lib/utilities"
)

func GenerateSecurityGroupStack() {
	// Create the stacks
	fmt.Println("Generating Mumbai Security Group stack template")

	// Generate the stack template
	var defaults  = securitygroup.CoreSecurityGroupDefaults{"Mumbai-Network","Mumbai-DS-SecurityGroup","DS Security Group",[]int{80,9090,22}}
	SecurityGroupTemplate := securitygroup.GenerateCoreSecurityGroupStack(defaults)	
	//Write the Json Templates
	SecurityGroupJsonTemplate, err := SecurityGroupTemplate.JSON()
	if err != nil {
		fmt.Println(err)
	} else {
		utils.WriteTemplatesToFile("./templates/mumbai/core","mumbai-securitygroup.json", SecurityGroupJsonTemplate)
	}
}
