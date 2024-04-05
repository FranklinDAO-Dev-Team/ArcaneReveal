extends Node2D

var draggable = false
var is_inside_dropable = false
var body_ref
var offset: Vector2
var initial_pos: Vector2
var is_animation_playing = false
var starting_pos

@onready var staff_position_top = get_parent().get_node("StaffPositionTop")
@onready var staff_position_bottom = get_parent().get_node("StaffPositionBottom")
@onready var staff_position_left = get_parent().get_node("StaffPositionLeft")
@onready var staff_position_right = get_parent().get_node("StaffPositionRight")
@onready var lightning_animation_right = $"SpellTop/WandAnimations/LightningStrikeTop"
@onready var lightning_animation_left = $"SpellTop/WandAnimations/LightningStrikeTop"
@onready var lightning_animation_top = $"SpellTop/WandAnimations/LightningStrikeTop"
@onready var lightning_animation_bottom = $"SpellTop/WandAnimations/LightningStrikeTop"

func _ready():
	starting_pos = self.position
	self.top_level = true
	
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
						position = get_parent().position + Vector2(16, 0)
						lightning_animation_top.play("play")
						is_animation_playing = true
					body_ref.Direction.BOTTOM:
						position = get_parent().position + Vector2(16, 32)
						rotation += PI 
						lightning_animation_bottom.play("play")
						is_animation_playing = true
					body_ref.Direction.LEFT:
						position = get_parent().position + Vector2(0, 16)
						rotation -= PI / 2 
						lightning_animation_left.play("play")
						is_animation_playing = true
					body_ref.Direction.RIGHT:
						position = get_parent().position + Vector2(32, 16)
						rotation += PI / 2
						lightning_animation_bottom.position = Vector2(20, -23)
						lightning_animation_right.play("play")
						is_animation_playing = true
			else:
				self.position = starting_pos

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
