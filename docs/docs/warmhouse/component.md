```puml
@startuml
title WarmHouse Component Diagram

!includeurl https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/master/C4_Component.puml

Person(user, "User", "Home owner")
Person(admin, "Admin", "System administrator")

System_Boundary(WarmHouseSystem, "Warm House Ecosystem") {
    Container(web_ui, "Web Interface") {
        Component(login, "Auth Service", "Go", "OAuth2 authentication")
        Component(dashboard, "Dashboard UI", "React", "Real-time monitoring")
        Component(controls, "Device Controls", "React", "Heating system management")
    }
    
    Container(backend, "Backend Service") {
        Component(api_gw, "API Gateway", "Go", "Routes requests to services")
        Component(device_mgr, "Device Manager", "Go", "MQTT device communication")
        Component(data_svc, "Data Service", "Go", "Telemetry processing")
        Component(auth_svc, "Auth Service", "Go", "JWT validation")
    }
    
    Container(postgres, "Database") {
        Component(users_db, "Users", "SQL", "Stores credentials")
        Component(telemetry_db, "Telemetry", "SQL", "Time-series data")
        Component(devices_db, "Devices", "SQL", "Device configurations")
    }
}

System_Ext(thermal_sensor, "Thermal Sensor", "IoT device")
System_Ext(heating_system, "Heating System", "Smart controller")

' ===== Связи =====
Rel(user, dashboard, "Views metrics", "HTTPS")
Rel(user, controls, "Adjusts settings", "HTTPS")
Rel(controls, api_gw, "Sends commands", "REST/HTTP")
Rel(api_gw, device_mgr, "Routes requests", "gRPC")
Rel(device_mgr, heating_system, "Controls", "MQTT")
Rel(device_mgr, thermal_sensor, "Reads data", "MQTT")
Rel(device_mgr, data_svc, "Stores metrics", "Protobuf")
Rel(data_svc, telemetry_db, "Persists data", "SQL")
Rel(login, auth_svc, "Validates tokens", "JWT")
Rel(auth_svc, users_db, "Checks credentials", "SQL")
Rel(admin, devices_db, "Manages devices", "SQL")

@enduml
```