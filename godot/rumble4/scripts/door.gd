extends StaticBody2D


# Called when the node enters the scene tree for the first time.
func _ready():
	$OpenDoor.visible = false
	$ClosedDoor.visible = true
	

func _on_key_door_opened():
	$OpenDoor.visible = true
	$ClosedDoor.visible = false
	
