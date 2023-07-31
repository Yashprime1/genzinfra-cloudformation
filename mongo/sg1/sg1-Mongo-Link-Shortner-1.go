package sg1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSg1MongoLinkShortnerTemplate() {
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
		SubnetCidrBlockSuffix:  "7.128/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "7.112/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "c",
		SubnetCidrBlockSuffix:  "7.144/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0d85a72517fc05f6b")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("sg1-SecurityGroup-MongoLinkShortenerEc2InstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "sg1-mongo-link-shortener-1-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance007133 := mongo.NewMongo(defaults)
	mongoReplicaInstance007133.EnableEc2instance = true
	mongoReplicaInstance007133.Ec2Instance.ImageId = cloudformation.String("ami-02df2edff7704b9fe")
	mongoReplicaInstance007133.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	mongoReplicaInstance007133.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance007133.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.7.133")
	mongoReplicaInstance007133.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance007133.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance007133.EnableXvdpGp3 = true
	mongoReplicaInstance007133.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance007133.StopServices = false
	mongoReplicaInstance007133.EnableMongoRegistryCache = true
	mongoReplicaInstance007133.MongoContainerTag = "bamboo-mongo-task-2652-mongo-1"
	mongoReplicaInstance007133.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance007117 := mongo.NewMongo(defaults)
	mongoReplicaInstance007117.EnableEc2instance = true
	mongoReplicaInstance007117.Ec2Instance.ImageId = cloudformation.String("ami-02df2edff7704b9fe")
	mongoReplicaInstance007117.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	mongoReplicaInstance007117.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance007117.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.7.117")
	mongoReplicaInstance007117.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance007117.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance007117.EnableXvdpGp3 = true
	mongoReplicaInstance007117.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance007117.StopServices = false
	mongoReplicaInstance007117.EnableMongoRegistryCache = true
	mongoReplicaInstance007117.MongoContainerTag = "bamboo-mongo-task-2652-mongo-1"
	mongoReplicaInstance007117.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance007149 := mongo.NewMongo(defaults)
	mongoReplicaInstance007149.EnableEc2instance = true
	mongoReplicaInstance007149.Ec2Instance.ImageId = cloudformation.String("ami-02df2edff7704b9fe")
	mongoReplicaInstance007149.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	mongoReplicaInstance007149.Ec2InstanceSubnet = subnetC
	mongoReplicaInstance007149.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.7.149") //primary
	mongoReplicaInstance007149.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance007149.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance007149.EnableXvdpGp3 = true
	mongoReplicaInstance007149.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance007149.MongoContainerTag = "bamboo-mongo-task-2652-mongo-1"
	mongoReplicaInstance007149.StopServices = false
	mongoReplicaInstance007149.EnableMongoRegistryCache = true
	mongoReplicaInstance007149.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sg1/Mongo-Link-Shortener", "sg1-Mongo-Link-Shortener-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sg1/Mongo-Link-Shortener", "sg1-Mongo-Link-Shortener-1-Service.json")
}
