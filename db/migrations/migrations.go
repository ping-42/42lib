package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
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
				err := tx.Migrator().CreateTable(
					&models.Organization{},
					&models.Sensor{},
					&models.LvTaskType{},
					&models.LvTaskStatus{},
					&models.LvProtocol{},
					&models.Subscription{},
					&models.Task{},
					&models.TsHostRuntimeStat{},
					&models.TsDnsResult{},
					&models.TsDnsResultAnswer{},
					&models.TsHttpResult{},
					&models.TsIcmpResult{},
					&models.TsTracerouteResultHop{},
					&models.TsTracerouteResult{},
					&models.SensorRank{},
					&models.LvUserGroup{},
					&models.LvPermission{},
					&models.User{},
					&models.PermissionToUserGroup{},
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
                    SELECT create_hypertable('ts_traceroute_results', by_range('time'));
                    SELECT create_hypertable('ts_traceroute_results_hop', by_range('time'));`).Error
				if err != nil {
					return err
				}

				// indices
				err = tx.Exec(`
                    CREATE INDEX idx_runtime_sensor_time ON ts_host_runtime_stats (sensor_id, time DESC);
                    CREATE INDEX idx_runtime_sensor_id   ON ts_host_runtime_stats (sensor_id);
                    --
                    CREATE INDEX idx_dns_results_sensor_time ON ts_dns_results (sensor_id, time DESC);
                    CREATE INDEX idx_dns_results_sensor_id   ON ts_dns_results (sensor_id);
                    CREATE INDEX idx_dns_results_answer_sensor_time ON ts_dns_results_answer (sensor_id, time DESC);
                    CREATE INDEX idx_dns_results_answer_sensor_id   ON ts_dns_results_answer (sensor_id);
                    --
                    CREATE INDEX idx_http_results_sensor_time ON ts_http_results (sensor_id, time DESC);
                    CREATE INDEX idx_http_results_sensor_id   ON ts_http_results (sensor_id);
                    --
                    CREATE INDEX idx_icmp_results_sensor_time ON ts_icmp_results (sensor_id, time DESC);
                    CREATE INDEX idx_icmp_results_sensor_id   ON ts_icmp_results (sensor_id);
					--
                    CREATE INDEX idx_traceroute_results_sensor_time ON ts_traceroute_results (sensor_id, time DESC);
                    CREATE INDEX idx_traceroute_results_sensor_id   ON ts_traceroute_results (sensor_id);
                    CREATE INDEX idx_traceroute_results_hop_sensor_time ON ts_traceroute_results_hop (sensor_id, time DESC);
                    CREATE INDEX idx_traceroute_results_hop_sensor_id   ON ts_traceroute_results_hop (sensor_id);
                    CREATE INDEX idx_traceroute_results_task     ON ts_traceroute_results     (task_id);
                    CREATE INDEX idx_traceroute_results_hop_task ON ts_traceroute_results_hop (task_id);
					`).Error
				if err != nil {
					return err
				}

				// lookup values
				err = tx.Exec(`
                    INSERT INTO lv_task_types(id, type) VALUES (1, 'DNS_TASK');
                    INSERT INTO lv_task_types(id, type) VALUES (2, 'ICMP_TASK');
                    INSERT INTO lv_task_types(id, type) VALUES (3, 'HTTP_TASK');
                    INSERT INTO lv_task_types(id, type) VALUES (4, 'TRACEROUTE_TASK');
					--
					INSERT INTO lv_protocols(id, type) VALUES (1, 'TCP');
                    INSERT INTO lv_protocols(id, type) VALUES (2, 'UDP');
					--
					INSERT INTO lv_task_statuses(id, status) VALUES (1, 'INITIATED_BY_SCHEDULER');
                    INSERT INTO lv_task_statuses(id, status) VALUES (2, 'PUBLISHED_TO_REDIS_BY_SCHEDULER');
                    INSERT INTO lv_task_statuses(id, status) VALUES (3, 'RECEIVED_BY_SERVER');
                    INSERT INTO lv_task_statuses(id, status) VALUES (4, 'SENT_TO_SENSOR_BY_SERVER');
                    INSERT INTO lv_task_statuses(id, status) VALUES (5, 'RECEIVED_BY_SENSOR');
                    INSERT INTO lv_task_statuses(id, status) VALUES (6, 'RESULTS_SENT_TO_SERVER_BY_SENSOR');
                    INSERT INTO lv_task_statuses(id, status) VALUES (7, 'RESULTS_RECEIVED_BY_SERVER');
                    INSERT INTO lv_task_statuses(id, status) VALUES (8, 'DONE');
                    INSERT INTO lv_task_statuses(id, status) VALUES (9, 'ERROR');
					--
					INSERT INTO lv_user_groups(id, group_name) VALUES (1, 'root');
					INSERT INTO lv_user_groups(id, group_name) VALUES (2, 'admin');
					INSERT INTO lv_user_groups(id, group_name) VALUES (3, 'user');
					--
					INSERT INTO lv_permissions(id, permission) VALUES (1, 'read');
					INSERT INTO lv_permissions(id, permission) VALUES (2, 'create');
					INSERT INTO lv_permissions(id, permission) VALUES (3, 'update');
					INSERT INTO lv_permissions(id, permission) VALUES (4, 'delete');
					INSERT INTO lv_permissions(id, permission) VALUES (5, 'create_organization_user');
					--
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (1, 1);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (1, 2);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (1, 3);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (1, 4);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (1, 5);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (2, 1);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (2, 2);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (2, 3);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (2, 4);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (2, 5);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (3, 1);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (3, 2);
					INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (3, 3);
					`).Error

				return err
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Rollback().Error
			},
		},
	}

	// dev demo data
	// if config.CurrentEnv() == config.Dev {
	// 	migrations = append(migrations, &gormigrate.Migration{
	// 		ID: "dev-seeds-01",
	// 		Migrate: func(tx *gorm.DB) error {

	// 			err := tx.Exec(devSeeds).Error
	// 			if err != nil {
	// 				return err
	// 			}
	// 			return err
	// 		},
	// 		Rollback: func(tx *gorm.DB) error {
	// 			return tx.Rollback().Error
	// 		},
	// 	})

	// 	migrations = append(migrations, &gormigrate.Migration{
	// 		ID: "dev-seeds-sensor-ranks",
	// 		Migrate: func(tx *gorm.DB) error {
	// 			return tx.Exec(`INSERT INTO sensor_ranks(id, sensor_id, rank, distribution_rank, created_at) VALUES (1, 'b9dc3d20-256b-4ac7-8cae-2f6dc962e183', 5, 0, now());`).Error
	// 		},
	// 		Rollback: func(tx *gorm.DB) error {
	// 			return tx.Rollback().Error
	// 		},
	// 	})

	// 	migrations = append(migrations, &gormigrate.Migration{
	// 		ID: "dev-seeds-organization-to-sensors-relation",
	// 		Migrate: func(tx *gorm.DB) error {

	// 			err := tx.Exec(`
	// 			INSERT INTO organizations(id, name) VALUES ('10e76fbd-77cf-4470-bcb0-25c72b09a511', 'test seed org');`)
	// 			if err.Error != nil {
	// 				return err.Error
	// 			}

	// 			return tx.Exec(`
	// 			INSERT INTO users(id, wallet_address, user_group_id, organization_id) VALUES ('63e76fbd-77cf-4470-bcb0-25c72b09a504', '0xd694cfc8c66e34371eae8ebe03d54867e5c6cec4', 1, '10e76fbd-77cf-4470-bcb0-25c72b09a511');
	// 			UPDATE sensors SET organization_id = '10e76fbd-77cf-4470-bcb0-25c72b09a511', is_active = true, created_at = now() WHERE id='b9dc3d20-256b-4ac7-8cae-2f6dc962e183';`).Error
	// 		},
	// 		Rollback: func(tx *gorm.DB) error {
	// 			return tx.Rollback().Error
	// 		},
	// 	})

	// 	migrations = append(migrations, &gormigrate.Migration{
	// 		ID: "for-squash-users-new-fields",
	// 		Migrate: func(tx *gorm.DB) error {
	// 			err := tx.Migrator().AddColumn(&models.User{}, "IsActive")
	// 			if err != nil {
	// 				return err
	// 			}
	// 			err = tx.Migrator().AddColumn(&models.User{}, "IsValidated")
	// 			if err != nil {
	// 				return err
	// 			}
	// 			err = tx.Migrator().AddColumn(&models.User{}, "CreatedAt")
	// 			if err != nil {
	// 				return err
	// 			}
	// 			err = tx.Migrator().AddColumn(&models.User{}, "LastLoginAt")
	// 			if err != nil {
	// 				return err
	// 			}
	// 			err = tx.Exec(`
	// 				-- add permissions to users group
	// 				INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (3, 1);
	// 				INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (3, 2);
	// 				INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (3, 3);
	// 				INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (2, 4);
	// 				-- create new create_organization_user permission
	// 				INSERT INTO lv_permissions(id, permission) VALUES (5, 'create_organization_user');
	// 				-- assign the new new pemission to roots and admins
	// 				INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (1, 5);
	// 				INSERT INTO permission_to_user_groups(user_group_id, permission_id) VALUES (2, 5);
	// 				`).Error
	// 			return err
	// 		},
	// 		Rollback: func(tx *gorm.DB) error {
	// 			return tx.Rollback().Error
	// 		},
	// 	})
	// }

	m := gormigrate.New(db, &gormigrate.Options{UseTransaction: true}, migrations)
	return m.Migrate()
}

