package eu1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateEu1MongoAuditTemplate() {
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
		Ecc2SubnetLogicalId:    "MongoAuditReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "10.0/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.15.7.48 - 10.15.7.64
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoAuditReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "10.16/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.49 - 10.14.6.62
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoAuditReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "10.32/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.65 - 10.14.6.78
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(512)
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("eu1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{
		"--dbpath",
		"/var/lib/mongo/data",
		"--replSet",
		"audit-rs0",
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

	MongoReplicaSetInstance010005 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance010005.EnableEc2instance = true
	MongoReplicaSetInstance010005.Ec2Instance.ImageId = cloudformation.String("ami-0f572ed10ef61301e")
	MongoReplicaSetInstance010005.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	MongoReplicaSetInstance010005.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance010005.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.10.5") //primary
	MongoReplicaSetInstance010005.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance010005.EnableXvdpGp3 = true
	MongoReplicaSetInstance010005.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance010005.MongoContainerTag = "bamboo-mongo-task-2652-mongo-1"
	MongoReplicaSetInstance010005.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance010021 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance010021.EnableEc2instance = true
	MongoReplicaSetInstance010021.Ec2Instance.ImageId = cloudformation.String("ami-0251986887b4fb951")
	MongoReplicaSetInstance010021.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	MongoReplicaSetInstance010021.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance010021.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.10.21")
	MongoReplicaSetInstance010021.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance010021.StopServices = false
	MongoReplicaSetInstance010021.EnableXvdpGp3 = true
	MongoReplicaSetInstance010021.EnableMongoRegistryCache = true
	MongoReplicaSetInstance010021.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance010021.MongoContainerTag = "bamboo-mongo-task-2652-mongo-1"
	MongoReplicaSetInstance010021.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance010037 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance010037.EnableEc2instance = true
	MongoReplicaSetInstance010037.Ec2Instance.ImageId = cloudformation.String("ami-0251986887b4fb951")
	MongoReplicaSetInstance010037.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	MongoReplicaSetInstance010037.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance010037.Ec2Instance.PrivateIpAddress = cloudformation.String("10.11.10.37")
	MongoReplicaSetInstance010037.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance010037.StopServices = false
	MongoReplicaSetInstance010037.EnableXvdpGp3 = true
	MongoReplicaSetInstance010037.EnableMongoRegistryCache = true
	MongoReplicaSetInstance010037.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance010037.MongoContainerTag = "bamboo-mongo-task-2652-mongo-1"
	MongoReplicaSetInstance010037.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/eu1/Mongo-Audit-1", "eu1-Mongo-Audit-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/eu1/Mongo-Audit-1", "eu1-Mongo-Audit-1-Service.json")
}
