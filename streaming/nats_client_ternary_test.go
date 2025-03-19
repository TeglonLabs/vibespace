package streaming

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/bmorphism/vibespace-mcp-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestPublishWithTernaryData tests that a world moment with balanced ternary data is published correctly
func TestPublishWithTernaryData(t *testing.T) {
	// Create a world moment with balanced ternary data
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: time.Now().Unix(),
		CreatorID: "test-creator",
		Sharing: models.SharingSettings{
			IsPublic: true,
		},
	}

	// Attach balanced ternary data
	// The sequence 1, 0, -1, 1, 0, -1 represents the decimal value 363
	ternaryData := models.BalancedTernaryData{1, 0, -1, 1, 0, -1}
	moment.AttachBalancedTernaryData(ternaryData)

	// Create a mock NATS connection
	mockConn := new(MockNatsConn)
	
	// We expect it to publish to at least the world and creator subjects
	// The data should include the encoded balanced ternary data
	mockConn.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	mockConn.On("IsConnected").Return(true)
	mockConn.On("ConnectedServerId").Return("test-server")
	mockConn.On("ConnectedUrl").Return("nats://localhost:4222")
	mockConn.On("RTT").Return(time.Duration(1*time.Millisecond), nil)
	mockConn.On("Close").Return()
	
	// Create a client with our mock connection
	client := &NATSClient{
		url:         "nats://localhost:4222",
		streamID:    "test-stream",
		connected:   true,
		conn:        mockConn,
		rateLimiter: NewRateLimiter(100, 10, 1000),
	}
	
	// Publish the moment
	err := client.PublishWorldMoment(moment, "test-user")
	assert.NoError(t, err)
	
	// Verify the mock was called at least once for publishing
	mockConn.AssertNumberOfCalls(t, "Publish", 2) // At least world and creator subjects
	
	// Get the arguments from the publish calls
	assert.True(t, len(mockConn.Calls) > 0, "No publish calls were made")
	
	// Find a call with the right subject pattern (world.moment)
	var subject string
	var data []byte
	for _, call := range mockConn.Calls {
		if len(call.Arguments) >= 2 {
			subj := call.Arguments.Get(0).(string)
			if subj != "" && strings.Contains(subj, "world.moment") {
				subject = subj
				data = call.Arguments.Get(1).([]byte)
				break
			}
		}
	}
	
	// Verify the subject is formatted correctly
	assert.NotEmpty(t, subject, "No suitable subject found")
	assert.Contains(t, subject, "test-stream")
	assert.Contains(t, subject, "world.moment")
	assert.Contains(t, subject, "test-world")
	
	// Verify that the data contains both balanced ternary and binary data
	assert.NotEmpty(t, data)
	
	// Decode the data back to a WorldMoment
	var decodedMoment models.WorldMoment
	assert.NoError(t, json.Unmarshal(data, &decodedMoment))
	
	// Verify the binary data was included
	assert.NotNil(t, decodedMoment.BinaryData)
	assert.Equal(t, models.EncodingBinary, decodedMoment.BinaryData.Encoding)
	assert.Equal(t, "application/balanced-ternary", decodedMoment.BinaryData.Format)
	
	// Verify we can convert the binary data back to balanced ternary
	numTrits := len(ternaryData)
	recoveredTernary := BytesToTernary(decodedMoment.BinaryData.Data, numTrits)
	assert.Equal(t, ternaryData, recoveredTernary)
}

// TestBinaryPublish tests publishing a world moment with binary data
func TestBinaryPublish(t *testing.T) {
	// Create a world moment
	moment := &models.WorldMoment{
		WorldID:   "test-world",
		Timestamp: time.Now().Unix(),
		CreatorID: "test-creator",
		Sharing: models.SharingSettings{
			IsPublic: true,
		},
	}

	// Attach binary data in different formats
	binaryFormats := []struct {
		data     []byte
		encoding models.DataEncoding
		format   string
	}{
		{[]byte{0x01, 0x02, 0x03, 0x04}, models.EncodingBinary, "application/octet-stream"},
		{[]byte("Hello, World!"), models.EncodingBase64, "text/plain"},
		{[]byte{0xDE, 0xAD, 0xBE, 0xEF}, models.EncodingHex, "application/x-hex"},
	}

	for _, bf := range binaryFormats {
		t.Run(string(bf.encoding), func(t *testing.T) {
			// Attach binary data
			err := moment.AttachBinaryData(bf.data, bf.encoding, bf.format)
			assert.NoError(t, err)

			// Create a mock NATS connection
			mockConn := new(MockNatsConn)
			mockConn.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
			mockConn.On("IsConnected").Return(true)
			mockConn.On("ConnectedServerId").Return("test-server")
			mockConn.On("ConnectedUrl").Return("nats://localhost:4222")
			mockConn.On("RTT").Return(time.Duration(1*time.Millisecond), nil)
			mockConn.On("Close").Return()
			
			// Create a client with our mock connection
			client := &NATSClient{
				url:         "nats://localhost:4222",
				streamID:    "test-stream",
				connected:   true,
				conn:        mockConn,
				rateLimiter: NewRateLimiter(100, 10, 1000),
			}
			
			// Publish the moment
			err = client.PublishWorldMoment(moment, "test-user")
			assert.NoError(t, err)
			
			// Verify the mock was called for publishing
			mockConn.AssertCalled(t, "Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"))
			
			// Get the arguments from the publish calls
			// Make sure we have calls
			assert.True(t, len(mockConn.Calls) > 0, "No publish calls were made")
			
			// Find a call with the right subject pattern (world.moment)
			var data []byte
			for _, call := range mockConn.Calls {
				if len(call.Arguments) >= 2 {
					subject := call.Arguments.Get(0).(string)
					if subject != "" && strings.Contains(subject, "world.moment") {
						data = call.Arguments.Get(1).([]byte)
						break
					}
				}
			}
			
			assert.NotEmpty(t, data, "No suitable publish data found")
			
			// Decode the data back to a WorldMoment
			var decodedMoment models.WorldMoment
			err = json.Unmarshal(data, &decodedMoment)
			assert.NoError(t, err)
			
			// Verify the binary data was included with correct encoding and format
			assert.NotNil(t, decodedMoment.BinaryData)
			assert.Equal(t, bf.encoding, decodedMoment.BinaryData.Encoding)
			assert.Equal(t, bf.format, decodedMoment.BinaryData.Format)
			
			// For binary encoding, verify the data is exactly the same
			if bf.encoding == models.EncodingBinary {
				assert.Equal(t, bf.data, decodedMoment.BinaryData.Data)
			}
			
			// Verify we can recover the original data
			recoveredData, err := decodedMoment.GetBinaryData()
			assert.NoError(t, err)
			assert.Equal(t, bf.data, recoveredData)
			
			// Reset mock for next test
			mockConn.ExpectedCalls = nil
			mockConn.Calls = nil
		})
	}
}