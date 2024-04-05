extends Area2D

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

func _ready():
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
	health -= 1
	if health < 0:
		health = max_health
	
	
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
		area.damage()
		
func _physics_process(delta):
	update_health()
	
func take_damage() -> void:
	health -= 1
	if health < 0:
		health = max_health
	#update_health_ui()
