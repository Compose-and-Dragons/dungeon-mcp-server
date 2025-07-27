# Dungeon Specification

## Map Structure
- A dungeon can be defined on a 6 x 6 grid map
- A cell can be:
  - a room
  - a corridor
  - an obstacle (therefore cannot be entered)
- There is an entrance room in the dungeon
- There is an exit room in the dungeon

## Rooms
- Each room can have between 1 and 4 doors
- Each door can lead to another room or a corridor
- A room has a description
- A room cannot have both a monster and an NPC
- A room can have a treasure without monster or NPC
- A room has coordinates on the map
- From a room, the player can:
  - move to another room
  - fight a monster
  - collect a treasure
  - interact with an NPC
- From a room you cand move to other rooms or corridors

## Monsters
- Each room can have a monster
- Each monster can be of type:
  - goblin
  - orc
  - dragon
- Each monster has:
  - A name and description
  - A difficulty level (1 to 10)
  - A number of hit points (1 to 100)
  - Can have a treasure

## NPCs
- Each room can have an NPC
- Each NPC can be of type:
  - merchant
  - healer
  - sage
- Each NPC has:
  - A name and description
  - Can answer the player's questions

## Treasures
- Each treasure can be of type:
  - gold
  - gem
  - artifact
- Each treasure has a value (1 to 1000)

## Items
- In a room you can find healing potions
- Each healing potion has a healing level (1 to 100)
- Each healing potion can be used by the player

## Player
### Character Creation
- The player can choose a class (warrior, mage, thief)
- The player can choose a name
- The player can choose an avatar

### Stats & Progression
- The player can have a level (1 to 100)
- The player can have skills (strength, agility, intelligence)
- The player can gain experience by fighting monsters
- The player can gain experience by interacting with NPCs
- The player can gain gold coins by collecting treasures

### Position & Movement
- The player has a position in the dungeon
- The player can move in the dungeon

### Inventory
- The player can have an inventory
- The inventory can contain:
  - weapons
  - armor
  - healing potions

## Player Actions
- The player can move in the dungeon
- The player can fight monsters
- The player can collect treasures
- The player can interact with NPCs

## Implementation
The dungeon description will be made in a YAML file