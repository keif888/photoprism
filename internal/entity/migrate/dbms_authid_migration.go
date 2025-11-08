package migrate

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ConvertDBMSAuthIDDataTypes applies the data type conversion for auth_id columns for SQLite
// Conversion is from VARBINARY(255) DEFAULT ” to TEXT NOT NULL COLLATE BINARY DEFAULT ”
// This can't be done in a script as the order of columns is not known, and is determined by the age of the database
func ConvertDBMSAuthIDDataTypes(db *gorm.DB) (err error) {
	switch db.Dialector.Name() {
	case SQLite3:
		// These create statements will get out of date, but that is ok, as the main migrate path will add any missing columns/indexes in later.
		authSessionsCreate := `CREATE TABLE "auth_sessions" ("id" VARBINARY(2048),"user_uid" VARBINARY(42) DEFAULT '',"user_name" varchar(200),"client_uid" VARBINARY(42) DEFAULT '',"client_name" varchar(200) DEFAULT '',"client_ip" varchar(64),"auth_provider" VARBINARY(128) DEFAULT '',"auth_method" VARBINARY(128) DEFAULT '',"auth_issuer" VARBINARY(255) DEFAULT '',"auth_id" TEXT NOT NULL COLLATE BINARY DEFAULT '',"auth_scope" varchar(1024) DEFAULT '',"grant_type" VARBINARY(64) DEFAULT '',"last_active" bigint,"sess_expires" bigint,"sess_timeout" bigint,"preview_token" VARBINARY(64) DEFAULT '',"download_token" VARBINARY(64) DEFAULT '',"access_token" VARBINARY(4096) DEFAULT '',"refresh_token" VARBINARY(2048) DEFAULT '',"id_token" VARBINARY(2048) DEFAULT '',"user_agent" varchar(512),"data_json" VARBINARY(4096),"ref_id" VARBINARY(16) DEFAULT '',"login_ip" varchar(64),"login_at" datetime,"created_at" datetime,"updated_at" datetime , PRIMARY KEY ("id"))`
		authUsersCreate := `CREATE TABLE "auth_users" ("id" integer primary key autoincrement,"user_uuid" VARBINARY(64),"user_uid" VARBINARY(42),"auth_provider" VARBINARY(128) DEFAULT '',"auth_method" VARBINARY(128) DEFAULT '',"auth_issuer" VARBINARY(255) DEFAULT '',"auth_id" TEXT NOT NULL COLLATE BINARY DEFAULT '',"user_name" varchar(200),"display_name" varchar(200),"user_email" varchar(255),"backup_email" varchar(255),"user_role" varchar(64) DEFAULT '',"user_scope" varchar(1024) DEFAULT '*',"user_attr" varchar(1024) DEFAULT '',"super_admin" bool,"can_login" bool,"login_at" datetime,"expires_at" datetime,"webdav" bool,"base_path" VARBINARY(1024),"upload_path" VARBINARY(1024),"can_invite" bool,"invite_token" VARBINARY(64),"invited_by" varchar(64),"verify_token" VARBINARY(64),"verified_at" datetime,"consent_at" datetime,"born_at" datetime,"reset_token" VARBINARY(64),"preview_token" VARBINARY(64),"download_token" VARBINARY(64),"thumb" VARBINARY(128) DEFAULT '',"thumb_src" VARBINARY(8) DEFAULT '',"ref_id" VARBINARY(16),"created_at" datetime,"updated_at" datetime,"deleted_at" datetime )`

		type resultIndex struct {
			Name string
		}

		type pragmaTable struct {
			Cid       int
			Name      string
			Type      string
			Notnull   int
			DfltValue string
			Pk        int
		}

		if !db.Migrator().HasTable("auth_sessions") {
			if err := db.Exec(authSessionsCreate).Error; err != nil {
				return fmt.Errorf("migrate: error creating auth_sessions %w", err)
			}
		} else {
			// Data Migration here, by rename, create new, data transfer, drop indexes
			if err := db.Exec(`ALTER TABLE "auth_sessions" RENAME TO "migrate_auth_sessions"`).Error; err != nil {
				return fmt.Errorf("migrate: error renaming auth_sessions %w", err)
			}
			if err := db.Exec(authSessionsCreate).Error; err != nil {
				return fmt.Errorf("migrate: error creating auth_sessions %w", err)
			}

			// Get the columns of both old and new table, and find the columns that are in the old and new table
			var oldPragmaColumns []pragmaTable
			var newPragmaColumns []pragmaTable
			oldColumns := make(map[string]bool)

			if err := db.Raw("PRAGMA table_info(migrate_auth_sessions)").Scan(&oldPragmaColumns).Error; err != nil {
				return fmt.Errorf("migrate: error getting column list for migrate_auth_sessions with %w", err)
			}
			for _, pragma := range oldPragmaColumns {
				oldColumns[pragma.Name] = false
			}

			if err := db.Raw("PRAGMA table_info(auth_sessions)").Scan(&newPragmaColumns).Error; err != nil {
				return fmt.Errorf("migrate: error getting column list for auth_sessions with %w", err)
			}
			for _, pragma := range newPragmaColumns {
				if _, present := oldColumns[pragma.Name]; present {
					oldColumns[pragma.Name] = true
				}
			}
			// Build the select into statement
			var columns []string
			for key, value := range oldColumns {
				if value {
					columns = append(columns, key)
				}
			}

			populateStmt := fmt.Sprintf("INSERT INTO auth_sessions (%s) SELECT %s FROM migrate_auth_sessions", strings.Join(columns, ", "), strings.Join(columns, ", "))

			if err := db.Exec(populateStmt).Error; err != nil {
				return fmt.Errorf("migrate: error migrating with stmt %s with %w", populateStmt, err)
			}

			var indexes []resultIndex
			if err := db.Raw("SELECT name FROM sqlite_master WHERE type = 'index' AND tbl_name = ? AND sql IS NOT NULL", "migrate_auth_sessions").Scan(&indexes).Error; err != nil {
				return fmt.Errorf("migrate: error getting index list %w", err)
			}
			for _, index := range indexes {
				dropStatement := fmt.Sprintf(`DROP INDEX IF EXISTS "%s"`, index.Name)
				if err := db.Exec(dropStatement).Error; err != nil {
					return fmt.Errorf("migrate: error dropping index %s was %w", index.Name, err)
				}
			}
		}
		if !db.Migrator().HasTable("auth_users") {
			if err := db.Exec(authUsersCreate).Error; err != nil {
				return fmt.Errorf("migrate: error creating auth_users %w", err)
			}
		} else {
			// Data Migration here, by rename, create new, data transfer, drop indexes
			if err := db.Exec(`ALTER TABLE "auth_users" RENAME TO "migrate_auth_users"`).Error; err != nil {
				return fmt.Errorf("migrate: error renaming auth_users %w", err)
			}
			if err := db.Exec(authUsersCreate).Error; err != nil {
				return fmt.Errorf("migrate: error creating auth_users %w", err)
			}

			// Get the columns of both old and new table, and find the columns that are in the old and new table
			var oldPragmaColumns []pragmaTable
			var newPragmaColumns []pragmaTable
			oldColumns := make(map[string]bool)

			if err := db.Raw("PRAGMA table_info(migrate_auth_users)").Scan(&oldPragmaColumns).Error; err != nil {
				return fmt.Errorf("migrate: error getting column list for migrate_auth_users with %w", err)
			}
			for _, pragma := range oldPragmaColumns {
				oldColumns[pragma.Name] = false
			}

			if err := db.Raw("PRAGMA table_info(auth_users)").Scan(&newPragmaColumns).Error; err != nil {
				return fmt.Errorf("migrate: error getting column list for auth_users with %w", err)
			}
			for _, pragma := range newPragmaColumns {
				if _, present := oldColumns[pragma.Name]; present {
					oldColumns[pragma.Name] = true
				}
			}
			// Build the select into statement
			var columns []string
			for key, value := range oldColumns {
				if value {
					columns = append(columns, key)
				}
			}

			populateStmt := fmt.Sprintf("INSERT INTO auth_users (%s) SELECT %s FROM migrate_auth_users", strings.Join(columns, ", "), strings.Join(columns, ", "))

			if err := db.Exec(populateStmt).Error; err != nil {
				return fmt.Errorf("migrate: error migrating with stmt %s with %w", populateStmt, err)
			}

			var indexes []resultIndex
			if err := db.Raw("SELECT name FROM sqlite_master WHERE type = 'index' AND tbl_name = ? AND sql IS NOT NULL", "migrate_auth_users").Scan(&indexes).Error; err != nil {
				return fmt.Errorf("migrate: error getting index list %w", err)
			}
			for _, index := range indexes {
				dropStatement := fmt.Sprintf(`DROP INDEX IF EXISTS "%s"`, index.Name)
				if err := db.Exec(dropStatement).Error; err != nil {
					return fmt.Errorf("migrate: error dropping index %s was %w", index.Name, err)
				}
			}
		}
	case MySQL:
		// These create statements will get out of date, but that is ok, as the main migrate path will add any missing columns/indexes in later.
		authSessionsCreate := "CREATE TABLE `auth_sessions` (`id` varbinary(2048),`user_uid` varbinary(42) DEFAULT '',`user_name` varchar(200),`client_uid` varbinary(42) DEFAULT '',`client_name` varchar(200) DEFAULT '',`client_ip` varchar(64),`auth_provider` varbinary(128) DEFAULT '',`auth_method` varbinary(128) DEFAULT '',`auth_issuer` varbinary(255) DEFAULT '',`auth_id` varbinary(255) DEFAULT '',`auth_scope` varchar(1024) DEFAULT '',`grant_type` varbinary(64) DEFAULT '',`last_active` bigint,`sess_expires` bigint,`sess_timeout` bigint,`preview_token` varbinary(64) DEFAULT '',`download_token` varbinary(64) DEFAULT '',`access_token` varbinary(4096) DEFAULT '',`refresh_token` varbinary(2048) DEFAULT '',`id_token` varbinary(2048) DEFAULT '',`user_agent` varchar(512),`data_json` varbinary(4096),`ref_id` varbinary(16) DEFAULT '',`login_ip` varchar(64),`login_at` datetime(3) NULL,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,PRIMARY KEY (`id`),INDEX `idx_auth_sessions_user_uid` (`user_uid`),INDEX `idx_auth_sessions_user_name` (`user_name`),INDEX `idx_auth_sessions_client_uid` (`client_uid`),INDEX `idx_auth_sessions_client_ip` (`client_ip`),INDEX `idx_auth_sessions_auth_id` (`auth_id`),INDEX `idx_auth_sessions_sess_expires` (`sess_expires`))"
		authUsersCreate := "CREATE TABLE `auth_users` (`id` bigint AUTO_INCREMENT,`user_uuid` varbinary(64),`user_uid` varbinary(42),`auth_provider` varbinary(128) DEFAULT '',`auth_method` varbinary(128) DEFAULT '',`auth_issuer` varbinary(255) DEFAULT '',`auth_id` varbinary(255) DEFAULT '',`user_name` varchar(200),`display_name` varchar(200),`user_email` varchar(255),`backup_email` varchar(255),`user_role` varchar(64) DEFAULT '',`user_scope` varchar(1024) DEFAULT '*',`user_attr` varchar(1024) DEFAULT '',`super_admin` boolean,`can_login` boolean,`login_at` datetime(3) NULL,`expires_at` datetime(3) NULL,`webdav` boolean,`base_path` varbinary(1024),`upload_path` varbinary(1024),`can_invite` boolean,`invite_token` varbinary(64),`invited_by` varchar(64),`verify_token` varbinary(64),`verified_at` datetime(3) NULL,`consent_at` datetime(3) NULL,`born_at` datetime(3) NULL,`reset_token` varbinary(64),`preview_token` varbinary(64),`download_token` varbinary(64),`thumb` varbinary(128) DEFAULT '',`thumb_src` varbinary(8) DEFAULT '',`ref_id` varbinary(16),`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,PRIMARY KEY (`id`),INDEX `idx_auth_users_uuid` (`user_uuid`),UNIQUE INDEX `idx_auth_users_user_uid` (`user_uid`),INDEX `idx_auth_users_auth_id` (`auth_id`),INDEX `idx_auth_users_user_name` (`user_name`),INDEX `idx_auth_users_user_email` (`user_email`),INDEX `idx_auth_users_invite_token` (`invite_token`),INDEX `idx_auth_users_thumb` (`thumb`))"
		if !db.Migrator().HasTable("auth_sessions") {
			if err := db.Exec(authSessionsCreate).Error; err != nil {
				return fmt.Errorf("migrate: error creating auth_sessions %w", err)
			}
		}
		if !db.Migrator().HasTable("auth_users") {
			if err := db.Exec(authUsersCreate).Error; err != nil {
				return fmt.Errorf("migrate: error creating auth_users %w", err)
			}
		}
		// There are no migration needs for MariaDB as the structure is not being manipulated.
	case Postgres:
		// These create statements will get out of date, but that is ok, as the main migrate path will add any missing columns/indexes in later.
		authSessionsCreate := "CREATE TABLE \"auth_sessions\" (\"id\" bytea,\"user_uid\" bytea DEFAULT '',\"user_name\" varchar(200),\"client_uid\" bytea DEFAULT '',\"client_name\" varchar(200) DEFAULT '',\"client_ip\" varchar(64),\"auth_provider\" bytea DEFAULT '',\"auth_method\" bytea DEFAULT '',\"auth_issuer\" bytea DEFAULT '',\"auth_id\" bytea DEFAULT '',\"auth_scope\" varchar(1024) DEFAULT '',\"grant_type\" bytea DEFAULT '',\"last_active\" bigint,\"sess_expires\" bigint,\"sess_timeout\" bigint,\"preview_token\" bytea DEFAULT '',\"download_token\" bytea DEFAULT '',\"access_token\" bytea DEFAULT '',\"refresh_token\" bytea DEFAULT '',\"id_token\" bytea DEFAULT '',\"user_agent\" varchar(512),\"data_json\" bytea,\"ref_id\" bytea DEFAULT '',\"login_ip\" varchar(64),\"login_at\" timestamptz,\"created_at\" timestamptz,\"updated_at\" timestamptz,PRIMARY KEY (\"id\"))"
		authUsersCreate := "CREATE TABLE \"auth_users\" (\"id\" bigserial,\"user_uuid\" bytea,\"user_uid\" bytea,\"auth_provider\" bytea DEFAULT '',\"auth_method\" bytea DEFAULT '',\"auth_issuer\" bytea DEFAULT '',\"auth_id\" bytea DEFAULT '',\"user_name\" varchar(200),\"display_name\" varchar(200),\"user_email\" varchar(255),\"backup_email\" varchar(255),\"user_role\" varchar(64) DEFAULT '',\"user_scope\" varchar(1024) DEFAULT '*',\"user_attr\" varchar(1024) DEFAULT '',\"super_admin\" boolean,\"can_login\" boolean,\"login_at\" timestamptz,\"expires_at\" timestamptz,\"webdav\" boolean,\"base_path\" bytea,\"upload_path\" bytea,\"can_invite\" boolean,\"invite_token\" bytea,\"invited_by\" varchar(64),\"verify_token\" bytea,\"verified_at\" timestamptz,\"consent_at\" timestamptz,\"born_at\" timestamptz,\"reset_token\" bytea,\"preview_token\" bytea,\"download_token\" bytea,\"thumb\" bytea DEFAULT '',\"thumb_src\" bytea DEFAULT '',\"ref_id\" bytea,\"created_at\" timestamptz,\"updated_at\" timestamptz,\"deleted_at\" timestamptz,PRIMARY KEY (\"id\"))"
		if !db.Migrator().HasTable("auth_sessions") {
			if err := db.Exec(authSessionsCreate).Error; err != nil {
				return fmt.Errorf("migrate: error creating auth_sessions %w", err)
			}
			if err := db.Exec("CREATE INDEX IF NOT EXISTS \"idx_auth_sessions_auth_id\" ON \"auth_sessions\" (\"auth_id\")").Error; err != nil {
				return fmt.Errorf("migrate: error creating idx_auth_sessions_auth_id %w", err)
			}
		}
		if !db.Migrator().HasTable("auth_users") {
			if err := db.Exec(authUsersCreate).Error; err != nil {
				return fmt.Errorf("migrate: error creating auth_users %w", err)
			}
			if err := db.Exec("CREATE INDEX IF NOT EXISTS \"idx_auth_users_auth_id\" ON \"auth_users\" (\"auth_id\")").Error; err != nil {
				return fmt.Errorf("migrate: error creating idx_auth_users_auth_id %w", err)
			}
		}
		// There are no migration needs for Postgres as the structure is not being manipulated.
	default:
	}
	return nil
}
