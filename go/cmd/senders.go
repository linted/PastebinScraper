package main

var senderFuncs = []func(chan pasteMatch){
	postToSlack,
}

var registerFuncs = []func(){
	registerSlackFlags,
}

var validationFuncs = []func(){
	validateSlackFlags,
}

func registerSenderFlags() {
	for _, setupFunc := range registerFuncs {
		setupFunc()
	}
}

func validateSenderFlags() {
	for _, validation := range validationFuncs {
		validation()
	}
}

func startSending(messages chan pasteMatch) {
	var channels = make([]chan pasteMatch, len(senderFuncs))
	for i, senders := range senderFuncs {
		channels[i] = make(chan pasteMatch)
		go senders(channels[i])
	}
	teeChannel(messages, channels)

	for _, chans := range channels {
		close(chans)
	}
}

func teeChannel(messages chan pasteMatch, channels []chan pasteMatch) {
teeForever:
	for {
		select {
		case <-stopFlag:
			break teeForever // leave the for loop
		case msg := <-messages:
			for _, sendQueue := range channels {
				sendQueue <- msg
			}
		}
	}
}