// const devSeeds = `
// -----sensors-----
// INSERT INTO sensors(id, name, location, secret) VALUES ('b9dc3d20-256b-4ac7-8cae-2f6dc962e183', 'Test Sensor', 'Sofia, Bulgaria', 'sensorSecret123!');
// -----client-----
// INSERT INTO clients(id, name, email) VALUES (1, 'Test Client', 'test_client@gmail.com');
// -----client_subscriptions-----
// INSERT INTO client_subscriptions(id, client_id, task_type_id, tests_count_subscribed, tests_count_executed, period, last_execution_completed, opts, is_active) VALUES (1, 1, 1, 9999, 0, 60, NULL, '{"Host":"https://google.com", "Proto":"udp"}', TRUE);
// INSERT INTO client_subscriptions(id, client_id, task_type_id, tests_count_subscribed, tests_count_executed, period, last_execution_completed, opts, is_active) VALUES (2, 1, 2, 9999, 0, 60, NULL, '{"TargetDomain":"","TargetIPs":["127.0.0.1"],"Count":3,"Payload":"MDgwOTBhMGIwYzBkMGUwZjEwMTExMjEzMTQxNTE2MTcxODE5MWExYjFjMWQxZTFmMjAyMTIyMjMyNDI1MjYyNzI4MjkyYTJiMmMyZDJlMmYzMDMxMzIzMzM0MzUzNjM3"}', TRUE);
// INSERT INTO client_subscriptions(id, client_id, task_type_id, tests_count_subscribed, tests_count_executed, period, last_execution_completed, opts, is_active) VALUES (3, 1, 3, 9999, 0, 60, NULL, '{"TargetDomain":"https://google.com","HttpMethod":"GET","RequestHeaders":{"Content-Type":["application/json"]},"RequestBody":"c29tZSB0ZXN0IGJvZHk="}', TRUE);
// INSERT INTO client_subscriptions(id, client_id, task_type_id, tests_count_subscribed, tests_count_executed, period, last_execution_completed, opts, is_active) VALUES (4, 1, 4, 9999, 0, 60, NULL, '{"Dest": "8.8.8.8", "Port": 33434, "MaxHops": 64, "Retries": 3, "Timeout": 500, "FirstHop": 1, "NetCapRaw": true, "PacketSize": 52}', TRUE);
// `

//INSERT INTO client_subscriptions(id, client_id, task_type_id, tests_count_subscribed, tests_count_executed, period, last_execution_completed, opts, is_active) VALUES (1, 1, 2, 9999, 0, 60, NULL, '{"TargetDomain":null,"TargetIPs":["127.0.0.1"],"Count":3,"Payload":"MDgwOTBhMGIwYzBkMGUwZjEwMTExMjEzMTQxNTE2MTcxODE5MWExYjFjMWQxZTFmMjAyMTIyMjMyNDI1MjYyNzI4MjkyYTJiMmMyZDJlMmYzMDMxMzIzMzM0MzUzNjM3"}', TRUE);
