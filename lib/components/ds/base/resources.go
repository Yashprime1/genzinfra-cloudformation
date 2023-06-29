package dsbase
import (
	"github.com/awslabs/goformation/v7/cloudformation/ec2"
	"github.com/awslabs/goformation/v7/cloudformation"
)

func AddResourcesForDsBaseStack(template *cloudformation.Template,defaults DsBaseDefaults){
	template.Resources["DSEc2Instance"] =  &ec2.Instance{
		InstanceType: cloudformation.String(cloudformation.Ref("InstanceType")),
	}
}

