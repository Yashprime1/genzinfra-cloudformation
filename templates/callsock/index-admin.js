const AWS = require('aws-sdk');
AWS.config.update({ region: 'eu-west-1' }); 
const ecs = new AWS.ECS();

const clusterName = 'eu1-Callsock-1-CallsockEcsCluster-dNuh3PWdH6I6'; // Replace with your ECS cluster name
const serviceName = 'eu1-Callsock-1-Service-CallSockV4EcsService-2upj4pQcckdC';
const checkInterval = 10 * 1000; // 5 minutes in milliseconds

const cloudwatch = new AWS.CloudWatch(); 

async function waitForMetricToReachZero(clusterName, serviceName, taskArn) {
    const metricName = "MyCallSockConnections"; // Replace with actual metric name
    const namespace = "AWS/ECS"; // Replace with actual namespace
    
    while (true) {
        const now = new Date();
        const past = new Date(now.getTime() - 5 * 60 * 1000); // Look at the last 5 minutes

        const params = {
            MetricDataQueries: [
                {
                    Id: "metricQuery",
                    MetricStat: {
                        Metric: {
                            Namespace: namespace,
                            MetricName: metricName,
                            Dimensions: [
                                { Name: "ClusterName", Value: clusterName },
                                { Name: "ServiceName", Value: serviceName },
                                { Name: "TaskId", Value: taskArn.split('/').pop() }
                            ],
                        },
                        Period: 60,
                        Stat: "Average",
                    },
                    ReturnData: true,
                },
            ],
            StartTime: past,
            EndTime: now,
        };

        console.log("Fetching metric data with params:", JSON.stringify(params, null, 2));

        try {
            const response = await cloudwatch.getMetricData(params).promise();
            const dataPoints = response.MetricDataResults[0]?.Values || [];
            const latestValue = dataPoints[0] || 0;
            console.log("Metric data points:", dataPoints);
            console.log("Latest metric value:", latestValue);

            if (latestValue === 0) {
                console.log("Metric has reached zero. Proceeding...");
                break;
            }
        } catch (error) {
            console.error("Error fetching metric data:", error);
        }

        console.log("Waiting for metric to reach zero...");
        await new Promise((resolve) => setTimeout(resolve, 30000)); // Wait 30 seconds before checking again
    }
}

async function checkServiceEvents() {
    try {
        const describeServicesResponse = await ecs.describeServices({
            cluster: clusterName,
            services: [serviceName],
        }).promise();

        const events = describeServicesResponse.services[0]?.events || [];
        const latestEvent = events[0]; // Events are sorted by timestamp (latest first)

        if (
            latestEvent &&
            latestEvent.message.includes(
                'was unable to reach steady state'
            ) &&
            latestEvent.message.includes('was unable to scale in due to')
        ) {
            console.log('Matching event found:', latestEvent.message);

            const taskSetIdMatch = latestEvent.message.match(/taskSet (ecs-svc\/\d+)/);
            if (!taskSetIdMatch) {
                console.error('TaskSet ID not found in the event message.');
                return;
            }

            const taskSetId = taskSetIdMatch[1];
            console.log('TaskSet ID:', taskSetId);

            // Pass the latest event to the handler
            await handleScalingIssue(taskSetId, latestEvent);
        } else {
            console.log('No matching event found.');
        }
    } catch (error) {
        console.error('Error checking service events:', error);
    }
}

