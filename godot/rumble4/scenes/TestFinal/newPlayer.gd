extends Area2D

const MAX_HEALTH = 5
var health = MAX_HEALTH
var x_pos
var y_pos
var id

const animation_speed = 3
var moving = false
const tile_size = 64

var inputs = {
	"right": Vector2.RIGHT,
	"left": Vector2.LEFT,
	"up": Vector2.UP,
	"down": Vector2.DOWN
}

@onready var ray = $RayCast2D
@onready var game_node = get_parent()

func _ready():
	$"../Player/LifeBar/Life1".play("hearts")
	$"../Player/LifeBar/Life2".play("hearts")
	$"../Player/LifeBar/Life3".play("hearts")
	$"../Player/LifeBar/Life4".play("hearts")
	$"../Player/LifeBar/Life5".play("hearts")
	update_health_ui()
	global_position = Vector2((x_pos - 1) * tile_size / 2, (y_pos - 1) * tile_size / 2)
	
	$StaffPositionTop.position = Vector2(16, 0)  # Adjust this offset
	$StaffPositionBottom.position = Vector2(16, 32)  # Adjust this offset
	$StaffPositionLeft.position = Vector2(0, 16)  # Adjust this offset
	$StaffPositionRight.position = Vector2(32, 16)  # Adjust this offset
	

func _process(delta):
	$Sprite.play("idle")
	update_health_ui()
	

func update_health_ui():
	if health == 0:
		queue_free()
		$"../GameOverLabel".visible = true  # Hide the GameOverLabel node
	for i in range(MAX_HEALTH):
		$"../Player/LifeBar".get_child(i).visible = health > i


func _unhandled_input(event):
	if moving:
		return
	for dir in inputs.keys():
		if event.is_action_pressed(dir):
			$RayCast2DEnemy.target_position = inputs[dir] * tile_size
			$RayCast2DEnemy.force_raycast_update()
			
			if $RayCast2DEnemy.is_colliding() and $RayCast2DEnemy.get_collider().name.begins_with("Enemy"):
				var resp = await game_node.client.rpc_async(game_node.session, "tx/game/player-turn", JSON.stringify({
					"GameIDStr": "2",
					"Action": "attack",
					"Direction": dir,
					"WandNum": "0",
					}))

			else:
				var resp = await game_node.client.rpc_async(game_node.session, "tx/game/player-turn", JSON.stringify({
					"GameIDStr": "2",
					"Action": "move",
					"Direction": dir,
					"WandNum": "0",
					}))


func move(x_curr, y_curr):
	var pos = Vector2((x_curr - 1) * tile_size / 2, (y_curr - 1) * tile_size / 2)
	var tween = get_tree().create_tween()
	tween.tween_property(self, "position", pos, 1.0 / animation_speed).set_trans(Tween.TRANS_SINE)
	moving = true
	await tween.finished
	moving = false
	x_pos = x_curr
	y_pos = y_curr


func hit_wall():
	$AnimationPlayer.play("hit_wall")

