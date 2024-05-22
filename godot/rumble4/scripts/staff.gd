extends Node2D

var selected_wand = 0
var is_animation_playing = false
var starting_pos

@onready var game_node = get_parent().get_parent()
@onready var staff_position_top = get_parent().get_node("StaffPositionTop")
@onready var staff_position_bottom = get_parent().get_node("StaffPositionBottom")
@onready var staff_position_left = get_parent().get_node("StaffPositionLeft")
@onready var staff_position_right = get_parent().get_node("StaffPositionRight")
@onready var lightning_animation_right = $"SpellTop/WandAnimations/LightningStrikeTop"
@onready var lightning_animation_left = $"SpellTop/WandAnimations/LightningStrikeTop"
@onready var lightning_animation_top = $"SpellTop/WandAnimations/LightningStrikeTop"
@onready var lightning_animation_bottom = $"SpellTop/WandAnimations/LightningStrikeTop"

var tile_size = 64
var inputs = {
	"right": Vector2.RIGHT,
	"left": Vector2.LEFT,
	"up": Vector2.UP,
	"down": Vector2.DOWN
}

func _ready():
	starting_pos = self.position
	self.top_level = true

	set_process_input(true)

	lightning_animation_right.connect("animation_finished", _on_animation_finished)
	lightning_animation_left.connect("animation_finished", _on_animation_finished)
	lightning_animation_top.connect("animation_finished", _on_animation_finished)
	lightning_animation_bottom.connect("animation_finished", _on_animation_finished)

func _process(delta):
	var hovered_wand = "Staff" + str(selected_wand)
	if Input.is_action_just_pressed("select_wand1"):
		selected_wand = 1
		if name == hovered_wand:
			print("selected wand" + str(selected_wand))
	elif Input.is_action_just_pressed("select_wand2"):
		selected_wand = 2
		if name == hovered_wand:
			print("selected wand" + str(selected_wand))
	elif Input.is_action_just_pressed("select_wand3"):
		selected_wand = 3
		if name == hovered_wand:
			print("selected wand" + str(selected_wand))
	elif Input.is_action_just_pressed("select_wand4"):
		selected_wand = 4
		if name == hovered_wand:
			print("selected wand" + str(selected_wand))

	if name == hovered_wand:
		$Icon.visible = true
	else:
		$Icon.visible = false

	if Input.is_action_just_pressed("wand_up"):
		cast_wand(hovered_wand, "up")
	elif Input.is_action_just_pressed("wand_down"):
		cast_wand(hovered_wand, "down")
	elif Input.is_action_just_pressed("wand_left"):
		cast_wand(hovered_wand, "left")
	elif Input.is_action_just_pressed("wand_right"):
		cast_wand(hovered_wand, "right")

func cast_wand(hovered_wand_selected, direction):
	if selected_wand == 0 or is_animation_playing:
		return
	var hovered_wand = "Staff" + str(selected_wand)
	if name == hovered_wand:
		print("casting wand" + str(selected_wand))

	if name == hovered_wand_selected:
		var gameID = await game_node.get_gameID_for_child()
		var resp = await game_node.client.rpc_async(game_node.session, "tx/game/player-turn", JSON.stringify({
			"GameIDStr": str(gameID),
			"Action": "wand",
			"Direction": direction,
			"WandNum": str(selected_wand - 1),
		}))

		var raycast = get_parent().get_node("RayCast2DMagic")
		raycast.target_position = inputs[direction] * tile_size
		raycast.force_raycast_update()

		if resp != null:
			match direction:
				"up":
					position = get_parent().position + Vector2(16, 0)
					lightning_animation_top.play("play")
				"down":
					position = get_parent().position + Vector2(16, 32)
					rotation += PI
					lightning_animation_bottom.play("play")
				"left":
					position = get_parent().position + Vector2(0, 16)
					rotation -= PI / 2
					lightning_animation_left.play("play")
				"right":
					position = get_parent().position + Vector2(32, 16)
					rotation += PI / 2
					lightning_animation_bottom.position = Vector2(20, -23)
					lightning_animation_right.play("play")

			is_animation_playing = true

func _on_animation_finished():
	is_animation_playing = false
	var casted_wand = "Staff" + str(selected_wand)
	if name == casted_wand:
		queue_free()  # Delete the staff node after the animation finishes
