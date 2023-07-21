package dsbase

import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func GenerateDsBaseStack(defaults DsBaseDefaults) *cloudformation.Template {
	// Create the DS Base Stack
	BaseTemplate := cloudformation.NewTemplate()
	BaseTemplate.Description = "DS Base Stack"
	AddParametersForDsBaseStack(BaseTemplate,defaults)
	AddResourcesForDsBaseStack(BaseTemplate,defaults)
	AddConditionsForDsBaseStack(BaseTemplate,defaults)
	AddOutputsForDsBaseStack(BaseTemplate)
	return BaseTemplate
}

