// MongoDB Replica Set Setup Script

// Wait for the server to be ready
sleep(2000);

// Initialize replica set
try {
  var result = rs.status();
  print("Replica set already initialized:", result.set);
} catch (e) {
  print("Initializing replica set...");
  var config = {
    _id: "rs0",
    members: [
      { _id: 0, host: "localhost:27017" }
    ]
  };

  var result = rs.initiate(config);
  print("Replica set initialization result:", result);

  // Wait for the replica set to become ready
  while (true) {
    try {
      var status = rs.status();
      if (status.members && status.members[0].state === 1) {
        print("Replica set is ready");
        break;
      }
    } catch (e) {
      print("Waiting for replica set to be ready...");
    }
    sleep(1000);
  }
}
