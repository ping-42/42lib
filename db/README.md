## DB Conventions

* Tables that contain constants, should be prefixed with "lv_" (lookup values). And the inserts for them should not change their order and be as part of the migration.
* IDs which are not for internal usage only, should be stored as UUID data type.
