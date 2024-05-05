extends Node

var enemy_state = {}
var player
var level
var wall_state = []
var game_over = false
var staff_nodes = []

@onready var client : NakamaClient
@onready var socket
@onready var session : NakamaSession
@onready var ray = $RayCast3D

var enemies_defeated = 0
const tile_size = 32
const grid_size = 11
const inputs = {
	"right": Vector2.RIGHT,
	"left": Vector2.LEFT,
	"up": Vector2.UP,
	"down": Vector2.DOWN
}

func _on_enemy_exited(enemy: Area2D):
	if not enemy.is_inside_tree():
		enemies_defeated += 1
		if enemies_defeated == 2:
			# Display WIN text or perform any other desired action
			$GameWinLabel.visible = true
	

func _ready():
	client = Nakama.create_client("defaultkey", "127.0.0.1", 7350, "http")
	socket = Nakama.create_socket_from(client)

	# Authenticate with the Nakama server using Device Authentication
	var device_id = OS.get_unique_id()
	session = await client.authenticate_device_async(device_id)
	if session.is_exception():
		print("An error occurred: %s" % session)
		return

	print("Successfully authenticated: %s" % session)


	var connected_result: NakamaAsyncResult = await socket.connect_async(session)
	if connected_result.is_exception():
		print("An error occurred: %s" % connected_result)
		return
	print("Socket connected.")
	
	socket.received_notification.connect(self._on_notification)
	
	# Check whether account already has persona
	var resp = await client.rpc_async(session, "nakama/show-persona")
	if resp.is_exception():
		print("An error occured: %s", % resp)
		# Create persona	
		resp = await client.rpc_async(session, "nakama/claim-persona", JSON.stringify({"personaTag": "CoolMage"}))
		if resp.is_exception():
			print("An error occured while claiming persona: %s", % resp)
			return
		print("Created PersonaTag: CoolMage")
	else:
		if JSON.parse_string(resp.payload)["status"] == "accepted":
			print("Device already has a persona, skipping creation")
	
	while true:
		# wait until show persona succeed
		var personaResp = await client.rpc_async(session, "nakama/show-persona")
		if personaResp.is_exception():
			print("persona is not registered yet, waiting")
			await get_tree().create_timer(0.25)
			continue
		else:
			if personaResp.payload != null and JSON.parse_string(personaResp.payload)["status"] == "accepted":
				print("Device has a persona, continuing")
				break
	
	# Make CreateGame TXN call
	var random = RandomNumberGenerator.new()
	resp = await client.rpc_async(session, "tx/game/request-game", JSON.stringify({"playerSource": str(random.randi_range(100000, 999999))}))
	#print(resp)
	if resp.is_exception():
		print("An error occured: %s", % resp)
		return
	print("Successfully created game request entity: %s", % resp)
	var payload = await wait_for_game_creation()
	var json = JSON.new()
	var state = json.parse_string(payload)
	for row in range(grid_size):
		wall_state.append([])
		for col in range(grid_size):
			wall_state[row].append(null)
	initialize_state(state)


func _on_rpc_response(result: NakamaAsyncResult):
	if result.is_successful():
		print("RPC succeeded: ", result.payload)
	else:
		print("RPC failed: ", result.error)


func _on_notification(p_notification : NakamaAPI.ApiNotification):
	var notification = JSON.new()
	notification.parse(p_notification.content)
	print("[Notification]: ", notification.data)
	if notification.data.has("event") and notification.data["event"] == "game-over":
		player.update_health_ui(true)
		game_over = true
		return  # Add this line to exit the function if the game is over
	if (not game_over) and notification.data.has("turnEvent"):
		print(game_over)
		process_event(notification.data)
		var payload = await handle_query()
		var json = JSON.new()
		if payload != null:
			var state = json.parse_string(payload)
			if state != null:  # Add this check
				process_state(state)
				print("caught wand turn event")

	if (not game_over) and notification.data.has("event") and notification.data["event"] == "player_turn":
		print("caught player turn event")
		print(game_over)
		var payload = await handle_query()
		var json = JSON.new()
		if payload != null:
			var state = json.parse_string(payload)
			if state != null:  # Add this check
				process_state(state)
				print("caught wand turn event")


func handle_query():
	var resp_getID = await client.rpc_async(session, "query/game/query-game-id-by-persona", JSON.stringify({
		"Persona": "CoolMage",
	}))
	print(resp_getID)  # This should show the response details including payload

	# Create a new JSON object and parse the response payload
	var json = JSON.new()
	var error = json.parse(resp_getID.payload)
	if error == OK:
		var response_dict = json.data  # Access the parsed data

		# Check if the 'Success' key is true and then access 'GameID'
		if response_dict and "Success" in response_dict and response_dict["Success"]:
			var game_id = response_dict["GameID"]
			print("Game ID: ", game_id)

			# Make another RPC call using the retrieved GameID
			var resp_getGameState = await client.rpc_async(session, "query/game/game-state", JSON.stringify({
				"GameID": game_id,  # Use the actual game ID retrieved
			}))
			print(resp_getGameState)  # Print the state response
			return resp_getGameState.payload
		else:
			print("Failed to get Game ID or the response did not indicate success.")
	else:
		print("JSON Parse Error:", json.get_error_message())


