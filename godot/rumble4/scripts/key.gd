extends StaticBody2D

signal door_opened

var key_taken = false
var next_to_door = false

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
	if key_taken and next_to_door:
		emit_signal("door_opened")

func _on_area_2d_body_entered(body: CharacterBody2D):
	if body.name == "Player":
		print("enter")
		key_taken = true
		$Sprite2D.queue_free()

