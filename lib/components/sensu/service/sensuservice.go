package sensuservice
import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func GenerateSensuServiceStack() *cloudformation.Template{
	// Create the DS Service Stack
	ServiceTemplate := cloudformation.NewTemplate()
	ServiceTemplate.Description = "Sensu Service Stack"
	AddParametersForSensuServiceStack(ServiceTemplate)
	AddResourcesForSensuServiceStack(ServiceTemplate)
	AddConditionsForSensuServiceStack(ServiceTemplate)
	return ServiceTemplate
}

