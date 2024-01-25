package sg1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSG1MongoExtendedStatsTemplate() {
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
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "42.0/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.42.0 - 10.15.42.15
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "42.16/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.42.16- 10.15.42.31
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "42.32/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.42.32 - 10.15.42.47
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0bbe76ab26239944e")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("sg1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "sg1-Mongo-Extended-Stats-1-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}

	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance042004 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance042004.EnableEc2instance = true
	MongoReplicaSetInstance042004.Ec2Instance.ImageId = cloudformation.String("ami-0a98453160fc9f10a")
	MongoReplicaSetInstance042004.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance042004.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance042004.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.42.4") //primary
	MongoReplicaSetInstance042004.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance042004.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance042004.EnableMongoRegistryCache = true
	MongoReplicaSetInstance042004.StopServices = false
	MongoReplicaSetInstance042004.EnableXvdpGp3 = true
	MongoReplicaSetInstance042004.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance042004.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance042004.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance042020 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance042020.EnableEc2instance = true
	MongoReplicaSetInstance042020.Ec2Instance.ImageId = cloudformation.String("ami-0a98453160fc9f10a")
	MongoReplicaSetInstance042020.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance042020.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance042020.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.42.20")
	MongoReplicaSetInstance042020.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance042020.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance042020.StopServices = false
	MongoReplicaSetInstance042020.EnableMongoRegistryCache = true
	MongoReplicaSetInstance042020.EnableXvdpGp3 = true
	MongoReplicaSetInstance042020.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance042020.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance042020.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance042036 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance042036.EnableEc2instance = true
	MongoReplicaSetInstance042036.Ec2Instance.ImageId = cloudformation.String("ami-0a98453160fc9f10a")
	MongoReplicaSetInstance042036.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance042036.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance042036.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.42.36")
	MongoReplicaSetInstance042036.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance042036.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance042036.EnableMongoRegistryCache = true
	MongoReplicaSetInstance042036.StopServices = false
	MongoReplicaSetInstance042036.EnableXvdpGp3 = true
	MongoReplicaSetInstance042036.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance042036.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance042036.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sg1/Mongo-Extended-Stats-1", "sg1-Mongo-Extended-Stats-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sg1/Mongo-Extended-Stats-1", "sg1-Mongo-Extended-Stats-1-Service.json")
}
