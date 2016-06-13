> Replace XXXXXXXXXXXX with your Account ID (i.e. 123456789023).

##### Getting started.

Create your queue.

```
aws sqs create-queue \
  --queue-name MyQueue
```

Create a role for autoscaling.

```
aws iam create-role \
  --role-name Test-Role \
  --assume-role-policy-document file://Test-Role-Trust-Policy.json
```

`Test-Role-Trust-Policy.json`

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "autoscaling.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

```
aws iam attach-role-policy \
  --policy-arn arn:aws:iam::aws:policy/service-role/AutoScalingNotificationAccessRole \
  --role-name Test-Role
```

Create a role for EC2 instance.

```
aws iam create-role \
  --role-name EC2-Role \
  --assume-role-policy-document file://EC2-Role-Trust-Policy.json
```

`EC2-Role-Trust-Policy.json`

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

```
aws iam put-role-policy \
  --role-name EC2-Role \
  --policy-name NotifyPolicy \
  --policy-document file://NotifyPolicy.json
```

`NotifyPolicy.json`

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Resource": "*",
            "Action": [
                "autoscaling:CompleteLifecycleAction",
                "autoscaling:RecordLifecycleActionHeartbeat",
                "ec2:DescribeInstances",
                "sqs:DeleteMessage",
                "sqs:GetQueueUrl",
                "sqs:ReceiveMessage",
                "sns:Publish"
            ]
        }
    ]
}
```

Create an instance profile.

```
aws iam create-instance-profile \
  --instance-profile-name EC2-Instance-Profile
```

Add role to an instance profile.

```
aws iam add-role-to-instance-profile \
  --role-name EC2-Role \
  --instance-profile-name EC2-Instance-Profile
```

Create a launch configuration using Amazon Linux AMI.

> Use your security group in place of sg-YYYYYYYY.

```
aws autoscaling create-launch-configuration \
  --launch-configuration-name my-launch-config \
  --security-groups sg-YYYYYYYY \
  --image-id ami-a4827dc9 \
  --iam-instance-profile arn:aws:iam::XXXXXXXXXXXX:instance-profile/EC2-Instance-Profile \
  --key-name AMI \
  --instance-type t2.nano
```

Create an auto scaling group.

> Use your subnet group in place of subnet-YYYYYYYY.

```
aws autoscaling create-auto-scaling-group \
  --auto-scaling-group-name my-auto-scaling-group \
  --launch-configuration-name my-launch-config \
  --min-size 0 --max-size 1 --desired-capacity 1 \
  --default-cooldown 60 \
  --vpc-zone-identifier subnet-YYYYYYYY
```

Create a lifecycle hook.

```
aws autoscaling put-lifecycle-hook \
  --lifecycle-hook-name my-lifecycle-hook \
  --auto-scaling-group-name my-auto-scaling-group \
  --lifecycle-transition autoscaling:EC2_INSTANCE_TERMINATING \
  --notification-target-arn arn:aws:sqs:us-east-1:XXXXXXXXXXXX:MyQueue \
  --role-arn arn:aws:iam::XXXXXXXXXXXX:role/Test-Role \
  --heartbeat-timeout 600 --default-result 'CONTINUE'
```

List your created hook. GlobalTimeout is the maximum time, in seconds, that an instance can remain in a Pending:Wait or Terminating:Wait state.

```
aws autoscaling describe-lifecycle-hooks \
  --auto-scaling-group-name my-auto-scaling-group
```

`Output`

```json
{
    "LifecycleHooks": [
        {
            "GlobalTimeout": 60000,
            "HeartbeatTimeout": 600,
            "RoleARN": "arn:aws:iam::XXXXXXXXXXXX:role/Test-Role",
            "AutoScalingGroupName": "my-auto-scaling-group",
            "LifecycleHookName": "my-lifecycle-hook",
            "DefaultResult": "CONTINUE",
            "NotificationTargetARN": "arn:aws:sqs:us-east-1:XXXXXXXXXXXX:MyQueue",
            "LifecycleTransition": "autoscaling:EC2_INSTANCE_TERMINATING"
        }
    ]
}
```

Retrieve a message from the queue. There will be a test messages sent from AWS.

```
aws sqs receive-message \
  --queue-url https://sqs.us-east-1.amazonaws.com/XXXXXXXXXXXX/MyQueue
