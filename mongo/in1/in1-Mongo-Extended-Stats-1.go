package in1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateIn1MongoExtendedStatsTemplate() {
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
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "58.0/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.12.58.0 - 10.12.58.15
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "in1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "58.16/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.12.58.16- 10.12.58.31
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "in1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "58.32/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.12.58.32 - 10.12.58.47
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-09334c0af4e60072e")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("in1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "in1-Mongo-Extended-Stats-1-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}

	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance058004 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance058004.EnableEc2instance = true
	MongoReplicaSetInstance058004.Ec2Instance.ImageId = cloudformation.String("ami-09334c0af4e60072e")
	MongoReplicaSetInstance058004.Ec2Instance.InstanceType = cloudformation.String("m5.xlarge")
	MongoReplicaSetInstance058004.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance058004.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.58.4") //primary
	MongoReplicaSetInstance058004.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance058004.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance058004.EnableMongoRegistryCache = true
	MongoReplicaSetInstance058004.StopServices = false
	MongoReplicaSetInstance058004.EnableXvdpGp3 = true
	MongoReplicaSetInstance058004.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance058004.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance058004.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance058020 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance058020.EnableEc2instance = true
	MongoReplicaSetInstance058020.Ec2Instance.ImageId = cloudformation.String("ami-09334c0af4e60072e")
	MongoReplicaSetInstance058020.Ec2Instance.InstanceType = cloudformation.String("m5.xlarge")
	MongoReplicaSetInstance058020.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance058020.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.58.20")
	MongoReplicaSetInstance058020.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance058020.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance058020.StopServices = false
	MongoReplicaSetInstance058020.EnableMongoRegistryCache = true
	MongoReplicaSetInstance058020.EnableXvdpGp3 = true
	MongoReplicaSetInstance058020.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance058020.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance058020.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance058036 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance058036.EnableEc2instance = true
	MongoReplicaSetInstance058036.Ec2Instance.ImageId = cloudformation.String("ami-09334c0af4e60072e")
	MongoReplicaSetInstance058036.Ec2Instance.InstanceType = cloudformation.String("m5.xlarge")
	MongoReplicaSetInstance058036.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance058036.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.58.36")
	MongoReplicaSetInstance058036.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance058036.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance058036.EnableMongoRegistryCache = true
	MongoReplicaSetInstance058036.StopServices = false
	MongoReplicaSetInstance058036.EnableXvdpGp3 = true
	MongoReplicaSetInstance058036.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance058036.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance058036.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/in1/Mongo-Extended-Stats-1", "in1-Mongo-Extended-Stats-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/in1/Mongo-Extended-Stats-1", "in1-Mongo-Extended-Stats-1-Service.json")
}
