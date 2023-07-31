package eu1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateEU1MongoExtendedStatsTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEbsDlmLifecyclePolicy"] = mongo.GetDlmLifeCyclePolicy()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("eu1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("eu1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "40.48/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.11.40.48 - 10.11.40.63
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "40.64/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.11.40.64- 10.11.40.79
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "40.80/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.11.40.80 - 10.11.40.95
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0e199854def4d1705")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("eu1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "eu1-Mongo-Extended-Stats-1-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}

	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance040052 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance040052.EnableEc2instance = true
	MongoReplicaSetInstance040052.Ec2Instance.ImageId = cloudformation.String("ami-0e199854def4d1705")
	MongoReplicaSetInstance040052.Ec2Instance.InstanceType = cloudformation.String("m5.2xlarge")
	MongoReplicaSetInstance040052.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance040052.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.40.52") //primary
	MongoReplicaSetInstance040052.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance040052.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance040052.EnableMongoRegistryCache = true
	MongoReplicaSetInstance040052.StopServices = false
	MongoReplicaSetInstance040052.EnableXvdpGp3 = true
	MongoReplicaSetInstance040052.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance040052.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance040052.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance040068 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance040068.EnableEc2instance = true
	MongoReplicaSetInstance040068.Ec2Instance.ImageId = cloudformation.String("ami-0e199854def4d1705")
	MongoReplicaSetInstance040068.Ec2Instance.InstanceType = cloudformation.String("m5.2xlarge")
	MongoReplicaSetInstance040068.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance040068.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.40.68")
	MongoReplicaSetInstance040068.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance040068.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance040068.StopServices = false
	MongoReplicaSetInstance040068.EnableMongoRegistryCache = true
	MongoReplicaSetInstance040068.EnableXvdpGp3 = true
	MongoReplicaSetInstance040068.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance040068.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance040068.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance040084 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance040084.EnableEc2instance = true
	MongoReplicaSetInstance040084.Ec2Instance.ImageId = cloudformation.String("ami-0e199854def4d1705")
	MongoReplicaSetInstance040084.Ec2Instance.InstanceType = cloudformation.String("m5.2xlarge")
	MongoReplicaSetInstance040084.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance040084.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.40.84")
	MongoReplicaSetInstance040084.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance040084.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance040084.EnableMongoRegistryCache = true
	MongoReplicaSetInstance040084.StopServices = false
	MongoReplicaSetInstance040084.EnableXvdpGp3 = true
	MongoReplicaSetInstance040084.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance040084.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance040084.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/eu1/Mongo-Extended-Stats-1", "eu1-Mongo-Extended-Stats-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/eu1/Mongo-Extended-Stats-1", "eu1-Mongo-Extended-Stats-1-Service.json")
}
