extends Area2D

const MAX_HEALTH = 5
var health = MAX_HEALTH
var x_pos
var y_pos
var id

const animation_speed = 3
var moving = false
const tile_size = 32

@onready var game_node = get_parent()
@onready var node_name = get_script().resource_name


func _ready():
	$"LifeBar/Life1".play("hearts")
	$"LifeBar/Life2".play("hearts")
	$"LifeBar/Life3".play("hearts")
	$"LifeBar/Life4".play("hearts")
	$"LifeBar/Life5".play("hearts")
	update_health_ui()
	global_position = Vector2((x_pos - 1) * tile_size, (y_pos - 1) * tile_size)


func _process(delta):
	$Sprite.play("idle")
	update_health_ui()


func update_health_ui():
	if health == 0:
		queue_free()
	for i in range(MAX_HEALTH):
		$"LifeBar".get_child(i).visible = health > i


func update_health():
	var healthbar = $healthbar  
	healthbar.value = health
	if health < MAX_HEALTH:
		healthbar.visible = false
	else:
		healthbar.visible = true


func attack(player_x, player_y):
	var new_pos = Vector2((player_x - 1) * tile_size, (player_y - 1) * tile_size)
	var curr_pos = Vector2((x_pos - 1) * tile_size, (y_pos - 1) * tile_size)
	var tween = get_tree().create_tween()
	tween.tween_property(self, "position", new_pos, 1.0 / (2 * animation_speed)).set_trans(Tween.TRANS_SINE)
	moving = true
	tween.tween_property(self, "position", curr_pos, 1.0 / (2 * animation_speed)).set_trans(Tween.TRANS_SINE)
	await tween.finished
	moving = false


func move(x_curr, y_curr):
	if moving:
		return
	var pos = Vector2((x_curr - 1) * tile_size, (y_curr - 1) * tile_size)
	var tween = get_tree().create_tween()
	tween.tween_property(self, "position", pos, 1.0 / animation_speed).set_trans(Tween.TRANS_SINE)
	moving = true
	await tween.finished
	moving = false
	x_pos = x_curr
	y_pos = y_curr
	

