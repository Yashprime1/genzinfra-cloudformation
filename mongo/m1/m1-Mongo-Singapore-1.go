package m1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func Generatem1MongoSingapore1Template() {
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
		SubnetCidrBlockSuffix:  "6.160/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "m1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "6.176/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(324)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0f8afc87a647567dd")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("m1-SecurityGroup-MongoCommonEc2InstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance006165 := mongo.NewMongo(defaults)
	mongoReplicaInstance006165.EnableEc2instance = false
	mongoReplicaInstance006165.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	mongoReplicaInstance006165.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006165.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.165")
	mongoReplicaInstance006165.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-045081a7c9bb95bff")
	mongoReplicaInstance006165.XvdpEc2Volume.Size = cloudformation.Int(324)
	mongoReplicaInstance006165.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006165.StopServices = false
	mongoReplicaInstance006165.EnableMongoRegistryCache = true
	mongoReplicaInstance006165.MongoContainerTag = "bamboo-mongo-task-sne-8311-5"
	mongoReplicaInstance006165.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006166 := mongo.NewMongo(defaults)
	mongoReplicaInstance006166.EnableEc2instance = false
	mongoReplicaInstance006166.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	mongoReplicaInstance006166.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006166.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.166")
	mongoReplicaInstance006166.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-045081a7c9bb95bff")
	mongoReplicaInstance006166.XvdpEc2Volume.Size = cloudformation.Int(324)
	mongoReplicaInstance006166.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006166.StopServices = false
	mongoReplicaInstance006166.EnableMongoRegistryCache = true
	mongoReplicaInstance006166.MongoContainerTag = "bamboo-mongo-task-sne-8311-5"
	mongoReplicaInstance006166.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006181 := mongo.NewMongo(defaults)
	mongoReplicaInstance006181.EnableEc2instance = false
	mongoReplicaInstance006181.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	mongoReplicaInstance006181.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006181.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.181")
	mongoReplicaInstance006181.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-045081a7c9bb95bff")
	mongoReplicaInstance006181.XvdpEc2Volume.Size = cloudformation.Int(324)
	mongoReplicaInstance006181.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006181.StopServices = false
	mongoReplicaInstance006181.EnableMongoRegistryCache = true
	mongoReplicaInstance006181.MongoContainerTag = "bamboo-mongo-task-sne-8311-5"
	mongoReplicaInstance006181.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/m1/Mongo-Singapore-1", "m1-Mongo-Singapore-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/m1/Mongo-Singapore-1", "m1-Mongo-Singapore-1-Service.json")
}
