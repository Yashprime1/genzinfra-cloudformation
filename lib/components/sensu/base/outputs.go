package sensubase
// FOr Sensubase
import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddOutputsForSensuBaseStack(template *cloudformation.Template) {
	template.Outputs["SensuElbTargetGroupArn"] = cloudformation.Output{
		Description: cloudformation.String("Sensu ELB Target Group Arn"),
		Value:       cloudformation.Ref("SensuElbTargetGroup"),
		Export: &cloudformation.Export{
			Name: cloudformation.Sub("${AWS::StackName}-Service-SensuElbTargetGroupArn"),
		},
	}

	template.Outputs["SensuBackendElbTargetGroupArn"] = cloudformation.Output{
		Description: cloudformation.String("Sensu Backend ELB Target Group Arn"),
		Value:       cloudformation.Ref("SensuBackendElbTargetGroup"),
		Export: &cloudformation.Export{
			Name: cloudformation.Sub("${AWS::StackName}-Service-SensuBackendElbTargetGroupArn"),
		},
	}

	template.Outputs["SensuApiElbTargetGroupArn"] = cloudformation.Output{
		Description: cloudformation.String("Sensu Api ELB Target Group Arn"),
		Value:       cloudformation.Ref("SensuApiElbTargetGroup"),
		Export: &cloudformation.Export{
			Name: cloudformation.Sub("${AWS::StackName}-Service-SensuApiElbTargetGroupArn"),
		},
	}
}
