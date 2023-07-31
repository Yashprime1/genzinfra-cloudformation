package eu1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateEu1MongoLegacyAccounts2Template() {
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
		SubnetCidrBlockSuffix:  "7.192/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.15.7.48 - 10.15.7.64
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoLegacyAccountsReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.208/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.49 - 10.14.6.62
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoLegacyAccountsReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.224/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.65 - 10.14.6.78
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(704)
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("eu1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "accounts-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007200 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007200.EnableEc2instance = true
	MongoReplicaSetInstance007200.Ec2Instance.ImageId = cloudformation.String("ami-0f572ed10ef61301e")
	MongoReplicaSetInstance007200.Ec2Instance.InstanceType = cloudformation.String("r5.8xlarge")
	MongoReplicaSetInstance007200.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007200.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.7.200") //Primary
	MongoReplicaSetInstance007200.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007200.EnableXvdpGp3 = true
	MongoReplicaSetInstance007200.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007200.MongoContainerTag = "bamboo-mongo-master-9"
	MongoReplicaSetInstance007200.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007215 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007215.EnableEc2instance = true
	MongoReplicaSetInstance007215.Ec2Instance.ImageId = cloudformation.String("ami-0251986887b4fb951")
	MongoReplicaSetInstance007215.Ec2Instance.InstanceType = cloudformation.String("r5.8xlarge")
	MongoReplicaSetInstance007215.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007215.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.7.215")
	MongoReplicaSetInstance007215.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007215.StopServices = false
	MongoReplicaSetInstance007215.EnableXvdpGp3 = true
	MongoReplicaSetInstance007215.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007215.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007215.MongoContainerTag = "bamboo-mongo-master-15"
	MongoReplicaSetInstance007215.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007230 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007230.EnableEc2instance = true
	MongoReplicaSetInstance007230.Ec2Instance.ImageId = cloudformation.String("ami-0251986887b4fb951")
	MongoReplicaSetInstance007230.Ec2Instance.InstanceType = cloudformation.String("r5.8xlarge")
	MongoReplicaSetInstance007230.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007230.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.7.230")
	MongoReplicaSetInstance007230.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007230.StopServices = false
	MongoReplicaSetInstance007230.EnableXvdpGp3 = true
	MongoReplicaSetInstance007230.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007230.XvdpEc2Volume.Iops = cloudformation.Int(5000)
	MongoReplicaSetInstance007230.MongoContainerTag = "bamboo-mongo-master-15"
	MongoReplicaSetInstance007230.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/eu1/eu1-Mongo-Legacy-Accounts", "eu1-Mongo-Legacy-Accounts-2.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/eu1/eu1-Mongo-Legacy-Accounts", "eu1-Mongo-Legacy-Accounts-2-Service.json")
}
