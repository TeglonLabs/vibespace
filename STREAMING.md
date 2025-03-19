# vibespace MCP Experience NATS Streaming

This document describes how to use the NATS streaming functionality in the vibespace MCP experience.

## Overview

The vibespace MCP experience now supports streaming "world moments" to a NATS server. A "world moment" is a snapshot of a world's state at a specific point in time, including its:

- World details (name, description, type, etc.)
- Current vibe (if any)
- Occupancy level
- Activity level (calculated from occupancy and other factors)
- Sensor data (temperature, humidity, light, sound, movement)
- Timestamp
- Additional data (extensible via the CustomData field)

This streaming capability allows external applications to receive real-time updates about worlds in the VibeSpace ecosystem.

## NATS Topics/Subjects

The streaming service publishes to the following NATS subjects:

### Public Topics
- `{streamID}.world.moment.{worldID}`: World moment updates for a specific world (only for public worlds)
- `{streamID}.world.vibe.{worldID}`: Vibe updates for a specific world (only for public vibes)

### User-Specific Topics
- `{streamID}.world.moment.{worldID}.user.{userID}`: World moment updates tailored for a specific user
- `{streamID}.world.vibe.{worldID}.user.{userID}`: Vibe updates tailored for a specific user

The content published to user-specific topics is filtered based on the user's access level and the sharing settings defined for the world or vibe.

> 📝 **Note**: The default Stream ID is "ies", so topics will typically be like "ies.world.moment.{worldID}".

## MCP Tools

The following MCP tools are available to control the streaming functionality:

> ⚠️ **Multiplayer Feature**: Users can now share their worlds with others through NATS streaming. All streaming operations now include user attribution and access control mechanisms to ensure privacy and enable collaboration.

### streaming_startStreaming

Starts streaming world moments to NATS.

Parameters:
- `interval` (optional): Stream interval in milliseconds (e.g., 5000 for 5 seconds)

Example:
```json
{
  "method": "streaming_startStreaming",
  "params": {
    "interval": 5000
  }
}
```

### streaming_stopStreaming

Stops streaming world moments.

Example:
```json
{
  "method": "streaming_stopStreaming",
  "params": {}
}
```

### streaming_status

Returns the current status of the streaming service.

Example:
```json
{
  "method": "streaming_status",
  "params": {}
}
```

Response:
```json
{
  "isStreaming": true,
  "streamInterval": "5s",
  "natsUrl": "nats://nonlocal.info:4222",
  "message": "Streaming active with interval 5s",
  "uiIndicators": {
    "streamActive": true,
    "streamIndicator": "ACTIVE",
    "statusColor": "#4CAF50",
    "connectionQuality": "EXCELLENT" 
  },
  "connectionStatus": {
    "isConnected": true,
    "url": "nats://nonlocal.info:4222",
    "reconnectCount": 0,
    "disconnectCount": 0,
    "lastConnectTime": "2025-03-19T15:22:34.123Z",
    "serverId": "NCHBCW3",
    "connectedUrl": "nats://nonlocal.info:4222",
    "rtt": "23.5ms"
  }
}
```

### streaming_streamWorld

Streams a single moment for a specific world (one-time).

Parameters:
- `worldId`: The ID of the world to stream
- `userId`: The ID of the user initiating the stream (required for multiplayer attribution)
- `sharing` (optional): Controls how this moment is shared with other users
  - `isPublic`: Boolean indicating if this world is visible to all users
  - `allowedUsers`: Array of user IDs who can access this world
  - `contextLevel`: Level of detail to share ("none", "partial", or "full")

Example:
```json
{
  "method": "streaming_streamWorld",
  "params": {
    "worldId": "office",
    "userId": "user123",
    "sharing": {
      "isPublic": false,
      "allowedUsers": ["user456", "user789"],
      "contextLevel": "partial"
    }
  }
}
```

### streaming_updateConfig

Updates the streaming configuration.

Parameters:
- `natsHost` (optional): The hostname of the NATS server (e.g., "nonlocal.info")
- `natsPort` (optional): The port of the NATS server (e.g., 4222)
- `natsUrl` (optional): The complete URL of the NATS server (overrides host/port if set)
- `streamId` (optional): The stream identifier (default: "ies")
- `streamInterval` (optional): Stream interval in milliseconds

