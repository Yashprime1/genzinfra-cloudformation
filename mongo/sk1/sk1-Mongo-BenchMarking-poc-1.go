package sk1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSk1MongoBenchMarkingPocTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("sk1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("sk1")

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("sk1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{
		"--dbpath",
		"/var/lib/mongo/data",
		"--logappend",
		"--auth",
		"--oplogSize",
		"51200",
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
	SubnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "c",
		SubnetCidrBlockSuffix:  "7.112/28",
	})
	SubnetC.AppendToTemplate(sTemplate)
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007231 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007231.EnableEc2instance = false
	MongoReplicaSetInstance007231.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance007231.Ec2Instance.InstanceType = cloudformation.String("r5.2xlarge")
	MongoReplicaSetInstance007231.Ec2InstanceSubnet = SubnetC
	MongoReplicaSetInstance007231.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.7.116")
	MongoReplicaSetInstance007231.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007231.StopServices = false
	MongoReplicaSetInstance007231.EnableXvdpGp3 = true
	MongoReplicaSetInstance007231.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007231.MongoContainerTag = "bamboo-mongo-task-2652-mongo-9"
	MongoReplicaSetInstance007231.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sk1/Mongo-Benchmarking-poc-1", "sk1-Mongo-Benchmarking-poc-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sk1/Mongo-Benchmarking-poc-1", "sk1-Mongo-Benchmarking-poc-1-Service.json")
}
