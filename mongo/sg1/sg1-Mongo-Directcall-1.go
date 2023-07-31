package sg1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSg1MongoDirectcallTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEbsDlmLifecyclePolicy"] = mongo.GetDlmLifeCyclePolicy()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("sg1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("sg1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoDirectcallReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "13.128/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.7.144 - 10.14.7.159
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoDirectcallReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "13.144/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.7.160 - 10.14.7.175
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoDirectcallReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "13.160/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.7.176 - 10.14.7.191
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0d85a72517fc05f6b")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("sg1-SecurityGroup-DirectcallMongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "dc-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}

	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance013132 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013132.EnableEc2instance = true
	MongoReplicaSetInstance013132.Ec2Instance.ImageId = cloudformation.String("ami-02df2edff7704b9fe")
	MongoReplicaSetInstance013132.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance013132.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance013132.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.13.132")
	MongoReplicaSetInstance013132.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013132.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013132.EnableXvdpGp3 = true
	MongoReplicaSetInstance013132.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013132.StopServices = false
	MongoReplicaSetInstance013132.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013132.MongoContainerTag = "bamboo-mongo-task-2652-mongo-1"
	MongoReplicaSetInstance013132.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance013148 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013148.EnableEc2instance = true
	MongoReplicaSetInstance013148.Ec2Instance.ImageId = cloudformation.String("ami-02df2edff7704b9fe")
	MongoReplicaSetInstance013148.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance013148.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance013148.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.13.148") //primary
	MongoReplicaSetInstance013148.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013148.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013148.EnableXvdpGp3 = true
	MongoReplicaSetInstance013148.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013148.MongoContainerTag = "bamboo-mongo-task-2652-mongo-1"
	MongoReplicaSetInstance013148.StopServices = false
	MongoReplicaSetInstance013148.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013148.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance013164 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013164.EnableEc2instance = true
	MongoReplicaSetInstance013164.Ec2Instance.ImageId = cloudformation.String("ami-02df2edff7704b9fe")
	MongoReplicaSetInstance013164.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance013164.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance013164.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.13.164")
	MongoReplicaSetInstance013164.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013164.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013164.EnableXvdpGp3 = true
	MongoReplicaSetInstance013164.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013164.StopServices = false
	MongoReplicaSetInstance013164.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013164.MongoContainerTag = "bamboo-mongo-task-2652-mongo-1"
	MongoReplicaSetInstance013164.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sg1/Mongo-Directcall-1", "sg1-Mongo-Directcall-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sg1/Mongo-Directcall-1", "sg1-Mongo-Directcall-1-Service.json")
}
