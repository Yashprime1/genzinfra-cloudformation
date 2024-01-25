package sg1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSg1MongoDelivery1Template() {
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
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.160/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.7.160 - 10.15.7.175
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.176/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.14.7.176 - 10.14.7.191
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.192/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.14.7.192 - 10.14.7.107
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(130)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-038ad0e960f5a6d6c")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableMongoLogger = true
	defaults.StackName = "sg1-Mongo-Delivery-1"
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("sg1-SecurityGroup-MongoInstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "sg1-mongo-delivery-1-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007164 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007164.EnableEc2instance = true
	MongoReplicaSetInstance007164.Ec2Instance.ImageId = cloudformation.String("ami-0a98453160fc9f10a")
	MongoReplicaSetInstance007164.Ec2Instance.InstanceType = cloudformation.String("t3.small")
	MongoReplicaSetInstance007164.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007164.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.7.164") // primary
	MongoReplicaSetInstance007164.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance007164.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007164.EnableXvdpGp3 = true
	MongoReplicaSetInstance007164.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007164.StopServices = false
	MongoReplicaSetInstance007164.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007164.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance007164.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007180 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007180.EnableEc2instance = true
	MongoReplicaSetInstance007180.Ec2Instance.ImageId = cloudformation.String("ami-0a98453160fc9f10a")
	MongoReplicaSetInstance007180.Ec2Instance.InstanceType = cloudformation.String("t3.small")
	MongoReplicaSetInstance007180.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007180.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.7.180")
	MongoReplicaSetInstance007180.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance007180.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007180.EnableXvdpGp3 = true
	MongoReplicaSetInstance007180.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007180.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance007180.StopServices = false
	MongoReplicaSetInstance007180.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007180.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007196 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007196.EnableEc2instance = true
	MongoReplicaSetInstance007196.Ec2Instance.ImageId = cloudformation.String("ami-0a98453160fc9f10a")
	MongoReplicaSetInstance007196.Ec2Instance.InstanceType = cloudformation.String("t3.small")
	MongoReplicaSetInstance007196.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007196.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.7.196")
	MongoReplicaSetInstance007196.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance007196.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007196.EnableXvdpGp3 = true
	MongoReplicaSetInstance007196.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007196.StopServices = false
	MongoReplicaSetInstance007196.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007196.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance007196.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sg1/sg1-Mongo-Delivery", "sg1-Mongo-Delivery-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sg1/sg1-Mongo-Delivery", "sg1-Mongo-Delivery-1-Service.json")
}
