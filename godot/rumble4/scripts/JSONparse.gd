#extends Node2D
#
#func readJSON(json_file_path):
	#var file = FileAccess.open(json_file_path, FileAccess.READ)
	#var content = file.get_as_text()
	#var json = JSON.new()
	#var finish = json.parse_string(content)
	#return finish
	#
## Called every frame. 'delta' is the elapsed time since the previous frame.
#func process_data():
	## cast spell: p _ _ m _ w
	## (3, 1, 0), (5, 1, 0), (7, 1, 1), (9, 1, 0) (10, 1, 2) 
	#print("entered process_data")
	#
	#var json_file_path = "res://testInput.json"
	#var data_received = readJSON(json_file_path);
	#print(data_received.data)
	#for data in data_received.data:
		#var x_pos = int(data[0])
		#var y_pos = int(data[1])
		#var action = int(data[2])
		#
		#print(str(x_pos) + " " + str(y_pos) + " " + str(action))
		#
		## Calculate position based on x_pos and y_pos, assuming each square has a size of 32
		#var position = Vector2(x_pos * 32, y_pos * 32)
		#
		## Instantiate animation player at the position
		#var animation_player = AnimationPlayer.new()
		#add_child(animation_player)
		##animation_player.global_position = position
		#
		## Initiate corresponding animation based on action
		#match action:
			#0:
				## Animate lightning bolt from the sky attack
				##animation_player.play("lightning_bolt_attack")
				#print("lightning at: " + str(position.x) + ", " + str(position.y))
			#1:
				## Animate explosion
				##animation_player.play("explosion")
				#print("explosion at: " + str(position.x) + ", " + str(position.y))
			#2:
				## Animate lightning bolt dissipating
				##animation_player.play("lightning_dissipate")
				#print("dissipate at: " + str(position.x) + ", " + str(position.y))
			#_:
				## Handle unexpected action
				#print("Unexpected action:", action)

extends Node

func _ready():
	# Sample stringified JSON data
	var jsonString = """
	{
		"payload": {
			"player": {
				"x": 1,
				"y": 5,
				"maxHealth": 5,
				"currHealth": 3
			},
			"wands": [
				{"number": 0, "isAvailable": true},
				{"number": 1, "isAvailable": true},
				{"number": 2, "isAvailable": true},
				{"number": 3, "isAvailable": true}
			],
			"walls": [
				{"x": 3, "y": 2, "type": 0},
				{"x": 5, "y": 2, "type": 0},
				{"x": 2, "y": 5, "type": 0},
				{"x": 2, "y": 7, "type": 0}
			],
			"monsters": [
				{"x": 9, "y": 5, "type": 2},
				{"x": 1, "y": 7, "type": 2},
				{"x": 3, "y": 5, "type": 2},
				{"x": 1, "y": 3, "type": 2}
			]
		}
	}
	"""

	# Parse JSON string into a Dictionary
	
	var json = JSON.new()
	var finish = json.parse_string(jsonString)

	# Access player data
	var player = json["payload"]["player"]
	var playerX = player["x"]
	var playerY = player["y"]
	var maxHealth = player["maxHealth"]
	var currHealth = player["currHealth"]
	print("Player Position:", playerX, ",", playerY)
	print("Max Health:", maxHealth)
	print("Current Health:", currHealth)

	# Access wands data
	var wands = json["payload"]["wands"]
	for wand in wands:
		var number = wand["number"]
		var isAvailable = wand["isAvailable"]
		print("Wand Number:", number)
		print("Is Available:", isAvailable)

	# Access walls data
	var walls = json["payload"]["walls"]
	for wall in walls:
		var wallX = wall["x"]
		var wallY = wall["y"]
		var wallType = wall["type"]
		print("Wall at:", wallX, ",", wallY)
		print("Type:", wallType)

	# Access monsters data
	var monsters = json["payload"]["monsters"]
	for monster in monsters:
		var monsterX = monster["x"]
		var monsterY = monster["y"]
		var monsterType = monster["type"]
		print("Monster at:", monsterX, ",", monsterY)
		print("Type:", monsterType)
