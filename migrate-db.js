const { Client } = require('pg');
const fs = require('fs');
const path = require('path');

// Your Railway PostgreSQL connection URL
const DATABASE_URL = 'postgresql://postgres:veQjTQKoMnjazNDruTCZcggOZUqLfVtf@caboose.proxy.rlwy.net:10803/railway';

async function runMigration() {
  console.log('🚀 Running Surplus Supper Database Migration...');
  
  const client = new Client({
    connectionString: DATABASE_URL,
  });

  try {
    // Connect to database
    await client.connect();
    console.log('✅ Connected to Railway PostgreSQL');

    // Read migration file
    const migrationPath = path.join(__dirname, 'backend', 'db', 'migrations', '001_initial_schema.sql');
    const migrationSQL = fs.readFileSync(migrationPath, 'utf8');

    // Run migration
    console.log('📊 Creating tables and inserting sample data...');
    await client.query(migrationSQL);

    console.log('✅ Migration completed successfully!');
    console.log('🎉 Database tables created with sample data');
    console.log('📋 Created tables: users, restaurants, inventory_items, offers, etc.');

    // Verify tables were created
    const result = await client.query(`
      SELECT table_name 
      FROM information_schema.tables 
      WHERE table_schema = 'public'
      ORDER BY table_name;
    `);

    console.log('\n📋 Tables in database:');
    result.rows.forEach(row => {
      console.log(`  - ${row.table_name}`);
    });

  } catch (error) {
    console.error('❌ Migration failed:', error.message);
    process.exit(1);
  } finally {
    await client.end();
  }
}

runMigration();
