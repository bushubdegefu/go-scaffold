package temps

var databasePglogicMigrationTemplate = `package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq"
)

func connectDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    return db, nil
}

func createReplicationSlot(db *sql.DB) {
    _, err := db.Exec("SELECT pg_create_logical_replication_slot('my_slot', 'pgoutput')")
    if err != nil {
        log.Fatalf("Failed to create replication slot: %v", err)
    }
}

func createReplicationSet(db *sql.DB) {
    _, err := db.Exec("SELECT pg_create_replication_set('my_replication_set')")
    if err != nil {
        log.Fatalf("Failed to create replication set: %v", err)
    }

    _, err = db.Exec("SELECT pg_add_table_to_replication_set('my_replication_set', 'my_table')")
    if err != nil {
        log.Fatalf("Failed to add table to replication set: %v", err)
    }
}

func createSubscription(db *sql.DB, sourceConnStr string) {
    _, err := db.Exec("CREATE SUBSCRIPTION my_subscription CONNECTION $1 PUBLICATION my_replication_set", sourceConnStr)
    if err != nil {
        log.Fatalf("Failed to create subscription: %v", err)
    }
}

func dropSubscription(db *sql.DB) {
    _, err := db.Exec("SELECT pg_drop_subscription('my_subscription')")
    if err != nil {
        log.Fatalf("Failed to drop subscription: %v", err)
    }
}

func dropReplicationSet(db *sql.DB) {
    _, err := db.Exec("SELECT pg_drop_replication_set('my_replication_set')")
    if err != nil {
        log.Fatalf("Failed to drop replication set: %v", err)
    }

    _, err = db.Exec("SELECT pg_drop_logical_replication_slot('my_slot')")
    if err != nil {
        log.Fatalf("Failed to drop replication slot: %v", err)
    }
}

func main() {
    sourceDSN := "host=source_host user=replication_user password=your_password dbname=source_db sslmode=disable"
    targetDSN := "host=target_host user=your_user password=your_password dbname=target_db sslmode=disable"

    sourceDB, err := connectDB(sourceDSN)
    if err != nil {
        log.Fatal(err)
    }
    defer sourceDB.Close()

    targetDB, err := connectDB(targetDSN)
    if err != nil {
        log.Fatal(err)
    }
    defer targetDB.Close()

    // Setup replication
    createReplicationSlot(sourceDB)
    createReplicationSet(sourceDB)
    createSubscription(targetDB, sourceDSN)

    log.Println("Replication setup completed.")

    // Optionally: cleanup after migration
    // dropSubscription(targetDB)
    // dropReplicationSet(sourceDB)
}
`
var databaseLLVMMigrationTemplate = ``
