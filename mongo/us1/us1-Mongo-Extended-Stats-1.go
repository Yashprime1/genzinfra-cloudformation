package us1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateUs1MongoExtendedStatsTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEbsDlmLifecyclePolicy"] = mongo.GetDlmLifeCyclePolicy()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("us1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("us1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "us1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "99.0/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.99.0 - 10.11.40.63
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "us1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "99.16/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.99.16- 10.16.99.31
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "us1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoExtendedStatsReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "99.32/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.99.32 - 10.16.99.47
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-04d527fb9bed69e16")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("us1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "us1-Mongo-Extended-Stats-1-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}

	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance099004 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance099004.EnableEc2instance = true
	MongoReplicaSetInstance099004.Ec2Instance.ImageId = cloudformation.String("ami-02726fee3f464392b")
	MongoReplicaSetInstance099004.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance099004.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance099004.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.99.4")
	MongoReplicaSetInstance099004.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance099004.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance099004.EnableMongoRegistryCache = true
	MongoReplicaSetInstance099004.StopServices = false
	MongoReplicaSetInstance099004.EnableXvdpGp3 = true
	MongoReplicaSetInstance099004.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance099004.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance099004.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance099020 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance099020.EnableEc2instance = true
	MongoReplicaSetInstance099020.Ec2Instance.ImageId = cloudformation.String("ami-02726fee3f464392b")
	MongoReplicaSetInstance099020.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance099020.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance099020.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.99.20") //primary
	MongoReplicaSetInstance099020.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance099020.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance099020.StopServices = false
	MongoReplicaSetInstance099020.EnableMongoRegistryCache = true
	MongoReplicaSetInstance099020.EnableXvdpGp3 = true
	MongoReplicaSetInstance099020.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance099020.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance099020.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance099036 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance099036.EnableEc2instance = true
	MongoReplicaSetInstance099036.Ec2Instance.ImageId = cloudformation.String("ami-02726fee3f464392b")
	MongoReplicaSetInstance099036.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance099036.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance099036.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.99.36")
	MongoReplicaSetInstance099036.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance099036.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance099036.EnableMongoRegistryCache = true
	MongoReplicaSetInstance099036.StopServices = false
	MongoReplicaSetInstance099036.EnableXvdpGp3 = true
	MongoReplicaSetInstance099036.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance099036.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance099036.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/us1/Mongo-Extended-Stats-1", "us1-Mongo-Extended-Stats-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/us1/Mongo-Extended-Stats-1", "us1-Mongo-Extended-Stats-1-Service.json")
}
