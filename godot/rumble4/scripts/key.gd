extends StaticBody2D

signal door_opened

var key_taken = false
var next_to_door = false

const TILE_SIZE = 32

func _ready():
	self.position = Vector2(TILE_SIZE * 1 - TILE_SIZE / 2, TILE_SIZE * 3 - TILE_SIZE / 2)

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
	if key_taken and next_to_door and Input.is_action_just_pressed("ui_accept"):
		emit_signal("door_opened")


func _on_area_2d_area_entered(area):
	if area.name == "Player" and not key_taken:
		key_taken = true
		$Sprite2D.queue_free()


func _on_door_entered_door_zone():
	next_to_door = true


func _on_door_exited_door_zone():
	next_to_door = false
