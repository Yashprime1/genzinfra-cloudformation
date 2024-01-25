package m1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func Generatem1MongoMetadata1Template() {
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
		SubnetCidrBlockSuffix:  "6.128/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "m1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "6.144/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0f8afc87a647567dd")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("m1-SecurityGroup-MongoCommonEc2InstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance006134 := mongo.NewMongo(defaults)
	mongoReplicaInstance006134.EnableEc2instance = false
	mongoReplicaInstance006134.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	mongoReplicaInstance006134.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006134.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.134")
	mongoReplicaInstance006134.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0f5ddafea48e22e91")
	mongoReplicaInstance006134.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006134.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006134.StopServices = false
	mongoReplicaInstance006134.EnableMongoRegistryCache = true
	mongoReplicaInstance006134.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	mongoReplicaInstance006134.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006135 := mongo.NewMongo(defaults)
	mongoReplicaInstance006135.EnableEc2instance = false
	mongoReplicaInstance006135.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	mongoReplicaInstance006135.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006135.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.135")
	mongoReplicaInstance006135.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0f5ddafea48e22e91")
	mongoReplicaInstance006135.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006135.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006135.StopServices = false
	mongoReplicaInstance006135.EnableMongoRegistryCache = true
	mongoReplicaInstance006135.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	mongoReplicaInstance006135.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006150 := mongo.NewMongo(defaults)
	mongoReplicaInstance006150.EnableEc2instance = false
	mongoReplicaInstance006150.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	mongoReplicaInstance006150.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006150.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.150")
	mongoReplicaInstance006150.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0f5ddafea48e22e91")
	mongoReplicaInstance006150.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006150.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006150.StopServices = false
	mongoReplicaInstance006150.EnableMongoRegistryCache = true
	mongoReplicaInstance006150.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	mongoReplicaInstance006150.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/m1/Mongo-Metadata-1", "m1-Mongo-Metadata-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/m1/Mongo-Metadata-1", "m1-Mongo-Metadata-1-Service.json")
}
