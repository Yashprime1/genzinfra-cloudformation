package dsbase
import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddParametersForDsBaseStack(template *cloudformation.Template,defaults DsBaseDefaults ){
	// Add the parameters
	template.Parameters["InstanceType"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The instance type for EC2 instances required for the DS service"),
		Default: cloudformation.String(defaults.InstanceType),
	}
		
}