```

`Output`

```json
{
    "Messages": [
        {
            "Body": "{\"AutoScalingGroupName\":\"my-auto-scaling-group\",\"Service\":\"AWS Auto Scaling\",\"Time\":\"2016-06-12T00:07:08.954Z\",\"AccountId\":\"XXXXXXXXXXXX\",\"Event\":\"autoscaling:TEST_NOTIFICATION\",\"RequestId\":\"GGGGGGGG-GGGG-GGGG-GGGG-GGGGGGGGGGGG\",\"AutoScalingGroupARN\":\"arn:aws:autoscaling:us-east-1:XXXXXXXXXXXX:autoScalingGroup:FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF:autoScalingGroupName/my-auto-scaling-group\"}",
            "ReceiptHandle": "AQEB4hFJn5yvhj9LjS7CALKZEMJXMa4PR8A1YICPWuD/jN4cpMqzfx0AEGzJMPRfy79CKdMaeEsgZwKqALYExIfpgFhLwotONelA6GeFlU1FzJoJZVy4/qA327IOSxPrsIvfOwYYV71OQvzJz2o5oH9r36WcAD5DmLEV5dCJJ5uUY0JjJtZNIMY4ZDBFMBdc+R/J0f0LLcZK0DxRF/qTYyCAdTb8ciYMa/JCP30vPm6I+kktF1vqnqtPV30pjKuZN6J2nqdqOIS/ZpIW5cWWZ+p49KsKMYyL3VkH3Zdahr6/GHOciXeoofs47DB/Giy91tBAdpVoRIYAv1lIuqwQanv7bOOb81DoktWAuLQOHeD4BIMi3nMZ0tdLAt9l4EvNx3jC",
            "MD5OfBody": "764efa883dda1e11db47671c4a3bbd9e",
            "MessageId": "DDDDDDDD-DDDD-DDDD-DDDD-DDDDDDDDDDDD"
        }
    ]
}
```

Delete an instance and poll the queue again. There will be terminating message for the EC2 instance.

```
aws sqs receive-message \
  --queue-url https://sqs.us-east-1.amazonaws.com/XXXXXXXXXXXX/MyQueue
```

`Output`

```json
{
    "Messages": [
        {
            "Body": "{\"AutoScalingGroupName\":\"my-auto-scaling-group\",\"Service\":\"AWS Auto Scaling\",\"Time\":\"2016-06-11T23:53:04.656Z\",\"AccountId\":\"XXXXXXXXXXXX\",\"LifecycleTransition\":\"autoscaling:EC2_INSTANCE_TERMINATING\",\"RequestId\":\"HHHHHHHH-HHHH-HHHH-HHHH-HHHHHHHHHHHH\",\"LifecycleActionToken\":\"ZZZZZZZZ-ZZZZ-ZZZZ-ZZZZ-ZZZZZZZZZZZZ\",\"EC2InstanceId\":\"i-59824f1c\",\"LifecycleHookName\":\"my-lifecycle-hook\"}",
            "ReceiptHandle": "AQEBm3WqzQtzWvclLpFdKAbqurUuPYrUyRYhFGuvxoL5poss23h90J1rIIPcIWh9W0zpXG8zWWGva4QDuT0G87dg6zF0qtjgUrbtrCek/5FpQTCZGAW4423fuKIsL6RTUWbQptaHfkv9N1ihXYWarTAYqLI8mpP4NzBJhSXVX03GpkDwxHYPMt7GPxPhCfpW4M0fDOFSpyCU9e0PoFf688Mt3rAfgQz6xBXhijmuHVDY3PzhhaYgWWKz5/PqX4bLnF1SbIhjE/TKcBZnLmrjItFUqQbLiwABJhXwqR9l+jWdxBu6NAQJsuXvKHQXJifau1xYnIYoUbBvHzEebS1Epek18Iy02F4MlXJSuo4UcD2zk/68l+Wy5bHTkPS5oTgxUpCg",
            "MD5OfBody": "b1946ac92492d2347c6235b4d2611184",
            "MessageId": "EEEEEEEE-EEEE-EEEE-EEEE-EEEEEEEEEEEE"
        }
    ]
}
```

Also, the autoscaling group will pause for about 10 minutes waiting for a heartbeat response. If you don't wish to wait 10 minutes, you can use 'CONTINUE' to complete the lifecycle action to terminate the instance.

> Use your lifecycle-action-token in place of ZZZZZZZZ-ZZZZ-ZZZZ-ZZZZ-ZZZZZZZZZZZZ.

```
aws autoscaling complete-lifecycle-action \
  --lifecycle-hook-name my-lifecycle-hook \
  --auto-scaling-group-name my-auto-scaling-group \
  --lifecycle-action-result CONTINUE \
  --lifecycle-action-token ZZZZZZZZ-ZZZZ-ZZZZ-ZZZZ-ZZZZZZZZZZZZ
