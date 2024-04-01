extends Node2D

func _ready():
	# Enable input events for the staff ability
	set_process_input(true)

func _input(event):
	if event is InputEventMouseButton and event.button_index == MOUSE_BUTTON_LEFT:
		if event.pressed:
			# Start dragging the staff ability
			set_process_input(false)
			var original_position = global_position
			while Input.is_mouse_button_pressed(MOUSE_BUTTON_LEFT):
				global_position = get_global_mouse_position()
				await get_tree().process_frame
			set_process_input(true)
			
			# Check if the staff ability is dropped on a valid tile
			var drop_area = $DropArea
			var overlapping_bodies = drop_area.get_overlapping_bodies()
			if overlapping_bodies.size() > 0:
				# Apply the ability effect on the overlapping tiles
				for body in overlapping_bodies:
					if body.is_in_group("tiles"):
						# Apply fire effect to the tile (implement this based on your game logic)
						body.apply_fire_effect()
			else:
				# Snap back to the original position if not dropped on a valid tile
				global_position = original_position
