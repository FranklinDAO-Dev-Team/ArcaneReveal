extends Area2D

@export var attack_damage = 1

var previous_move
var prev_x = 1
var prev_y = 1
var x_pos
var y_pos

var animation_speed = 4
var moving = false
var tile_size = 64
const max_health = 5
var health = max_health
var inputs = {
	"enemy_right": Vector2.RIGHT,
	"enemy_left": Vector2.LEFT,
	"enemy_up": Vector2.UP,
	"enemy_down": Vector2.DOWN
}

@onready var ray = $RayCast2D
@onready var player = $Player

@onready var game_node = get_parent()

@onready var node_name = get_script().resource_name


func _ready():

	$"LifeBar/Life1".play("hearts")
	$"LifeBar/Life2".play("hearts")
	$"LifeBar/Life3".play("hearts")
	$"LifeBar/Life4".play("hearts")
	$"LifeBar/Life5".play("hearts")
	
	update_health_ui()
	
	if node_name == "Enemy1":
		var resp = await game_node.client.rpc_async(game_node.session, "query/game/game-state", JSON.stringify({}))
		#print(resp)
		update_health()

	position = position.snapped(Vector2.ONE * tile_size)


func update_health_ui():
	if health == 0:
		queue_free()
	for i in range(max_health):
		$"LifeBar".get_child(i).visible = health > i
		
func damage() -> void:	
	health -= 1
	if health == 0:
		queue_free()
	update_health_ui()
	
func update_health():
	var healthbar = $healthbar  
	healthbar.value = health
	if health < max_health:
		healthbar.visible = false
	else:
		healthbar.visible = true
	
func _process(delta):
	$Sprite.play("idle")
	#if (prev_x != x_pos or prev_y != y_pos):
		#move(prev_x, prev_y, x_pos, y_pos)
		#prev_x = x_pos
		#prev_y = y_pos
	global_position = Vector2((x_pos - 1) * tile_size / 2, (y_pos - 1) * tile_size / 2)
	update_health_ui()
			
			
func attack_animation():
	$Sprite.stop()
	$Sprite.play("attack")
	$Sprite.play("idle")

func move(prev_x, prev_y, curr_x, curr_y):
	var delta = Vector2(curr_x - prev_x, prev_y - curr_y)
	var tween = get_tree().create_tween()
	tween.tween_property(self, "position", position + delta * tile_size, 1.0/animation_speed).set_trans(Tween.TRANS_SINE)
	moving = true
	await tween.finished
	moving = false
		
func _on_area_entered(area):
	if area.name == "Player" && moving == true:
		area.damage(attack_damage)
		attack_animation()
