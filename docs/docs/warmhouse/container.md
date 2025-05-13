```puml
@startuml
title WarmHouse Container Diagram
!includeurl https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/master/C4_Container.puml

Person(user, "User", "Manages home heating via web interface")
Person(admin, "Administrator", "Manages devices and users")

System_Boundary(WarmHouseSystem, "Warm House Ecosystem") {
    Container(web_ui, "Web Interface", "Go + React", "User dashboard and controls")
    Container(backend, "Backend Service", "Go", "Business logic and device communication")
    Container(postgres, "Database", "PostgreSQL", "Stores users, devices and telemetry")
}

System_Ext(thermal_sensor, "Thermal Sensor", "IoT temperature sensor")
System_Ext(heating_system, "Heating System", "Smart HVAC controller")

' ===== Связи =====
Rel(user, web_ui, "Uses", "HTTPS")
Rel(web_ui, backend, "API calls", "REST/HTTP")
Rel(backend, postgres, "Reads/writes data", "SQL")
Rel(backend, heating_system, "Controls", "MQTT")
Rel(backend, thermal_sensor, "Receives data", "MQTT")
Rel(admin, web_ui, "Manages system", "HTTPS")
Rel(admin, thermal_sensor, "Registers devices", "MQTT-over-TLS")

@enduml
``` 