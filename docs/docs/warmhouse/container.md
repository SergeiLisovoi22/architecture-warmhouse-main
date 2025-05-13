```puml
@startuml
title SmartHome Container Diagram
!includeurl https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/master/C4_Component.puml

Person(User, "User", "Controls devices via UI")
Person(Admin, "Admin", "Manages config")

System_Boundary(SmartHome_System, "SmartHome Ecosystem") {
    Container(WebUI, "Web Interface", "React + Go") {
        Component(Auth_UI, "Auth UI", "Login/register")
        Component(Device_Control, "Device Control", "Real-time UI")
        Component(Admin_Console, "Admin Console", "System config")
    }
    
    Container(Backend, "Backend", "Go") {
        Component(Controllers, "Controllers", "REST/gRPC")
        Component(Services, "Services", "Business logic")
        Component(TelemetryService, "Telemetry Service", "Data processing")
        Component(Repositories, "Repositories", "Data access")
        Component(AutomationService, "Automation Service", "Rules engine")
    }
    
    Container(PostgreSQL, "Database", "PostgreSQL", "Stores data")
}

System_Ext(IoT_Devices, "IoT Devices", "MQTT")
System_Ext(Notification_Service, "Notification Service", "Email/SMS")

Rel(User, WebUI, "HTTPS/WebSocket")
Rel(WebUI, Controllers, "API calls", "REST/HTTP")
Rel(Controllers, Services, "Process requests", "gRPC")
Rel(Services, Repositories, "Data operations", "gRPC")
Rel(Repositories, PostgreSQL, "SQL queries", "PGX")
Rel(Services, IoT_Devices, "Device comms", "MQTT")
Rel(TelemetryService, AutomationService, "Triggers rules", "Kafka")
Rel(AutomationService, Notification_Service, "Alerts", "API")
Rel(Admin, Admin_Console, "HTTPS")
@enduml


``` 