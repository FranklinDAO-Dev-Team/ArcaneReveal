extends Area2D

@export var attack_damage = 1

var previous_move
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
	global_position = Vector2((x_pos - 1) * tile_size / 2, (y_pos - 1) * tile_size / 2)
		
func _unhandled_input(event):
	if moving:
		return
	for dir in inputs.keys():
		if event.is_action_pressed(dir):
			move(dir)
			
			
func attack_animation():
	$Sprite.stop()
	$Sprite.play("attack")
	$Sprite.play("idle")

func move(dir):
	ray.target_position = inputs[dir] * tile_size
	ray.force_raycast_update()
	if !ray.is_colliding():
		previous_move = dir 
		var tween = get_tree().create_tween()
		tween.tween_property(self, "position", position + inputs[dir] * tile_size, 1.0/animation_speed).set_trans(Tween.TRANS_SINE)
		moving = true
		await tween.finished
		moving = false
		
func _on_area_entered(area):
	if area.name == "Player" && moving == true:
		area.damage(attack_damage)
		#print(previous_move)
		attack_animation()
