package sensubase

import (
	"github.com/awslabs/goformation/v7/cloudformation/elasticloadbalancingv2"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/autoscaling"
	"github.com/awslabs/goformation/v7/cloudformation/iam"
)

func AddResourcesForSensuBaseStack(template *cloudformation.Template, defaults SensuBaseDefaults) {
	template.Resources["SensuEc2IamRole"] = &iam.Role{
		AssumeRolePolicyDocument: map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": []string{
							"ec2.amazonaws.com",
						},
					},
				},
			},
		},
	}
	template.Resources["SensuEc2RolePolicy"] = &iam.Policy{
		PolicyDocument: map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Action": "*",
					"Effect": "Allow",
					"Resource": []string{
						"*",
					},
				},
			},
		},
		PolicyName: cloudformation.Join("-", []string{
			cloudformation.Ref("AWS::StackName"),
			"SensuEc2RolePolicy",
		}),
		Roles: []string{
			cloudformation.Ref("SensuEc2IamRole"),
		},
	}
	template.Resources["SensuEc2InstanceProfile"] = &iam.InstanceProfile{
		Path: cloudformation.String("/"),
		Roles: []string{
			cloudformation.Ref("SensuEc2IamRole"),
		},
	}
	template.Resources["SensuLaunchConfiguration"] = &autoscaling.LaunchConfiguration{
		AssociatePublicIpAddress: cloudformation.Bool(true),
		ImageId:                  cloudformation.Ref("SensuAmiId"),
		InstanceType:             cloudformation.Ref("SensuInstanceType"),
		InstanceMonitoring:       cloudformation.Bool(false),
		SecurityGroups: []string{
			cloudformation.ImportValue(defaults.SecurityGroupStack + "-DS2SecurityGroupId"),
		},
		IamInstanceProfile: cloudformation.String(cloudformation.Ref("SensuEc2InstanceProfile")),
	}
	template.Resources["SensuAsg"] = &autoscaling.AutoScalingGroup{
		DesiredCapacity:         cloudformation.String(cloudformation.Ref("SensuDesiredSize")),
		LaunchConfigurationName: cloudformation.String(cloudformation.Ref("SensuLaunchConfiguration")),
		MaxSize:                 cloudformation.Ref("SensuMaxSize"),
		MinSize:                 cloudformation.Ref("SensuMinSize"),
		VPCZoneIdentifier: []string{
			cloudformation.ImportValue(defaults.NetworkStack + "-AppPublicSubnet1Id"),
			cloudformation.ImportValue(defaults.NetworkStack + "-AppPublicSubnet2Id"),
		},
		TargetGroupARNs: []string{
			cloudformation.Ref("SensuElbTargetGroup"),
		},
	}
	template.Resources["SensuElb"] = &elasticloadbalancingv2.LoadBalancer{
		Scheme:        cloudformation.String("internet-facing"),
		Type:          cloudformation.String("application"),
		IpAddressType: cloudformation.String("ipv4"),
		LoadBalancerAttributes: []elasticloadbalancingv2.LoadBalancer_LoadBalancerAttribute{
			{
				Key:   cloudformation.String("idle_timeout.timeout_seconds"),
				Value: cloudformation.String("1800"),
			},
		},
		SecurityGroups: []string{
			cloudformation.ImportValue(defaults.SecurityGroupStack + "-DS2SecurityGroupId"),
		},
		Subnets: []string{
			cloudformation.ImportValue(defaults.NetworkStack + "-AppPublicSubnet1Id"),
			cloudformation.ImportValue(defaults.NetworkStack + "-AppPublicSubnet2Id"),
		},
	}
	template.Resources["SensuElbTargetGroup"] = &elasticloadbalancingv2.TargetGroup{
		HealthCheckIntervalSeconds: cloudformation.Int(30),
		HealthCheckPath:            cloudformation.String("/"),
		HealthCheckProtocol:        cloudformation.String("HTTP"),
		HealthCheckTimeoutSeconds:  cloudformation.Int(10),
		HealthyThresholdCount:      cloudformation.Int(5),
		Matcher: &elasticloadbalancingv2.TargetGroup_Matcher{
			HttpCode: cloudformation.String("200"),
		},
		Port:                    cloudformation.Int(8080),
		Protocol:                cloudformation.String("HTTP"),
		TargetType:              cloudformation.String("instance"),
		UnhealthyThresholdCount: cloudformation.Int(2),
		VpcId:                   cloudformation.String(cloudformation.ImportValue(defaults.NetworkStack + "-AppVPCId")),
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
	template.Resources["SensuElbListener"] = &elasticloadbalancingv2.Listener{
		DefaultActions: []elasticloadbalancingv2.Listener_Action{
			{
				TargetGroupArn: cloudformation.String(cloudformation.Ref("SensuElbTargetGroup")),
				Type:           "forward",
			},
		},
		LoadBalancerArn: cloudformation.Ref("SensuElb"),
		Port:            cloudformation.Int(8080),
		Protocol:        cloudformation.String("HTTP"),
	}
}
