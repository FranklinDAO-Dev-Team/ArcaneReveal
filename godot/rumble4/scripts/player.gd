extends CharacterBody2D

var currPos = [0, 0]
#
#
##func _physics_process(delta):
	##move()
	##
##func move():
	##input_movement = Input.get_vector("left", "right", "up", "down")
	##
	##if input_movement != Vector2.ZERO:
		##velocity = input_movement * speed
	##
	##if input_movement == Vector2.ZERO:
		##velocity = Vector2.ZERO
		##
	##move_and_slide()
#
func _input(event):
	if event.is_action_pressed("right"):
		if currPos[0] <= 96:
			currPos[0] += 32
	elif event.is_action_pressed("left"):
		if currPos[0] >= 32:
			currPos[0] -= 32	
	elif event.is_action_pressed("up"):
		if currPos[1] >= 32:
			currPos[1] -= 32
	elif event.is_action_pressed("down"):
		if currPos[1] <= 96:	
			currPos[1] += 32

	self.position = Vector2(currPos[0], currPos[1])


