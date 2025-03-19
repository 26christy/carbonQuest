# 🚨 Alarm System Microservices

This repository contains three microservices for managing alarms, notifications, and acknowledgment (ACK) states. The services are:

1. **Alarm Service:** Create and manage alarms.  
2. **ACK Service:** Handle alarm acknowledgment to control notification frequency.  
3. **Notification Service:** Manage alarm notifications and webhook integrations.



## 📁 **Microservices Overview**
| Service             | Port | Purpose                             |
|---------------------|------|-------------------------------------|
| Alarm Service       | 8080 | Manage alarms and states            |
| ACK Service         | 8082 | Handle acknowledgment and schedules |
| Notification Service| 8081 | Send notifications                  |


# 📦 **Installation and Setup**

### 🔹 **Step 1: Clone Repository**
```bash
git clone https://github.com/26christy/carbonQuest.git
```

### 🔹 **Step 2: Set Up Environment Variables**
Create a `.env` file for each service in the service folder:

- **Alarm Service:** `.env.alarm-service`
    ```env
    SERVICE_NAME=alarm-service
    PORT=8080
    HOST=localhost
    NOTIFICATION_SERVICE_PORT=8081
    ```
- **ACK Service:** `.env.ack-service`
    ```env
    SERVICE_NAME=ack-service
    ALARM_SERVICE_PORT=8080
    HOST=localhost
    ACK_DURATION=3
    ```
- **Notification Service:** `.env.notification-service`
    ```env
    SERVICE_NAME=notification-service
    NOTIFICATION_SERVICE_PORT=8081
    ACK_DURATION=3
    UNACK_DURATION=3
    HOST=localhost
    ALARM_SERVICE_PORT=8080
    ACK_SERVICE_PORT=8082
    NOTIFIER_TYPE=log
    NOTIFIER_PARAMS=""
    ```


### 🔹 **Step 3: Run Microservices**
```bash
# Run Alarm Service
cd alarm-service
go run main.go

# Run ACK Service
cd ack-service
go run main.go

# Run Notification Service
cd notification-service
go run main.go
```

## 📁 **Directory Structure**
```
alarm-system/
│── alarm-service/          # Manages alarms & states
│   ├── handlers/           # REST API handlers
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── routes.go
│   ├── storage/            # Storage logic & interface
│   │   ├── memory.go       # In-memory implementation
│   │   ├── memory_test.go 
│   │   ├── iface.go        # Defines the storage interface
│   ├── service/            # Business logic (Alarm processing)
|   │   ├── service.go
│   │   ├── service_test.go
│   │   ├── iface.go        # Defines the service interface
│   ├── .env.alarm-service          
│   ├── main.go             # Entry point
│
│── notification-service/   # Handles notifications
│   ├── handlers/           # REST API handlers
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── routes.go
│   ├── service/            # Notification logic (Webhook, log, etc.)
|   │   ├── service.go
│   │   ├── service_test.go
│   │   ├── iface.go        # Defines the service interface
│   │── notifiers/
│   │   ├── iface.go        # Defines the Notifier interface
│   │   ├── createNotifier.go # Factory function to create notifiers
│   │   ├── log.go          # Implements a logger notifier
│   │   ├── webhook.go      # Implements a webhook notifier
│   ├── .env.notification-service
│   ├── main.go             # Entry point
│
│── ack-service/            # Handles acknowledgments
│   ├── handlers/           # REST API handlers
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── routes.go
│   ├── service/            # ACK logic
|   │   ├── service.go
│   │   ├── service_test.go
│   │   ├── iface.go        # Defines the service interface
│   ├── storage/            # Storage for ACK states
│   │   ├── memory.go       # In-memory implementation
│   │   ├── memory_test.go 
│   │   ├── iface.go        # Defines the storage interface
│   ├── .env.ack-service 
│   ├── main.go             # Entry point
│
│── middleware/             
│   ├── middleware.go
│── common/                 # Shared logic & models
│   ├── models/             # Shared data models (Alarm, ACK, Notification)
│   ├── utils/              # Helper functions
│
│── README.md               # API Documentation
```

# 📄 **API Documentation**

---

## 🚨 **1. Alarm Service (Port 8080)**

### **Endpoints:**
| Method | Endpoint           | Description                          |
|--------|---------------------|--------------------------------------|
| POST   | `/alarms`           | Create a new alarm                   |
| GET    | `/alarms`           | Retrieve all alarms                  |
| GET    | `/alarms/{id}`      | Get alarm details by ID              |
| PUT    | `/alarms/{id}`      | Update alarm state by ID             |
| DELETE | `/alarms/{id}`      | Delete an alarm by ID                |

### **Request/Response:**
#### ➔ Create an Alarm
This endpoint allows you to add a new alarm.
```bash
curl --location 'localhost:8080/alarms' \
--header 'Content-Type: application/json' \
--data '{
    "name": "morning-alarm",
    "timestamp": "2025-03-18T22:01:00.000000+05:30"
}'
```
#### ➔ Response:
```json
{
    "id": "e70a2066-7bec-4999-a688-62da73b6187f",
    "name": "morning-alarm",
    "timestamp": "2025-03-18T22:01:00+05:30",
    "status": "triggered",
    "created_at": "2025-03-18T21:59:54.318719+05:30",
    "updated_at": "2025-03-18T21:59:54.318719+05:30"
}
```

