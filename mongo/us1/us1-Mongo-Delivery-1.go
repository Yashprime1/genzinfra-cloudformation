package us1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateUs1MongoDelivery1Template() {
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
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.160/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.7.160 - 10.16.7.175
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "us1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.176/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.7.176 - 10.16.7.191
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "us1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.192/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.7.192 - 10.16.7.207
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(130)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0245178d1c349ad3b")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("us1-SecurityGroup-MongoInstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "us1-Mongo-Delivery-1-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007164 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007164.EnableEc2instance = true
	MongoReplicaSetInstance007164.Ec2Instance.ImageId = cloudformation.String("ami-02726fee3f464392b")
	MongoReplicaSetInstance007164.Ec2Instance.InstanceType = cloudformation.String("t3.micro")
	MongoReplicaSetInstance007164.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007164.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.7.164")
	MongoReplicaSetInstance007164.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance007164.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007164.StopServices = false
	MongoReplicaSetInstance007164.EnableXvdpGp3 = true
	MongoReplicaSetInstance007164.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007164.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	MongoReplicaSetInstance007164.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007164.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007180 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007180.EnableEc2instance = true
	MongoReplicaSetInstance007180.Ec2Instance.ImageId = cloudformation.String("ami-0030f75d7382838b8")
	MongoReplicaSetInstance007180.Ec2Instance.InstanceType = cloudformation.String("t3.micro")
	MongoReplicaSetInstance007180.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007180.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.7.180") //primary
	MongoReplicaSetInstance007180.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance007180.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007180.StopServices = false
	MongoReplicaSetInstance007180.EnableXvdpGp3 = true
	MongoReplicaSetInstance007180.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007180.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	MongoReplicaSetInstance007180.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007180.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007196 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007196.EnableEc2instance = true
	MongoReplicaSetInstance007196.Ec2Instance.ImageId = cloudformation.String("ami-02726fee3f464392b")
	MongoReplicaSetInstance007196.Ec2Instance.InstanceType = cloudformation.String("t3.micro")
	MongoReplicaSetInstance007196.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007196.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.7.196")
	MongoReplicaSetInstance007196.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance007196.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007196.StopServices = false
	MongoReplicaSetInstance007196.EnableXvdpGp3 = true
	MongoReplicaSetInstance007196.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007196.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	MongoReplicaSetInstance007196.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007196.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/us1/us1-Mongo-Delivery", "us1-Mongo-Delivery-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/us1/us1-Mongo-Delivery", "us1-Mongo-Delivery-1-Service.json")
}
