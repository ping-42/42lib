package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"github.com/ping-42/42lib/config"
	"github.com/ping-42/42lib/db/models"
	"gorm.io/gorm"
)

func MigrateAndSeed(gormClient *gorm.DB) {
	err := migrate(gormClient)
	if err != nil {
		panic(err)
	}
}

// https://docs.timescale.com/getting-started/latest/
/*
-- Drop everything from public schema:
DO $$
DECLARE
    current_table text;
BEGIN
    FOR current_table IN (SELECT table_name FROM information_schema.tables WHERE table_schema = 'public')
    LOOP
        EXECUTE 'DROP TABLE IF EXISTS public.' || current_table || ' CASCADE';
    END LOOP;
END $$;
*/

func migrate(db *gorm.DB) error {
	// TODOs:
	// 1.prior to going live, the structs referenced here should
	//   be copy-pasted explicitly to preserve a changelog. Additionally,
	//   all migrations should be squashed into one.
	// 2. Add indexes for rank queries
	migrations := []*gormigrate.Migration{
		{
			// IDs must be unique
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
				type Task struct {
					ID                   uuid.UUID                 `gorm:"type:uuid;primary_key;" json:"id"`
					TaskTypeID           uint64                    //FK to TaskType.id
					TaskType             models.LvTaskType         `gorm:"foreignKey:TaskTypeID"`
					TaskStatusID         uint8                     //FK to TaskType.id
					TaskStatus           models.LvTaskStatus       `gorm:"foreignKey:TaskStatusID"`
					SensorID             uuid.UUID                 //FK to Sensor.id
					Sensor               models.Sensor             `gorm:"foreignKey:SensorID"`
					ClientSubscriptionID uint64                    //FK to ClientSubscription.id
					ClientSubscription   models.ClientSubscription `gorm:"foreignKey:ClientSubscriptionID"`
					Opts                 []byte                    `gorm:"type:jsonb"`
				}
				err := tx.Migrator().CreateTable(
					&models.Sensor{},
					&models.LvTaskType{},
					&models.LvTaskStatus{},
					&models.LvProtocol{},
					&models.Client{},
					&models.ClientSubscription{},
					&Task{},
					&models.TsHostRuntimeStat{},
					&models.TsDnsResult{},
					&models.TsDnsResultAnswer{},
					&models.TsHttpResult{},
					&models.TsIcmpResult{},
					&models.TsHopResult{},
					&models.TsTracerouteResult{},
				)
				if err != nil {
					return err
				}

				// hypertables for timeseries data
				err = tx.Exec(`
					SELECT create_hypertable('ts_host_runtime_stats', by_range('time'));
					--
					SELECT create_hypertable('ts_dns_results', by_range('time'));
					SELECT create_hypertable('ts_dns_results_answer', by_range('time'));
					--
					SELECT create_hypertable('ts_http_results', by_range('time'));
					--
					SELECT create_hypertable('ts_icmp_results', by_range('time'));
					--
					SELECT create_hypertable('ts_traceroute_results', by_range('time'));`).Error
				if err != nil {
					return err
				}

				// indices
				err = tx.Exec(`
					CREATE INDEX idx_runtime_sensor_time ON ts_host_runtime_stats (sensor_id, time DESC);
					--
					CREATE INDEX idx_dns_results_sensor_time ON ts_dns_results (sensor_id, time DESC);
					CREATE INDEX idx_dns_results_answer_sensor_time ON ts_dns_results_answer (sensor_id, time DESC);
					-- CREATE INDEX idx_dns_results_task ON ts_dns_results (task_id);
					-- CREATE INDEX idx_dns_results_answer_task ON ts_dns_results_answer (task_id);
					--
					CREATE INDEX idx_http_results_sensor_time ON ts_http_results (sensor_id, time DESC);
					--
					CREATE INDEX idx_icmp_results_sensor_time ON ts_icmp_results (sensor_id, time DESC);
					--
					CREATE INDEX idx_traceroute_results_sensor_time ON ts_traceroute_results (sensor_id, time DESC);`).Error
				if err != nil {
					return err
				}

				// lookup values
				err = tx.Exec(`
					INSERT INTO lv_task_types(id, type) VALUES (1, 'DNS_TASK');
					INSERT INTO lv_task_types(id, type) VALUES (2, 'ICMP_TASK');
					INSERT INTO lv_task_types(id, type) VALUES (3, 'HTTP_TASK');
					INSERT INTO lv_task_types(id, type) VALUES (4, 'TRACEROUTE_TASK');`).Error
				if err != nil {
					return err
				}

				err = tx.Exec(`
					INSERT INTO lv_protocols(id, type) VALUES (1, 'TCP');
					INSERT INTO lv_protocols(id, type) VALUES (2, 'UDP');`).Error
				if err != nil {
					return err
				}

				err = tx.Exec(`
					INSERT INTO lv_task_statuses(id, status) VALUES (1, 'INITIATED_BY_SCHEDULER');
					INSERT INTO lv_task_statuses(id, status) VALUES (2, 'PUBLISHED_TO_REDIS_BY_SCHEDULER');
					INSERT INTO lv_task_statuses(id, status) VALUES (3, 'RECEIVED_BY_SERVER');
					INSERT INTO lv_task_statuses(id, status) VALUES (4, 'SENT_TO_SENSOR_BY_SERVER');
					INSERT INTO lv_task_statuses(id, status) VALUES (5, 'RECEIVED_BY_SENSOR');
					INSERT INTO lv_task_statuses(id, status) VALUES (6, 'RESULTS_SENT_TO_SERVER_BY_SENSOR');
					INSERT INTO lv_task_statuses(id, status) VALUES (7, 'RESULTS_RECEIVED_BY_SERVER');
					INSERT INTO lv_task_statuses(id, status) VALUES (8, 'DONE');
					INSERT INTO lv_task_statuses(id, status) VALUES (9, 'ERROR');`).Error
				return err
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Rollback().Error
			},
		},
		{
			ID: "for-squash-1",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().CreateTable(&models.SensorRank{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Rollback().Error
			},
		},
		{
			ID: "for-squash-2",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().AddColumn(&models.Task{}, "CreatedAt")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Rollback().Error
			},
		},
	}

	// dev demo data
	if config.CurrentEnv() == config.Dev {
		migrations = append(migrations, &gormigrate.Migration{
			ID: "dev-seeds-01",
			Migrate: func(tx *gorm.DB) error {

				err := tx.Exec(devSeeds).Error
				if err != nil {
					return err
				}
				return err
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Rollback().Error
			},
		})

		migrations = append(migrations, &gormigrate.Migration{
			ID: "dev-seeds-sensor-ranks",
			Migrate: func(tx *gorm.DB) error {
				return tx.Exec(`INSERT INTO sensor_ranks(id, sensor_id, rank, distribution_rank, created_at) VALUES (1, 'b9dc3d20-256b-4ac7-8cae-2f6dc962e183', 5, 0, now());`).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Rollback().Error
			},
		})
	}

	m := gormigrate.New(db, &gormigrate.Options{UseTransaction: true}, migrations)
	return m.Migrate()
}

