package us1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateUs1MongoLinkShortnerTemplate() {
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
		SubnetCidrBlockSuffix:  "6.128/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "us1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "6.224/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "us1",
		AvailabilityZoneSuffix: "c",
		SubnetCidrBlockSuffix:  "6.144/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-083a57bca71fb2945")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("us1-SecurityGroup-MongoLinkShortenerEc2InstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "us1-mongo-link-shortener-1-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance006133 := mongo.NewMongo(defaults)
	mongoReplicaInstance006133.EnableEc2instance = true
	mongoReplicaInstance006133.Ec2Instance.ImageId = cloudformation.String("ami-02726fee3f464392b")
	mongoReplicaInstance006133.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	mongoReplicaInstance006133.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006133.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.6.133") //primary
	mongoReplicaInstance006133.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006133.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006133.StopServices = false
	mongoReplicaInstance006133.EnableMongoRegistryCache = true
	mongoReplicaInstance006133.EnableXvdpGp3 = true
	mongoReplicaInstance006133.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance006133.MongoContainerTag = "bamboo-mongo-task-2652-mongo-9"
	mongoReplicaInstance006133.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006229 := mongo.NewMongo(defaults)
	mongoReplicaInstance006229.EnableEc2instance = true
	mongoReplicaInstance006229.Ec2Instance.ImageId = cloudformation.String("ami-02726fee3f464392b")
	mongoReplicaInstance006229.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	mongoReplicaInstance006229.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006229.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.6.229") 
	mongoReplicaInstance006229.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006229.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006229.StopServices = true
	mongoReplicaInstance006229.EnableMongoRegistryCache = true
	mongoReplicaInstance006229.EnableXvdpGp3 = true
	mongoReplicaInstance006229.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance006229.MongoContainerTag = "bamboo-mongo-task-2652-mongo-9"
	mongoReplicaInstance006229.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006149 := mongo.NewMongo(defaults)
	mongoReplicaInstance006149.EnableEc2instance = true
	mongoReplicaInstance006149.Ec2Instance.ImageId = cloudformation.String("ami-02726fee3f464392b")
	mongoReplicaInstance006149.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	mongoReplicaInstance006149.Ec2InstanceSubnet = subnetC
	mongoReplicaInstance006149.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.6.149")
	mongoReplicaInstance006149.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006149.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006149.StopServices = false
	mongoReplicaInstance006149.EnableMongoRegistryCache = true
	mongoReplicaInstance006149.EnableXvdpGp3 = true
	mongoReplicaInstance006149.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance006149.MongoContainerTag = "bamboo-mongo-task-2652-mongo-9"
	mongoReplicaInstance006149.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/us1/Mongo-Link-Shortener", "us1-Mongo-Link-Shortener-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/us1/Mongo-Link-Shortener", "us1-Mongo-Link-Shortener-1-Service.json")
}
