package sensor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Creds struct {
	SensorId uuid.UUID
	Secret   string
}

func (sc Creds) GetSensorEnvToken() (token string, err error) {
	j, err := json.Marshal(sc)
	if err != nil {
		err = fmt.Errorf("getSensorEnvToken Marshal err:%v", err)
		return
	}
	token = base64.StdEncoding.EncodeToString(j)
	return
}