async function handleScalingIssue(taskSetId, latestEvent) {
    try {
        const listTasksResponse = await ecs.listTasks({
            cluster: clusterName,
            serviceName: serviceName,
            desiredStatus: 'RUNNING',
        }).promise();

        const taskArns = listTasksResponse.taskArns;
        console.log('Running tasks:', taskArns);

        // Extract the number of tasks under protection from the latest event message
        const tasksUnderProtectionMatch = latestEvent.message.match(/reason (\d+) tasks under protection/);
        if (!tasksUnderProtectionMatch) {
            console.error('Number of tasks under protection not found in the event message.');
            return;
        }

        const tasksUnderProtection = parseInt(tasksUnderProtectionMatch[1], 10);
        console.log(`Number of tasks under protection: ${tasksUnderProtection}`);

        // Get task details to find the oldest tasks
        const taskDetails = await Promise.all(
            taskArns.map(async (taskArn) => {
                const taskResponse = await ecs.describeTasks({
                    cluster: clusterName,
                    tasks: [taskArn],
                }).promise();
                const task = taskResponse.tasks[0];
                const protectionResponse = await ecs.getTaskProtection({
                    cluster: clusterName,
                    tasks: [taskArn],
                }).promise();

                const protectionEnabled = protectionResponse.protectedTasks[0]?.protectionEnabled || false;

                return {
                    taskArn,
                    startedAt: task.startedAt,
                    protectionEnabled,
                };
            })
        );
        console.log('Task details:', taskDetails); // Debugging log
        // Sort tasks by start time (oldest first)
        const sortedTasks = taskDetails
            .filter((task) => task.protectionEnabled) // Only consider protected tasks
            .sort((a, b) => a.startedAt - b.startedAt);

        // Remove protection from the required number of oldest tasks
        const tasksToUnprotect = sortedTasks.slice(0, tasksUnderProtection);
        console.log('Tasks to unprotect:', tasksToUnprotect); // Debugging log
        if (tasksToUnprotect.length === 0) {
            console.log('No tasks to unprotect. Exiting.');
            return;
        }
        for (const task of tasksToUnprotect) {
            console.log(`Removing protection from task ${task.taskArn}...`);

            await waitForMetricToReachZero(clusterName, serviceName, task.taskArn);
            await ecs.updateTaskProtection({
                cluster: clusterName,
                tasks: [task.taskArn],
                protectionEnabled: false,
            }).promise();
            console.log(`Protection removed from task ${task.taskArn}.`);
        }
    } catch (error) {
        console.error('Error handling scaling issue:', error);
    }
}

// Start checking service events at regular intervals
setInterval(checkServiceEvents, checkInterval);const AWS = require('aws-sdk');
AWS.config.update({ region: 'eu-west-1' }); 
const ecs = new AWS.ECS();

const clusterName = 'eu1-Callsock-1-CallsockEcsCluster-dNuh3PWdH6I6'; // Replace with your ECS cluster name
const serviceName = 'eu1-Callsock-1-Service-CallSockV4EcsService-2upj4pQcckdC';
const checkInterval = 10 * 1000; // 5 minutes in milliseconds

const cloudwatch = new AWS.CloudWatch(); 

async function waitForMetricToReachZero(clusterName, serviceName, taskArn) {
    const metricName = "MyCallSockConnections"; // Replace with actual metric name
    const namespace = "AWS/ECS"; // Replace with actual namespace
    
    while (true) {
        const now = new Date();
        const past = new Date(now.getTime() - 5 * 60 * 1000); // Look at the last 5 minutes

        const params = {
            MetricDataQueries: [
                {
                    Id: "metricQuery",
                    MetricStat: {
                        Metric: {
                            Namespace: namespace,
                            MetricName: metricName,
                            Dimensions: [
                                { Name: "ClusterName", Value: clusterName },
                                { Name: "ServiceName", Value: serviceName },
                                { Name: "TaskId", Value: taskArn.split('/').pop() }
                            ],
                        },
                        Period: 60,
                        Stat: "Average",
                    },
                    ReturnData: true,
                },
            ],
            StartTime: past,
            EndTime: now,
        };

        console.log("Fetching metric data with params:", JSON.stringify(params, null, 2));

        try {
            const response = await cloudwatch.getMetricData(params).promise();
            const dataPoints = response.MetricDataResults[0]?.Values || [];
            const latestValue = dataPoints[0] || 0;
            console.log("Metric data points:", dataPoints);
            console.log("Latest metric value:", latestValue);

            if (latestValue === 0) {
                console.log("Metric has reached zero. Proceeding...");
                break;
            }
        } catch (error) {
            console.error("Error fetching metric data:", error);
        }

        console.log("Waiting for metric to reach zero...");
        await new Promise((resolve) => setTimeout(resolve, 30000)); // Wait 30 seconds before checking again
    }
}

