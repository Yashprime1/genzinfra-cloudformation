package sk1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSk1MongoIntegrationDebuggerTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("sk1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("sk1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoIntegrationDebuggerReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "58.48/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.14.58.0 - 10.14.58.15
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoIntegrationDebuggerReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "58.64/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.14.58.16- 10.14.58.31
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoIntegrationDebuggerReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "58.80/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.14.58.32 - 10.14.58.47
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-01cc09728d28c3961")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("sk1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "mongo-integration-debugger-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}

	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance058052 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance058052.EnableEc2instance = true
	MongoReplicaSetInstance058052.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance058052.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance058052.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance058052.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.58.52") //primary
	MongoReplicaSetInstance058052.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance058052.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance058052.EnableMongoRegistryCache = true
	MongoReplicaSetInstance058052.StopServices = false
	MongoReplicaSetInstance058052.EnableXvdpGp3 = true
	MongoReplicaSetInstance058052.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance058052.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	MongoReplicaSetInstance058052.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance058068 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance058068.EnableEc2instance = true
	MongoReplicaSetInstance058068.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance058068.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance058068.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance058068.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.58.68")
	MongoReplicaSetInstance058068.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance058068.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance058068.StopServices = false
	MongoReplicaSetInstance058068.EnableMongoRegistryCache = true
	MongoReplicaSetInstance058068.EnableXvdpGp3 = true
	MongoReplicaSetInstance058068.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance058068.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	MongoReplicaSetInstance058068.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance058084 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance058084.EnableEc2instance = true
	MongoReplicaSetInstance058084.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance058084.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance058084.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance058084.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.58.84")
	MongoReplicaSetInstance058084.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance058084.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance058084.EnableMongoRegistryCache = true
	MongoReplicaSetInstance058084.StopServices = false
	MongoReplicaSetInstance058084.EnableXvdpGp3 = true
	MongoReplicaSetInstance058084.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance058084.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	MongoReplicaSetInstance058084.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sk1/Mongo-Integration-Debugger-1", "sk1-Mongo-Integration-Debugger-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sk1/Mongo-Integration-Debugger-1", "sk1-Mongo-Integration-Debugger-1-Service.json")
}
