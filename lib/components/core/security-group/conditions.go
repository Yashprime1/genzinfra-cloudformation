package securitygroup
import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddConditionsForCoreSecurityGroupStack(template *cloudformation.Template){
	// Create the DS Service Stack
	ServiceTemplate := cloudformation.NewTemplate()
	ServiceTemplate.Description = "DS Service Stack"
}
