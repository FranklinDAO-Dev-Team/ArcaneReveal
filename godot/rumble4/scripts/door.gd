extends StaticBody2D

signal entered_door_zone
signal exited_door_zone

# Called when the node enters the scene tree for the first time.
func _ready():
	$OpenDoor.visible = false
	$ClosedDoor.visible = true
	

func _on_key_door_opened():
	$OpenDoor.visible = true
	$ClosedDoor.visible = false
	

func _on_door_zone_area_entered(area):
	if area.name == "Player":
		emit_signal("entered_door_zone")


func _on_door_zone_area_exited(area):
	if area.name == "Player":
		emit_signal("exited_door_zone")
