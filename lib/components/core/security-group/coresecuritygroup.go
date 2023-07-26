package securitygroup

import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func GenerateCoreSecurityGroupStack(NetworkStackName string) *cloudformation.Template {
	// Create the DS Service Stack
	ServiceTemplate := cloudformation.NewTemplate()
	ServiceTemplate.Description = "Core Security Group Stack"
	defaults := CoreSecurityGroupDefaults{}
	defaults.NetworkStack = NetworkStackName
	AddParametersForCoreSecurityGroupStack(ServiceTemplate)
	AddResourcesForCoreSecurityGroupStack(ServiceTemplate, defaults)
	AddConditionsForCoreSecurityGroupStack(ServiceTemplate)
	AddOutputsForCoreSecurityGroupStack(ServiceTemplate)
	return ServiceTemplate
}
