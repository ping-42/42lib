package models

import (
	"time"
)

type ClientSubscription struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ClientID uint64 //FK to Client.id
	Client   Client `gorm:"foreignKey:ClientID"`

	TaskTypeID uint64     //FK to TaskType.id
	TaskType   LvTaskType `gorm:"foreignKey:TaskTypeID"`

	TestsCountSubscribed int
	TestsCountExecuted   int

	// the tasks should be assined on each X period
	Period time.Duration

	LastExecutionCompleted time.Time

	Opts []byte `gorm:"type:jsonb"`

	IsActive bool
}

// // JsonbCustom represents the custom type for Opts/Jsonb field
// type JsonbCustom map[string]interface{}

// // Value Marshal
// func (jsonField JsonbCustom) Value() (driver.Value, error) {
// 	return json.Marshal(jsonField)
// }

// // Scan Unmarshal
// func (jsonField *JsonbCustom) Scan(value interface{}) error {
// 	fmt.Printf("Type of value: %T\n", value)
// 	switch v := value.(type) {
// 	case []byte:
// 		// If value is []byte, unmarshal it directly
// 		return json.Unmarshal(v, jsonField)
// 	case string:
// 		// If value is string, convert it to []byte and then unmarshal
// 		return json.Unmarshal([]byte(v), jsonField)
// 	case map[string]interface{}:
// 		// If value is map[string]interface{}, assign it directly
// 		*jsonField = v
// 	case int:
// 		// If value is int, create a map with a single key-value pair
// 		*jsonField = JsonbCustom{"value": v}
// 	default:
// 		// Unsupported type
// 		return fmt.Errorf("unsupported type: %T", value)
// 	}
// 	return nil
// 	// return json.Unmarshal(data, &jsonField)
// }

// func (jsonField JsonbCustom) ToOpts() (Opts, error) {
// 	opts := Opts{}

// 	// Convert JsonbCustom to JSON bytes
// 	jsonBytes, err := json.Marshal(jsonField)
// 	if err != nil {
// 		return opts, err
// 	}

// 	// Unmarshal JSON bytes to Opts struct
// 	err = json.Unmarshal(jsonBytes, &opts)
// 	if err != nil {
// 		return opts, err
// 	}

// 	return opts, nil
// }
