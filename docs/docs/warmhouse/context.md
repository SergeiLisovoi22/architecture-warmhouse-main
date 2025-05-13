```puml
@startuml
title WarmHouse Context Diagram

top to bottom direction

!includeurl https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/master/C4_Component.puml

Person(user, "User", "A user of the warmhouse system")
Person(Admin, "Administrator", "An administrator managing the system")
System(Web, "Web interface for users", "System managing of the heating")
System(Backend, "Backend for web interface", "System managing of the heating")
System(postgres, "PostgreSQL", "Реляционная СУБД")
System_Ext(house, "House", "A user's house") 
System_Ext(ThermalSensor, "Thermal sensor", "Thermal sensor for house ")
System_Ext(HeatingSystem, "Heating System", "System managing of the heating")

Rel(user, Web, "Users monitor the temperature from houses", "HTTPS")
Rel(user, Web, "Turn on and turn off the system")
Rel(Web, Backend, "Turn on and turn off the system", "REST/HTTP")
Rel(Backend, HeatingSystem, "Turn on and turn off the system", "MQTT")
Rel(Backend, ThermalSensor, "Gets the data from thermal sensors", "MQTT")
Rel(ThermalSensor, house, "Sensors mounted in houses")
Rel(Web, Backend, "Web reads data form backend")
Rel(Backend, postgres, "Backend writes data to postgres")
Rel(postgres, Backend, "Backend reads data form postgres")
Rel(Admin, Web, "Gives privileges to users", "HTTPS")
Rel(Admin, ThermalSensor, "Installing and removes the sensors", "MQTT")

@enduml
``` 