const devSeeds = `
-----sensors-----
INSERT INTO sensors(id, name, location, secret) VALUES ('b9dc3d20-256b-4ac7-8cae-2f6dc962e183', 'Test Sensor', 'Sofia, Bulgaria', 'sensorSecret123!');
-----client-----
INSERT INTO clients(id, name, email) VALUES (1, 'Test Client', 'test_client@gmail.com');
-----client_subscriptions-----
INSERT INTO client_subscriptions(id, client_id, task_type_id, tests_count_subscribed, tests_count_executed, period, last_execution_completed, opts, is_active) VALUES (1, 1, 1, 9999, 0, 60, NULL, '{"Host":"https://google.com", "Proto":"udp"}', TRUE);
INSERT INTO client_subscriptions(id, client_id, task_type_id, tests_count_subscribed, tests_count_executed, period, last_execution_completed, opts, is_active) VALUES (2, 1, 2, 9999, 0, 60, NULL, '{"TargetDomain":"","TargetIPs":["127.0.0.1"],"Count":3,"Payload":"MDgwOTBhMGIwYzBkMGUwZjEwMTExMjEzMTQxNTE2MTcxODE5MWExYjFjMWQxZTFmMjAyMTIyMjMyNDI1MjYyNzI4MjkyYTJiMmMyZDJlMmYzMDMxMzIzMzM0MzUzNjM3"}', TRUE);
INSERT INTO client_subscriptions(id, client_id, task_type_id, tests_count_subscribed, tests_count_executed, period, last_execution_completed, opts, is_active) VALUES (3, 1, 3, 9999, 0, 60, NULL, '{"TargetDomain":"https://google.com","HttpMethod":"GET","RequestHeaders":{"Content-Type":["application/json"]},"RequestBody":"c29tZSB0ZXN0IGJvZHk="}', TRUE);
`

//INSERT INTO client_subscriptions(id, client_id, task_type_id, tests_count_subscribed, tests_count_executed, period, last_execution_completed, opts, is_active) VALUES (1, 1, 2, 9999, 0, 60, NULL, '{"TargetDomain":null,"TargetIPs":["127.0.0.1"],"Count":3,"Payload":"MDgwOTBhMGIwYzBkMGUwZjEwMTExMjEzMTQxNTE2MTcxODE5MWExYjFjMWQxZTFmMjAyMTIyMjMyNDI1MjYyNzI4MjkyYTJiMmMyZDJlMmYzMDMxMzIzMzM0MzUzNjM3"}', TRUE);
