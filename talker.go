// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

import (
	"math/rand"
	"strings"
	"time"
)

type Conversation struct {
	speaker string
	phrase  []string
}

var mouse = Conversation{
	speaker: "mouse",
	phrase: []string{
		"Ahem.\nHello.",
		"Do you have any chewing gum?",
		"Do you smell something?\nThere's supposed to be a pony.\nI haven't found a pony yet.",
		"The answer is 42, of course.",
		"Excuse me.",
		"They say it snows in the summer here sometimes.",
		"I think I've seen you before. With the cat.",
		"I'm ever so hungry.\nI wonder what's for dinner?",
		"Do you like poetry?",
		"Oh, gross.\n\nI'm pretty sure I stepped in something.",
		"sniff",
		"boo",
		"Pssst! Try /go home",
		"Cats make me nervous.",
	},
}

var cat = Conversation{
	speaker: "cat",
	phrase: []string{
		"Pfffttt!!!",
		"Zzzzzzz",
		"zzzzzzzzzzzzz",
		"Purrrrrrrr",
		"Meoooowwwww!",
	},
}

var conversation = []*Conversation{&cat, &mouse, &mouse, &cat, &mouse}

func RunTalker() {
	pause := 2
	priorText := ""
	for {
		seconds := rand.Intn(65)
		c := conversation[rand.Intn(len(conversation))]
		text := c.phrase[rand.Intn(len(c.phrase))]
		if priorText == text {
			continue
		}
		priorText = text
		speaker := c.speaker
		lines := strings.Split(text, "\n")
		time.Sleep(time.Duration(seconds) * time.Second)
		for _, line := range lines {
			if len(line) > 0 {
				MakeSmalltalk(line, speaker)
			}
			time.Sleep(time.Duration(pause) * time.Second)
			speaker = ""
		}
	}
}
