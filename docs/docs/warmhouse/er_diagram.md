@startuml
' ER-диаграмма для экосистемы "Умный Дом"

hide circle
skinparam linetype ortho

' ===== Сущности =====
entity "User" as user {
*id: INT <<PK>>
--
*email: VARCHAR(255)
*password_hash: VARCHAR(255)
created_at: TIMESTAMP
}

entity "House" as house {
*id: INT <<PK>>
--
*user_id: INT <<FK>>
address: VARCHAR(500)
description: TEXT
}

entity "DeviceType" as device_type {
*id: INT <<PK>>
--
*name: VARCHAR(100)
description: TEXT
protocol: VARCHAR(50) [MQTT, HTTP, Zigbee]
}

entity "Device" as device {
*id: INT <<PK>>
--
*type_id: INT <<FK>>
*house_id: INT <<FK>>
serial_number: VARCHAR(100) <<UQ>>
status: ENUM('online', 'offline', 'error')
last_seen: TIMESTAMP
}

entity "TelemetryData" as telemetry {
*id: BIGINT <<PK>>
--
*device_id: INT <<FK>>
*value: DOUBLE
*timestamp: TIMESTAMP
unit: VARCHAR(20)
}

entity "AutomationRule" as automation {
*id: INT <<PK>>
--
*device_id: INT <<FK>>
condition: TEXT
action: TEXT
is_active: BOOLEAN
}

' ===== Связи =====
user ||--o{ house : "owns"
house ||--o{ device : "contains"
device_type ||--o{ device : "defines_type"
device ||--o{ telemetry : "generates"
device ||--o{ automation : "has_rules"

' ===== Комментарии =====
note top of user
Пользователи системы
email - уникальный идентификатор
end note

note right of device_type
Типы устройств (термостат, датчик движения и т.д.)
protocol - поддерживаемый протокол связи
end note

note bottom of automation
Правила автоматизации:
Пример условия: "temperature > 25"
Пример действия: "turn_on AC"
end note
@enduml