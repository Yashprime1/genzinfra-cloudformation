package sensubase
import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddParametersForSensuBaseStack(template *cloudformation.Template,defaults SensuBaseDefaults ){
	// Add the parameters
	// Add the parameters
	template.Parameters["SensuAmiId"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The AMI ID for the Sensu"),
	}
	
	template.Parameters["SensuInstanceType"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The instance type for the Sensu"),
	}
	
	template.Parameters["SensuMaxSize"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The maximum size for the Sensu"),
	}

	template.Parameters["SensuMinSize"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The minimum size for the Sensu"),
	}

	template.Parameters["SensuDesiredSize"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The desired capacity for the Sensu"),
	}
}

