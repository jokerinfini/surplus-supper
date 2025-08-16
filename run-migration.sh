#!/bin/bash

# Database Migration Script for Railway PostgreSQL
# Usage: ./run-migration.sh

echo "🚀 Running Surplus Supper Database Migration..."

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ Error: DATABASE_URL environment variable not set"
    echo "Please set your Railway PostgreSQL connection URL:"
    echo "export DATABASE_URL=postgresql://username:password@host:port/database"
    exit 1
fi

echo "✅ Connecting to database..."
echo "📊 Running migration..."

# Run the migration SQL
psql "$DATABASE_URL" -f backend/db/migrations/001_initial_schema.sql

if [ $? -eq 0 ]; then
    echo "✅ Migration completed successfully!"
    echo "🎉 Database tables created with sample data"
    echo "📋 Created tables: users, restaurants, inventory_items, offers, etc."
else
    echo "❌ Migration failed. Please check your DATABASE_URL and try again."
    exit 1
fi
