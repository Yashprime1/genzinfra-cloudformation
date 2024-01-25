package m1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func Generatem1Mongo1Template() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("m1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("m1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "m1",
		AvailabilityZoneSuffix: "a",
		SubnetCidrBlockSuffix:  "6.96/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "m1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "6.112/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(4608)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0f8afc87a647567dd")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("m1-SecurityGroup-MongoCommonEc2InstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance006101 := mongo.NewMongo(defaults)
	mongoReplicaInstance006101.EnableEc2instance = false
	mongoReplicaInstance006101.Ec2Instance.InstanceType = cloudformation.String("r5.24xlarge")
	mongoReplicaInstance006101.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006101.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.101")
	mongoReplicaInstance006101.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0b78eada87c6da5bf")
	mongoReplicaInstance006101.XvdpEc2Volume.Size = cloudformation.Int(4608)
	mongoReplicaInstance006101.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006101.StopServices = false
	mongoReplicaInstance006101.EnableMongoRegistryCache = true
	mongoReplicaInstance006101.MongoContainerTag = "bamboo-mongo-master-15"
	mongoReplicaInstance006101.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006102 := mongo.NewMongo(defaults)
	mongoReplicaInstance006102.EnableEc2instance = false
	mongoReplicaInstance006102.Ec2Instance.InstanceType = cloudformation.String("r5.24xlarge")
	mongoReplicaInstance006102.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006102.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.102")
	mongoReplicaInstance006102.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0b78eada87c6da5bf")
	mongoReplicaInstance006102.XvdpEc2Volume.Size = cloudformation.Int(4608)
	mongoReplicaInstance006102.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006102.StopServices = false
	mongoReplicaInstance006102.EnableMongoRegistryCache = true
	mongoReplicaInstance006102.MongoContainerTag = "bamboo-mongo-master-15"
	mongoReplicaInstance006102.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006117 := mongo.NewMongo(defaults)
	mongoReplicaInstance006117.EnableEc2instance = false
	mongoReplicaInstance006117.Ec2Instance.InstanceType = cloudformation.String("r5.24xlarge")
	mongoReplicaInstance006117.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006117.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.117")
	mongoReplicaInstance006117.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0b78eada87c6da5bf")
	mongoReplicaInstance006117.XvdpEc2Volume.Size = cloudformation.Int(4608)
	mongoReplicaInstance006117.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006117.StopServices = false
	mongoReplicaInstance006117.EnableMongoRegistryCache = true
	mongoReplicaInstance006117.MongoContainerTag = "bamboo-mongo-master-15"
	mongoReplicaInstance006117.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/m1/Mongo-1", "m1-Mongo-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/m1/Mongo-1", "m1-Mongo-1-Service.json")
}
