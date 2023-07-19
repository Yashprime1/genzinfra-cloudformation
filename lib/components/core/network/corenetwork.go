package network

import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func GenerateCoreNetworkStack(defaults CoreNetworkDefaults) *cloudformation.Template {
	// Create the DS Base Stack
	BaseTemplate := cloudformation.NewTemplate()
	BaseTemplate.Description = "Core Network Stack"
	AddParametersForCoreNetworkStack(BaseTemplate,defaults)
	AddResourcesForCoreNetworkStack(BaseTemplate,defaults)
	AddConditionsForCoreNetworkStack(BaseTemplate,defaults)
	AddOutputsForCoreNetworkStack(BaseTemplate)
	return BaseTemplate
}

