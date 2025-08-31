.PHONY: migrate clean

# Run migrations up
migrate:
	goose up

# Clean database (down + up for fresh start)
clean:
	goose down
	goose up

# Just rollback migrations
rollback:
	goose down

# Check migration status
status:
	goose status