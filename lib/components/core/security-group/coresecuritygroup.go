package securitygroup

import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func GenerateCoreSecurityGroupStack(defaults CoreSecurityGroupDefaults) *cloudformation.Template {
	// Create the DS Service Stack
	ServiceTemplate := cloudformation.NewTemplate()
	ServiceTemplate.Description = defaults.SecurityGroupDescription
	AddParametersForCoreSecurityGroupStack(ServiceTemplate)
	AddResourcesForCoreSecurityGroupStack(ServiceTemplate, defaults)
	AddConditionsForCoreSecurityGroupStack(ServiceTemplate)
	AddOutputsForCoreSecurityGroupStack(ServiceTemplate)
	return ServiceTemplate
}
