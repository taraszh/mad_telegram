package emoji

import (
	"math/rand"
	"time"
)

type EndStringEmojifier struct {
	rng *rand.Rand
}

func NewEndStringEmojifier() *EndStringEmojifier {
	return &EndStringEmojifier{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (e *EndStringEmojifier) Emojify(message string) (string, error) {
	randomEmoji := themedEmojis[e.rng.Intn(len(themedEmojis))]

	return message + " " + randomEmoji, nil
}

var themedEmojis = []string{
	// Mechanisms
	"⚙️",
	"🛞",
	"🛠️", // Hammer and Wrench
	"🔧",  // Wrench
	"🔨",  // Hammer
	"🪚",  // Carpentry Saw
	"🔩",  // Nut and Bolt
	"⛓️", // Chains
	"🪝",  // Hook
	"🔗",  // Link
	"🧰",  // Toolbox
	"⚒️", // Hammer and Pick
	"🕰️", // Mantelpiece Clock (for clockwork)
	"⌛",  // Hourglass
	"⏳",  // Hourglass with Flowing Sand
	"⚖️", // Balance Scale
	"🪤",  // Mouse Trap
	"🧲",  // Magnet
	"🪚",  // Carpentry Saw
	"🔒",  // Locked
	"🗝️", // Old Key
	"🔭",  // Telescope
	"🧪",  // Test Tube (for alchemical/mechanical experiments)
	"⚗️", // Alembic (for early mechanical chemistry)

	// Medieval
	"👑",
	"🧛",
	"🏰", // Castle
	"👑", // Crown
	"🧝", // Elf
	"🕍", // Synagogue (for medieval religious buildings)
	"⛪", // Church
	"🔔", // Bell
	"🏺", // Amphora (for medieval pottery)
	"🦁", // Lion (heraldic symbol)

	// Weapons
	"🗡️", "⚔️", "🏹", "🛡️", "💣", "🪓",

	// Nature
	"🌚",
	"🎑",
	"🔥",
	"🌊",
	"🍄‍🟫",
	"🌿", // Herb
	"🌱", // Seedling
	"🌳", // Deciduous Tree
	"🌲", // Evergreen Tree
	//"🌴",  // Palm Tree
	//"🌵",  // Cactus
	"🌾",  // Sheaf of Rice
	"🍀",  // Four Leaf Clover
	"🍁",  // Maple Leaf
	"🍂",  // Fallen Leaf
	"🍃",  // Leaf Fluttering in Wind
	"🌺",  // Hibiscus
	"🌸",  // Cherry Blossom
	"🌹",  // Rose
	"🌼",  // Blossom
	"🍄",  // Mushroom
	"🐚",  // Spiral Shell
	"🌊",  // Water Wave
	"🏞️", // National Park
	"⛰️", // Mountain
	"🌄",  // Sunrise Over Mountains
	"🌅",  // Sunrise
	"🌙",  // Crescent Moon
	"⭐",  // Star
	"🌈",  // Rainbow
	"❄️", // Snowflake
	//"☁️", // Cloud
	"🌧️", // Cloud with Rain
	"⛈️", // Cloud with Lightning and Rain
	"🌋",  // Volcano
	"🪨",  // Rock
	"🌍",  // Globe Showing Europe-Africa

	// Space
	"🪐", "🌠", "🌟", "☄️",

	// Magic
	"✨",
	"🪬",
	"🧿",
	"🪄",    // Magic Wand
	"✨",    // Sparkles
	"🔮",    // Crystal Ball
	"💫",    // Dizzy
	"🧙‍♂️", // Man Mage
	"🧚‍♀️", // Woman Fairy
	"🧝",    // Elf
	"🧞‍♂️", // Man Genie
	"⭐",    // Star
	"🌟",    // Glowing Star
	"🎇",    // Sparkler
	"🎆",    // Fireworks
	"🪐",    // Ringed Planet
	"🌌",    // Milky Way
	"🦄",    // Unicorn
	"🐉",    // Dragon
	"🧪",    // Test Tube (for potions)
	"📜",    // Scroll
	"⚗️",   // Alembic (for alchemy)
}
