```puml
@startuml
title SmartHome Component Diagram
!includeurl https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/master/C4_Component.puml

Person(User, "User", "Home owner")
Person(Admin, "Admin", "System admin")

System_Boundary(SmartHome_System, "SmartHome Ecosystem") {
    Container(WebUI, "Web Interface") {
        Component(LoginComponent, "Auth Component", "OAuth2/JWT")
        Component(DashboardComponent, "Dashboard UI", "Real-time metrics")
        Component(ControlPanel, "Control Panel", "Device management")
    }
    
    Container(Backend, "Backend Service") {
        Component(AuthController, "AuthController", "Handles auth")
        Component(DeviceController, "DeviceController", "Device commands")
        Component(TelemetryController, "TelemetryController", "Data ingestion")
        
        Component(DeviceService, "DeviceService", "State management")
        Component(TelemetryService, "TelemetryService", "Data processing")
        Component(AutomationService, "AutomationService", "Rules engine")
        
        Component(DeviceRepo, "DeviceRepository", "SQL operations")
        Component(TelemetryRepo, "TelemetryRepository", "Timeseries")
        Component(UserRepo, "UserRepository", "User management")
    }
    
    Container(Database, "Database") {
        Component(DevicesTable, "devices", "Device metadata")
        Component(TelemetryTable, "telemetry", "Sensor data")
        Component(UsersTable, "users", "User data")
    }
}

System_Ext(MQTT_Broker, "MQTT Broker", "IoT")
System_Ext(Heating_System, "Heating System", "Smart HVAC")
System_Ext(Thermal_Sensor, "Thermal Sensor", "IoT")
System_Ext(Notification_Service, "Notification Service", "Email/SMS")

' ===== Связи =====
Rel(User, DashboardComponent, "Views metrics", "WebSocket")
Rel(User, ControlPanel, "Controls devices", "HTTPS")
Rel(Admin, ControlPanel, "Manages config", "HTTPS")

Rel(LoginComponent, AuthController, "API calls", "REST")
Rel(ControlPanel, DeviceController, "Commands", "REST")
Rel(DeviceController, DeviceService, "Processes", "gRPC")
Rel(TelemetryController, TelemetryService, "Analyzes", "gRPC")

Rel(AuthController, UserRepo, "Validates users", "gRPC")
Rel(DeviceService, DeviceRepo, "CRUD", "gRPC")
Rel(TelemetryService, TelemetryRepo, "Stores", "gRPC")
Rel(DeviceRepo, DevicesTable, "SQL", "PGX")
Rel(TelemetryRepo, TelemetryTable, "SQL", "PGX")

Rel(DeviceService, MQTT_Broker, "Controls", "MQTT")
Rel(MQTT_Broker, Heating_System, "Commands", "MQTT")
Rel(Thermal_Sensor, MQTT_Broker, "Sends data", "MQTT")
Rel(MQTT_Broker, TelemetryService, "Streams data", "MQTT")

Rel(TelemetryService, AutomationService, "Triggers rules", "Kafka")
Rel(AutomationService, DeviceService, "Executes actions", "gRPC")
Rel(AutomationService, Notification_Service, "Sends alerts", "API")
@enduml
```