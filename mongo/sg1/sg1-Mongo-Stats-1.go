package sg1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSg1MongoStatsTemplate() {
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
		SubnetCidrBlockSuffix:  "10.112/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "10.128/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sg1",
		AvailabilityZoneSuffix: "c",
		SubnetCidrBlockSuffix:  "10.144/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(100)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-038ad0e960f5a6d6c")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("sg1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "sg1-Mongo-Stats-1-rs0", "--logpath", "/va2r/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance006005 := mongo.NewMongo(defaults)
	mongoReplicaInstance006005.EnableEc2instance = true
	mongoReplicaInstance006005.Ec2Instance.ImageId = cloudformation.String("ami-0a98453160fc9f10a")
	mongoReplicaInstance006005.Ec2Instance.InstanceType = cloudformation.String("r5.xlarge")
	mongoReplicaInstance006005.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006005.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.10.117")
	mongoReplicaInstance006005.XvdpEc2Volume.Size = cloudformation.Int(164)
	mongoReplicaInstance006005.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006005.EnableXvdpGp3 = true
	mongoReplicaInstance006005.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance006005.MongoContainerTag = "github-sne-6117-16"
	mongoReplicaInstance006005.StopServices = false
	mongoReplicaInstance006005.EnableMongoRegistryCache = true
	mongoReplicaInstance006005.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006006 := mongo.NewMongo(defaults)
	mongoReplicaInstance006006.EnableEc2instance = true
	mongoReplicaInstance006006.Ec2Instance.ImageId = cloudformation.String("ami-0a98453160fc9f10a")
	mongoReplicaInstance006006.Ec2Instance.InstanceType = cloudformation.String("r5.xlarge")
	mongoReplicaInstance006006.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006006.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.10.133") //primary
	mongoReplicaInstance006006.XvdpEc2Volume.Size = cloudformation.Int(164)
	mongoReplicaInstance006006.EnableXvdpGp3 = true
	mongoReplicaInstance006006.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance006006.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006006.StopServices = false
	mongoReplicaInstance006006.EnableMongoRegistryCache = true
	mongoReplicaInstance006006.MongoContainerTag = "github-sne-6117-16"
	mongoReplicaInstance006006.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006020 := mongo.NewMongo(defaults)
	mongoReplicaInstance006020.EnableEc2instance = true
	mongoReplicaInstance006020.Ec2Instance.ImageId = cloudformation.String("ami-0a98453160fc9f10a")
	mongoReplicaInstance006020.Ec2Instance.InstanceType = cloudformation.String("r5.xlarge")
	mongoReplicaInstance006020.Ec2InstanceSubnet = subnetC
	mongoReplicaInstance006020.Ec2Instance.PrivateIpAddress = cloudformation.String("10.15.10.149")
	mongoReplicaInstance006020.XvdpEc2Volume.Size = cloudformation.Int(164)
	mongoReplicaInstance006020.EnableXvdpGp3 = true
	mongoReplicaInstance006020.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	mongoReplicaInstance006020.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006020.StopServices = false
	mongoReplicaInstance006020.EnableMongoRegistryCache = true
	mongoReplicaInstance006020.MongoContainerTag = "github-sne-6117-16"
	mongoReplicaInstance006020.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sg1/Mongo-Stats", "sg1-Mongo-Stats-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sg1/Mongo-Stats", "sg1-Mongo-Stats-1-Service.json")
}
