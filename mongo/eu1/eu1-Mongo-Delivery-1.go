package eu1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateEu1MongoDeliveryTemplate() {
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
		SubnetCidrBlockSuffix:  "3.0/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "3.16/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "c",
		SubnetCidrBlockSuffix:  "3.32/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(625)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0c0475d803735333b")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableMongoLogger = true
	defaults.StackName = "eu1-Mongo-Delivery-1"
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("eu1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "eu1-Mongo-Delivery-1-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance003005 := mongo.NewMongo(defaults)
	mongoReplicaInstance003005.EnableEc2instance = true
	mongoReplicaInstance003005.Ec2Instance.ImageId = cloudformation.String("ami-0d7a65c5a518a12c3")
	mongoReplicaInstance003005.Ec2Instance.InstanceType = cloudformation.String("r5.xlarge")
	mongoReplicaInstance003005.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance003005.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.3.5") //primary
	mongoReplicaInstance003005.XvdpEc2Volume.Size = cloudformation.Int(625)
	mongoReplicaInstance003005.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance003005.EnableXvdpGp3 = true
	mongoReplicaInstance003005.StopServices = false
	mongoReplicaInstance003005.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance003005.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	mongoReplicaInstance003005.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance003021 := mongo.NewMongo(defaults)
	mongoReplicaInstance003021.EnableEc2instance = true
	mongoReplicaInstance003021.Ec2Instance.ImageId = cloudformation.String("ami-0d7a65c5a518a12c3")
	mongoReplicaInstance003021.Ec2Instance.InstanceType = cloudformation.String("r5.xlarge")
	mongoReplicaInstance003021.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance003021.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.3.21")
	mongoReplicaInstance003021.XvdpEc2Volume.Size = cloudformation.Int(625)
	mongoReplicaInstance003021.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance003021.StopServices = false
	mongoReplicaInstance003021.EnableXvdpGp3 = true
	mongoReplicaInstance003021.EnableMongoRegistryCache = true
	mongoReplicaInstance003021.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance003021.MongoContainerTag = "github-sne-6117-16"
	mongoReplicaInstance003021.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance003037 := mongo.NewMongo(defaults)
	mongoReplicaInstance003037.EnableEc2instance = true
	mongoReplicaInstance003037.Ec2Instance.ImageId = cloudformation.String("ami-0d7a65c5a518a12c3")
	mongoReplicaInstance003037.Ec2Instance.InstanceType = cloudformation.String("r5.xlarge")
	mongoReplicaInstance003037.Ec2InstanceSubnet = subnetC
	mongoReplicaInstance003037.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.3.37")
	mongoReplicaInstance003037.XvdpEc2Volume.Size = cloudformation.Int(625)
	mongoReplicaInstance003037.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance003037.StopServices = false
	mongoReplicaInstance003037.EnableXvdpGp3 = true
	mongoReplicaInstance003037.EnableMongoRegistryCache = true
	mongoReplicaInstance003037.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance003037.MongoContainerTag = "github-sne-6117-16"
	mongoReplicaInstance003037.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/eu1/Mongo-Delivery", "eu1-Mongo-Delivery-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/eu1/Mongo-Delivery", "eu1-Mongo-Delivery-1-Service.json")
}
