extends Area2D

const MAX_HEALTH = 5
var health = MAX_HEALTH
var x_pos
var y_pos
var id
const animation_speed = 3
var moving = false
const tile_size = 32
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
	update_health_ui(false)
	global_position = Vector2((x_pos - 1) * tile_size, (y_pos - 1) * tile_size)
	$StaffPositionTop.position = Vector2(16, 0)  # Adjust this offset
	$StaffPositionBottom.position = Vector2(16, 32)  # Adjust this offset
	$StaffPositionLeft.position = Vector2(0, 16)  # Adjust this offset
	$StaffPositionRight.position = Vector2(32, 16)  # Adjust this offset

func _process(delta):
	$Sprite.play("idle")
	update_health_ui(false)

func update_health_ui(game_over):
	if health == 0 or game_over:
		for child in game_node.get_children():
			if child.name.begins_with("Ability"):
				child.queue_free()
				
		for id in game_node.enemy_state.keys():
			game_node.enemy_state[id].queue_free()
			game_node.enemy_state.erase(id)
			
		# Clear existing walls
		for row in range(game_node.grid_size):
			for col in range(game_node.grid_size):
				if game_node.wall_state[row][col] != null:
					var curr_wall = game_node.wall_state[row][col]
					curr_wall.queue_free()
					game_node.wall_state[row][col] = null
		
		queue_free()
		$"../GameOverLabel".visible = true  # Hide the GameOverLabel node
		# Prompt the user for a new username
		var username_input_screen = preload("res://scenes/TestFinal/username.tscn").instantiate()
		game_node.add_child(username_input_screen)
		username_input_screen.connect("username_submitted", Callable(game_node, "_on_username_submitted"))
	for i in range(MAX_HEALTH):
		$"../Player/LifeBar".get_child(i).visible = health > i

func move(x_curr, y_curr):
	var pos = Vector2((x_curr - 1) * tile_size, (y_curr - 1) * tile_size)
	var tween = get_tree().create_tween()
	tween.tween_property(self, "position", pos, 1.0 / animation_speed).set_trans(Tween.TRANS_SINE)
	moving = true
	await tween.finished
	moving = false
	x_pos = x_curr
	y_pos = y_curr


func attack(dir):
	if dir == "left":
		$AnimationPlayer.play("attack_left")
	elif dir == "right":
		$AnimationPlayer.play("attack_right")
	elif dir == "up":
		$AnimationPlayer.play("attack_up")
	elif dir == "down":
		$AnimationPlayer.play("attack_down")
	else:
		return


func hit_wall():
	$AnimationPlayer.play("hit_wall")

