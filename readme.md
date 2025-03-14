```
alarm-system/
│── alarm-service/          # Manages alarms & states
│   ├── handlers/           # REST API handlers
│   ├── storage/            # Storage logic & interface
│   │   ├── memory.go       # In-memory implementation
│   │   ├── storage.go      # Defines the storage interface
│   ├── service/            # Business logic (Alarm processing)
│   ├── events/             # Publishes events (alarm triggered)
│   ├── main.go             # Entry point
│
│── notification-service/   # Handles notifications
│   ├── handlers/           # REST API handlers
│   ├── service/            # Notification logic (Webhook, Email, etc.)
│   ├── events/             # Listens to alarm events
│   ├── main.go             # Entry point
│
│── ack-service/            # Handles acknowledgments
│   ├── handlers/           # REST API handlers
│   ├── service/            # ACK logic (notification throttling)
│   ├── storage/            # Storage for ACK states
│   │   ├── memory.go       # In-memory implementation
│   │   ├── storage.go      # Defines the storage interface
│   ├── events/             # Listens to alarm events
│   ├── main.go             # Entry point
│
│── common/                 # Shared logic & models
│   ├── models/             # Shared data models (Alarm, ACK, Notification)
│   ├── utils/              # Helper functions (Logging, Time)
│
│── tests/                  # Unit & integration tests
│── README.md               # API Documentation
```