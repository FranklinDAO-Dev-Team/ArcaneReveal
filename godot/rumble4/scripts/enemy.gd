extends Area2D

@export var attack_damage = 1

var previous_move

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
@onready var health_bar = $healthbar
@onready var player = $Player

@onready var game_node = get_parent()

@onready var node_name = get_script().resource_name


func _ready():
	if node_name == "Enemy1":
		var resp = await game_node.client.rpc_async(game_node.session, "query/game/game-state", JSON.stringify({}))
		print(resp)
	update_health()
	position = position.snapped(Vector2.ONE * tile_size)
		
	
func update_health():
	var healthbar = $healthbar  
	healthbar.value = health
	if health < max_health:
		healthbar.visible = false
	else:
		healthbar.visible = true
		
func damage() -> void:
	print("hello")
	print(health)
	health -= 1
	if health == 0:
		queue_free()
	$healthbar.value = health
	print(health)
	
func _process(delta):
	$Sprite.play("idle")
		
func _unhandled_input(event):
	if moving:
		return
	for dir in inputs.keys():
		if event.is_action_pressed(dir):
			move(dir)
			
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
	print("sugma")
	if area.name == "Player" and moving == true:
		area.damage(attack_damage)
		print(previous_move)
		match area.previous_move:
			"right": area.move("left")
			"left": area.move("right")
			"up": area.move("down")
			"down": area.move("up")
