package securitygroup
import (
	"github.com/awslabs/goformation/v7/cloudformation/ec2"
	"github.com/awslabs/goformation/v7/cloudformation/tags"
	"github.com/awslabs/goformation/v7/cloudformation"

)

func AddResourcesForCoreSecurityGroupStack(template *cloudformation.Template,defaults CoreSecurityGroupDefaults){
	template.Resources["DSSecurityGroup"] =  &ec2.SecurityGroup{
		GroupDescription: "DS Security Group",
		SecurityGroupIngress: []ec2.SecurityGroup_Ingress{
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(-1),
				ToPort: cloudformation.Int(-1),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
		},
		VpcId: cloudformation.String(cloudformation.ImportValue(defaults.NetworkStack+"-AppVPCId")),
		Tags: []tags.Tag{
			{
				Key: "Name",
				Value: "DS Security Groups",
			},
		},
	}
}
