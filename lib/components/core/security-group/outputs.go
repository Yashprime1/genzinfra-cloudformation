package securitygroup

import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddOutputsForCoreSecurityGroupStack(template *cloudformation.Template) {
	template.Outputs["DS2SecurityGroupId"] = cloudformation.Output{
		Description: cloudformation.String("DS2 Security Group Id"),
		Value:       cloudformation.GetAtt("DSSecurityGroup", "GroupId"),
		Export: &cloudformation.Export{
			Name: cloudformation.Sub("${AWS::StackName}-DS2SecurityGroupId"),
		},
	}
}
