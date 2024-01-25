package eu2Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateEu2MongoAccountsTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("eu2")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("eu2")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu2",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoAccountsReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.48/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.15.7.48 - 10.15.7.64
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu2",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoAccountsReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.64/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.49 - 10.14.6.62
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "eu2",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoAccountsReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.96/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.65 - 10.14.6.78
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("eu2-MongoInstanceSecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "accounts-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007007 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007007.EnableEc2instance = false
	MongoReplicaSetInstance007007.Ec2Instance.ImageId = cloudformation.String("ami-0f572ed10ef61301e")
	MongoReplicaSetInstance007007.Ec2Instance.InstanceType = cloudformation.String("t3.medium")
	MongoReplicaSetInstance007007.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007007.Ec2Instance.PrivateIpAddress = cloudformation.String("10.13.7.52")
	MongoReplicaSetInstance007007.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007007.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007007.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance007007.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007022 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007022.EnableEc2instance = false
	MongoReplicaSetInstance007022.Ec2Instance.ImageId = cloudformation.String("ami-0f572ed10ef61301e")
	MongoReplicaSetInstance007022.Ec2Instance.InstanceType = cloudformation.String("t3.medium")
	MongoReplicaSetInstance007022.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007022.Ec2Instance.PrivateIpAddress = cloudformation.String("10.13.7.68")
	MongoReplicaSetInstance007022.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007022.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007022.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance007022.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007038 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007038.EnableEc2instance = false
	MongoReplicaSetInstance007038.Ec2Instance.ImageId = cloudformation.String("ami-0f572ed10ef61301e")
	MongoReplicaSetInstance007038.Ec2Instance.InstanceType = cloudformation.String("t3.medium")
	MongoReplicaSetInstance007038.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007038.Ec2Instance.PrivateIpAddress = cloudformation.String("10.13.7.100")
	MongoReplicaSetInstance007038.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007038.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007038.MongoContainerTag = "bamboo-mongo-sne-6117-1"
	MongoReplicaSetInstance007038.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/eu2/eu2-Mongo-Accounts-1", "eu2-Mongo-Accounts-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/eu2/eu2-Mongo-Accounts-1", "eu2-Mongo-Accounts-1-Service.json")
}
