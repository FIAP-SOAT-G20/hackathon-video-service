// MongoDB initialization script for video service
// This script sets up the initial database structure for the video service

print('Starting MongoDB initialization for video service...');

// Switch to the video_service database
db = db.getSiblingDB('video_service');

// Create a user for the application (if needed)
// Note: In production, you'd want more secure credentials
db.createUser({
  user: 'video_service_user',
  pwd: 'video_service_password',
  roles: [
    {
      role: 'readWrite',
      db: 'video_service'
    }
  ]
});

// Create the videos collection
db.createCollection('videos');

// Create indexes for optimal performance
db.videos.createIndex({ "video_id": 1 }, { unique: true });
db.videos.createIndex({ "customer_id": 1 });
db.videos.createIndex({ "status": 1 });
db.videos.createIndex({ "status": 1, "customer_id": 1 });
db.videos.createIndex({ "created_at": -1 });
db.videos.createIndex({ "updated_at": -1 });

// Insert some sample data for testing (optional)
db.videos.insertMany([
  {
    video_id: 1,
    customer_id: 100,
    status: "OPEN",
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    video_id: 2,
    customer_id: 101,
    status: "PENDING",
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    video_id: 3,
    customer_id: 102,
    status: "COMPLETED",
    created_at: new Date(),
    updated_at: new Date()
  }
]);

print('MongoDB initialization completed successfully!');
print('Database: video_service');
print('Collection: videos');
print('Sample records inserted: 3');
print('Indexes created: 6');
