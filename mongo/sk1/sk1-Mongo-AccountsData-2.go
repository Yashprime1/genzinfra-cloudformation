package sk1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSk1MongoWiredTigerAccountsData24Dot2Template() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEbsDlmLifecyclePolicy"] = mongo.GetDlmLifeCyclePolicy()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("sk1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("sk1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoAccountsData2ReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.192/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.14.7.144 - 10.14.7.159
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoAccountsData2ReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.208/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.14.7.160 - 10.14.7.175
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoAccountsData2ReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.224/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.14.7.176 - 10.14.7.191
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0a91977023aa59ed0")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("sk1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{
		"--dbpath",
		"/var/lib/mongo/data",
		"--replSet",
		"accounts-rs0",
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

	MongoReplicaSetInstance007197 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007197.EnableEc2instance = true
	MongoReplicaSetInstance007197.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance007197.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance007197.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007197.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.7.197") //primary
	MongoReplicaSetInstance007197.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007197.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007197.StopServices = false
	MongoReplicaSetInstance007197.EnableXvdpGp3 = true
	MongoReplicaSetInstance007197.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007197.MongoContainerTag = "github-task-2652-mongo-17"
	MongoReplicaSetInstance007197.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007197.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007213 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007213.EnableEc2instance = true
	MongoReplicaSetInstance007213.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance007213.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance007213.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007213.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.7.213") //primary
	MongoReplicaSetInstance007213.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007213.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007213.StopServices = false
	MongoReplicaSetInstance007213.EnableXvdpGp3 = true
	MongoReplicaSetInstance007213.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007213.MongoContainerTag = "github-task-2652-mongo-17"
	MongoReplicaSetInstance007213.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007213.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007229 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007229.EnableEc2instance = true
	MongoReplicaSetInstance007229.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance007229.Ec2Instance.InstanceType = cloudformation.String("c5.large")
	MongoReplicaSetInstance007229.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007229.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.7.229")
	MongoReplicaSetInstance007229.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007229.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007229.StopServices = false
	MongoReplicaSetInstance007229.EnableXvdpGp3 = true
	MongoReplicaSetInstance007229.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007229.MongoContainerTag = "github-task-2652-mongo-17"
	MongoReplicaSetInstance007229.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007229.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sk1/Mongo-AccountsData-2", "sk1-Mongo-AccountsData-2.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sk1/Mongo-AccountsData-2", "sk1-Mongo-AccountsData-2-Service.json")
}
