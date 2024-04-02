extends CharacterBody2D

var input_movement = Vector2.ZERO
var speed = 70

func _physics_process(delta):
	move()
	
func move():
	input_movement = Input.get_vector("left", "right", "up", "down")
	
	if input_movement != Vector2.ZERO:
		velocity = input_movement * speed
	
	if input_movement == Vector2.ZERO:
		velocity = Vector2.ZERO
		
	move_and_slide()
