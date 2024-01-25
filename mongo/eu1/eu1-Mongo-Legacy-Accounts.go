package eu1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateEu1MongoLegacyAccountsTemplate() {
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
		Ecc2SubnetLogicalId:    "MongoLegacyAccountsReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.80/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.15.7.48 - 10.15.7.64
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoLegacyAccountsReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.112/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.49 - 10.14.6.62
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoLegacyAccountsReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.128/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.65 - 10.14.6.78
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("eu1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "accounts-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007085 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007085.EnableEc2instance = true
	MongoReplicaSetInstance007085.Ec2Instance.ImageId = cloudformation.String("ami-0d7a65c5a518a12c3")
	MongoReplicaSetInstance007085.Ec2Instance.InstanceType = cloudformation.String("r5.4xlarge")
	MongoReplicaSetInstance007085.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007085.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.7.85")
	MongoReplicaSetInstance007085.XvdpEc2Volume.Size = cloudformation.Int(1536)
	MongoReplicaSetInstance007085.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007085.EnableXvdpGp3 = true
	MongoReplicaSetInstance007085.StopServices = false
	MongoReplicaSetInstance007085.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007085.MongoContainerTag = "github-master-18"
	MongoReplicaSetInstance007085.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007117 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007117.EnableEc2instance = true
	MongoReplicaSetInstance007117.Ec2Instance.ImageId = cloudformation.String("ami-0d7a65c5a518a12c3")
	MongoReplicaSetInstance007117.Ec2Instance.InstanceType = cloudformation.String("r5.4xlarge")
	MongoReplicaSetInstance007117.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007117.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.7.117") //primary
	MongoReplicaSetInstance007117.XvdpEc2Volume.Size = cloudformation.Int(1536)
	MongoReplicaSetInstance007117.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007117.StopServices = false
	MongoReplicaSetInstance007117.EnableXvdpGp3 = true
	MongoReplicaSetInstance007117.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007117.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007117.MongoContainerTag = "bamboo-mongo-master-19"
	MongoReplicaSetInstance007117.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007133 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007133.EnableEc2instance = true
	MongoReplicaSetInstance007133.Ec2Instance.ImageId = cloudformation.String("ami-0d7a65c5a518a12c3")
	MongoReplicaSetInstance007133.Ec2Instance.InstanceType = cloudformation.String("r5.4xlarge")
	MongoReplicaSetInstance007133.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007133.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.7.133")
	MongoReplicaSetInstance007133.XvdpEc2Volume.Size = cloudformation.Int(1536)
	MongoReplicaSetInstance007133.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007133.StopServices = false
	MongoReplicaSetInstance007133.EnableXvdpGp3 = true
	MongoReplicaSetInstance007133.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007133.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007133.MongoContainerTag = "github-master-18"
	MongoReplicaSetInstance007133.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/eu1/eu1-Mongo-Legacy-Accounts", "eu1-Mongo-Legacy-Accounts.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/eu1/eu1-Mongo-Legacy-Accounts", "eu1-Mongo-Legacy-Accounts-Service.json")
}
