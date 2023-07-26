package securitygroup
import (
	"github.com/awslabs/goformation/v7/cloudformation/ec2"
	"github.com/awslabs/goformation/v7/cloudformation/tags"
	"github.com/awslabs/goformation/v7/cloudformation"

)

func AddResourcesForCoreSecurityGroupStack(template *cloudformation.Template,defaults CoreSecurityGroupDefaults){

	ingress_rules := []ec2.SecurityGroup_Ingress{}
	for _,port := range defaults.SecurityGroupIngressPorts{
		ingress_rules = append(ingress_rules,ec2.SecurityGroup_Ingress{
			IpProtocol: "tcp",
			FromPort: cloudformation.Int(port),
			ToPort: cloudformation.Int(port),
			CidrIp: cloudformation.String("0.0.0.0/0"),
		})
	}
	template.Resources[defaults.SecurityGroupStack] =  &ec2.SecurityGroup{
		GroupDescription: defaults.SecurityGroupDescription,
		SecurityGroupIngress: ingress_rules,
		VpcId: cloudformation.String(cloudformation.ImportValue(defaults.NetworkStack+"-AppVPCId")),
		Tags: []tags.Tag{
			{
				Key: "Name",
				Value: defaults.SecurityGroupStack,
			},
		},
	}
}
