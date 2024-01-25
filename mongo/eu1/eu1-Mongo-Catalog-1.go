package eu1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateEu1MongoCatalogTemplate() {
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
		Ecc2SubnetLogicalId:    "MongoCatalogReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.144/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.11.7.144 - 10.11.7.159
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoCatalogReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.160/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.11.7.160 - 10.11.7.175
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoCatalogReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.176/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.11.7.176 - 10.11.7.191
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-06ec1aa1d318a028f")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("eu1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{
		"--dbpath",
		"/var/lib/mongo/data",
		"--replSet",
		"catalog-rs0",
		"--logappend",
		"--auth",
		"--oplogSize",
		"51200",
		"--keyFile",
		"/var/lib/mongodb-keyfile",
		"--journal",
		"--directoryperdb",
		"--storageEngine",
		"wiredTiger",
		"--wiredTigerEngineConfigString",
		"file_manager=(close_handle_minimum=100,close_idle_time=30,close_scan_interval=30)",
		"--port",
		"27017",
		"--bind_ip_all",
		"--setParameter",
		"enableLocalhostAuthBypass=true",
		"--logpath",
		"/var/log/mongodb/mongod.log",
	}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007150 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007150.EnableEc2instance = true
	MongoReplicaSetInstance007150.Ec2Instance.ImageId = cloudformation.String("ami-0d7a65c5a518a12c3")
	MongoReplicaSetInstance007150.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	MongoReplicaSetInstance007150.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007150.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.7.150") //primary
	MongoReplicaSetInstance007150.XvdpEc2Volume.Size = cloudformation.Int(384)
	MongoReplicaSetInstance007150.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007150.EnableXvdpGp3 = true
	MongoReplicaSetInstance007150.StopServices = false
	MongoReplicaSetInstance007150.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007150.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007150.MongoContainerTag = "bamboo-mongo-task-2652-mongo-9"
	MongoReplicaSetInstance007150.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007165 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007165.EnableEc2instance = true
	MongoReplicaSetInstance007165.Ec2Instance.ImageId = cloudformation.String("ami-0d7a65c5a518a12c3")
	MongoReplicaSetInstance007165.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	MongoReplicaSetInstance007165.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007165.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.7.165")
	MongoReplicaSetInstance007165.XvdpEc2Volume.Size = cloudformation.Int(384)
	MongoReplicaSetInstance007165.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007165.EnableXvdpGp3 = true
	MongoReplicaSetInstance007165.StopServices = false
	MongoReplicaSetInstance007165.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007165.MongoContainerTag = "github-task-2652-mongo-17"
	MongoReplicaSetInstance007165.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007182 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007182.EnableEc2instance = true
	MongoReplicaSetInstance007182.Ec2Instance.ImageId = cloudformation.String("ami-0d7a65c5a518a12c3")
	MongoReplicaSetInstance007182.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	MongoReplicaSetInstance007182.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007182.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.7.182")
	MongoReplicaSetInstance007182.XvdpEc2Volume.Size = cloudformation.Int(384)
	MongoReplicaSetInstance007182.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007182.EnableXvdpGp3 = true
	MongoReplicaSetInstance007182.StopServices = false
	MongoReplicaSetInstance007182.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007182.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007182.MongoContainerTag = "github-task-2652-mongo-17"
	MongoReplicaSetInstance007182.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/eu1/Mongo-Catalog-1", "eu1-Mongo-Catalog-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/eu1/Mongo-Catalog-1", "eu1-Mongo-Catalog-1-Service.json")
}
