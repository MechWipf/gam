package actor

//Spawn an actor with an auto generated id
func Spawn(props Props) *PID {
	id := ProcessRegistry.getAutoId()
	pid := spawn(id, props, nil)
	return pid
}

//SpawnNamed spawns a named actor
func SpawnNamed(props Props, name string) *PID {
	pid := spawn(name, props, nil)
	return pid
}

func spawn(id string, props Props, parent *PID) *PID {
	if props.RouterConfig() != nil {
		return spawnRouter(props.RouterConfig(), props, parent)
	}

	cell := NewActorCell(props, parent)
	mailbox := props.ProduceMailbox()
	ref := NewLocalActorRef(mailbox)
	pid, new := ProcessRegistry.registerPID(ref, id)

	if new {
		mailbox.RegisterHandlers(cell.invokeUserMessage, cell.invokeSystemMessage)
		cell.self = pid
		cell.invokeUserMessage(UserMessage{Message: &Started{}})
	}

	return pid
}