func wait_for_game_creation():
	var created = false
	while not created:
		var resp_getID = await client.rpc_async(session, "query/game/query-game-id-by-persona", JSON.stringify({
		"Persona": "CoolMage",
		}))
		print(resp_getID)  # This should show the response details including payload

		# Create a new JSON object and parse the response payload
		var json = JSON.new()
		var error = json.parse(resp_getID.payload)
		if error == OK:
			var response_dict = json.data  # Access the parsed data

			# Check if the 'Success' key is true and then access 'GameID'
			if response_dict and "Success" in response_dict and response_dict["Success"]:
				var game_id = response_dict["GameID"]
				print("Game ID: ", game_id)

				# Make another RPC call using the retrieved GameID
				var resp_getGameState = await client.rpc_async(session, "query/game/game-state", JSON.stringify({
					"GameID": game_id,  # Use the actual game ID retrieved
				}))
				print(resp_getGameState)  # Print the state response
				created = false
				return resp_getGameState.payload
			else:
				print("Failed to get Game ID or the response did not indicate success.")
		else:
			print("JSON Parse Error:", json.get_error_message())
	return


func initialize_state(state : Dictionary):
	var player_init = state["player"]
	var wands = state["wands"]
	var walls = state["walls"]
	var monsters = state["monsters"]
	level = state["level"]
	
	if player == null:
		var player_scene = load("res://scenes/TestFinal/newPlayer.tscn")
		var player_instance = player_scene.instantiate()
		player_instance.x_pos = int(player_init["x"])
		player_instance.y_pos = int(player_init["y"])
		player_instance.health = int(player_init["currHealth"])
		player_instance.id = int(player_init["id"])
		player = player_instance
		add_child(player)
	else:
		player.x_pos = int(player_init["x"])
		player.y_pos = int(player_init["y"])
		player.health = int(player_init["currHealth"])
		player.id = int(player_init["id"])
		
		
	# Remove existing staff nodes
	for staff_node in staff_nodes:
		staff_node.queue_free()
	staff_nodes.clear()
	
	# Create new staff nodes
	for i in range(1, 5):
		var staff_scene = load("res://scenes/Gama/staff_1.tscn")
		var staff_instance = staff_scene.instantiate()
		staff_instance.position = Vector2(41 + (i - 1) * 65, -63)
		staff_instance.scale = Vector2(1.2, 1.2)
		staff_instance.name = "Staff" + str(i)
		player.add_child(staff_instance)
		staff_nodes.append(staff_instance)
		
	for row in range(grid_size):
		for col in range(grid_size):
			if wall_state[row][col] != null:
				var curr_wall = wall_state[row][col]
				curr_wall.queue_free()
				wall_state[row][col] = null
			
	for wall in walls:
		var x_pos = int(wall["x"])
		var y_pos = int(wall["y"])
		
		var position = Vector2((x_pos - 1) * tile_size, (y_pos - 1) * tile_size)
			
		var wall_scene = load("res://scenes/Gama/wall.tscn")
		var wall_instance = wall_scene.instantiate()
				
		wall_instance.global_position = position
			
		# Add the instance as a child to the main scene
		wall_state[x_pos][y_pos] = wall_instance
		add_child(wall_instance)
	
	for monster in monsters:
		# Load the BasicLightning scene
		var enemy_scene = load("res://scenes/TestFinal/newEnemy.tscn")
			
		# Create an instance of the BasicLightning scene
		var enemy_instance = enemy_scene.instantiate()
			
		# Set the global position of the instance to the specified position
		enemy_instance.x_pos = int(monster["x"])
		enemy_instance.y_pos = int(monster["y"])
		enemy_instance.health = int(monster["currHealth"])
		enemy_instance.id = int(monster["id"])
			
		# Add the instance as a child to the main scene
		add_child(enemy_instance)
		enemy_state[enemy_instance.id] = enemy_instance


