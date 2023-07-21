package dsbase
import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddParametersForDsBaseStack(template *cloudformation.Template,defaults DsBaseDefaults ){
	// Add the parameters
	// Add the parameters
	// template.Parameters["DsAmiId"] = cloudformation.Parameter{
	// 	Type: "String",
	// 	Description: cloudformation.String("The AMI ID for the DS"),
	// }
	
	// template.Parameters["DsInstanceType"] = cloudformation.Parameter{
	// 	Type: "String",
	// 	Description: cloudformation.String("The instance type for the DS"),
	// }
	
	// template.Parameters["DsMaxSize"] = cloudformation.Parameter{
	// 	Type: "String",
	// 	Description: cloudformation.String("The maximum size for the DS"),
	// }

	// template.Parameters["DsMinSize"] = cloudformation.Parameter{
	// 	Type: "String",
	// 	Description: cloudformation.String("The minimum size for the DS"),
	// }

	// template.Parameters["DsDesiredSize"] = cloudformation.Parameter{
	// 	Type: "String",
	// 	Description: cloudformation.String("The desired capacity for the DS"),
	// }
}

