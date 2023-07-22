package network

import (
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ec2"
	"github.com/awslabs/goformation/v7/cloudformation/tags"
)

func AddResourcesForCoreNetworkStack(template *cloudformation.Template, defaults CoreNetworkDefaults) {
	template.Resources["AppVPC"] = &ec2.VPC{
		CidrBlock:          cloudformation.String(cloudformation.Ref("AppVpcCidr")),
		EnableDnsHostnames: cloudformation.Bool(true),
		EnableDnsSupport:   cloudformation.Bool(true),
		Tags: []tags.Tag{
			{
				Key:  "Name",
				Value: "AppVPC",
			},
		},
	}

	template.Resources["AppPublicSubnet1"] = &ec2.Subnet{
		AvailabilityZone: cloudformation.SelectPtr(0, cloudformation.GetAZsPtr(cloudformation.Ref("AWS::Region"))),
		CidrBlock:        cloudformation.String(cloudformation.Ref("AppPublicSubnet1Cidr")),
		VpcId:            cloudformation.Ref("AppVPC"),
		Tags: []tags.Tag{
			{
				Key:  "Name",
				Value: "AppPublicSubnet1",
			},
		},
	}

	template.Resources["AppPublicSubnet2"] = &ec2.Subnet{
		AvailabilityZone: cloudformation.SelectPtr(1, cloudformation.GetAZsPtr(cloudformation.Ref("AWS::Region"))),
		CidrBlock:        cloudformation.String(cloudformation.Ref("AppPublicSubnet2Cidr")),
		VpcId:            cloudformation.Ref("AppVPC"),
		Tags: []tags.Tag{
			{
				Key:  "Name",
				Value: "AppPublicSubnet2",
			},
		},
	}

	template.Resources["AppPrivateSubnet1"] = &ec2.Subnet{
		CidrBlock: cloudformation.String(cloudformation.Ref("AppPrivateSubnet1Cidr")),
		VpcId:     cloudformation.Ref("AppVPC"),
	}

	template.Resources["AppPrivateSubnet2"] = &ec2.Subnet{
		CidrBlock: cloudformation.String(cloudformation.Ref("AppPrivateSubnet2Cidr")),
		VpcId:     cloudformation.Ref("AppVPC"),
	}

	template.Resources["AppPublicRouteTable1"] = &ec2.RouteTable{
		VpcId: cloudformation.Ref("AppVPC"),
		
	}

	template.Resources["AppPublicRouteTable2"] = &ec2.RouteTable{
		VpcId: cloudformation.Ref("AppVPC"),
	}

	template.Resources["AppInternetGateway"] = &ec2.InternetGateway{}

	template.Resources["AppVpcInternetGatewayAttachment"] = &ec2.VPCGatewayAttachment{
		VpcId:             cloudformation.Ref("AppVPC"),
		InternetGatewayId: cloudformation.String(cloudformation.GetAtt("AppInternetGateway", "InternetGatewayId")),
	}

	template.Resources["AppPublicRouteTable1IGRoute"] = &ec2.Route{
		RouteTableId:         cloudformation.Ref("AppPublicRouteTable1"),
		GatewayId:            cloudformation.String(cloudformation.GetAtt("AppInternetGateway", "InternetGatewayId")),
		DestinationCidrBlock: cloudformation.String("0.0.0.0/0"),
	}

	template.Resources["AppPublicRouteTable2IGRoute"] = &ec2.Route{
		RouteTableId:         cloudformation.Ref("AppPublicRouteTable2"),
		GatewayId:            cloudformation.String(cloudformation.GetAtt("AppInternetGateway", "InternetGatewayId")),
		DestinationCidrBlock: cloudformation.String("0.0.0.0/0"),
	}

}