Example:
```json
{
  "method": "streaming_updateConfig",
  "params": {
    "natsHost": "nonlocal.info",
    "natsPort": 4222,
    "streamId": "ies",
    "streamInterval": 10000
  }
}
```

## Binary and Balanced Ternary Data Support

The streaming service now supports attaching arbitrary binary data and balanced ternary data to world moments. This enables efficient transmission of specialized data formats and provides native support for the balanced ternary numeral system.

### Binary Data Support

World moments can include arbitrary binary data with different encoding formats:

```go
// Attach binary data directly
moment.AttachBinaryData([]byte{0x01, 0x02, 0x03, 0x04}, models.EncodingBinary, "application/octet-stream")

// Attach text data with base64 encoding
moment.AttachBinaryData([]byte("Hello, World!"), models.EncodingBase64, "text/plain")

// Attach binary data with hex encoding
moment.AttachBinaryData([]byte{0xDE, 0xAD, 0xBE, 0xEF}, models.EncodingHex, "application/x-hex")
```

The binary data is encoded according to the specified format and included in the published world moment. The `BinaryData` struct includes:

- `Data`: The actual binary payload
- `Encoding`: The encoding format (binary, base64, hex)
- `Format`: MIME type or format description

Clients can retrieve the original data using:

```go
// Get the original binary data, regardless of encoding
originalData, err := moment.GetBinaryData()
```

### Balanced Ternary Support

The system provides native support for balanced ternary, a non-standard numeral system with digits {-1, 0, 1} (represented as 'T', '0', '1'). This offers computational advantages for certain operations and can represent signed values without a separate sign bit.

```go
// Create balanced ternary data from a string representation
// 'T' or '-' represents -1, '0' represents 0, '1' represents 1
ternary := models.NewBalancedTernaryFromString("10T01")

// Attach balanced ternary data to a moment
moment.AttachBalancedTernaryData(ternary)

// Attach balanced ternary from a decimal value
moment.AttachBalancedTernaryFromDecimal(42)
```

The balanced ternary data is automatically encoded to a compact binary representation (2 bits per trit) when published and can be reconstructed by clients:

```go
// Retrieve the original balanced ternary data from binary
ternaryData := streaming.BytesToTernary(moment.BinaryData.Data, numTrits)

// Convert ternary to decimal
decimal := ternaryData.ToDecimal()

// Get string representation (using 'T' for -1)
ternaryString := ternaryData.String() // e.g. "10T01"
```

### Advanced Use Cases

Binary and balanced ternary data support enables:

1. **Compact Data Transmission**: Efficient encoding of specialized data formats
2. **Specialized Computation**: Native support for ternary operations
3. **Custom Data Formats**: Transmission of application-specific binary formats
4. **Improved Precision**: Balanced ternary can provide advantages for certain mathematical operations

Example: Publishing a world moment with sensor data and ternary encoding of signal processing results:

```go
moment := &models.WorldMoment{
    WorldID:   "laboratory",
    Timestamp: time.Now().Unix(),
    CreatorID: "researcher1",
    Sharing: models.SharingSettings{
        IsPublic: true,
    },
    SensorData: models.SensorData{
        Temperature: &temp,
        Humidity:    &humidity,
    },
}

// Attach balanced ternary data representing signal processing results
moment.AttachBalancedTernaryFromString("10T01T0110T")

// Publish the moment
client.PublishWorldMoment(moment, "researcher1")
```

Balanced ternary is particularly useful for:
- Efficient representation of signed values
- Certain mathematical algorithms (like multiplication)
- Signal processing with three-state transitions
- Specialized computational models

## Privacy, Access Control and Advanced Features

The multiplayer streaming implementation follows Model Context Protocol (MCP) best practices for shared experiences:

1. **User Attribution**: All world moments are attributed to their creator via the `creatorId` field
2. **Permission Model**: Access to streams is controlled through `SharingSettings`
   - `isPublic`: When true, the world is visible to all users
   - `allowedUsers`: Specific users who can access private worlds
   - `contextLevel`: Controls how much information is shared with others

