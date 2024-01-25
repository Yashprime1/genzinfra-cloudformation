package in1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateIn1MongoDeliveryTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEbsDlmLifecyclePolicy"] = mongo.GetDlmLifeCyclePolicy()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("in1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("in1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "in1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoDirectcallReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "10.112/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.7.144 - 10.14.7.159
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "in1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoDirectcallReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "10.128/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.7.160 - 10.14.7.175
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "in1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoDirectcallReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "10.144/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.7.176 - 10.14.7.191
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(130)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0de320fb812039314")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.EnableMongoLogger = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.StackName = "in1-Mongo-Delivery-1"
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("in1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "Delivery-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}

	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance013132 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013132.EnableEc2instance = true
	MongoReplicaSetInstance013132.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	MongoReplicaSetInstance013132.Ec2Instance.InstanceType = cloudformation.String("r5.xlarge")
	MongoReplicaSetInstance013132.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance013132.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.10.116") //primary
	MongoReplicaSetInstance013132.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance013132.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013132.EnableXvdpGp3 = true
	MongoReplicaSetInstance013132.StopServices = false
	MongoReplicaSetInstance013132.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013132.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	MongoReplicaSetInstance013132.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance013148 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013148.EnableEc2instance = true
	MongoReplicaSetInstance013148.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	MongoReplicaSetInstance013148.Ec2Instance.InstanceType = cloudformation.String("r5.xlarge")
	MongoReplicaSetInstance013148.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance013148.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.10.132")
	MongoReplicaSetInstance013148.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance013148.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013148.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013148.StopServices = false
	MongoReplicaSetInstance013148.EnableXvdpGp3 = true
	MongoReplicaSetInstance013148.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013148.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance013148.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance013164 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013164.EnableEc2instance = true
	MongoReplicaSetInstance013164.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	MongoReplicaSetInstance013164.Ec2Instance.InstanceType = cloudformation.String("r5.xlarge")
	MongoReplicaSetInstance013164.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance013164.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.10.148")
	MongoReplicaSetInstance013164.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance013164.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013164.StopServices = false
	MongoReplicaSetInstance013164.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013164.EnableXvdpGp3 = true
	MongoReplicaSetInstance013164.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013164.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance013164.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/in1/Mongo-Delivery-1", "in1-Mongo-Delivery-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/in1/Mongo-Delivery-1", "in1-Mongo-Delivery-1-Service.json")
}
