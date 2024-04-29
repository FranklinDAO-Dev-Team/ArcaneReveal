extends StaticBody2D

signal entered_door_zone
signal exited_door_zone

const FILE_BEGIN = "res://scenes/Beta/level_"


# Called when the node enters the scene tree for the first time.
func _ready():
	$OpenDoor.visible = false
	$ClosedDoor.visible = true
	

func _on_key_door_opened():
	$OpenDoor.visible = true
	$ClosedDoor.visible = false
	var current_scene_file = get_tree().current_scene.scene_file_path
	var next_level_number = current_scene_file.to_int() + 1 
	
	#print(next_level_number)
	
	var next_level_path = FILE_BEGIN + str(next_level_number) +  ".tscn"
	#print(next_level_path)
	get_tree().change_scene_to_file(next_level_path)
	#var next_level = load("res://scenes/JASON_V1/level_2.tscn").instantiate()
	#get_tree().root.add_child(next_level)
	#print("stuff happened")
	#

func _on_door_zone_area_entered(area):
	if area.name == "Player":
		emit_signal("entered_door_zone")


func _on_door_zone_area_exited(area):
	if area.name == "Player":
		emit_signal("exited_door_zone")
