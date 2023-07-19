package network

import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddOutputsForCoreNetworkStack(template *cloudformation.Template) {
	template.Outputs["AppVPCId"] = cloudformation.Output{
		Description: cloudformation.String("App Vpc Id"),
		Value:       cloudformation.Ref("AppVPC"),
		Export: &cloudformation.Export{
			Name: cloudformation.Sub("${AWS::StackName}-AppVPCId"),
		},
	}
	template.Outputs["AppPrivateSubnet1Id"] = cloudformation.Output{
		Description: cloudformation.String("App Private Subnet 1 Id"),
		Value:       cloudformation.Ref("AppPrivateSubnet1"),
		Export: &cloudformation.Export{
			Name: cloudformation.Sub("${AWS::StackName}-AppPrivateSubnet1Id"),
		},
	}
	template.Outputs["AppPrivateSubnet2Id"] = cloudformation.Output{
		Description: cloudformation.String("App Private Subnet 2 Id"),
		Value:       cloudformation.Ref("AppPrivateSubnet2"),
		Export: &cloudformation.Export{
			Name: cloudformation.Sub("${AWS::StackName}-AppPrivateSubnet2Id"),
		},
	}
	template.Outputs["AppPublicSubne1tId"] = cloudformation.Output{
		Description: cloudformation.String("App Public Subnet 1 Id"),
		Value:       cloudformation.Ref("AppPublicSubnet1"),
		Export: &cloudformation.Export{
			Name: cloudformation.Sub("${AWS::StackName}-AppPublicSubnet1Id"),
		},
	}
	template.Outputs["AppPublicSubnet2Id"] = cloudformation.Output{
		Description: cloudformation.String("App Public Subnet 2 Id"),
		Value:       cloudformation.Ref("AppPublicSubnet2"),
		Export: &cloudformation.Export{
			Name: cloudformation.Sub("${AWS::StackName}-AppPublicSubnet2Id"),
		},
	}

}
