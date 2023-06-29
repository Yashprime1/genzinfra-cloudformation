package dsservice
import (
	"github.com/awslabs/goformation/v7/cloudformation/ec2"
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddResourcesForDsServiceStack(template *cloudformation.Template){
	template.Resources["DSEc2Instance"] =  &ec2.Instance{
		InstanceType: cloudformation.String(cloudformation.Ref("InstanceType")),
	}
}

