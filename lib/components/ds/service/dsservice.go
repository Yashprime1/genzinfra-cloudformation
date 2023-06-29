package dsservice
import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func GenerateDsServiceStack() *cloudformation.Template{
	// Create the DS Service Stack
	ServiceTemplate := cloudformation.NewTemplate()
	ServiceTemplate.Description = "DS Service Stack"
	AddParametersForDsServiceStack(ServiceTemplate)
	AddResourcesForDsServiceStack(ServiceTemplate)
	AddConditionsForDsServiceStack(ServiceTemplate)
	return ServiceTemplate
}

