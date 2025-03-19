package streaming

import (
	"testing"
)

func TestRunToolTests(t *testing.T) {
	t.Run("TestToolsStartStreaming", TestToolsStartStreaming)
	t.Run("TestToolsStopStreaming", TestToolsStopStreaming)
	t.Run("TestToolsStatus", TestToolsStatus)
	t.Run("TestToolsStreamWorld", TestToolsStreamWorld)
	t.Run("TestToolsUpdateConfig", TestToolsUpdateConfig)
}