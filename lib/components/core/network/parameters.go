package network
import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddParametersForCoreNetworkStack(template *cloudformation.Template,defaults CoreNetworkDefaults ){
	// Add the parameters
	template.Parameters["AppPublicSubnet1Cidr"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The CIDR block for the public subnet1 for App"),
	}
	template.Parameters["AppPublicSubnet2Cidr"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The CIDR block for the public subnet2 for App"),
	}
	template.Parameters["AppPrivateSubnet1Cidr"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The CIDR block for the private subnet 1 for App"),
	}
	template.Parameters["AppPrivateSubnet2Cidr"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The CIDR block for the private subnet 2 for App"),
	}	
		
	template.Parameters["AppVpcCidr"] = cloudformation.Parameter{
		Type: "String",
		Description: cloudformation.String("The CIDR block for the VPC for App"),
	}
	
}

