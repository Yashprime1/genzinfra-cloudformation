package us1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ec2"
	"github.com/awslabs/goformation/v7/cloudformation/tags"
)

func GetDefaultMongoConfiguration() mongo.Mongo {
	defaults := mongo.Mongo{
		Ec2Instance: ec2.Instance{
			DisableApiTermination:             cloudformation.Bool(false),
			IamInstanceProfile:                cloudformation.String(cloudformation.Ref("MongoEc2InstanceIamInstanceProfile")),
			ImageId:                           cloudformation.String("ami-04b18275a0f59862d"),
			InstanceInitiatedShutdownBehavior: cloudformation.String("stop"),
			InstanceType:                      cloudformation.String("t2.small"),
			Monitoring:                        cloudformation.Bool(true),
			SecurityGroupIds: []string{
				cloudformation.ImportValue("us1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
			},
			Tags: []tags.Tag{
				{
					Key:   "customer",
					Value: "clevertap",
				},
				{
					Key:   "role",
					Value: "mongo",
				},
				{
					Key:   "ecs_cluster",
					Value: cloudformation.Ref("MongoEcsCluster"),
				},
			},
		},
		XvdpEc2Volume: ec2.Volume{
			Encrypted:  cloudformation.Bool(true),
			KmsKeyId:   cloudformation.String(cloudformation.GetAtt("MongoVolumeXvdpKmsKey", "Arn")),
			Tags:       []tags.Tag{},
			VolumeType: cloudformation.String("gp2"),
		},
		StackPrefix:       "us1",
		EnableEc2instance: true,
	}

	tagSnapshotIdentifierValue := make([]string, 0)
	tagSnapshotIdentifierValue = append(tagSnapshotIdentifierValue, cloudformation.Ref("AWS::StackName"), "MongoVolumeXvdp")
	defaults.XvdpEc2Volume.Tags = append(defaults.XvdpEc2Volume.Tags, tags.Tag{Key: "DlmIdentifier", Value: cloudformation.Join("-", tagSnapshotIdentifierValue)})
	defaults.XvdpEc2Volume.AWSCloudFormationDeletionPolicy = "Snapshot"

	return defaults
}
