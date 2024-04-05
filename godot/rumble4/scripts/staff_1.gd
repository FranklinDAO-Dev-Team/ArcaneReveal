extends Node2D

var draggable = false
var is_inside_dropable = false
var body_ref
var offset: Vector2
var initial_pos: Vector2
var is_animation_playing = false

@onready var staff_position_top = get_parent().get_node("StaffPositionTop")
@onready var staff_position_bottom = get_parent().get_node("StaffPositionBottom")
@onready var staff_position_left = get_parent().get_node("StaffPositionLeft")
@onready var staff_position_right = get_parent().get_node("StaffPositionRight")
@onready var lightning_animation_right = $"SpellRight/WandAnimations/LightningStrikeRight"
@onready var lightning_animation_left = $"SpellLeft/WandAnimations/LightningStrikeLeft"
@onready var lightning_animation_top = $"SpellTop/WandAnimations/LightningStrikeTop"
@onready var lightning_animation_bottom = $"SpellBottom/WandAnimations/LightningStrikeBottom"

func _ready():
	set_process_input(true)
	var tile_pieces = get_parent().get_parent().get_node("TilePiece").get_children()
	for tile_piece in tile_pieces:
		tile_piece.connect("direction_changed", _on_direction_changed)
		
	lightning_animation_right.connect("animation_finished", _on_animation_finished)
	lightning_animation_left.connect("animation_finished", _on_animation_finished)
	lightning_animation_top.connect("animation_finished", _on_animation_finished)
	lightning_animation_bottom.connect("animation_finished", _on_animation_finished)
		
func _on_direction_changed(new_direction):
	if is_inside_dropable:
		body_ref.current_direction = new_direction

func _process(delta):
	if draggable:
		if Input.is_action_just_pressed("click"):
			offset = get_global_mouse_position() - global_position
			global.is_dragging = true

		if Input.is_action_pressed("click"):
			global_position = get_global_mouse_position() - offset

		elif Input.is_action_just_released("click"):
			global.is_dragging = false
			
			if is_inside_dropable:
				match body_ref.current_direction:
					body_ref.Direction.TOP:
						position = staff_position_top.position  # Assign the position of the StaffPositionTop node
						lightning_animation_top.position = Vector2(20, -23)
						lightning_animation_top.play("play")
						is_animation_playing = true
					body_ref.Direction.BOTTOM:
						position = staff_position_bottom.position  # Assign the position of the StaffPositionBottom node
						lightning_animation_bottom.position = Vector2(20, -23)
						lightning_animation_bottom.play("play")
						is_animation_playing = true
					body_ref.Direction.LEFT:
						position = staff_position_left.position  # Assign the position of the StaffPositionLeft node
						lightning_animation_left.position = Vector2(20, -23)
						lightning_animation_left.play("play")
						is_animation_playing = true
					body_ref.Direction.RIGHT:
						position = staff_position_right.position  # Assign the position of the StaffPositionRight node
						lightning_animation_right.position = Vector2(20, -23)
						lightning_animation_right.play("play")
						is_animation_playing = true

func _on_animation_finished():
	is_animation_playing = false
	queue_free()  # or visible = false, or modulate.a = 0

func _on_area_2d_mouse_entered(): 
	if not global.is_dragging:
		draggable = true
		scale = Vector2(1.25, 1.25)


func _on_area_2d_mouse_exited():
	if not global.is_dragging:
		draggable = false
		scale = Vector2(1, 1)


func _on_area_2d_body_entered(body):
	if body.is_in_group('dropable'):
		is_inside_dropable = true
		body.modulate = Color(Color.REBECCA_PURPLE, 1)
		body_ref = body


func _on_area_2d_body_exited(body):
	if body.is_in_group('dropable'):
		is_inside_dropable = false
		body.modulate = Color(Color.MEDIUM_PURPLE, 0.7)
