## Overview
Mongo is one of our core platform databases and there are various purposes for which mongo is being used in every region.

- Mongo is being utilised in a variety of ways in our production environment -
  - mongo-common 
      description:  NA
       
  - mongo-metadata 
      description: holds regional level accounts and user metadata ie. commons, billing database. 
      
  - mongo-directcall
      description: holds regional directcall database.
      
  - mongo-link-shortener 
      description: hold regional link_shortener_*  collections.
      
  - mongo-catalog 
      description: hold regional catalogs and recommendations collections.
      
  - mongo-accounts 
      description: hold regional accounts (new_arch) database.
      
  - mongo-legacy-accounts 
      description: hold regional account legacy databases in old architecture. 
      
  - mongo-audit 
      description: hold regional dashboard audited actions info in audit database.
      
  - mongo-whatsapp
      description: hold regional whatsapp conversations and users data.
      
  - mongo-stats 
      description: hold regional stats related collections like targeting_stats , segment_stats, realtime_stats.
      
## Architecture
- A replica set in MongoDB is a group of mongod processes that maintain the same data set. Replica sets provide redundancy and high availability, and are the basis for all production deployments.
- So, here we have a replica setup for our production environment too.
- There are 3 mongo containers running on 3 instances using placement constraint for ECS service and a SSM document that's available in every region `<NetworkPrefix>-SharedResources-primarynodeBootstrapMongo-*` once ran on any one of the instances will help setup the replication.   
- For more information regarding mongo bootstrap and initializing a replica setup can be found here -
`https://wizrocket.atlassian.net/wiki/spaces/SNE/pages/2018477590/MongoDB+bootstrap`  


## Auto Scaling
NA
  

## Containers
There are 4 various containers that run on a particular EC2 instance here -
- Mongo Container (3 Mongo ECS Services running in replica mode - dedicated to run on a particluar instance)
- Splunk Forwarder Container (1 ECS Service running in daemon mode)
- Sensu (1 ECS Service running in daemon mode)
- Cadvisor (1 ECS Service running in daemon mode)

## Data Storage
This is done on the xvdp EBS volume provisioned for mongo.

## Data Backup
The retention period of a mongo snapshot is 14 days. We take snapshots of EBS volume every 4 hours. The detailed process to restore DB from it can be found here - `https://wizrocket.atlassian.net/wiki/spaces/SNE/pages/3976757607/Document+to+restore+data+from+mongo+snapshot+and+configure+it`

## Monitoring
- Sensu
- Splunk
- Grafana

  
## Application logs
Mongo logs are forwarded by splunk forwarder to splunk at index "mongo".

## Stdout
NA

## Alarm
Not yet configured

## Debugging tools
- Splunk Cloud
- Prometheus
- Grafana
- Sensu

## Frequently Asked Questions
No FAQ
