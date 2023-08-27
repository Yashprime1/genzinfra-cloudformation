package sensuservice
import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddConditionsForSensuServiceStack(template *cloudformation.Template){
	// Create the DS Service Stack
	ServiceTemplate := cloudformation.NewTemplate()
	ServiceTemplate.Description = "Sensu Service Stack"
}
