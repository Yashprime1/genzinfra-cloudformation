package dsbase
// FOr dsbase
import (
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddOutputsForDsBaseStack(template *cloudformation.Template) {
	template.Outputs["DsElbTargetGroupArn"] = cloudformation.Output{
		Description: cloudformation.String("DS ELB Target Group Arn"),
		Value:       cloudformation.Ref("DsElbTargetGroup"),
		Export: &cloudformation.Export{
			Name: cloudformation.Sub("${AWS::StackName}-Service-DsElbTargetGroupArn"),
		},
	}
}
