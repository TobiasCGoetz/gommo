package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventLogAccessControl(t *testing.T) {
	// Setup test environment
	ts := setupTestSuite(t)
	defer func() {
		// Reset global event logger after test
		eventLogger = NewEventLogger()
	}()

	// Create test players
	player1ID := ts.createPlayerAt(1, 1)
	player2ID := ts.createPlayerAt(2, 2)

	// Create test events
	now := time.Now()
	testEvents := []GameEvent{
		{
			ID:        "1",
			Type:      EventPlayerMove,
			PlayerID:  player1ID,
			Timestamp: now.Add(-5 * time.Minute),
			Details:   map[string]interface{}{"x": 1, "y": 2},
		},
		{
			ID:        "2",
			Type:      EventCardUsage,
			PlayerID:  player1ID,
			Timestamp: now.Add(-4 * time.Minute),
			Details:   map[string]interface{}{"card": "weapon", "target": "zombie"},
		},
		{
			ID:        "3",
			Type:      EventPlayerMove,
			PlayerID:  player2ID,
			Timestamp: now.Add(-3 * time.Minute),
			Details:   map[string]interface{}{"x": 2, "y": 3},
		},
		{
			ID:        "4",
			Type:      EventPlayerDeath,
			PlayerID:  player2ID,
			Timestamp: now.Add(-2 * time.Minute),
			Details:   map[string]interface{}{"cause": "zombies", "location": map[string]int{"x": 2, "y": 3}},
		},
		{
			ID:        "5",
			Type:      EventCombatResult,
			PlayerID:  "",
			Timestamp: now.Add(-1 * time.Minute),
			Details: map[string]interface{}{
				"tile_x": 2, "tile_y": 3,
				"involved_players": []string{player1ID, player2ID},
			},
		},
	}

	// Add events to logger
	el := NewEventLogger().(*EventLoggerImpl)
	for _, event := range testEvents {
		el.events = append(el.events, event)
	}
	eventLogger = el

	tests := []struct {
		name            string
		requestingPlayer string
		expectedEvents  []string // Expected event IDs
		filterType     EventType
		description    string
	}{
		{
			name:            "player1 sees own events and global deaths",
			requestingPlayer: player1ID,
			expectedEvents:  []string{"1", "2", "4", "5"}, // Own events + death event + combat
			description:     "Should see own events and all player deaths",
		},
		{
			name:            "player2 sees own events and global deaths",
			requestingPlayer: player2ID,
			expectedEvents:  []string{"3", "4", "5"}, // Own events + death event + combat
			description:     "Should see own events and all player deaths",
		},
		{
			name:            "filter player1 moves only",
			requestingPlayer: player1ID,
			filterType:      EventPlayerMove,
			expectedEvents:  []string{"1"},
			description:     "Should only see player1's movement events",
		},
		{
			name:            "filter death events",
			requestingPlayer: player1ID,
			filterType:      EventPlayerDeath,
			expectedEvents:  []string{"4"},
			description:     "Should see death events from any player",
		},
		{
			name:            "nonexistent player sees only global events",
			requestingPlayer: "nonexistent",
			expectedEvents:  []string{"4"}, // Only global death event
			description:     "Should only see global events",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var events []GameEvent
			if tt.filterType != "" {
				events = eventLogger.GetPlayerEvents(tt.requestingPlayer, EventFilters{EventType: tt.filterType})
			} else {
				events = eventLogger.GetPlayerEvents(tt.requestingPlayer, EventFilters{})
			}

			// Verify expected number of events
			assert.Equal(t, len(tt.expectedEvents), len(events), tt.description)

			// Verify expected event IDs
			eventIDs := make([]string, 0, len(events))
			for _, e := range events {
				eventIDs = append(eventIDs, e.ID)
				
				// Verify access control
				if e.Type != EventPlayerDeath && e.Type != EventCombatResult {
					assert.Equal(t, tt.requestingPlayer, e.PlayerID, 
						"Player should only see their own non-global events")
				}
			}
			assert.ElementsMatch(t, tt.expectedEvents, eventIDs, "Unexpected events returned")
		})
	}
}

func TestEventLogger_ConcurrentAccess(t *testing.T) {
	el := NewEventLogger().(*EventLoggerImpl)
	playerID := "testPlayer"
	
	// Test concurrent writes
	t.Run("concurrent writes", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Run("", func(t *testing.T) {
				t.Parallel()
				el.LogEvent(EventPlayerMove, playerID, map[string]interface{}{"test": "data"})
			})
		}
	})

	// Test concurrent reads during writes
	t.Run("concurrent read during write", func(t *testing.T) {
		go func() {
			for i := 0; i < 100; i++ {
				el.LogEvent(EventPlayerMove, playerID, map[string]interface{}{"test": "data"})
			}
		}()

		for i := 0; i < 100; i++ {
			events := el.GetPlayerEvents(playerID, EventFilters{})
			assert.NotNil(t, events, "Should not panic during concurrent access")
		}
	})
}
