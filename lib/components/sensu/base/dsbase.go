package sensubase

import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func GenerateSensuBaseStack(defaults SensuBaseDefaults) *cloudformation.Template {
	// Create the Sensu Base Stack
	BaseTemplate := cloudformation.NewTemplate()
	BaseTemplate.Description = "Sensu's Base Stack"
	AddParametersForSensuBaseStack(BaseTemplate,defaults)
	AddResourcesForSensuBaseStack(BaseTemplate,defaults)
	AddConditionsForSensuBaseStack(BaseTemplate,defaults)
	AddOutputsForSensuBaseStack(BaseTemplate)
	return BaseTemplate
}

