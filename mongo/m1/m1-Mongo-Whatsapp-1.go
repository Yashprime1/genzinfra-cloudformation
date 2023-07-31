package m1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func Generatem1MongoWhatsappTemplate() {
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
		SubnetCidrBlockSuffix:  "6.0/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "m1",
		AvailabilityZoneSuffix: "b",
		SubnetCidrBlockSuffix:  "6.16/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs
	})
	subnetB.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-09ea601f56e47ab9e")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{cloudformation.ImportValue("m1-SecurityGroup-MongoWhatsappEc2InstanceEC2SecurityGroupId")}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "wa-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	mongoReplicaInstance006005 := mongo.NewMongo(defaults)
	mongoReplicaInstance006005.EnableEc2instance = false
	mongoReplicaInstance006005.Ec2Instance.ImageId = cloudformation.String("ami-0f4af75f590b4738e")
	mongoReplicaInstance006005.Ec2Instance.InstanceType = cloudformation.String("c5a.24xlarge")
	mongoReplicaInstance006005.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006005.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.5")
	mongoReplicaInstance006005.XvdpEc2Volume.Size = cloudformation.Int(3072)
	mongoReplicaInstance006005.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006005.StopServices = false
	mongoReplicaInstance006005.EnableMongoRegistryCache = true
	mongoReplicaInstance006005.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	mongoReplicaInstance006005.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006006 := mongo.NewMongo(defaults)
	mongoReplicaInstance006006.EnableEc2instance = false
	mongoReplicaInstance006006.Ec2Instance.ImageId = cloudformation.String("ami-0f4af75f590b4738e")
	mongoReplicaInstance006006.Ec2Instance.InstanceType = cloudformation.String("c5a.24xlarge")
	mongoReplicaInstance006006.Ec2InstanceSubnet = subnetA
	mongoReplicaInstance006006.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.6")
	mongoReplicaInstance006006.XvdpEc2Volume.Size = cloudformation.Int(3072)
	mongoReplicaInstance006006.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006006.StopServices = false
	mongoReplicaInstance006006.EnableMongoRegistryCache = true
	mongoReplicaInstance006006.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	mongoReplicaInstance006006.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006020 := mongo.NewMongo(defaults)
	mongoReplicaInstance006020.EnableEc2instance = false
	mongoReplicaInstance006020.Ec2Instance.ImageId = cloudformation.String("ami-0f4af75f590b4738e")
	mongoReplicaInstance006020.Ec2Instance.InstanceType = cloudformation.String("c5a.24xlarge")
	mongoReplicaInstance006020.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006020.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.20")
	mongoReplicaInstance006020.XvdpEc2Volume.Size = cloudformation.Int(3072)
	mongoReplicaInstance006020.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006020.StopServices = false
	mongoReplicaInstance006020.EnableMongoRegistryCache = true
	mongoReplicaInstance006020.MongoContainerTag = "bamboo-mongo-sne-6117-2"
	mongoReplicaInstance006020.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006022 := mongo.NewMongo(defaults)
	mongoReplicaInstance006022.EnableEc2instance = false
	mongoReplicaInstance006022.Ec2Instance.ImageId = cloudformation.String("ami-0f4af75f590b4738e")
	mongoReplicaInstance006022.Ec2Instance.InstanceType = cloudformation.String("r5.8xlarge")
	mongoReplicaInstance006022.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006022.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.22")
	mongoReplicaInstance006022.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-07d7228ca1f24e795")
	mongoReplicaInstance006022.XvdpEc2Volume.Size = cloudformation.Int(1024)
	mongoReplicaInstance006022.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006022.StopServices = false
	mongoReplicaInstance006022.EnableMongoRegistryCache = true
	mongoReplicaInstance006022.MongoContainerTag = "bamboo-mongo-master-15"
	mongoReplicaInstance006022.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006023 := mongo.NewMongo(defaults)
	mongoReplicaInstance006023.EnableEc2instance = false
	mongoReplicaInstance006023.Ec2Instance.ImageId = cloudformation.String("ami-0f4af75f590b4738e")
	mongoReplicaInstance006023.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	mongoReplicaInstance006023.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006023.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.23")
	mongoReplicaInstance006023.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-01e00d0cc1cc3a81c")
	mongoReplicaInstance006023.XvdpEc2Volume.Size = cloudformation.Int(196)
	mongoReplicaInstance006023.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006023.StopServices = false
	mongoReplicaInstance006023.EnableMongoRegistryCache = true
	mongoReplicaInstance006023.MongoContainerTag = "bamboo-mongo-master-15"
	mongoReplicaInstance006023.AppendToTemplate(sTemplate, serviceTemplate)

	mongoReplicaInstance006024 := mongo.NewMongo(defaults)
	mongoReplicaInstance006024.EnableEc2instance = false
	mongoReplicaInstance006024.Ec2Instance.ImageId = cloudformation.String("ami-0f4af75f590b4738e")
	mongoReplicaInstance006024.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	mongoReplicaInstance006024.Ec2InstanceSubnet = subnetB
	mongoReplicaInstance006024.Ec2Instance.PrivateIpAddress = cloudformation.String("10.18.6.24")
	mongoReplicaInstance006024.XvdpEc2Volume.SnapshotId = cloudformation.String("snap-0067c919d918f2936")
	mongoReplicaInstance006024.XvdpEc2Volume.Size = cloudformation.Int(256)
	mongoReplicaInstance006024.EnableMongoArtifactoryRepository = true
	mongoReplicaInstance006024.StopServices = false
	mongoReplicaInstance006024.EnableMongoRegistryCache = true
	mongoReplicaInstance006024.MongoContainerTag = "bamboo-mongo-master-15"
	mongoReplicaInstance006024.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/m1/Mongo-Whatsapp", "m1-Mongo-Whatsapp-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/m1/Mongo-Whatsapp", "m1-Mongo-Whatsapp-1-Service.json")
}
