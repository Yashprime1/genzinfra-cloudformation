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
				FromPort: cloudformation.Int(80),
				ToPort: cloudformation.Int(80),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(22),
				ToPort: cloudformation.Int(22),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(23),
				ToPort: cloudformation.Int(23),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(9090),
				ToPort: cloudformation.Int(9090),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(17001),
				ToPort: cloudformation.Int(17001),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(27017),
				ToPort: cloudformation.Int(27017),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(9216),
				ToPort: cloudformation.Int(9216),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(3000),
				ToPort: cloudformation.Int(3000),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(8080),
				ToPort: cloudformation.Int(8080),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(8081),
				ToPort: cloudformation.Int(8081),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(2379),
				ToPort: cloudformation.Int(2379),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(2380),
				ToPort: cloudformation.Int(2380),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(80),
				ToPort: cloudformation.Int(80),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(8000),
				ToPort: cloudformation.Int(8000),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},	
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(9997),
				ToPort: cloudformation.Int(9997),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(8088),
				ToPort: cloudformation.Int(8088),
				CidrIp: cloudformation.String("0.0.0.0/0"),
			},
			{
				IpProtocol: "tcp",
				FromPort: cloudformation.Int(443),
				ToPort: cloudformation.Int(443),
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
