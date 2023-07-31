package sk1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSk1MongoWiredTiger3Dot6Template() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("sk1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("sk1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "WiredTigerMongoReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.0/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.33 - 10.14.6.46
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "WiredTigerMongoReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.16/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.49 - 10.14.6.62
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "WiredTigerMongoReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.32/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.65 - 10.14.6.78
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0a91977023aa59ed0")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("sk1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "wa-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007007 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007007.EnableEc2instance = true
	MongoReplicaSetInstance007007.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance007007.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance007007.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007007.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.7.7") //primary
	MongoReplicaSetInstance007007.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007007.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007007.StopServices = false
	MongoReplicaSetInstance007007.EnableXvdpGp3 = true
	MongoReplicaSetInstance007007.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007007.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	MongoReplicaSetInstance007007.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007007.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007022 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007022.EnableEc2instance = true
	MongoReplicaSetInstance007022.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance007022.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance007022.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007022.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.7.22")
	MongoReplicaSetInstance007022.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007022.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007022.StopServices = false
	MongoReplicaSetInstance007022.EnableXvdpGp3 = true
	MongoReplicaSetInstance007022.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007022.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	MongoReplicaSetInstance007022.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007022.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007038 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007038.EnableEc2instance = true
	MongoReplicaSetInstance007038.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance007038.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance007038.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007038.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.7.38")
	MongoReplicaSetInstance007038.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007038.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007038.StopServices = false
	MongoReplicaSetInstance007038.EnableXvdpGp3 = true
	MongoReplicaSetInstance007038.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007038.MongoContainerTag = "bamboo-mongo-sne-6117-3"
	MongoReplicaSetInstance007038.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007038.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sk1/Mongo-MetaData-1", "sk1-Mongo-MetaData-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sk1/Mongo-MetaData-1", "sk1-Mongo-MetaData-1-Service.json")
}
