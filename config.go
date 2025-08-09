package main

// Config holds all game configuration values
type Config struct {
	Map struct {
		Width  int
		Height int
	}
	Game struct {
		BotNumber        int
		TurnLength       int
		MaxTurns         int
		VictoryNumber    int
		DefaultDirection Direction
	}
	Combat struct {
		ZombieCutoff    int
		WeaponStrength  int
		PlayerMinAttack int
		PlayerMaxAttack int
	}
	Player struct {
		NameMaxLength int
	}
	Server struct {
		IDSalt string
	}
	TerrainResources map[Terrain]TerrainReward
}

// TerrainReward defines what resources a terrain type provides
type TerrainReward struct {
	amount    int
	givesCard Card
}

// NewDefaultConfig creates a new configuration with default values
func NewDefaultConfig() *Config {
	config := &Config{}

	// Map configuration
	config.Map.Width = 100
	config.Map.Height = 100

	// Game configuration
	config.Game.BotNumber = 0
	config.Game.TurnLength = 15
	config.Game.MaxTurns = 500
	config.Game.VictoryNumber = 2
	config.Game.DefaultDirection = South

	// Combat configuration
	config.Combat.ZombieCutoff = 13
	config.Combat.WeaponStrength = 6
	config.Combat.PlayerMinAttack = 1
	config.Combat.PlayerMaxAttack = 6

	// Player configuration
	config.Player.NameMaxLength = 20

	// Server configuration
	config.Server.IDSalt = "6LIBN8OWPzTKctUvbZtXV2mFn2tCq3qZKjHYbTTnLWtu6oGTU3ow3tuNx9SBTuND"

	// Terrain resources configuration
	config.TerrainResources = map[Terrain]TerrainReward{
		City:       {amount: 1, givesCard: Weapon},
		Forest:     {amount: 2, givesCard: Wood},
		Farm:       {amount: 1, givesCard: Food},
		Laboratory: {amount: 1, givesCard: Research},
	}

	return config
}

// Global configuration instance
var gameConfig = NewDefaultConfig()

// Legacy global variables for backward compatibility
// These will be removed in future refactoring phases
const mapWidth = 100
const mapHeight = 100

var botNumber = gameConfig.Game.BotNumber
var zombieCutoff = gameConfig.Combat.ZombieCutoff
var weaponStrength = gameConfig.Combat.WeaponStrength
var playerNameMaxLength = gameConfig.Player.NameMaxLength
var playerMinAttack = gameConfig.Combat.PlayerMinAttack
var playerMaxAttack = gameConfig.Combat.PlayerMaxAttack
var turnLength = gameConfig.Game.TurnLength
var maxTurns = gameConfig.Game.MaxTurns
var victoryNumber = gameConfig.Game.VictoryNumber
var defaultDirection = gameConfig.Game.DefaultDirection
var terrainResources = gameConfig.TerrainResources
var idSalt = gameConfig.Server.IDSalt

// Legacy global variables that weren't in config
var hasWon bool = false
