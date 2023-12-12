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
						Value: cloudformation.String(cloudformation.ImportValue(cloudformation.Sub("${AWS::StackName}-DsElbTargetGroupArn"))),
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
				Name: "Mongo",
				Image: "mongo:latest",
				Environment: []ecs.TaskDefinition_KeyValuePair{
					{
						Name:  cloudformation.String("App Name"),
						Value: cloudformation.String("Mongo"),
					},
				},
				Essential: cloudformation.Bool(false),
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
				Command: []string{
					"--replSet",
					"rs0",
				},
				PortMappings: []ecs.TaskDefinition_PortMapping{
					{
						ContainerPort: cloudformation.Int(27017),
						HostPort:      cloudformation.Int(27017),
						Protocol:      cloudformation.String("tcp"),
					},
				},
			},
			{
				Name: "MongoExporter",
				Image: "percona/mongodb_exporter:0.40.0",
				Environment: []ecs.TaskDefinition_KeyValuePair{
					{
						Name:  cloudformation.String("App Name"),
						Value: cloudformation.String("MongoExporter"),
					},
				},
				Command: []string{
					"---mongodb.uri",
					"http://myUserAdmin:abc123@localhost:17001/admin?ssl=false",
				},
				Essential: cloudformation.Bool(false),
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
						ContainerPort: cloudformation.Int(9216),
						HostPort:      cloudformation.Int(9216),
						Protocol:      cloudformation.String("tcp"),
					},
					{
						ContainerPort: cloudformation.Int(17001),
						HostPort:      cloudformation.Int(17001),
						Protocol:      cloudformation.String("tcp"),
					},
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
