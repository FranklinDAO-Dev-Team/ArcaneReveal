extends Node

var client : NakamaClient

func _ready():
	client = Nakama.create_client("defaultkey", "127.0.0.1", 7350, "http")

	# Get the System's unique device identifier
	var device_id = OS.get_unique_id()

	# Authenticate with the Nakama server using Device Authentication
	var session : NakamaSession = await client.authenticate_device_async(device_id)
	if session.is_exception():
		print("An error occurred: %s" % session)
		return
	print("Successfully authenticated: %s" % session)
