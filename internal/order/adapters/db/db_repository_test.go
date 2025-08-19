package db

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/uuid"

	"r2-challenge/cmd/envs"
	orderdomain "r2-challenge/internal/order/domain"
	appdb "r2-challenge/pkg/db"
	"r2-challenge/pkg/observability"
)

// applyMigrations applies all SQL files in db/migrations in lexicographic order.
func applyMigrations(t *testing.T, gdb *appdb.Database) {
	t.Helper()

	migrationsDir, err := findMigrationsDir()
	if err != nil {
		t.Fatalf("locate migrations dir: %v", err)
	}

	dirEntries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("read migrations dir: %v", err)
	}

	files := make([]string, 0, len(dirEntries))
	for _, e := range dirEntries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".sql" {
			continue
		}
		files = append(files, filepath.Join(migrationsDir, e.Name()))
	}
	sort.Strings(files)

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read migration %s: %v", f, err)
		}
		if err := gdb.Exec(string(content)).Error; err != nil {
			t.Fatalf("apply migration %s: %v", f, err)
		}
	}
}

// truncateAll removes data to ensure test isolation.
func truncateAll(t *testing.T, gdb *appdb.Database) {
	t.Helper()
	stmts := []string{
		"TRUNCATE order_items RESTART IDENTITY CASCADE",
		"TRUNCATE orders RESTART IDENTITY CASCADE",
		"TRUNCATE products RESTART IDENTITY CASCADE",
		"TRUNCATE users RESTART IDENTITY CASCADE",
		"TRUNCATE payments RESTART IDENTITY CASCADE",
	}
	for _, s := range stmts {
		if err := gdb.Exec(s).Error; err != nil {
			t.Fatalf("truncate failed for %s: %v", s, err)
		}
	}
}

func setupDatabase(t *testing.T) (*appdb.Database, observability.Tracer) {
	t.Helper()

	tracer, err := observability.SetupTracer()
	if err != nil {
		t.Fatalf("setup tracer: %v", err)
	}

	env := envs.Envs{
		DBHost:            getenvOr("DB_HOST", "localhost"),
		DBPort:            getenvOr("DB_PORT", "5432"),
		DBUser:            getenvOr("DB_USER", "postgres"),
		DBPassword:        getenvOr("DB_PASSWORD", "postgres"),
		DBName:            getenvOr("DB_NAME", "r2_db"),
		DBSSLMode:         getenvOr("DB_SSLMODE", "disable"),
		DBMaxOpenConns:    5,
		DBMaxIdleConns:    5,
		DBConnMaxLifetime: "5m",
	}

	database, err := appdb.Setup(env)
	if err != nil {
		t.Skipf("database not available: %v", err)
	}

	applyMigrations(t, database)
	truncateAll(t, database)

	return database, tracer
}

func getenvOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// findMigrationsDir walks up from CWD until it finds db/migrations.
func findMigrationsDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for i := 0; i < 6; i++ {
		candidate := filepath.Join(cwd, "db", "migrations")
		if st, err := os.Stat(candidate); err == nil && st.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}

	return "", errors.New("db/migrations not found walking up directories")
}

func TestOrderRepository_Save_And_ListByUser_ReturnsItems(t *testing.T) {
	database, tracer := setupDatabase(t)
	repo, err := NewDBRepository(database, tracer)
	if err != nil {
		t.Fatalf("new repo: %v", err)
	}

	ctx := context.Background()

	// Create supporting records
	userID := uuid.NewString()
	productID := uuid.NewString()

	if err := database.Exec(`INSERT INTO users (id, email, password_hash, name, role) VALUES (?, 't@example.com', 'x', 'Test', 'user')`, userID).Error; err != nil {
		t.Fatalf("insert user: %v", err)
	}
	if err := database.Exec(`INSERT INTO products (id, name, description, category, price_cents, inventory) VALUES (?, 'P', 'D', 'c', 1234, 10)`, productID).Error; err != nil {
		t.Fatalf("insert product: %v", err)
	}

	// Save order with one item
	saved, err := repo.Save(ctx, orderdomain.Order{
		UserID:     userID,
		Status:     "created",
		TotalCents: 1234,
		Items: []orderdomain.OrderItem{
			{ProductID: productID, Quantity: 1, PriceCents: 1234},
		},
	})
	if err != nil {
		t.Fatalf("save order: %v", err)
	}
	if saved.ID == "" {
		t.Fatalf("expected saved order to have ID")
	}

	// List orders by user and ensure items are populated
	list, err := repo.ListByUser(ctx, userID, OrderFilter{Limit: 10})
	if err != nil {
		t.Fatalf("list by user: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("expected at least one order in list")
	}

	var found *orderdomain.Order
	for i := range list {
		if list[i].ID == saved.ID {
			found = &list[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("saved order not found in list")
	}
	if len(found.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(found.Items))
	}
	if found.Items[0].ProductID != productID {
		t.Fatalf("unexpected product id: %s", found.Items[0].ProductID)
	}
	if found.Items[0].Quantity != 1 {
		t.Fatalf("unexpected quantity: %d", found.Items[0].Quantity)
	}
	if found.Items[0].PriceCents != 1234 {
		t.Fatalf("unexpected price: %d", found.Items[0].PriceCents)
	}

	// Optional: payload snapshot using JSON to ensure shape includes items
	payload, err := json.Marshal(found)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	if !containsJSONKey(payload, "items") {
		t.Fatalf("payload missing items key: %s", string(payload))
	}
}

// containsJSONKey checks if a top-level key exists in a JSON object.
func containsJSONKey(b []byte, key string) bool {
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return false
	}
	_, ok := m[key]
	return ok
}