#### ➔ Get all Alarms
This endpoint fetches the list of all the alarms
```bash
curl --location 'localhost:8080/alarms'
```
#### ➔ Response:
```json
{
    "alarms": [
        {
            "id": "e70a2066-7bec-4999-a688-62da73b6187f",
            "name": "morning-alarm",
            "timestamp": "2025-03-18T22:01:00+05:30",
            "status": "acknowledged",
            "created_at": "2025-03-18T21:59:54.318719+05:30",
            "updated_at": "2025-03-18T22:12:25.25511+05:30"
        }
    ]
}
```

#### ➔ Get an Alarm by ID
This endpoint allows you to get the alarm details by ID
```bash
curl --location 'localhost:8080/alarms/7f7ad0ae-ab13-42fe-8829-fcdcb3c78c67'
```
#### ➔ Response:
```json
{
    "id": "7f7ad0ae-ab13-42fe-8829-fcdcb3c78c67",
    "name": "morning-alarm",
    "timestamp": "2025-03-17T14:15:40.263543+05:30",
    "status": "triggered",
    "created_at": "2025-03-17T14:24:34.739461+05:30",
    "updated_at": "2025-03-17T14:24:34.739461+05:30"
}
```
#### ➔ Update an Alarm by ID
This endpoint allows you to update the alarm status by ID
```bash
curl --location --request PUT 'localhost:8080/alarms/e70a2066-7bec-4999-a688-62da73b6187f' \
--header 'Content-Type: application/json' \
--data '{
    "status": "active"
}
'
```
#### ➔ Response:
```json
{
    "id": "e70a2066-7bec-4999-a688-62da73b6187f",
    "name": "morning-alarm",
    "timestamp": "2025-03-17T14:15:40.263543+05:30",
    "status": "active",
    "created_at": "2025-03-17T14:24:34.739461+05:30",
    "updated_at": "2025-03-17T14:24:34.739461+05:30"
}
```

#### ➔ Delete an Alarm by ID
This endpoint allows you to delete an alarm by ID
```bash
curl --location --request DELETE 'localhost:8080/alarms/7f7ad0ae-ab13-42fe-8829-fcdcb3c78c67'
```
#### ➔ Response:
```json
```
---

## ✅ **2. ACK Service (Port 8082)**

### **Endpoints:**
| Method | Endpoint               | Description                                |
|--------|------------------------|--------------------------------------------|
| POST   | `/ack/{alarm_id}`      | Acknowledge an alarm                       |
| GET    | `/ack/{alarm_id}`      | Get the acked alarm details                |


### **Request/Response:**
#### ➔ Ack an Alarm by ID
This endpoint allows you to ack an alarm after getting a notification. The duration between each notification post ack can be set in the environment variable (ACK_DURATION).
```bash
curl --location --request POST 'http://localhost:8082/ack/e70a2066-7bec-4999-a688-62da73b6187f'
```
#### ➔ Response:
```json
{
    "message": "alarm successfully acknowledged"
}
```

#### ➔ Get an Acked Alarm by ID
This endpoint fetches the details of an acked alarm by ID
```bash
curl --location 'http://localhost:8082/ack/e70a2066-7bec-4999-a688-62da73b6187f'
```
#### ➔ Response:
```json
{
    "acked_at": "2025-03-18T22:12:25.256537+05:30",
    "alarm_id": "e70a2066-7bec-4999-a688-62da73b6187f",
    "next_notification_at": "2025-03-19T22:12:25.256537+05:30",
    "should_notify": true
}
```

---

## 🔔 **3. Notification Service (Port 8081)**

### **Endpoints:**
| Method | Endpoint                     | Description                           |
|--------|-------------------------------|--------------------------------------|
| POST   | `/notify/register-notifier`   | Trigger a manual notification method |
| POST   | `/notify`                     | Send an alarm notification           |

### **Request/Response:**
#### ➔ Register Notifier
This API allows user to register a notification type (logs, webhook etc). Notifier is registered when the service is started by default. Which ever type is mentioned in the environment variable (NOTIFIER_TYPE) is registered as the method of notification. However if the user wants to change it without restarting the service then this API can be called.
```bash
curl --location 'http://localhost:8081/notify/register-notifier' \
--header 'Content-Type: application/json' \
--data '{"type": "log", "param": ""}'
```
#### ➔ Response:
```json
{
    "message": "Notifier registered successfully"
}
```

#### ➔ Send an alarm notification
This API allows user to send an alarm notication using the registered notification method
```bash
curl --location 'http://localhost:8081/notify/' \
--header 'Content-Type: application/json' \
--data '{
    "alarm_id": "554e76e4-fc22-4eea-b6c6-616e5d4c8caf",
    "name": "morning-alarm",
    "type": "active",
    "timestamp": "2025-03-19T11:07:00+05:30"
    }'
```
#### ➔ Response:
```json
{
    "message": "Notification received successfully"
}
```

### 🔹 **Unit Tests**
Run tests using:
```bash
go test ./... -v
```