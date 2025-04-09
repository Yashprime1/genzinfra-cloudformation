// Go through this for async explanantion - https://www.youtube.com/watch?v=PoRJizFvM7s
const express = require('express');
const AWS = require('aws-sdk');
const axios = require('axios');
const { randomInt } = require('crypto');

const app = express();
const port = 3000;
var switchToZero = false;

enableTaskProtection();
// Configure AWS SDK
AWS.config.update({ region: process.env.AWS_REGION }); // Change the region as needed
const cloudwatch = new AWS.CloudWatch();
const ecs = new AWS.ECS();

// Function to fetch service name from ECS
async function getServiceName(clusterName,taskArn) {
    try {
      const clusterName = 'eu1-Callsock-1-CallsockEcsCluster-dNuh3PWdH6I6';
  
      // Describe the task using the ECS client
      const response = await ecs.describeTasks({
        cluster: clusterName,
        tasks: [taskArn]
      }).promise();
  
      const taskGroup = response.tasks[0]?.group;
      
      if (!taskGroup) {
        throw new Error('No task group found');
      }
  
      // Extract service name from task group
      const serviceName = taskGroup.replace('service:', '');
      
      console.log('Service Name:', serviceName);
      return serviceName;
    } catch (error) {
      console.error('Error fetching service name:', error);
      throw error;
    }
  }

// Function to fetch ECS metadata
async function fetchECSMetadata() {
    try {
        const metadataUri = process.env.ECS_CONTAINER_METADATA_URI_V4;
        if (!metadataUri) {
            throw new Error('ECS_CONTAINER_METADATA_URI_V4 is not set');
        }

        const response = await axios.get(`${metadataUri}/task`);
        const { Cluster, TaskARN,Family } = response.data;

        const clusterName = Cluster.split('/').pop();
        const taskId = TaskARN.split('/').pop();
        const serviceName = await getServiceName(clusterName,TaskARN);

        return { clusterName, serviceName, taskId };
    } catch (error) {
        console.error('Error fetching ECS metadata:', error.message);
        throw error;
    }
}

// Function to publish a constant metric
async function publishMetric() {
    const { clusterName, serviceName, taskId } = await fetchECSMetadata();

    const dimensions = [
        { Name: 'ClusterName', Value: clusterName },
        { Name: 'ServiceName', Value: serviceName },
        { Name: 'TaskId', Value: taskId },
    ];

    var metricValue = randomInt(500,600);
    if (switchToZero) {
        metricValue = 0;
    }
    const params = {
        MetricData: [
            {
                MetricName: 'MyCallSockConnections', // Replace with your metric name
                Dimensions: dimensions,
                Unit: 'Count',
                Value: metricValue, // Constant value
            },
        ],
        Namespace: 'AWS/ECS', // Replace with your namespace
    };

    try {
        await cloudwatch.putMetricData(params).promise();
        console.log('Metric published successfully');
    } catch (error) {
        console.error('Error publishing metric:', error.message);
    }
}


// Endpoint to trigger metric publishing
app.get('/publish-metric', async (req, res) => {
    try {
        switchToZero = true;
        res.send('Metric published successfully');
    } catch (error) {
        res.status(500).send('Error publishing metric');
    }
});



app.get('/', async (req, res) => {
    try {
        res.send("Yo Welcome To My Callsock");
    } catch (error) {
        res.status(500).send('Error'+error.message);
    }
});

app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});


let initiateShutdown = async function () {
    console.log("Sigterm received. Initiating graceful shutdown. Server IP: ");
    try {
        const { clusterName, serviceName, taskId } = await fetchECSMetadata();
        // Submit shutdown details to SQS
        await submitShutdownMessageToSQS({ clusterName, serviceName, taskId });
        console.log("Server shutdown complete");
    } catch (error) {
        console.error("Error during shutdown process:", error.message);
    }

    process.exit(0);
};

async function submitShutdownMessageToSQS({ clusterName, serviceName, taskId }) {
    const sqs = new AWS.SQS();
    const queueUrl = "https://sqs.eu-west-1.amazonaws.com/736548753645/eu1-SharedResources-callsock-to-admin-signedcall-sqs-queue"; // Ensure this environment variable is set

    if (!queueUrl) {
        console.error("SQS_QUEUE_URL is not set");
        return;
    }

    const messageBody = JSON.stringify({
        event: "TaskShuttingDown",
        clusterName,
        serviceName,
        taskId,
        timestamp: new Date().toISOString(),
    });

    const params = {
        QueueUrl: queueUrl,
        MessageBody: messageBody,
    };

    try {
        await sqs.sendMessage(params).promise();
        console.log("Shutdown message submitted to SQS successfully");
    } catch (error) {
        console.error("Error submitting message to SQS:", error.message);
    }
}
  
process.on('SIGINT', initiateShutdown);
process.on('SIGTERM', initiateShutdown);


async function enableTaskProtection(taskArn) {
    try {
        const ecsAgentUri = process.env.ECS_AGENT_URI;
        if (!ecsAgentUri) {
            throw new Error('ECS_CONTAINER_METADATA_URI_V4 is not set');
        }

        const endpoint = `${ecsAgentUri}/task-protection/v1/state`;
        const response = await axios.put(endpoint, {
            "ProtectionEnabled": true
        });
        console.log('Task protection enabled:', response.data);
    } catch (error) {
        console.error('Error enabling task protection:', error.message);
    }
}

setInterval(publishMetric, 10000);