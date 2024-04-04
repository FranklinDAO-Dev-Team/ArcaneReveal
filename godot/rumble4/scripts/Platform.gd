extends StaticBody2D
enum Direction { TOP, BOTTOM, LEFT, RIGHT }
var current_direction
signal direction_changed(new_direction)

# Called when the node enters the scene tree for the first time.
func _ready():
	modulate = Color(Color.MEDIUM_PURPLE, 0.7)
	if name == "collisionTop":
		current_direction = Direction.TOP
	elif name == "collisionBottom":
		current_direction = Direction.BOTTOM
	elif name == "collisionLeft":
		current_direction = Direction.LEFT
	elif name == "collisionRight":
		current_direction = Direction.RIGHT
		
	emit_signal("direction_changed", current_direction)



# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
	if global.is_dragging:
		visible = true
	else:
		visible = false
