package sensuservice

import (
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
	"github.com/awslabs/goformation/v7/cloudformation/iam"
)

func AddResourcesForSensuServiceStack(template *cloudformation.Template) {
	template.Resources["SensuEcsTaskRole"] = &iam.Role{
		AssumeRolePolicyDocument: map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": []string{
							"ecs-tasks.amazonaws.com",
						},
					},
				},
			},
		},
	}
	template.Resources["SensuEcsTaskRolePolicy"] = &iam.Policy{
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
		PolicyName: cloudformation.Join("-",[]string{
										cloudformation.Ref("AWS::StackName"),
										"SensuEcsTaskRolePolicy",	
									}),
		Roles: []string{
			cloudformation.Ref("SensuEcsTaskRole"),
		},
	}
	template.Resources["SensuEcsTaskExecutionRole"] = &iam.Role{
		AssumeRolePolicyDocument: map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": []string{
							"ecs-tasks.amazonaws.com",
						},
					},
				},
			},
		},
	}
	template.Resources["SensuEcsTaskExecutionRolePolicy"] = &iam.Policy{
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
		PolicyName: cloudformation.Join("-",[]string{
			cloudformation.Ref("AWS::StackName"),
			"SensuEcsTaskExecutionRolePolicy",	
		}),
		Roles: []string{
			cloudformation.Ref("SensuEcsTaskExecutionRole"),
		},
	}

	template.Resources["SensuEcsCluster"] = &ecs.Cluster{
		ClusterName: cloudformation.String(cloudformation.Ref("AWS::StackName")),
	}
	template.Resources["SensuEcsTaskDefinition"] = &ecs.TaskDefinition{
		NetworkMode:      cloudformation.String("host"),
		Family:           cloudformation.String(cloudformation.Ref("AWS::StackName")),
		ExecutionRoleArn: cloudformation.String(cloudformation.Ref("SensuEcsTaskExecutionRole")),
		ContainerDefinitions: []ecs.TaskDefinition_ContainerDefinition{
			{
				Name: "sensu",
				Image: "sensu/sensu:6.10.0",
				Environment: []ecs.TaskDefinition_KeyValuePair{
					{
						Name:  cloudformation.String("App Name"),
						Value: cloudformation.String("Sensu"),
					},
				},
				Essential: cloudformation.Bool(true),
				MemoryReservation: cloudformation.Int(256),
				Privileged:        cloudformation.Bool(false),
				ReadonlyRootFilesystem: cloudformation.Bool(false),
				Ulimits: []ecs.TaskDefinition_Ulimit{
					{
						Name:      "nofile",
						SoftLimit: 65536,
						HardLimit: 65536,
					},
				},
				PortMappings: []ecs.TaskDefinition_PortMapping{
					{
						ContainerPort: cloudformation.Int(3000),
						HostPort:      cloudformation.Int(3000),
						Protocol:      cloudformation.String("tcp"),
					},
					{
						ContainerPort: cloudformation.Int(8080),
						HostPort:      cloudformation.Int(8080),
						Protocol:      cloudformation.String("tcp"),
					},
					{
						ContainerPort: cloudformation.Int(8081),
						HostPort:      cloudformation.Int(8081),
						Protocol:      cloudformation.String("tcp"),
					},
				},
				MountPoints: []ecs.TaskDefinition_MountPoint{
					{
						ContainerPath: cloudformation.String("/var/lib/sensu"),
						SourceVolume:  cloudformation.String("sensu-backend"),
					},
				},
				Command: []string{
					"sensu-backend",
					"start",
					"--state-dir",
					"/var/lib/sensu/sensu-backend",
					"--log-level",
					"info",
				},
			},
		},
		Volumes: []ecs.TaskDefinition_Volume{
			{
				Name: cloudformation.String("sensu-backend"),
				Host: &ecs.TaskDefinition_HostVolumeProperties{
					SourcePath: cloudformation.String("/var/lib/sensu"),
				},
			},
		},
	}
	template.Resources["SensuEcsService"] = &ecs.Service{
		Cluster: cloudformation.String(cloudformation.Ref("SensuEcsCluster")),
		DeploymentConfiguration: &ecs.Service_DeploymentConfiguration{
			MaximumPercent:        cloudformation.Int(100),
			MinimumHealthyPercent: cloudformation.Int(0),
		},
		HealthCheckGracePeriodSeconds: cloudformation.Int(60),
		LaunchType:                    cloudformation.String("EC2"),
		LoadBalancers:                 []ecs.Service_LoadBalancer{
			{
				ContainerName:  cloudformation.String("sensu"),
				ContainerPort:  cloudformation.Int(3000),
				TargetGroupArn: cloudformation.String(cloudformation.ImportValue(
					cloudformation.Join("-", []string{
						cloudformation.Ref("AWS::StackName"),
						"SensuElbTargetGroupArn",
					}),
				)),
			},
			{
				ContainerName:  cloudformation.String("sensu"),
				ContainerPort:  cloudformation.Int(8081),
				TargetGroupArn: cloudformation.String(cloudformation.ImportValue(
					cloudformation.Join("-", []string{
						cloudformation.Ref("AWS::StackName"),
						"SensuBackendElbTargetGroupArn",
					}),
				)),
			},
			{
				ContainerName:  cloudformation.String("sensu"),
				ContainerPort:  cloudformation.Int(8080),
				TargetGroupArn: cloudformation.String(cloudformation.ImportValue(
					cloudformation.Join("-", []string{
						cloudformation.Ref("AWS::StackName"),
						"SensuApiElbTargetGroupArn",
					}),
				)),
			},
		},
		PlacementConstraints: []ecs.Service_PlacementConstraint{
			{
				Type: "distinctInstance",
			},
		},
		SchedulingStrategy: cloudformation.String("DAEMON"),
		TaskDefinition:     cloudformation.String(cloudformation.Ref("SensuEcsTaskDefinition")),
	}
	// template.Resources["PrometheusEcsTaskDefinition"] = &ecs.TaskDefinition{
	// 	NetworkMode:      cloudformation.String("host"),
	// 	Family:           cloudformation.String(cloudformation.Ref("AWS::StackName")),
	// 	ExecutionRoleArn: cloudformation.String(cloudformation.Ref("SensuEcsTaskExecutionRole")),
	// 	ContainerDefinitions: []ecs.TaskDefinition_ContainerDefinition{
	// 		{
	// 			Name: "prometheus",
	// 			Image: "bitnami/prometheus",
	// 			Environment: []ecs.TaskDefinition_KeyValuePair{
	// 				{
	// 					Name:  cloudformation.String("App Name"),
	// 					Value: cloudformation.String("prometheuss"),
	// 				},
	// 			},
	// 			Essential: cloudformation.Bool(true),
	// 			MemoryReservation: cloudformation.Int(64),
	// 			Privileged:        cloudformation.Bool(false),
	// 			ReadonlyRootFilesystem: cloudformation.Bool(false),
	// 			Ulimits: []ecs.TaskDefinition_Ulimit{
	// 				{
	// 					Name:      "nofile",
	// 					SoftLimit: 65536,
	// 					HardLimit: 65536,
	// 				},
	// 			},
	// 			PortMappings: []ecs.TaskDefinition_PortMapping{
	// 				{
	// 					ContainerPort: cloudformation.Int(9090),
	// 					HostPort:      cloudformation.Int(9090),
	// 					Protocol:      cloudformation.String("tcp"),
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// template.Resources["PrometheusEcsService"] = &ecs.Service{
	// 	Cluster: cloudformation.String(cloudformation.Ref("SensuEcsCluster")),
	// 	DeploymentConfiguration: &ecs.Service_DeploymentConfiguration{
	// 		MaximumPercent:        cloudformation.Int(100),
	// 		MinimumHealthyPercent: cloudformation.Int(0),
	// 	},
	// 	HealthCheckGracePeriodSeconds: cloudformation.Int(65),
	// 	LaunchType:                    cloudformation.String("EC2"),
	// 	LoadBalancers:                 []ecs.Service_LoadBalancer{
	// 		{
	// 			ContainerName:  cloudformation.String("prometheus"),
	// 			ContainerPort:  cloudformation.Int(9090),
	// 			TargetGroupArn: cloudformation.String(cloudformation.ImportValue(
	// 				cloudformation.Join("-", []string{
	// 					cloudformation.Ref("AWS::StackName"),
	// 					"SensuPrometheusElbTargetGroupArn",
	// 				}),
	// 			)),
	// 		},
	// 	},
	// 	PlacementConstraints: []ecs.Service_PlacementConstraint{
	// 		{
	// 			Type: "distinctInstance",
	// 		},
	// 	},
	// 	SchedulingStrategy: cloudformation.String("DAEMON"),
	// 	TaskDefinition:     cloudformation.String(cloudformation.Ref("PrometheusEcsTaskDefinition")),
	// }
}
