package eu2Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateEu2MongoWhatsappTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("eu2")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("eu2")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu2",
		AvailabilityZoneSuffix: "a",
		SubnetCidrBlockSuffix:  "6.0/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu2",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "6.16/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0f572ed10ef61301e")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("eu2-SecurityGroup-MongoWhatsappEc2InstanceEc2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "wa-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance006005 := mongo.NewMongo(defaults)
	mongoReplicaInstance006005.EnableEc2instance = false
	mongoReplicaInstance006005.Ec2Instance.ImageId = cloudformation.String("ami-0f572ed10ef61301e")
	mongoReplicaInstance006005.Ec2Instance.InstanceType = cloudformation.String("r4.xlarge")
	mongoReplicaInstance006005.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006005.Ec2Instance.PrivateIpAddress = cloudformation.String("10.13.6.5")
	mongoReplicaInstance006005.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006005.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006005.MongoContainerTag = "bamboo-mongo-task-2652-mongo-8"
	mongoReplicaInstance006005.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006006 := mongo.NewMongo(defaults)
	mongoReplicaInstance006006.EnableEc2instance = false
	mongoReplicaInstance006006.Ec2Instance.ImageId = cloudformation.String("ami-0f572ed10ef61301e")
	mongoReplicaInstance006006.Ec2Instance.InstanceType = cloudformation.String("r4.xlarge")
	mongoReplicaInstance006006.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006006.Ec2Instance.PrivateIpAddress = cloudformation.String("10.13.6.6")
	mongoReplicaInstance006006.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006006.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006006.MongoContainerTag = "bamboo-mongo-task-2652-mongo-8"
	mongoReplicaInstance006006.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006020 := mongo.NewMongo(defaults)
	mongoReplicaInstance006020.EnableEc2instance = false
	mongoReplicaInstance006020.Ec2Instance.ImageId = cloudformation.String("ami-0f572ed10ef61301e")
	mongoReplicaInstance006020.Ec2Instance.InstanceType = cloudformation.String("r4.xlarge")
	mongoReplicaInstance006020.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006020.Ec2Instance.PrivateIpAddress = cloudformation.String("10.13.6.20")
	mongoReplicaInstance006020.XvdpEc2Volume.Size = cloudformation.Int(64)
	mongoReplicaInstance006020.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006020.MongoContainerTag = "bamboo-mongo-task-2652-mongo-8"
	mongoReplicaInstance006020.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/eu2/Mongo-Whatsapp", "eu2-Mongo-Whatsapp-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/eu2/Mongo-Whatsapp", "eu2-Mongo-Whatsapp-1-Service.json")
}
