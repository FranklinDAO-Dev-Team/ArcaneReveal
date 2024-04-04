extends Node


# Called when the node enters the scene tree for the first time.
func parse():
	var data_to_send = ["a", "b", "c"]
	var json_string = JSON.stringify(data_to_send)
	# Save data
	# ...
	# Retrieve data
	var json = JSON.new()
	var error = json.parse(json_string)
	if error == OK:
		var data_received = json.data
		if typeof(data_received) == TYPE_ARRAY:
			print(data_received) # Prints array
		else:
			print("Unexpected data")
	else:
		print("JSON Parse Error: ", json.get_error_message(), " in ", json_string, " at line ", json.get_error_line())

func readJSON(json_file_path):
	var file = FileAccess.open(json_file_path, FileAccess.READ)
	var content = file.get_as_text()
	var json = JSON.new()
	var finish = json.parse_string(content)
	return finish
	
# Called every frame. 'delta' is the elapsed time since the previous frame.
func process_data(delta):
	# cast spell: p _ _ m _ w
	# (3, 1, 0), (5, 1, 0), (7, 1, 1), (9, 1, 0) (10, 1, 2) 
	print("entered process_data")
	
	var json_file_path = "res://testInput.json"
	var data_received = readJSON(json_file_path);
	for data in data_received:
		var x_pos = data[0]
		var y_pos = data[1]
		var action = data[2]
		
		# Calculate position based on x_pos and y_pos, assuming each square has a size of 32
		var position = Vector2(x_pos * 32, y_pos * 32)
		
		# Instantiate animation player at the position
		var animation_player = AnimationPlayer.new()
		add_child(animation_player)
		animation_player.global_position = position
		
		# Initiate corresponding animation based on action
		match action:
			0:
				# Animate lightning bolt from the sky attack
				#animation_player.play("lightning_bolt_attack")
				print("lightning at: " + position)
			1:
				# Animate explosion
				#animation_player.play("explosion")
				print("explosion at: " + position)
			2:
				# Animate lightning bolt dissipating
				#animation_player.play("lightning_dissipate")
				print("dissipate at: " + position)
			_:
				# Handle unexpected action
				print("Unexpected action:", action)



