package m1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func Generatem1MongoStats1Template() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("m1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("m1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "m1",
		AvailabilityZoneSuffix: "a",
		SubnetCidrBlockSuffix:  "6.32/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "m1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "6.48/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0f8afc87a647567dd")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("m1-SecurityGroup-MongoCommonEc2InstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{
		"--dbpath",
		"/var/lib/mongo/data",
		"--replSet",
		"stats-rs0",
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

	mongoReplicaInstance006037 := mongo.NewMongo(defaults)
	mongoReplicaInstance006037.EnableEc2instance = false
	mongoReplicaInstance006037.Ec2Instance.InstanceType = cloudformation.String("c5.12xlarge")
	mongoReplicaInstance006037.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006037.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.37")
	mongoReplicaInstance006037.XvdpEc2Volume.Size = cloudformation.Int(320)
	mongoReplicaInstance006037.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006037.StopServices = false
	mongoReplicaInstance006037.EnableMongoRegistryCache = true
	mongoReplicaInstance006037.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	mongoReplicaInstance006037.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006038 := mongo.NewMongo(defaults)
	mongoReplicaInstance006038.EnableEc2instance = false
	mongoReplicaInstance006038.Ec2Instance.InstanceType = cloudformation.String("c5.12xlarge")
	mongoReplicaInstance006038.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006038.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.38")
	mongoReplicaInstance006038.XvdpEc2Volume.Size = cloudformation.Int(320)
	mongoReplicaInstance006038.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006038.StopServices = false
	mongoReplicaInstance006038.EnableMongoRegistryCache = true
	mongoReplicaInstance006038.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	mongoReplicaInstance006038.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006053 := mongo.NewMongo(defaults)
	mongoReplicaInstance006053.EnableEc2instance = false
	mongoReplicaInstance006053.Ec2Instance.InstanceType = cloudformation.String("c5.12xlarge")
	mongoReplicaInstance006053.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006053.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.53")
	mongoReplicaInstance006053.XvdpEc2Volume.Size = cloudformation.Int(320)
	mongoReplicaInstance006053.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006053.StopServices = false
	mongoReplicaInstance006053.EnableMongoRegistryCache = true
	mongoReplicaInstance006053.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	mongoReplicaInstance006053.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/m1/Mongo-Stats", "m1-Mongo-Stats-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/m1/Mongo-Stats", "m1-Mongo-Stats-1-Service.json")
}
