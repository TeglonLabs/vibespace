package streaming

import ()

// This file contains helper functions for testing

// SetRateLimiter sets the rate limiter for testing
func (c *NATSClient) SetRateLimiter(limiter *RateLimiter) {
	c.rateLimiter = limiter
}

// SetConnectedState sets the connected flag for testing
func (c *NATSClient) SetConnectedState(connected bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connected = connected
}

// InjectConnection allows injecting a mock connection for testing
func (c *NATSClient) InjectConnection(conn NatsConnection) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn = conn
	c.connected = true // Also need to set connected state
}

// SetConfig allows setting the config directly for testing
func (s *StreamingService) SetConfig(config *StreamingConfig) {
	s.config = config
}

// SetClient sets the NATS client for testing
func (s *StreamingService) SetClient(client NATSClientInterface) {
	s.natsClient = client
}

// GetNATSClient returns the NATS client for testing
func (s *StreamingService) GetNATSClient() NATSClientInterface {
	return s.natsClient
}

// GetConfig returns the config for testing
func (s *StreamingService) GetConfig() *StreamingConfig {
	return s.config
}

// SetStreamingActive sets the streaming active flag for testing
func (s *StreamingService) SetStreamingActive(active bool) {
	s.streamingActive = active
}

// SetMomentGenerator sets the moment generator for testing
func (s *StreamingService) SetMomentGenerator(generator MomentGeneratorInterface) {
	s.momentGenerator = generator
}

// SetRepository sets the repository for testing
func (s *StreamingService) SetRepository(repo RepositoryInterface) {
	s.repo = repo
}