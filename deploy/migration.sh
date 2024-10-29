#!/bin/bash
source ../auth.env


where migrations
until goose -dir "${MIGRATION_DIR}" postgres "${DSN}" up -v; do
  echo "Migration failed, retrying in 5 seconds..."
  sleep 5
done

echo "Migrations applied successfully!"