3. **Content Filtering**: The content shared with non-creator users is filtered based on the `contextLevel`:
   - `none`: Minimal information (just ID and basic metadata)
   - `partial`: Moderate information (core data without sensitive/custom data)
   - `full`: Complete information (everything visible)

4. **User Interface Recommendations**: Clients should implement:
   - Visual indicators when streaming is active
   - Clear display of who is viewing a world
   - Obvious opt-in controls for sharing
   
   The `streaming_status` tool now provides UI indicator data that can be used to implement visual cues:
   
   ```json
   "uiIndicators": {
     "streamActive": true,                // Boolean flag for active streaming
     "streamIndicator": "ACTIVE",         // "ACTIVE", "READY", or "OFFLINE"
     "statusColor": "#4CAF50",            // Suggested color for indicators (green, blue, or red)
     "connectionQuality": "EXCELLENT"     // "EXCELLENT", "GOOD", "FAIR", "REMOTE", "LOCAL", "CUSTOM", or "DISCONNECTED"
   }
   ```
   
   Clients should use these indicators to:
   - Show a colored status badge when streaming is active
   - Display the connection quality (e.g., "Connected to REMOTE server")
   - Provide visual feedback when streaming starts or stops

5. **Advanced Connection Monitoring**: The streaming service now provides detailed connection status information:

   ```json
   "connectionStatus": {
     "isConnected": true,
     "url": "nats://nonlocal.info:4222",
     "reconnectCount": 2,
     "disconnectCount": 1,
     "lastConnectTime": "2025-03-19T14:30:45.123Z",
     "lastErrorMessage": "",
     "serverId": "NCHXEV2",
     "connectedUrl": "nats://nonlocal.info:4222",
     "rtt": "45.2ms"
   }
   ```

   This detailed information can be used to:
   - Monitor connection stability
   - Display reconnection statistics
   - Show real-time connection performance (RTT)
   - Diagnose connection issues

6. **Rate Limiting**: The streaming service now includes automatic rate limiting to prevent overwhelming the NATS server:
   - Burst capacity of 100 messages
   - Sustained rate of 10 messages per second
   - Automatic throttling when limits are exceeded
   - Error feedback when rate limits are hit

   This ensures optimal performance and prevents client applications from accidentally flooding the system with too many updates.

## Using with NATS Clients

To receive world moments in your application, subscribe to the relevant NATS subjects using a NATS client library. For multiplayer awareness, subscribe to your user-specific topics:

```go
// Go example - for a specific user
userID := "user123"
streamID := "ies" // Default stream ID
nc, _ := nats.Connect("nats://nonlocal.info:4222")

// Subscribe to public world moments with stream ID
sub1, _ := nc.Subscribe(fmt.Sprintf("%s.world.moment.*", streamID), func(msg *nats.Msg) {
    var moment models.WorldMoment
    json.Unmarshal(msg.Data, &moment)
    fmt.Printf("Received public moment for world %s at %v\n", moment.WorldID, moment.Timestamp)
})

// Subscribe to user-specific world moments
userSubject := fmt.Sprintf("%s.world.moment.*.user.%s", streamID, userID)
sub2, _ := nc.Subscribe(userSubject, func(msg *nats.Msg) {
    var moment models.WorldMoment
    json.Unmarshal(msg.Data, &moment)
    fmt.Printf("Received user-specific moment for world %s at %v\n", moment.WorldID, moment.Timestamp)
    
    // Check if you are viewing this world
    isViewing := false
    for _, viewer := range moment.Viewers {
        if viewer == userID {
            isViewing = true
            break
        }
    }
    
    if isViewing {
        fmt.Println("You are currently viewing this world")
    }
    
    // Show other viewers
    fmt.Printf("Other viewers: %d people\n", len(moment.Viewers))
})
```

