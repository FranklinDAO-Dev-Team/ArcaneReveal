extends StaticBody2D

signal door_opened

var key_taken = false
var next_to_door = false

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
	if key_taken and next_to_door and Input.is_action_just_pressed("ui_accept"):
		emit_signal("door_opened")

func _on_area_2d_body_entered(body: CharacterBody2D):
	if body.name == "Player" and not key_taken:
		print("enter")
		key_taken = true
		$Sprite2D.queue_free()

func _on_door_zone_body_entered(body: CharacterBody2D):
	if body.name == "Player":
		next_to_door = true


func _on_door_zone_body_exited(body: CharacterBody2D):
	if body.name == "Player":
		next_to_door = false
