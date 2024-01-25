package in1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateIn1MongoAccountsTemplate() {
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
		Ecc2SubnetLogicalId:    "MongoAccountsReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.48/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.15.7.48 - 10.15.7.64
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "in1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoAccountsReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.64/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.49 - 10.14.6.62
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "in1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoAccountsReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.96/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.65 - 10.14.6.78
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableVividCortexEcsService = true
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("in1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "accounts-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007052 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007052.EnableEc2instance = true
	MongoReplicaSetInstance007052.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	MongoReplicaSetInstance007052.Ec2Instance.InstanceType = cloudformation.String("c5.9xlarge")
	MongoReplicaSetInstance007052.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007052.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.7.52")
	MongoReplicaSetInstance007052.XvdpEc2Volume.Size = cloudformation.Int(192)
	MongoReplicaSetInstance007052.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007052.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance007052.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007052.StopServices = false
	MongoReplicaSetInstance007052.EnableXvdpGp3 = true
	MongoReplicaSetInstance007052.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007052.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007053 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007053.EnableEc2instance = false
	MongoReplicaSetInstance007053.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	MongoReplicaSetInstance007053.Ec2Instance.InstanceType = cloudformation.String("c5.9xlarge")
	MongoReplicaSetInstance007053.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007053.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.7.53")
	MongoReplicaSetInstance007053.XvdpEc2Volume.Size = cloudformation.Int(192)
	MongoReplicaSetInstance007053.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0436300a954b5a27b")
	MongoReplicaSetInstance007053.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007053.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance007053.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007053.EnableXvdpGp3 = true
	MongoReplicaSetInstance007053.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007053.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007054 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007054.EnableEc2instance = false
	MongoReplicaSetInstance007054.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	MongoReplicaSetInstance007054.Ec2Instance.InstanceType = cloudformation.String("c5.9xlarge")
	MongoReplicaSetInstance007054.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007054.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.7.54")
	MongoReplicaSetInstance007054.XvdpEc2Volume.Size = cloudformation.Int(192)
	MongoReplicaSetInstance007054.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0f8b081d592b8f73c")
	MongoReplicaSetInstance007054.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007054.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance007054.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007054.EnableXvdpGp3 = true
	MongoReplicaSetInstance007054.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007054.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007055 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007055.EnableEc2instance = false
	MongoReplicaSetInstance007055.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	MongoReplicaSetInstance007055.Ec2Instance.InstanceType = cloudformation.String("c5.9xlarge")
	MongoReplicaSetInstance007055.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007055.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.7.55")
	MongoReplicaSetInstance007055.XvdpEc2Volume.Size = cloudformation.Int(192)
	MongoReplicaSetInstance007055.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-048bd80dfd850d1d8")
	MongoReplicaSetInstance007055.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007055.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance007055.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007055.EnableXvdpGp3 = true
	MongoReplicaSetInstance007055.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007055.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007056 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007056.EnableEc2instance = false
	MongoReplicaSetInstance007056.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	MongoReplicaSetInstance007056.Ec2Instance.InstanceType = cloudformation.String("c5.9xlarge")
	MongoReplicaSetInstance007056.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007056.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.7.56")
	MongoReplicaSetInstance007056.XvdpEc2Volume.Size = cloudformation.Int(192)
	MongoReplicaSetInstance007056.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-09978ff9cfe1c540a")
	MongoReplicaSetInstance007056.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007056.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance007056.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007056.EnableXvdpGp3 = true
	MongoReplicaSetInstance007056.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007056.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007068 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007068.EnableEc2instance = true
	MongoReplicaSetInstance007068.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	MongoReplicaSetInstance007068.Ec2Instance.InstanceType = cloudformation.String("c5.9xlarge")
	MongoReplicaSetInstance007068.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007068.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.7.68")
	MongoReplicaSetInstance007068.XvdpEc2Volume.Size = cloudformation.Int(192)
	MongoReplicaSetInstance007068.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007068.StopServices = false
	MongoReplicaSetInstance007068.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance007068.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007068.EnableXvdpGp3 = true
	MongoReplicaSetInstance007068.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007068.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007100 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007100.EnableEc2instance = true
	MongoReplicaSetInstance007100.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	MongoReplicaSetInstance007100.Ec2Instance.InstanceType = cloudformation.String("c5.9xlarge")
	MongoReplicaSetInstance007100.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007100.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.7.100") //primary
	MongoReplicaSetInstance007100.XvdpEc2Volume.Size = cloudformation.Int(192)
	MongoReplicaSetInstance007100.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007100.StopServices = false
	MongoReplicaSetInstance007100.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	MongoReplicaSetInstance007100.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007100.EnableXvdpGp3 = true
	MongoReplicaSetInstance007100.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007100.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/in1/in1-Mongo-Accounts-1", "in1-Mongo-Accounts-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/in1/in1-Mongo-Accounts-1", "in1-Mongo-Accounts-1-Service.json")
}