```javascript
// JavaScript example - with multiplayer awareness
const NATS = require('nats');
const userID = "user123";
const streamID = "ies"; // Default stream ID

async function setupSubscriptions() {
    const nc = await NATS.connect({servers: ["nats://nonlocal.info:4222"]});
    
    // Subscribe to public world moments
    const publicSub = nc.subscribe(`${streamID}.world.moment.*`);
    (async () => {
        for await (const msg of publicSub) {
            const moment = JSON.parse(msg.data);
            console.log(`Received public moment for world ${moment.worldID} at ${moment.timestamp}`);
        }
    })();
    
    // Subscribe to user-specific world moments
    const userSubject = `${streamID}.world.moment.*.user.${userID}`;
    const userSub = nc.subscribe(userSubject);
    (async () => {
        for await (const msg of userSub) {
            const moment = JSON.parse(msg.data);
            console.log(`Received user-specific moment for world ${moment.worldID}`);
            
            // Display multiplayer info
            const isCreator = moment.creatorID === userID;
            console.log(isCreator ? "You created this world" : `Created by: ${moment.creatorID}`);
            console.log(`Sharing level: ${moment.sharing.contextLevel}`);
            console.log(`Current viewers: ${moment.viewers.length}`);
        }
    })();
}

setupSubscriptions();
```

### Example NATS Subscriber

The repository includes a complete example NATS subscriber application in the `examples/` directory:

- `examples/nats_subscriber.go`: A Go application that subscribes to both world moment and vibe update streams
- `examples/run_nats_subscriber.sh`: A script to build and run the subscriber

To run the example:

```bash
cd examples
./run_nats_subscriber.sh [nats-server-url] [stream-id] [user-id]
```

Parameters:
- `nats-server-url`: NATS server URL (defaults to `nats://nonlocal.info:4222` if not specified)
- `stream-id`: Stream ID for subject namespacing (defaults to `ies` if not specified)
- `user-id`: User ID for user-specific subscriptions (optional)

The example subscriber:
1. Connects to the specified NATS server
2. Subscribes to all world moments (`{streamID}.world.moment.*`)
3. Subscribes to all vibe updates (`{streamID}.world.vibe.*`) 
4. If a user ID is provided, subscribes to user-specific messages (`{streamID}.world.moment.*.user.{userID}`)
5. Prints formatted information about each received message
6. Gracefully handles disconnections and reconnections

## Running a NATS Server

To use the streaming functionality, you need a running NATS server. You can run one using Docker:

```bash
docker run -p 4222:4222 -p 8222:8222 nats
```

Or install and run natively:

```bash
# Install NATS Server
go install github.com/nats-io/nats-server/v2@latest

# Run the server
nats-server
```

## Configuration

The streaming service can be configured by updating the following values in the `main.go` file:

```go
streamingConfig := &streaming.StreamingConfig{
    NATSHost:       "nonlocal.info",   // NATS server hostname
    NATSPort:       4222,              // NATS server port
    StreamID:       "ies",             // Stream identifier for subject namespacing
    StreamInterval: 5 * time.Second,   // Interval between streaming moments
    AutoStart:      false,             // Whether to start streaming automatically
}
```

These values can also be modified at runtime using the `streaming_updateConfig` tool.

## Stream ID Namespacing

The Stream ID feature allows for better organization and isolation of NATS subjects. This is particularly useful in environments where multiple MCP servers or applications share the same NATS server.

### Subject Pattern With Stream ID

All NATS subjects now follow this pattern:

- `{streamID}.world.moment.{worldID}` - Public world moments
- `{streamID}.world.moment.{worldID}.user.{userID}` - User-specific world moments
- `{streamID}.world.vibe.{worldID}` - Public vibe updates
- `{streamID}.world.vibe.{worldID}.user.{userID}` - User-specific vibe updates

### Default Stream ID

The default Stream ID is "ies", resulting in subjects like:

- `ies.world.moment.office`
- `ies.world.moment.office.user.user123`
- `ies.world.vibe.office`

### Changing the Stream ID

You can change the Stream ID at runtime using the `streaming_updateConfig` tool:

```json
{
  "method": "streaming_updateConfig",
  "params": {
    "streamId": "custom-stream"
  }
}
```

This would result in subjects like:

- `custom-stream.world.moment.office`
- `custom-stream.world.moment.office.user.user123`

### Using Multiple Stream IDs

When subscribing to subjects, you need to match the Stream ID used by the server. The example NATS subscriber now supports this with a command-line parameter:

```bash
./run_nats_subscriber.sh nats://nonlocal.info:4222 custom-stream user123
```

This will subscribe to subjects with the "custom-stream" Stream ID for user "user123".