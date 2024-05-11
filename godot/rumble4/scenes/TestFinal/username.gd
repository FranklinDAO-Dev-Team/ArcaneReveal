extends Control

signal username_submitted(username)

var username = ""

func _on_button_pressed():
	username = $TextEdit.text
	print(username)
	emit_signal("username_submitted", username)
	queue_free()
