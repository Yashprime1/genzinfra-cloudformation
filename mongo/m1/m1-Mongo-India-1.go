package m1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func Generatem1MongoIndia1Template() {
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
		SubnetCidrBlockSuffix:  "6.224/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "m1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "6.240/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(1024)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0f8afc87a647567dd")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("m1-SecurityGroup-MongoCommonEc2InstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance006230 := mongo.NewMongo(defaults)
	mongoReplicaInstance006230.EnableEc2instance = false
	mongoReplicaInstance006230.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	mongoReplicaInstance006230.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006230.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.230")
	mongoReplicaInstance006230.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0c777b9c5dd3fce0e")
	mongoReplicaInstance006230.XvdpEc2Volume.Size = cloudformation.Int(1024)
	mongoReplicaInstance006230.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006230.StopServices = false
	mongoReplicaInstance006230.EnableMongoRegistryCache = true
	mongoReplicaInstance006230.MongoContainerTag = "bamboo-mongo-master-15"
	mongoReplicaInstance006230.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006231 := mongo.NewMongo(defaults)
	mongoReplicaInstance006231.EnableEc2instance = false
	mongoReplicaInstance006231.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	mongoReplicaInstance006231.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006231.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.231")
	mongoReplicaInstance006231.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0c777b9c5dd3fce0e")
	mongoReplicaInstance006231.XvdpEc2Volume.Size = cloudformation.Int(1024)
	mongoReplicaInstance006231.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006231.StopServices = false
	mongoReplicaInstance006231.EnableMongoRegistryCache = true
	mongoReplicaInstance006231.MongoContainerTag = "bamboo-mongo-master-15"
	mongoReplicaInstance006231.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006246 := mongo.NewMongo(defaults)
	mongoReplicaInstance006246.EnableEc2instance = false
	mongoReplicaInstance006246.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	mongoReplicaInstance006246.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006246.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.246")
	mongoReplicaInstance006246.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0c777b9c5dd3fce0e")
	mongoReplicaInstance006246.XvdpEc2Volume.Size = cloudformation.Int(1024)
	mongoReplicaInstance006246.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006246.StopServices = false
	mongoReplicaInstance006246.EnableMongoRegistryCache = true
	mongoReplicaInstance006246.MongoContainerTag = "bamboo-mongo-master-15"
	mongoReplicaInstance006246.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/m1/Mongo-India-1", "m1-Mongo-India-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/m1/Mongo-India-1", "m1-Mongo-India-1-Service.json")
}