async function checkServiceEvents() {
    try {
        const describeServicesResponse = await ecs.describeServices({
            cluster: clusterName,
            services: [serviceName],
        }).promise();

        const events = describeServicesResponse.services[0]?.events || [];
        const latestEvent = events[0]; // Events are sorted by timestamp (latest first)

        if (
            latestEvent &&
            latestEvent.message.includes(
                'was unable to reach steady state'
            ) &&
            latestEvent.message.includes('was unable to scale in due to')
        ) {
            console.log('Matching event found:', latestEvent.message);

            const taskSetIdMatch = latestEvent.message.match(/taskSet (ecs-svc\/\d+)/);
            if (!taskSetIdMatch) {
                console.error('TaskSet ID not found in the event message.');
                return;
            }

            const taskSetId = taskSetIdMatch[1];
            console.log('TaskSet ID:', taskSetId);

            // Pass the latest event to the handler
            await handleScalingIssue(taskSetId, latestEvent);
        } else {
            console.log('No matching event found.');
        }
    } catch (error) {
        console.error('Error checking service events:', error);
    }
}

async function handleScalingIssue(taskSetId, latestEvent) {
    try {
        const listTasksResponse = await ecs.listTasks({
            cluster: clusterName,
            serviceName: serviceName,
            desiredStatus: 'RUNNING',
        }).promise();

        const taskArns = listTasksResponse.taskArns;
        console.log('Running tasks:', taskArns);

        // Extract the number of tasks under protection from the latest event message
        const tasksUnderProtectionMatch = latestEvent.message.match(/reason (\d+) tasks under protection/);
        if (!tasksUnderProtectionMatch) {
            console.error('Number of tasks under protection not found in the event message.');
            return;
        }

        const tasksUnderProtection = parseInt(tasksUnderProtectionMatch[1], 10);
        console.log(`Number of tasks under protection: ${tasksUnderProtection}`);

        // Get task details to find the oldest tasks
        const taskDetails = await Promise.all(
            taskArns.map(async (taskArn) => {
                const taskResponse = await ecs.describeTasks({
                    cluster: clusterName,
                    tasks: [taskArn],
                }).promise();
                const task = taskResponse.tasks[0];
                const protectionResponse = await ecs.getTaskProtection({
                    cluster: clusterName,
                    tasks: [taskArn],
                }).promise();

                const protectionEnabled = protectionResponse.protectedTasks[0]?.protectionEnabled || false;

                return {
                    taskArn,
                    startedAt: task.startedAt,
                    protectionEnabled,
                };
            })
        );
        console.log('Task details:', taskDetails); // Debugging log
        // Sort tasks by start time (oldest first)
        const sortedTasks = taskDetails
            .filter((task) => task.protectionEnabled) // Only consider protected tasks
            .sort((a, b) => a.startedAt - b.startedAt);

        // Remove protection from the required number of oldest tasks
        const tasksToUnprotect = sortedTasks.slice(0, tasksUnderProtection);
        console.log('Tasks to unprotect:', tasksToUnprotect); // Debugging log
        if (tasksToUnprotect.length === 0) {
            console.log('No tasks to unprotect. Exiting.');
            return;
        }
        for (const task of tasksToUnprotect) {
            console.log(`Removing protection from task ${task.taskArn}...`);

            await waitForMetricToReachZero(clusterName, serviceName, task.taskArn);
            await ecs.updateTaskProtection({
                cluster: clusterName,
                tasks: [task.taskArn],
                protectionEnabled: false,
            }).promise();
            console.log(`Protection removed from task ${task.taskArn}.`);
        }
    } catch (error) {
        console.error('Error handling scaling issue:', error);
    }
}

// Start checking service events at regular intervals
setInterval(checkServiceEvents, checkInterval);