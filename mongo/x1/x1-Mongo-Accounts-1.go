package x1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateX1MongoAccountsTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKey()
	sTemplate.Resources["MongoEbsDlmLifecyclePolicy"] = mongo.GetDlmLifeCyclePolicy()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("x1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("x1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "x1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoAccountsReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.48/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.17.7.48 - 10.15.7.64
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "x1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoAccountsReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.64/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.16.6.49 - 10.14.6.62
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "x1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoAccountsReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.96/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.16.6.65 - 10.14.6.78
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableVividCortexEcsService = true
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("x1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "accounts-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007052 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007052.EnableEc2instance = true
	MongoReplicaSetInstance007052.Ec2Instance.ImageId = cloudformation.String("ami-00ac255069dc94e3c")
	MongoReplicaSetInstance007052.Ec2Instance.InstanceType = cloudformation.String("t2.small")
	MongoReplicaSetInstance007052.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007052.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.7.52") //primary
	MongoReplicaSetInstance007052.XvdpEc2Volume.Size = cloudformation.Int(128)
	MongoReplicaSetInstance007052.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007052.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance007052.EnableXvdpGp3 = true
	MongoReplicaSetInstance007052.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007068 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007068.EnableEc2instance = true
	MongoReplicaSetInstance007068.Ec2Instance.ImageId = cloudformation.String("ami-00ac255069dc94e3c")
	MongoReplicaSetInstance007068.Ec2Instance.InstanceType = cloudformation.String("t2.small")
	MongoReplicaSetInstance007068.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007068.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.7.68")
	MongoReplicaSetInstance007068.XvdpEc2Volume.Size = cloudformation.Int(128)
	MongoReplicaSetInstance007068.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007068.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance007068.EnableXvdpGp3 = true
	MongoReplicaSetInstance007068.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007100 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007100.EnableEc2instance = true
	MongoReplicaSetInstance007100.Ec2Instance.ImageId = cloudformation.String("ami-00ac255069dc94e3c")
	MongoReplicaSetInstance007100.Ec2Instance.InstanceType = cloudformation.String("t2.small")
	MongoReplicaSetInstance007100.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007100.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.7.100")
	MongoReplicaSetInstance007100.XvdpEc2Volume.Size = cloudformation.Int(128)
	MongoReplicaSetInstance007100.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007100.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance007100.EnableXvdpGp3 = true
	MongoReplicaSetInstance007100.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/x1/x1-Mongo-Accounts-1", "x1-Mongo-Accounts-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/x1/x1-Mongo-Accounts-1", "x1-Mongo-Accounts-1-Service.json")
}
