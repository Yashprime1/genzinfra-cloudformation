package sg1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSg1MongoIntegrationDebuggerTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEbsDlmLifecyclePolicy"] = mongo.GetDlmLifeCyclePolicy()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("sg1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("sg1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoIntegrationDebuggerReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "13.192/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.13.192 - 10.15.13.207
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoIntegrationDebuggerReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "13.208/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.13.208- 10.15.13.223
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoIntegrationDebuggerReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "13.224/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.15.13.224 - 10.15.13.239
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0df43911a5efd456b")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableSplunkPersistentState = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("sg1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "mongo-integration-debugger-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}

	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance013196 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013196.EnableEc2instance = true
	MongoReplicaSetInstance013196.Ec2Instance.ImageId = cloudformation.String("ami-0df43911a5efd456b")
	MongoReplicaSetInstance013196.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance013196.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance013196.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.13.196")
	MongoReplicaSetInstance013196.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013196.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013196.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013196.StopServices = false
	MongoReplicaSetInstance013196.EnableXvdpGp3 = true
	MongoReplicaSetInstance013196.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013196.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance013196.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance013212 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013212.EnableEc2instance = true
	MongoReplicaSetInstance013212.Ec2Instance.ImageId = cloudformation.String("ami-0df43911a5efd456b")
	MongoReplicaSetInstance013212.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance013212.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance013212.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.13.212") //primary
	MongoReplicaSetInstance013212.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013212.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013212.StopServices = false
	MongoReplicaSetInstance013212.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013212.EnableXvdpGp3 = true
	MongoReplicaSetInstance013212.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013212.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance013212.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance013228 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013228.EnableEc2instance = true
	MongoReplicaSetInstance013228.Ec2Instance.ImageId = cloudformation.String("ami-0df43911a5efd456b")
	MongoReplicaSetInstance013228.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance013228.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance013228.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.13.228")
	MongoReplicaSetInstance013228.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013228.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013228.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013228.StopServices = false
	MongoReplicaSetInstance013228.EnableXvdpGp3 = true
	MongoReplicaSetInstance013228.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013228.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance013228.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sg1/Mongo-Integration-Debugger-1", "sg1-Mongo-Integration-Debugger-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sg1/Mongo-Integration-Debugger-1", "sg1-Mongo-Integration-Debugger-1-Service.json")
}