```

Also scaling activities can be shown to inspect pauses.

```
aws autoscaling describe-scaling-activities \
  --auto-scaling-group-name my-auto-scaling-group
```

`Output`

```json
{
    "Activities": [
        {
            "Description": "Terminating EC2 instance: i-59824f1c",
            "AutoScalingGroupName": "my-auto-scaling-group",
            "ActivityId": "IIIIIIII-IIII-IIII-IIII-IIIIIIIIIIII",
            "Details": "{\"Availability Zone\":\"us-east-1c\",\"Subnet ID\":\"subnet-YYYYYYYY\"}",
            "StartTime": "2016-06-12T02:55:41.406Z",
            "Progress": 60,
            "Cause": "At 2016-06-12T02:55:14Z a user request update of AutoScalingGroup constraints to min: 0, max: 1,
 desired: 0 changing the desired capacity from 1 to 0.  At 2016-06-12T02:55:41Z an instance was taken out of service i
n response to a difference between desired and actual capacity, shrinking the capacity from 1 to 0.  At 2016-06-12T02:
55:41Z instance i-59824f1c was selected for termination.",
            "StatusCode": "MidTerminatingLifecycleAction"
        }
    ]
}
```

You will also notice the instance have been put into the "Terminating:Wait" lifecycle state.

```
aws autoscaling describe-auto-scaling-groups \
  --auto-scaling-group-name my-auto-scaling-group
```

`Output`

```json
{
    "AutoScalingGroups": [
        {
            "AutoScalingGroupARN": "arn:aws:autoscaling:us-east-1:XXXXXXXXXXXX:autoScalingGroup:FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF:autoScalingGroupName/my-auto-scaling-group",
            "HealthCheckGracePeriod": 0,
            "SuspendedProcesses": [],
            "DesiredCapacity": 0,
            "Tags": [],
            "EnabledMetrics": [],
            "LoadBalancerNames": [],
            "AutoScalingGroupName": "my-auto-scaling-group",
            "DefaultCooldown": 60,
            "MinSize": 0,
            "Instances": [
                {
                    "ProtectedFromScaleIn": false,
                    "AvailabilityZone": "us-east-1c",
                    "InstanceId": "i-59824f1c",
                    "HealthStatus": "Healthy",
                    "LifecycleState": "Terminating:Wait",
                    "LaunchConfigurationName": "my-launch-config"
                }
            ],
            "MaxSize": 1,
            "VPCZoneIdentifier": "subnet-YYYYYYYY",
            "TerminationPolicies": [
                "Default"
            ],
            "LaunchConfigurationName": "my-launch-config",
            "CreatedTime": "2016-06-12T01:00:35.418Z",
            "AvailabilityZones": [
                "us-east-1c"
            ],
            "HealthCheckType": "EC2",
            "NewInstancesProtectedFromScaleIn": false
        }
    ]
}
```

##### Cleaning up.

Delete your lifecycle hook.

```
aws autoscaling delete-lifecycle-hook \
  --lifecycle-hook-name my-lifecycle-hook \
  --auto-scaling-group-name my-auto-scaling-group
```

TODO: Need to add more commands to clean up.