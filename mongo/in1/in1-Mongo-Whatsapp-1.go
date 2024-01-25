package in1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateIn1MongoWhatsappTemplate() {
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
		SubnetCidrBlockSuffix:  "6.0/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "in1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "6.16/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetB.AppendToTemplate(sTemplate)
	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "in1",
		AvailabilityZoneSuffix: "c",
		SubnetCidrBlockSuffix:  "13.16/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-041db52f9e857be3c")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("in1-SecurityGroup-MongoWhatsappEc2InstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "wa-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance006005 := mongo.NewMongo(defaults)
	mongoReplicaInstance006005.EnableEc2instance = true
	mongoReplicaInstance006005.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	mongoReplicaInstance006005.Ec2Instance.InstanceType = cloudformation.String("r5.4xlarge")
	mongoReplicaInstance006005.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006005.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.6.5")
	mongoReplicaInstance006005.XvdpEc2Volume.Size = cloudformation.Int(1664)
	mongoReplicaInstance006005.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006005.EnableXvdpGp3 = true
	mongoReplicaInstance006005.StopServices = false
	mongoReplicaInstance006005.XvdpEc2Volume.Iops = cloudformation.Int(9000)
	mongoReplicaInstance006005.MongoContainerTag = "github-task-2652-mongo-17"
	mongoReplicaInstance006005.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006006 := mongo.NewMongo(defaults)
	mongoReplicaInstance006006.EnableEc2instance = true
	mongoReplicaInstance006006.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	mongoReplicaInstance006006.Ec2Instance.InstanceType = cloudformation.String("r5.4xlarge")
	mongoReplicaInstance006006.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006006.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.6.6")
	mongoReplicaInstance006006.XvdpEc2Volume.Size = cloudformation.Int(1664)
	mongoReplicaInstance006006.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006006.EnableMongoRegistryCache = true
	mongoReplicaInstance006006.StopServices = false
	mongoReplicaInstance006006.EnableXvdpGp3 = true
	mongoReplicaInstance006006.XvdpEc2Volume.Iops = cloudformation.Int(9000)
	mongoReplicaInstance006006.MongoContainerTag = "github-task-2652-mongo-17"
	mongoReplicaInstance006006.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006020 := mongo.NewMongo(defaults)
	mongoReplicaInstance006020.EnableEc2instance = true
	mongoReplicaInstance006020.Ec2Instance.ImageId = cloudformation.String("ami-03f1bfde83c761b99")
	mongoReplicaInstance006020.Ec2Instance.InstanceType = cloudformation.String("r5.4xlarge")
	mongoReplicaInstance006020.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006020.Ec2Instance.PrivateIpAddress = cloudformation.String("10.12.6.20") //primary
	mongoReplicaInstance006020.XvdpEc2Volume.Size = cloudformation.Int(1664)
	mongoReplicaInstance006020.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006020.StopServices = false
	mongoReplicaInstance006020.EnableXvdpGp3 = true
	mongoReplicaInstance006020.XvdpEc2Volume.Iops = cloudformation.Int(9000)
	mongoReplicaInstance006020.MongoContainerTag = "bamboo-mongo-task-2652-mongo-9"
	mongoReplicaInstance006020.EnableMongoRegistryCache = true
	mongoReplicaInstance006020.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/in1/Mongo-Whatsapp", "in1-Mongo-Whatsapp-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/in1/Mongo-Whatsapp", "in1-Mongo-Whatsapp-1-Service.json")
}
