```puml
@startuml
title SmartHome Context Diagram

top to bottom direction
!includeurl https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/master/C4_Component.puml

Person(User, "User", "Controls smart home devices")
Person(Admin, "Administrator", "Manages system config")
System(WebUI, "Web Interface", "React + Go")
System(Backend, "Backend Service", "Go")
System(PostgreSQL, "Database", "PostgreSQL")
System_Ext(MQTT_Broker, "MQTT Broker", "IoT")
System_Ext(Heating_System, "Heating System", "Smart HVAC")
System_Ext(Thermal_Sensor, "Thermal Sensor", "IoT")
System_Ext(Notification_Service, "Notification Service", "SendGrid/Twilio")

Rel(User, WebUI, "Uses", "HTTPS/WebSocket")
Rel(WebUI, Backend, "API calls", "REST/HTTP")
Rel(Backend, MQTT_Broker, "Device control", "MQTT")
Rel(MQTT_Broker, Heating_System, "Sends commands", "MQTT")
Rel(Thermal_Sensor, MQTT_Broker, "Sends data", "MQTT")
Rel(Backend, PostgreSQL, "Stores data", "SQL")
Rel(Backend, Notification_Service, "Sends alerts", "API")
Rel(Admin, WebUI, "Manages", "HTTPS")
@enduml
``` 