```puml
@startuml
title SmartHome Code Level Diagram (C4 Level 4)

!includeurl https://raw.githubusercontent.com/RicardoNiepel/C4-PlantUML/master/C4_Component.puml

class AuthController {
  +ValidateToken(token: string): User
  +Login(credentials: Credentials): Session
  +Register(userData: UserDTO): User
}

class DeviceController {
  +SendCommand(deviceId: string, command: Command): void
  +GetDeviceStatus(deviceId: string): DeviceStatus
  +UpdateDeviceConfig(config: DeviceConfig): void
}

class TelemetryController {
  +GetHistoricalData(deviceId: string, period: TimeRange): TelemetryData[]
  +GetRealtimeStream(deviceId: string): IObservable<TelemetryData>
}

class DeviceService {
  -_mqttClient: MQTTClient
  -_deviceRepo: DeviceRepo
  +ExecuteCommand(command: Command): void
  +UpdateDeviceState(deviceId: string, state: DeviceState): void
  +SyncWithIoTBroker(): void
  +GetDevice(deviceId: string): Device
  +UpdateDevice(device: Device): void
}

class TelemetryService {
  -_telemetryRepo: TelemetryRepo
  +ProcessSensorData(data: TelemetryRawData): void
  +AnalyzeTrends(deviceId: string): AnalyticsReport
  +TriggerAnomalyAlerts(): void
  +SaveData(data: TelemetryData): void
}

class AutomationService {
  -_ruleEngine: RuleEngine
  +EvaluateRules(context: RuleContext): void
  +ExecuteAutomationAction(action: AutomationAction): void
  +NotifyMaintenance(alert: Alert): void
}

class DeviceRepo {
  +GetDevice(deviceId: string): Device
  +UpdateDevice(device: Device): void
  +GetAllDevices(): Device[]
}

class MQTTClient {
  +Connect(brokerUrl: string): void
  +Subscribe(topic: string): void
  +Publish(topic: string, payload: string): void
}

class Device {
  +Id: string
  +Name: string
  +Type: DeviceType
  +Status: DeviceStatus
  +Config: DeviceConfig
  +LastSeen: DateTime
}

class TelemetryData {
  +DeviceId: string
  +Timestamp: DateTime
  +Value: float
  +Unit: string
}

class AutomationRule {
  +Condition: string
  +Actions: AutomationAction[]
  +IsEnabled: bool
  +Evaluate(context: RuleContext): bool
}

class NotificationService {
  +SendEmail(to: string, template: EmailTemplate): void
  +SendSMS(phone: string, message: string): void
}

AuthController "1" --> "1" UserRepo : Validates via >
DeviceController "1" --> "1" DeviceService : Uses >
TelemetryController "1" --> "1" TelemetryService : Uses >

DeviceService "1" --> "1" MQTTClient : Manages >
DeviceService "1" --> "1" DeviceRepo : Persists data >
TelemetryService "1" --> "1" TelemetryRepo : Stores >
TelemetryService "1" --> "1" AutomationService : Triggers via Kafka >

AutomationService "1" --> "1" DeviceService : Controls >
AutomationService "1" --> "1" NotificationService : Triggers via API >

MQTTClient "1" ..> "n" Device : Listens/Updates
@enduml
```