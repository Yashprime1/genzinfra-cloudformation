package sk1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSk1MongoWhatsappTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("sk1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("sk1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "a",
		SubnetCidrBlockSuffix:  "6.0/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})

	subnetA.AppendToTemplate(sTemplate)
	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "13.16/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetB.AppendToTemplate(sTemplate)
	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "c",
		SubnetCidrBlockSuffix:  "6.16/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0a91977023aa59ed0")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("sk1-SecurityGroup-MongoWhatsappEc2InstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "wa-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance006005 := mongo.NewMongo(defaults)
	mongoReplicaInstance006005.EnableEc2instance = true
	mongoReplicaInstance006005.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	mongoReplicaInstance006005.Ec2Instance.InstanceType = cloudformation.String("t2.large")
	mongoReplicaInstance006005.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006005.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.6.5") //primary
	mongoReplicaInstance006005.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006005.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006005.StopServices = false
	mongoReplicaInstance006005.EnableMongoRegistryCache = true
	mongoReplicaInstance006005.EnableXvdpGp3 = true
	mongoReplicaInstance006005.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance006005.MongoContainerTag = "bamboo-mongo-task-2652-mongo-9"
	mongoReplicaInstance006005.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006006 := mongo.NewMongo(defaults)
	mongoReplicaInstance006006.EnableEc2instance = true
	mongoReplicaInstance006006.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	mongoReplicaInstance006006.Ec2Instance.InstanceType = cloudformation.String("t2.large")
	mongoReplicaInstance006006.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006006.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.6.6")
	mongoReplicaInstance006006.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006006.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006006.StopServices = false
	mongoReplicaInstance006006.EnableXvdpGp3 = true
	mongoReplicaInstance006006.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance006006.MongoContainerTag = "bamboo-mongo-task-2652-mongo-9"
	mongoReplicaInstance006006.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006020 := mongo.NewMongo(defaults)
	mongoReplicaInstance006020.EnableEc2instance = true
	mongoReplicaInstance006020.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	mongoReplicaInstance006020.Ec2Instance.InstanceType = cloudformation.String("t2.large")
	mongoReplicaInstance006020.Ec2InstanceSubnet = subnetC
	mongoReplicaInstance006020.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.6.20")
	mongoReplicaInstance006020.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006020.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006020.StopServices = false
	mongoReplicaInstance006020.EnableMongoRegistryCache = true
	mongoReplicaInstance006020.EnableXvdpGp3 = true
	mongoReplicaInstance006020.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance006020.MongoContainerTag = "bamboo-mongo-task-2652-mongo-9"
	mongoReplicaInstance006020.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance013020 := mongo.NewMongo(defaults)
	mongoReplicaInstance013020.EnableEc2instance = true
	mongoReplicaInstance013020.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	mongoReplicaInstance013020.Ec2Instance.InstanceType = cloudformation.String("t3.large")
	mongoReplicaInstance013020.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance013020.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.13.20")
	mongoReplicaInstance013020.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance013020.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance013020.EnableMongoRegistryCache = true
	mongoReplicaInstance013020.StopServices = false
	mongoReplicaInstance013020.EnableXvdpGp3 = true
	mongoReplicaInstance013020.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance013020.MongoContainerTag = "bamboo-mongo-task-2652-mongo-9"
	mongoReplicaInstance013020.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sk1/Mongo-Whatsapp", "sk1-Mongo-Whatsapp-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sk1/Mongo-Whatsapp", "sk1-Mongo-Whatsapp-1-Service.json")
}