package main

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of game event
type EventType string

// Event types
const (
	// EventPlayerJoin is triggered when a new player joins the game
	EventPlayerJoin EventType = "player_join"
	// EventPlayerMove is triggered when a player moves to a new tile
	EventPlayerMove EventType = "player_move"
	// EventCardUsage is triggered when a player uses a card
	EventCardUsage EventType = "card_usage"
	// EventPlayerDeath is triggered when a player dies
	EventPlayerDeath EventType = "player_death"
	// EventCombatResult is triggered after combat resolution
	EventCombatResult EventType = "combat_result"
	// EventResourceGained is triggered when a player gains resources
	EventResourceGained EventType = "resource_gained"
	// EventGameTick is triggered on each game tick
	EventGameTick EventType = "game_tick"
	// EventCardPlayed is triggered when a player plays a card
	EventCardPlayed EventType = "card_played"
	// EventCardUsed is triggered when a player uses a card in combat
	EventCardUsed EventType = "card_used"
	// EventCardSelected is triggered when a player selects a card to consume
	EventCardSelected EventType = "card_selected"
	// EventCardConsumed is triggered when a player consumes a card
	EventCardConsumed EventType = "card_consumed"
	// EventDiceRoll is triggered when a player rolls dice for combat
	EventDiceRoll EventType = "dice_roll"
	// EventCombatStart is triggered when combat starts on a tile
	EventCombatStart EventType = "combat_start"
	// EventZombieSpawn is triggered when zombies spawn on a tile
	EventZombieSpawn EventType = "zombie_spawn"
	// EventCardDrawn is triggered when a player draws a card
	EventCardDrawn EventType = "card_drawn"
	// EventCardDiscarded is triggered when a player discards a card
	EventCardDiscarded EventType = "card_discarded"
)

// EventTypeList returns all valid event types
func EventTypeList() []EventType {
	return []EventType{
		EventPlayerMove,
		EventCardUsage,
		EventPlayerDeath,
		EventCombatResult,
		EventResourceGained,
		EventGameTick,
	}
}

// GameEvent represents a single game event
type GameEvent struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	PlayerID  string                 `json:"player_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details"`
	Turn      int64                  `json:"turn"`
}

// EventFilters defines filtering options for event queries
type EventFilters struct {
	EventType EventType    // Filter by event type
	Since     time.Time    // Only return events after this time
	Limit     int          // Maximum number of events to return (0 for no limit)
	LastTurns int          // Number of most recent turns to return events for (0 means no limit)
}

// EventLogger defines the interface for logging game events
type EventLogger interface {
	// LogEvent logs a new game event
	LogEvent(eventType EventType, playerID string, details map[string]interface{})
	// GetPlayerEvents returns all events visible to the specified player
	GetPlayerEvents(playerID string, filters EventFilters) []GameEvent
	// GetEventCount returns the total number of events in the log
	GetEventCount() int64
	// GetEventTypeCount returns the number of events of a specific type
	GetEventTypeCount(eventType EventType) int64
}

// EventLoggerImpl is the in-memory implementation of EventLogger
type EventLoggerImpl struct {
	events      []GameEvent
	mu          sync.RWMutex
	eventCounts map[EventType]*int64
	totalEvents int64
	currentTurn int64
}

// NewEventLogger creates a new instance of EventLogger
func NewEventLogger() EventLogger {
	// Initialize event counts
	eventCounts := make(map[EventType]*int64)
	for _, et := range EventTypeList() {
		var count int64
		eventCounts[et] = &count
	}

	return &EventLoggerImpl{
		events:      make([]GameEvent, 0, 1000), // Pre-allocate some capacity
		eventCounts: eventCounts,
	}
}

// reverseSlice reverses the order of elements in a slice of GameEvent
func reverseSlice(events []GameEvent) {
	for i := len(events)/2 - 1; i >= 0; i-- {
		opp := len(events) - 1 - i
		events[i], events[opp] = events[opp], events[i]
	}
}

// LogEvent adds a new event to the log
func (el *EventLoggerImpl) LogEvent(eventType EventType, playerID string, details map[string]interface{}) {
	el.mu.Lock()
	defer el.mu.Unlock()

	event := GameEvent{
		ID:        uuid.New().String(),
		Type:      eventType,
		PlayerID:  playerID,
		Timestamp: time.Now(),
		Details:   details,
		Turn:      el.currentTurn,
	}

	el.events = append(el.events, event)
	atomic.AddInt64(&el.totalEvents, 1)
	if count, exists := el.eventCounts[eventType]; exists {
		atomic.AddInt64(count, 1)
	}
}

// GetPlayerEvents returns all events visible to the specified player
func (e *EventLoggerImpl) GetPlayerEvents(playerID string, filters EventFilters) []GameEvent {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// If LastTurns is set, find the turn number N turns ago
	var minTurn int64 = 0
	if filters.LastTurns > 0 && e.currentTurn > int64(filters.LastTurns) {
		minTurn = e.currentTurn - int64(filters.LastTurns)
	}

	var result []GameEvent

	// Iterate in reverse chronological order (newest first)
	for i := len(e.events) - 1; i >= 0; i-- {
		event := e.events[i]

		// Skip events from before the minimum turn
		if filters.LastTurns > 0 && event.Turn < minTurn {
			continue
		}

		// Check if event should be visible to this player
		visible := false

		switch event.Type {
		case EventPlayerMove, EventCardUsage:
			// Only show player's own actions
			visible = event.PlayerID == playerID

		case EventPlayerDeath:
			// Show all player deaths
			visible = true

		case EventCombatResult:
			// Show combat events where player was involved
			if participants, ok := event.Details["involved_players"].([]interface{}); ok {
				for _, p := range participants {
					if pID, ok := p.(string); ok && pID == playerID {
						visible = true
						break
					}
				}
			} else if participants, ok := event.Details["involved_players"].([]string); ok {
				// Handle case where participants is a []string
				for _, pID := range participants {
					if pID == playerID {
						visible = true
						break
					}
				}
			} else if event.PlayerID == playerID {
				// If no participants but player is the event owner, show it
				visible = true
			}

		default:
			// For any other event type, only show if it's the player's own event
			visible = event.PlayerID == playerID
		}

		if !visible {
			continue
		}

		// Apply type filter if specified
		if filters.EventType != "" && event.Type != filters.EventType {
			continue
		}

		// Apply time filter if specified
		if !filters.Since.IsZero() && event.Timestamp.Before(filters.Since) {
			continue
		}

		result = append(result, event)

		// Apply limit if specified
		if filters.Limit > 0 && len(result) >= filters.Limit {
			break
		}
	}

	// Return in chronological order (oldest first)
	reverseSlice(result)
	return result
}

// GetEventCount returns the total number of events in the log
func (l *EventLoggerImpl) GetEventCount() int64 {
	return atomic.LoadInt64(&l.totalEvents)
}

// GetEventTypeCount returns the number of events of a specific type
func (l *EventLoggerImpl) GetEventTypeCount(eventType EventType) int64 {
	if count, exists := l.eventCounts[eventType]; exists {
		return atomic.LoadInt64(count)
	}
	return 0
}

// Helper function to check if a slice contains a player ID
func containsPlayer(players interface{}, playerID string) bool {
	if players == nil {
		return false
	}
	playerSlice, ok := players.([]string)
	if !ok {
		return false
	}
	
	for _, id := range playerSlice {
		if id == playerID {
			return true
		}
	}
	return false
}
