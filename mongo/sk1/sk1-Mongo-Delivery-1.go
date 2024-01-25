package sk1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSk1MongoDeliveryTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEbsDlmLifecyclePolicy"] = mongo.GetDlmLifeCyclePolicy()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("sk1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("sk1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "8.128/28",
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "8.144/28",
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoDeliveryReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "8.160/28",
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(30)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-002623b6d0b3d66a3")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("sk1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "Delivery-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}

	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance008132 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance008132.EnableEc2instance = true
	MongoReplicaSetInstance008132.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance008132.Ec2Instance.InstanceType = cloudformation.String("t3.micro")
	MongoReplicaSetInstance008132.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance008132.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.8.132") //primary
	MongoReplicaSetInstance008132.XvdpEc2Volume.Size = cloudformation.Int(30)
	MongoReplicaSetInstance008132.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance008132.StopServices = false
	MongoReplicaSetInstance008132.EnableMongoRegistryCache = true
	MongoReplicaSetInstance008132.EnableXvdpGp3 = true
	MongoReplicaSetInstance008132.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance008132.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance008132.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance008148 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance008148.EnableEc2instance = true
	MongoReplicaSetInstance008148.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance008148.Ec2Instance.InstanceType = cloudformation.String("t3.micro")
	MongoReplicaSetInstance008148.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance008148.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.8.148")
	MongoReplicaSetInstance008148.XvdpEc2Volume.Size = cloudformation.Int(30)
	MongoReplicaSetInstance008148.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance008148.StopServices = false
	MongoReplicaSetInstance008148.EnableXvdpGp3 = true
	MongoReplicaSetInstance008148.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance008148.EnableMongoRegistryCache = true
	MongoReplicaSetInstance008148.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance008148.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance008164 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance008164.EnableEc2instance = true
	MongoReplicaSetInstance008164.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance008164.Ec2Instance.InstanceType = cloudformation.String("t3.micro")
	MongoReplicaSetInstance008164.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance008164.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.8.164")
	MongoReplicaSetInstance008164.XvdpEc2Volume.Size = cloudformation.Int(30)
	MongoReplicaSetInstance008164.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance008164.StopServices = false
	MongoReplicaSetInstance008164.EnableXvdpGp3 = true
	MongoReplicaSetInstance008164.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance008164.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance008164.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sk1/Mongo-Delivery-1", "sk1-Mongo-Delivery-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sk1/Mongo-Delivery-1", "sk1-Mongo-Delivery-1-Service.json")
}
