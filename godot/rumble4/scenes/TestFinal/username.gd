extends Control

signal username_submitted(username)

var username = ""

func _on_button_pressed():
	# Get the text from the TextEdit node
	var text = $TextEdit.text
	
	# Sanitize the text by removing newline characters
	text = text.replace("\n", "").replace("\r", "")
	
	# Update the username variable
	username = text
	
	# Print and emit the signal with the sanitized username
	print(username)
	emit_signal("username_submitted", username)
	
	# Free the node
	queue_free()
