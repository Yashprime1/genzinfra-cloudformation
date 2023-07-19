package dsbase

import (
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
	"github.com/awslabs/goformation/v7/cloudformation/elasticloadbalancingv2"
)

func AddResourcesForDsBaseStack(template *cloudformation.Template, defaults DsBaseDefaults) {
	template.Resources["DsEcsCluster"] = &ecs.Cluster{}
	template.Resources["DsEcsCluster2"] = &ecs.Cluster{}
	template.Resources["DsElb"] = &elasticloadbalancingv2.LoadBalancer{
		Scheme:        cloudformation.String("internet-facing"),
		Type:          cloudformation.String("application"),
		IpAddressType: cloudformation.String("ipv4"),
		SecurityGroups: []string{
			cloudformation.ImportValue(defaults.SecurityGroupStack + "-DS2SecurityGroupId"),
		},
		Subnets: []string{
			cloudformation.ImportValue(defaults.NetworkStack + "-AppPublicSubnet1Id"),
			cloudformation.ImportValue(defaults.NetworkStack + "-AppPublicSubnet2Id"),
		},
	}

	template.Resources["DsElbTargetGroup"] = &elasticloadbalancingv2.TargetGroup{
		HealthCheckIntervalSeconds: cloudformation.Int(30),
		HealthCheckPath:            cloudformation.String("/"),
		HealthCheckProtocol:        cloudformation.String("HTTP"),
		HealthCheckTimeoutSeconds:  cloudformation.Int(5),
		HealthyThresholdCount:      cloudformation.Int(5),
		Matcher: &elasticloadbalancingv2.TargetGroup_Matcher{
			HttpCode: cloudformation.String("200"),
		},
		Port:                    cloudformation.Int(80),
		Protocol:                cloudformation.String("HTTP"),
		TargetType:              cloudformation.String("ip"),
		UnhealthyThresholdCount: cloudformation.Int(2),
		VpcId:                   cloudformation.String(cloudformation.ImportValue(defaults.NetworkStack+"-AppVPCId")),
		TargetGroupAttributes: []elasticloadbalancingv2.TargetGroup_TargetGroupAttribute{
			{
				Key:   cloudformation.String("stickiness.enabled"),
				Value: cloudformation.String("false"),
			},
			{
				Key:   cloudformation.String("deregistration_delay.timeout_seconds"),
				Value: cloudformation.String("300"),
			},
		},
	}

	template.Resources["DsElbListener"] = &elasticloadbalancingv2.Listener{
		DefaultActions: []elasticloadbalancingv2.Listener_Action{
			{
				TargetGroupArn: cloudformation.String(cloudformation.Ref("DsElbTargetGroup")),
				Type:          "forward",
			},
		},
		LoadBalancerArn: cloudformation.Ref("DsElb"),
		Port:            cloudformation.Int(80),
		Protocol:        cloudformation.String("HTTP"),
	}
}
