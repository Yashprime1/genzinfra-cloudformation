package in1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateIn1MongoStatsTemplate() {
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
		SubnetCidrBlockSuffix:  "7.112/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "in1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "7.128/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "in1",
		AvailabilityZoneSuffix: "c",
		SubnetCidrBlockSuffix:  "7.144/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(500)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0de320fb812039314")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("in1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "in1-Mongo-Stats-1-rs0", "--logpath", "/va2r/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance007116 := mongo.NewMongo(defaults)
	mongoReplicaInstance007116.EnableEc2instance = true
	mongoReplicaInstance007116.Ec2Instance.ImageId = cloudformation.String("ami-0de320fb812039314")
	mongoReplicaInstance007116.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	mongoReplicaInstance007116.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance007116.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.7.116")
	mongoReplicaInstance007116.XvdpEc2Volume.Size = cloudformation.Int(2048)
	mongoReplicaInstance007116.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance007116.EnableXvdpGp3 = true
	mongoReplicaInstance007116.StopServices = false
	mongoReplicaInstance007116.XvdpEc2Volume.Iops = cloudformation.Int(8000)
	mongoReplicaInstance007116.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	mongoReplicaInstance007116.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance007132 := mongo.NewMongo(defaults)
	mongoReplicaInstance007132.EnableEc2instance = true
	mongoReplicaInstance007132.Ec2Instance.ImageId = cloudformation.String("ami-076622203c9703088")
	mongoReplicaInstance007132.Ec2Instance.InstanceType = cloudformation.String("m5.2xlarge")
	mongoReplicaInstance007132.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance007132.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.7.132") //primary
	mongoReplicaInstance007132.XvdpEc2Volume.Size = cloudformation.Int(2048)
	mongoReplicaInstance007132.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance007132.EnableMongoRegistryCache = true
	mongoReplicaInstance007132.StopServices = false
	mongoReplicaInstance007132.EnableXvdpGp3 = true
	mongoReplicaInstance007132.XvdpEc2Volume.Iops = cloudformation.Int(8000)
	mongoReplicaInstance007132.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	mongoReplicaInstance007132.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance007148 := mongo.NewMongo(defaults)
	mongoReplicaInstance007148.EnableEc2instance = true
	mongoReplicaInstance007148.Ec2Instance.ImageId = cloudformation.String("ami-076622203c9703088")
	mongoReplicaInstance007148.Ec2Instance.InstanceType = cloudformation.String("m5.2xlarge")
	mongoReplicaInstance007148.Ec2InstanceSubnet = subnetC
	mongoReplicaInstance007148.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.7.148")
	mongoReplicaInstance007148.XvdpEc2Volume.Size = cloudformation.Int(2048)
	mongoReplicaInstance007148.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance007148.StopServices = false
	mongoReplicaInstance007148.EnableXvdpGp3 = true
	mongoReplicaInstance007148.XvdpEc2Volume.Iops = cloudformation.Int(8000)
	mongoReplicaInstance007148.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	mongoReplicaInstance007148.EnableMongoRegistryCache = true
	mongoReplicaInstance007148.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/in1/Mongo-Stats", "in1-Mongo-Stats-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/in1/Mongo-Stats", "in1-Mongo-Stats-1-Service.json")
}
