package mec1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateMec1MongoStatsTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEbsDlmLifecyclePolicy"] = mongo.GetDlmLifeCyclePolicy()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("mec1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("mec1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "mec1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoStatsReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.112/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.7.112 - 10.16.7.128
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "mec1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoStatsReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.128/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.7.128 - 10.16.13.144
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "mec1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoStatsReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.144/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.7.144 - 10.16.7.160
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-022cae5767ebb7d1d")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("mec1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "mec1-Mongo-Stats-1-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}

	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance013197 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013197.EnableEc2instance = true
	MongoReplicaSetInstance013197.Ec2Instance.ImageId = cloudformation.String("ami-0c3ddec888ab9cd5e")
	MongoReplicaSetInstance013197.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance013197.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance013197.Ec2Instance.PrivateIpAddress = cloudformation.String("10.19.7.118") //primary
	MongoReplicaSetInstance013197.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013197.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013197.StopServices = false
	MongoReplicaSetInstance013197.EnableXvdpGp3 = true
	MongoReplicaSetInstance013197.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013197.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013197.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	MongoReplicaSetInstance013197.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance013213 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013213.EnableEc2instance = true
	MongoReplicaSetInstance013213.Ec2Instance.ImageId = cloudformation.String("ami-0c3ddec888ab9cd5e")
	MongoReplicaSetInstance013213.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance013213.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance013213.Ec2Instance.PrivateIpAddress = cloudformation.String("10.19.7.134")
	MongoReplicaSetInstance013213.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013213.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013213.StopServices = false
	MongoReplicaSetInstance013213.EnableXvdpGp3 = true
	MongoReplicaSetInstance013213.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013213.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013213.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance013213.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance013229 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013229.EnableEc2instance = true
	MongoReplicaSetInstance013229.Ec2Instance.ImageId = cloudformation.String("ami-0c3ddec888ab9cd5e")
	MongoReplicaSetInstance013229.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance013229.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance013229.Ec2Instance.PrivateIpAddress = cloudformation.String("10.19.7.150")
	MongoReplicaSetInstance013229.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013229.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013229.StopServices = false
	MongoReplicaSetInstance013229.EnableXvdpGp3 = true
	MongoReplicaSetInstance013229.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013229.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013229.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance013229.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/mec1/Mongo-Stats-1", "mec1-Mongo-Stats-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/mec1/Mongo-Stats-1", "mec1-Mongo-Stats-1-Service.json")
}
