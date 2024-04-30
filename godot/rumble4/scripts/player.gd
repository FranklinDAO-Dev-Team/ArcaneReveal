#var currPos = [0, 0]
#
#
##func _physics_process(delta):
	##move()
	##
##func move():
	##input_movement = Input.get_vector("left", "right", "up", "down")
	##
	##if input_movement != Vector2.ZERO:
		##velocity = input_movement * speed
	##
	##if input_movement == Vector2.ZERO:
		##velocity = Vector2.ZERO
		##
	##move_and_slide()
#
#func _input(event):
	#if event.is_action_pressed("right"):
		#if currPos[0] <= 96:
			#currPos[0] += 32
	#elif event.is_action_pressed("left"):
		#if currPos[0] >= 32:
			#currPos[0] -= 32	
	#elif event.is_action_pressed("up"):
		#if currPos[1] >= 32:
			#currPos[1] -= 32
	#elif event.is_action_pressed("down"):
		#if currPos[1] <= 96:	
			#currPos[1] += 32
#
	#self.position = Vector2(currPos[0], currPos[1])
	
extends Area2D



const MAX_HEALTH = 5
var health = MAX_HEALTH
var previous_move
var previous_position

var animation_speed = 4



@export var moving = false
var tile_size = 64
var inputs = {
	"right": Vector2.RIGHT,
	"left": Vector2.LEFT,
	"up": Vector2.UP,
	"down": Vector2.DOWN
}

@onready var ray = $RayCast2D
@onready var game_node = get_parent()
@onready var animation_player = $"../BasicLightning"

	
func _ready():
	$"../Player/LifeBar/Life1".play("hearts")
	$"../Player/LifeBar/Life2".play("hearts")
	$"../Player/LifeBar/Life3".play("hearts")
	$"../Player/LifeBar/Life4".play("hearts")
	$"../Player/LifeBar/Life5".play("hearts")
	update_health_ui()

	position = position.snapped(Vector2.ONE * tile_size)
	$StaffPositionTop.position = Vector2(16, 0)  # Adjust this offset
	$StaffPositionBottom.position = Vector2(16, 32)  # Adjust this offset
	$StaffPositionLeft.position = Vector2(0, 16)  # Adjust this offset
	$StaffPositionRight.position = Vector2(32, 16)  # Adjust this offset

	if game_node == null:
		print("game_node is null, cannot access session")
		return

	if game_node.session == null:
		print("game_node.session is null, cannot proceed")
		return

func readJSON(json_file_path):
	var file = FileAccess.open(json_file_path, FileAccess.READ)
	var content = file.get_as_text()
	var json = JSON.new()
	var finish = json.parse_string(content)
	return finish
	
# Called every frame. 'delta' is the elapsed time since the previous frame.
func process_data():
	var json_file_path = "res://testInput.json"
	var data_received = readJSON(json_file_path);
	
	for data in data_received.data:
		var x_pos = int(data[0])
		var y_pos = int(data[1])
		var action = int(data[2])
		
		#print(str(x_pos) + " " + str(y_pos) + " " + str(action))
		
		# Calculate position based on x_pos and y_pos, assuming each square has a size of 32
		var position = Vector2(x_pos * 32 - 32, y_pos * 32 - 32)
		
		# Load the BasicLightning scene
		var basic_lightning_scene = load("res://scenes/Gama/BasicLightning.tscn")
		
		# Create an instance of the BasicLightning scene
		var basic_lightning_instance = basic_lightning_scene.instantiate()
		
		# Set the global position of the instance to the specified position
		basic_lightning_instance.global_position = position
		
		# Add the instance as a child to the main scene
		$"../".add_child(basic_lightning_instance)
		
		var animation_player = basic_lightning_instance.get_node("Blank")
		
		if animation_player != null:
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
				_:
					# Handle unexpected action
					print("")


func _process(delta):
	$Sprite.play("idle")

func update_health_ui():
	for i in range(MAX_HEALTH):
		$"../Player/LifeBar".get_child(i).visible = health > i


func _input(event: InputEvent) -> void:
	if event.is_action_pressed("query"):
		handle_query()

func handle_query():
	var resp_getID = await game_node.client.rpc_async(game_node.session, "query/game/query-game-id-by-persona", JSON.stringify({
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
			var resp_getGameState = await game_node.client.rpc_async(game_node.session, "query/game/game-state", JSON.stringify({
				"GameID": game_id,  # Use the actual game ID retrieved
			}))
			print(resp_getGameState)  # Print the state response
		else:
			print("Failed to get Game ID or the response did not indicate success.")
	else:
		print("JSON Parse Error:", json.get_error_message())

	
	
	


func damage(damage) -> void:
	health -= damage
	if health == 0:
		queue_free()
		$"../GameOverLabel".visible = true  # Hide the GameOverLabel node
	update_health_ui()


func _unhandled_input(event):
	if moving:
		return
	for dir in inputs.keys():
		if event.is_action_pressed(dir):
			$RayCast2DEnemy.target_position = inputs[dir] * tile_size
			$RayCast2DEnemy.force_raycast_update()
			
			if $RayCast2DEnemy.is_colliding() and $RayCast2DEnemy.get_collider().name.begins_with("Enemy"):
				var resp = await game_node.client.rpc_async(game_node.session, "tx/game/player-turn", JSON.stringify({
					"GameIDStr": "2",
					"Action": "attack",
					"Direction": dir,
					"WandNum": "0",
					}))
				#print(resp)
				if resp != null:
					move(dir)
			else:
				var resp = await game_node.client.rpc_async(game_node.session, "tx/game/player-turn", JSON.stringify({
					"GameIDStr": "2",
					"Action": "move",
					"Direction": dir,
					"WandNum": "0",
					}))
		
				if resp != null:
					move(dir)
					


func move(dir):
	#print("move")
	previous_move = dir
	previous_position = position
	#print(previous_position)
	ray.target_position = inputs[dir] * tile_size
	ray.force_raycast_update()
	if !ray.is_colliding():
		#print("no ray collide")
		previous_move = dir 
		#position += inputs[dir] * tile_size
		var tween = get_tree().create_tween()
		tween.tween_property(self, "position", position + inputs[dir] * tile_size, 1.0/animation_speed).set_trans(Tween.TRANS_SINE)
		moving = true 
		await tween.finished
		moving = false
	else:
		$AnimationPlayer.play("hit_wall")
	
func recoil():
	#print("recoil triggered")
	var tween = get_tree().create_tween()
	tween.tween_property(self, "position", previous_position, 1.0/animation_speed).set_trans(Tween.TRANS_SINE)
	await tween.finished
	#print(position)


func _on_area_entered(area):
	#print("do stuff")
	if area.name.begins_with("Enemy") && moving == true:
		self.recoil()		
		area.damage()

		
		match area.previous_move:
			"right": area.move("left")
			"left": area.move("right")
			"up": area.move("down")
			"down": area.move("up")
		
		
		

