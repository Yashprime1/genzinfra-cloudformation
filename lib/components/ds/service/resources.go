package dsservice

import (
	"github.com/awslabs/goformation/v7/cloudformation/iam"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
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
		Roles: []string{
			cloudformation.Ref("DsEcsTaskExecutionRole"),
		},
	}


	template.Resources["DsEcsCluster"] = &ecs.Cluster{}
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
						Name:  cloudformation.String(""),
						Value: cloudformation.String(""),
					},
				},
				LogConfiguration: &ecs.TaskDefinition_LogConfiguration{
					LogDriver: "awslogs",
				},
				Essential: cloudformation.Bool(true),
				MemoryReservation: cloudformation.Int(300),
				Privileged:        cloudformation.Bool(false),
				ReadonlyRootFilesystem: cloudformation.Bool(false),
				User:                   cloudformation.String("ds"),
				Ulimits: []ecs.TaskDefinition_Ulimit{
					{
						Name:      "nofile",
						SoftLimit: 65536,
						HardLimit: 65536,
					},
				},
			},
		},
	}
	template.Resources["DsEcsService"] = &ecs.Service{
		Cluster: cloudformation.String(cloudformation.Ref("DsEcsCluster")),
		DeploymentConfiguration: &ecs.Service_DeploymentConfiguration{
			MaximumPercent:        cloudformation.Int(200),
			MinimumHealthyPercent: cloudformation.Int(50),
		},
		HealthCheckGracePeriodSeconds: cloudformation.Int(60),
		LaunchType:                    cloudformation.String("EC2"),
		LoadBalancers:                 []ecs.Service_LoadBalancer{
			{
				ContainerName:  cloudformation.String("ds"),
				ContainerPort:  cloudformation.Int(80),
				TargetGroupArn: cloudformation.String(cloudformation.ImportValue(cloudformation.Sub("${AWS::StackName}-DsElbTargetGroupArn"))),
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
}
