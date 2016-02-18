// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

import (
	"math/rand"
	"time"
)

type Conversation struct {
	speaker string
	phrase  []string
}

var ignatz = Conversation{
	speaker: "Ignatz",
	phrase:  []string{"Yo. Buddy", "Get outta here."},
}

var anonymous = Conversation{
	speaker: "",
	phrase: []string{
		"Ahem.  Hello.",
		"Do you have any chewing gum?",
		"I haven't found the pony yet.",
		"Do you smell that?",
		"The answer is 42, of course.",
		"Oh. Excuse me.",
		"They say it snows in the summer here sometimes.",
		"Wait. Don't I know you?",
		"I wonder what's for dinner?",
		"I'm ever so hungry",
		"Do you like poetry?",
		"Oh, gross. I'm pretty sure I stepped in something.",
		"sniff",
		"boo",
	},
}

var conversation = []*Conversation{&ignatz, &anonymous, &anonymous, &anonymous, &anonymous}

func RunTalker() {
	for {
		seconds := rand.Intn(65)
		c := conversation[rand.Intn(len(conversation))]
		text := c.phrase[rand.Intn(len(c.phrase))]
		speaker := c.speaker
		time.Sleep(time.Duration(seconds) * time.Second)
		MakeSmalltalk(text, speaker)
	}
}
