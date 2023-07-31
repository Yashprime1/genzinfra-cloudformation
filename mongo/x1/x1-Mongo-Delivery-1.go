package x1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateX1MongoDelivery1Template() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKey()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("x1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("x1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "x1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.160/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.7.160 - 10.16.7.175
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "x1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.176/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.7.176 - 10.16.7.191
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "x1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.192/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.7.192 - 10.16.7.207
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(130)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-00ac255069dc94e3c")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("x1-SecurityGroup-MongoInstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "x1-Mongo-Delivery-1-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007164 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007164.EnableEc2instance = true
	MongoReplicaSetInstance007164.Ec2Instance.ImageId = cloudformation.String("ami-00ac255069dc94e3c")
	MongoReplicaSetInstance007164.Ec2Instance.InstanceType = cloudformation.String("t2.small")
	MongoReplicaSetInstance007164.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007164.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.7.164")
	MongoReplicaSetInstance007164.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance007164.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007164.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance007164.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007180 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007180.EnableEc2instance = true
	MongoReplicaSetInstance007180.Ec2Instance.ImageId = cloudformation.String("ami-00ac255069dc94e3c")
	MongoReplicaSetInstance007180.Ec2Instance.InstanceType = cloudformation.String("t2.small")
	MongoReplicaSetInstance007180.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007180.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.7.180") //primary
	MongoReplicaSetInstance007180.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance007180.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007180.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance007180.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007196 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007196.EnableEc2instance = true
	MongoReplicaSetInstance007196.Ec2Instance.ImageId = cloudformation.String("ami-00ac255069dc94e3c")
	MongoReplicaSetInstance007196.Ec2Instance.InstanceType = cloudformation.String("t2.small")
	MongoReplicaSetInstance007196.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007196.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.7.196")
	MongoReplicaSetInstance007196.XvdpEc2Volume.Size = cloudformation.Int(130)
	MongoReplicaSetInstance007196.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007196.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance007196.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/x1/x1-Mongo-Delivery", "x1-Mongo-Delivery-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/x1/x1-Mongo-Delivery", "x1-Mongo-Delivery-1-Service.json")
}
