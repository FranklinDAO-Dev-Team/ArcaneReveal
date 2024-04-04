extends Area2D

var animation_speed = 4
var moving = false
var tile_size = 64
var inputs = {
	"enemy_right": Vector2.RIGHT,
	"enemy_left": Vector2.LEFT,
	"enemy_up": Vector2.UP,
	"enemy_down": Vector2.DOWN
}

@onready var ray = $RayCast2D

func _ready():
	position = position.snapped(Vector2.ONE * tile_size)
	#position += Vector2.ONE * tile_size / 

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
		#position += inputs[dir] * tile_size
		var tween = get_tree().create_tween()
		tween.tween_property(self, "position", position + inputs[dir] * tile_size, 1.0/animation_speed).set_trans(Tween.TRANS_SINE)
		moving = true
		await tween.finished
		moving = false


func _on_area_entered(area):
	if area.name == "Player":
		area.damage_heavy()