func process_state(state : Dictionary):
	if level != state["level"]:
		initialize_state(state)
	var player_state = state["player"]
	var player_x = int(player_state["x"])
	var player_y = int(player_state["y"])
	
	player.move(player_x, player_y)
	player.health = int(player_state["currHealth"])
	
	# Update staff nodes' positions
	for staff_node in staff_nodes:
		if is_instance_valid(staff_node):  # Check if the staff node still exists
			staff_node.global_position = player.global_position + staff_node.position
	
	var monsters = state["monsters"]
	var monster_ids = []
	for i in range(monsters.size()):
		var monster = monsters[i]
		var x_pos = int(monster["x"])
		var y_pos = int(monster["y"])
		var health = int(monster["currHealth"])
		var id = int(monster["id"])

		# Set the global position of the instance to the specified position
		if id not in enemy_state.keys():
			var enemy_scene = load("res://scenes/TestFinal/newEnemy.tscn")
			var enemy_instance = enemy_scene.instantiate()
				
			enemy_instance.x_pos = x_pos
			enemy_instance.y_pos = y_pos
			enemy_instance.health = health
			enemy_instance.id = id
				
			# Add the instance as a child to the main scene
			add_child(enemy_instance)
			enemy_state[enemy_instance.id] = enemy_instance
		else:
			var enemy_instance = enemy_state[id]
			enemy_instance.move(x_pos, y_pos)
			enemy_instance.health = health
		monster_ids.append(id)
		
	# Check if an enemy died between turns
	for id in enemy_state.keys():
		if id not in monster_ids:
			enemy_state[id].queue_free()
			enemy_state.erase(id)
		
	
# Called every frame. 'delta' is the elapsed time since the previous frame.
func process_event(notification : Dictionary):
	var event_log = notification["turnEvent"]
	for event in event_log:
		var action = int(event["Event"])
		var x_pos = int(event["X"])
		var y_pos = int(event["Y"])
			
		# Calculate position based on x_pos and y_pos, assuming each square has a size of 32
		var position = Vector2((x_pos - 1) * tile_size, (y_pos - 1) * tile_size)
			
		# Load the BasicLightning scene
		var basic_lightning_scene = load("res://scenes/Gama/BasicLightning.tscn")
			
		# Create an instance of the BasicLightning scene
		var basic_lightning_instance = basic_lightning_scene.instantiate()
			
		# Set the global position of the instance to the specified position
		basic_lightning_instance.global_position = position
			
		# Add the instance as a child to the main scene
		add_child(basic_lightning_instance)
			
		var animation_player = basic_lightning_instance.get_node("Blank")
			
		if animation_player != null and x_pos > 0 and x_pos < 10 and y_pos > 0 and y_pos < 10:
			# Initiate corresponding animation based on action
			match action:
				0:
					# Animate lightning bolt from the sky attack
					# Access the AnimationPlayer in the BasicLightning scene
					animation_player = basic_lightning_instance.get_node("Blank")
					animation_player.play("default")
				1:
					# Animate explosion
					animation_player = basic_lightning_instance.get_node("Explosion")
					animation_player.play("default")
				2:
					# Animate lightning bolt dissipating
					animation_player = basic_lightning_instance.get_node("Spark")
					animation_player.play("default")
				3:
					# Animate lightning bolt dissipating
					animation_player = basic_lightning_instance.get_node("WallActivation")
					animation_player.play("default")
				4: 
					var payload = await handle_query()
					var json = JSON.new()
					if payload != null:
						var state = json.parse_string(payload)
						process_state(state)
					for enemy in enemy_state.values():
						if x_pos == enemy.x_pos and y_pos == enemy.y_pos and player != null:
							enemy.attack(player.x_pos, player.y_pos)
				_:
					# Handle unexpected action
					print("")

func has_player_attacked(dir):
	var new_player_pos = Vector2(player.x_pos, player.y_pos) + 2 * inputs[dir]
	for enemy in enemy_state.values():
		var monster_curr_pos = Vector2(enemy.x_pos, enemy.y_pos)
		if new_player_pos.x == monster_curr_pos.x and new_player_pos.y == monster_curr_pos.y:
			player.attack(dir)
			return true
	return false


func has_wall_collision(dir):
	var new_player_pos = Vector2(player.x_pos, player.y_pos) + inputs[dir]
	if wall_state[new_player_pos.x][new_player_pos.y] != null:
		player.hit_wall()
		return true
	return false


func _input(event: InputEvent) -> void:
	if event.is_action_pressed("query"):
		handle_query()


func _unhandled_input(event):
	for dir in inputs.keys():
		if event.is_action_pressed(dir) and not has_wall_collision(dir):
			if has_player_attacked(dir):
				var resp = await client.rpc_async(session, "tx/game/player-turn", JSON.stringify({
					"GameIDStr": "2",
					"Action": "attack",
					"Direction": dir,
					"WandNum": "0",
					}))
			else:
				var resp = await client.rpc_async(session, "tx/game/player-turn", JSON.stringify({
					"GameIDStr": "2",
					"Action": "move",
					"Direction": dir,
					"WandNum": "0",
					}))



