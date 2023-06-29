package dsservice
import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddParametersForDsServiceStack(template *cloudformation.Template ){
	// Add the parameters
	template.Parameters["InstanceType"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The instance type for EC2 instances required for the DS service"),
	}
		
}

