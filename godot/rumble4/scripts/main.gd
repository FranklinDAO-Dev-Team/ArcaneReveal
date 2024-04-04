extends Node2D

const MAX_HEALTH = 5
var health = MAX_HEALTH

func _ready() -> void:
	#$LifeBar/Life1.play("hearts")
	#$LifeBar/Life2.play("hearts")
	#$LifeBar/Life3.play("hearts")
	#$LifeBar/Life4.play("hearts")
	#$LifeBar/Life5.play("hearts")
	#update_health_ui()
	#$HealthBar.max_value = MAX_HEALTH
	print("todo")

func update_health_ui():
	set_health_bar()

func set_health_bar() -> void:
	for i in range(MAX_HEALTH):
		$LifeBar.get_child(i).visible = health > i
		

func _input(event: InputEvent) -> void:
	if event.is_action_pressed("ui_accept"):
		damage()

func damage() -> void:
	health -= 1
	if health < 0:
		health = MAX_HEALTH
	update_health_ui()
	
