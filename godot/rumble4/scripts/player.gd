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
var previous_position = Vector2()
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
@onready var animation_player = $"../BasicLightning"

	
func _ready():
	$"../LifeBar/Life1".play("hearts")
	$"../LifeBar/Life2".play("hearts")
	$"../LifeBar/Life3".play("hearts")
	$"../LifeBar/Life4".play("hearts")
	$"../LifeBar/Life5".play("hearts")
	update_health_ui()

	position = position.snapped(Vector2.ONE * tile_size)
	$StaffPositionTop.position = Vector2(16, 0)  # Adjust this offset
	$StaffPositionBottom.position = Vector2(16, 32)  # Adjust this offset
	$StaffPositionLeft.position = Vector2(0, 16)  # Adjust this offset
	$StaffPositionRight.position = Vector2(32, 16)  # Adjust this offset
	#position += Vector2.ONE * tile_size / 





func readJSON(json_file_path):
	var file = FileAccess.open(json_file_path, FileAccess.READ)
	var content = file.get_as_text()
	var json = JSON.new()
	var finish = json.parse_string(content)
	return finish
	
# Called every frame. 'delta' is the elapsed time since the previous frame.
func process_data():
	# cast spell: p _ _ m _ w
	# (3, 1, 0), (5, 1, 0), (7, 1, 1), (9, 1, 0) (10, 1, 2) 
	print("entered process_data")
	
	var json_file_path = "res://testInput.json"
	var data_received = readJSON(json_file_path);
	print(data_received.data)
	for data in data_received.data:
		var x_pos = int(data[0])
		var y_pos = int(data[1])
		var action = int(data[2])
		
		print(str(x_pos) + " " + str(y_pos) + " " + str(action))
		
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
		
		# Access the AnimationPlayer in the BasicLightning scene
		var animation_player = basic_lightning_instance.get_node("AnimationPlayer")
		if animation_player != null:
			# Initiate corresponding animation based on action
			match action:
				0:
					# Animate lightning bolt from the sky attack
					animation_player.play("default")
					print("lightning at: " + str(position.x) + ", " + str(position.y))
				1:
					# Animate explosion
					animation_player.play("explosion")
					print("explosion at: " + str(position.x) + ", " + str(position.y))
				2:
					# Animate lightning bolt dissipating
					animation_player.play("lightning_dissipate")
					print("dissipate at: " + str(position.x) + ", " + str(position.y))
				_:
					# Handle unexpected action
					print("Unexpected action:", action)
			
			# Queue the instance for deletion after the animation finishes
			#animation_player.queue_free()
			#animation_player.connect("animation_finished", basic_lightning_instance, "_on_animation_finished")
		else:
			print("AnimationPlayer not found in BasicLightning scene")

# Callback function to delete the instance after the animation finishes
#func _on_animation_finished():
	#var instance = get_parent()
	#instance.queue_free()








func _process(delta):
	$Sprite.play("idle")


func update_health_ui():
	for i in range(MAX_HEALTH):
		$"../LifeBar".get_child(i).visible = health > i


func _input(event: InputEvent) -> void:
	#if event.is_action_pressed("ui_accept"):
		#damage()
	if event.is_action_pressed("enemy_left"):
		$AnimationPlayer.play("attack_left")
	if event.is_action_pressed("enemy_right"):
		$AnimationPlayer.play("attack_right")
	if event.is_action_pressed("enemy_down"):
		$AnimationPlayer.play("attack_down")
	if event.is_action_pressed("enemy_up"):
		$AnimationPlayer.play("attack_up")
	


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
			move(dir)


func move(dir):
	var initial_position = position  # Store the current position before moving
	previous_move = dir

	var tween = get_tree().create_tween()
	tween.tween_property(self, "position", position + inputs[dir] * tile_size, 1.0/animation_speed)
	tween.set_trans(Tween.TRANS_SINE)
	moving = true
	await tween.finished
	moving = false

	# Check collision after moving
	if is_colliding_with_enemy():
		position = initial_position  # Move back to the previous position if collided with an enemy

func is_colliding_with_enemy() -> bool:
	for area in get_overlapping_areas():
		if area.name == "Enemy1" or area.name == "Enemy2":
			return true
	return false



func _on_area_entered(area):
	print("do stuff")
	if (area.name == "Enemy1" or area.name == "Enemy2") && moving == true:
		area.damage()
		
		
		
		

