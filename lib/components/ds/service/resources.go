package dsservice

import (
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
	"github.com/awslabs/goformation/v7/cloudformation/iam"
)

func AddResourcesForDsServiceStack(template *cloudformation.Template) {
	template.Resources["DsEcsTaskRole"] = &iam.Role{
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
	template.Resources["DsEcsTaskRolePolicy"] = &iam.Policy{
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
										"DsEcsTaskRolePolicy",	
									}),
		Roles: []string{
			cloudformation.Ref("DsEcsTaskRole"),
		},
	}
	template.Resources["DsEcsTaskExecutionRole"] = &iam.Role{
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
	template.Resources["DsEcsTaskExecutionRolePolicy"] = &iam.Policy{
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
			"DsEcsTaskExecutionRolePolicy",	
		}),
		Roles: []string{
			cloudformation.Ref("DsEcsTaskExecutionRole"),
		},
	}

	template.Resources["DsEcsCluster"] = &ecs.Cluster{
		ClusterName: cloudformation.String(cloudformation.Ref("AWS::StackName")),
	}
	template.Resources["DsEcsTaskDefinition"] = &ecs.TaskDefinition{
		NetworkMode:      cloudformation.String("host"),
		Family:           cloudformation.String(cloudformation.Ref("AWS::StackName")),
		ExecutionRoleArn: cloudformation.String(cloudformation.Ref("DsEcsTaskExecutionRole")),
		ContainerDefinitions: []ecs.TaskDefinition_ContainerDefinition{
			{
				Name: "ds",
				Image: "httpd",
				Environment: []ecs.TaskDefinition_KeyValuePair{
					{
						Name:  cloudformation.String("App Name"),
						Value: cloudformation.String("DS"),
					},
					{
						Name:  cloudformation.String("TestImport"),
						Value: cloudformation.ImportValue(cloudformation.String(cloudformation.Sub("${AWS::StackName}-DsElbTargetGroupArn"))),
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
						ContainerPort: cloudformation.Int(80),
						HostPort:      cloudformation.Int(80),
						Protocol:      cloudformation.String("tcp"),
					},
				},
			},
			{
				Name: "sensu",
				Image: "sensu/sensu:6.10.0",
				Environment: []ecs.TaskDefinition_KeyValuePair{
					{
						Name:  cloudformation.String("SENSU_BACKEND_URL"),
						Value: cloudformation.String("ws://mu-sen-sensu-voab7tutj8ly-464414231.ap-south-1.elb.amazonaws.com:8081/"),
					},
					{
						Name:  cloudformation.String("SENSU_INSECURE_SKIP_TLS_VERIFY"),
						Value: cloudformation.String("false"),
					},
					{
						Name:  cloudformation.String("SENSU_LOG_LEVEL"),
						Value: cloudformation.String("info"),
					},
					{
						Name:  cloudformation.String("SENSU_STRIP_NETWORKS"),
						Value: cloudformation.String("system"),
					},
					{
						Name:  cloudformation.String("SENSU_API_HOST"),
						Value: cloudformation.String("0.0.0.0"),
					},
					{
						Name:  cloudformation.String("SENSU_CACHE_DIR"),
						Value: cloudformation.String("/var/lib/sensu"),
					},
					{
						Name:  cloudformation.String("SENSU_SUBSCRIPTIONS"),
						Value: cloudformation.String("ds"),
					},
					{
						Name:  cloudformation.String("SENSU_ANNOTATIONS"),
						Value: cloudformation.String("{\"maintainer\": \"Team A\",\"priority\": \"P1\"}"),
					},
					{
						Name:  cloudformation.String("SENSU_LABELS"),
						Value: cloudformation.String("{\"stack\": \"ding\"}"),
					},
					{
						Name:  cloudformation.String("SENSU_NAMESPACE"),
						Value: cloudformation.String("production"),
					},
					{
						Name:  cloudformation.String("SENSU_DETECT_CLOUD_PROVIDER"),
						Value: cloudformation.String("true"),
					},
					{
						Name:  cloudformation.String("SENSU_KEEPALIVE_CRITICAL_TIMEOUT"),
						Value: cloudformation.String("180"),
					},
					{
						Name:  cloudformation.String("SENSU_KEEPALIVE_INTERNAL"),
						Value: cloudformation.String("5"),
					},
					{
						Name:  cloudformation.String("SENSU_KEEPALIVE_WARNING_TIMEOUT"),
						Value: cloudformation.String("120"),
					},
					{
						Name:  cloudformation.String("SENSU_AGENT_MANAGED_ENTITY"),
						Value: cloudformation.String("true"),
					},
					{
						Name:  cloudformation.String("SENSU_DEREGISTER"),
						Value: cloudformation.String("true"),
					},
					{
						Name:  cloudformation.String("SENSU_USER"),
						Value: cloudformation.String("admin"),
					},
					{
						Name:  cloudformation.String("SENSU_PASSWORD"),
						Value: cloudformation.String("P@ssw0rd!"),
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
					"sensu-agent",
					"start",
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
	template.Resources["DsEcsService"] = &ecs.Service{
		Cluster: cloudformation.String(cloudformation.Ref("DsEcsCluster")),
		DeploymentConfiguration: &ecs.Service_DeploymentConfiguration{
			MaximumPercent:        cloudformation.Int(100),
			MinimumHealthyPercent: cloudformation.Int(0),
		},
		HealthCheckGracePeriodSeconds: cloudformation.Int(60),
		LaunchType:                    cloudformation.String("EC2"),
		LoadBalancers:                 []ecs.Service_LoadBalancer{
			{
				ContainerName:  cloudformation.String("ds"),
				ContainerPort:  cloudformation.Int(80),
				TargetGroupArn: cloudformation.String(cloudformation.ImportValue(
					cloudformation.Join("-", []string{
						cloudformation.Ref("AWS::StackName"),
						"DsElbTargetGroupArn",
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
		TaskDefinition:     cloudformation.String(cloudformation.Ref("DsEcsTaskDefinition")),
	}
	// template.Resources["PrometheusEcsTaskDefinition"] = &ecs.TaskDefinition{
	// 	NetworkMode:      cloudformation.String("host"),
	// 	Family:           cloudformation.String(cloudformation.Ref("AWS::StackName")),
	// 	ExecutionRoleArn: cloudformation.String(cloudformation.Ref("DsEcsTaskExecutionRole")),
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
	// 	Cluster: cloudformation.String(cloudformation.Ref("DsEcsCluster")),
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
	// 					"DsPrometheusElbTargetGroupArn",
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
