// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Small, random conversations which are injected into rooms
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Conversation struct {
	// Speaker is the actor given credit for saying things.  Phrases
	// which are segmented only give credit for speaking the the
	// first segment of the phrase. Subsequent segments are spoken
	// without credit since they will appear as continuations of the
	// initial segment.
	speaker string
	// Phrases spoken during a conversation. A phrase with embedded
	// newlines ("\n") will be split into segments and each segment will
	// be spoken separately with a brief pause between spoken segments.
	phrase []string
	// Unsaid controls the sequence of phrases spoken during a conversation.
	// Before a conversation is started, this should be properly initialized
	// (randomly) with the indicies of the strings in phrase. Elements are
	// removed from the beginning until the conversation is finished.
	unsaid []int
}

const (
	pauseBetweenSegments = 2
)

var (
	mouseSays = []string{
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
	}

	mouse = Conversation{
		speaker: "mouse",
		phrase:  mouseSays,
	}

	catSays = []string{
		"Pfffttt!!!",
		"Zzzzzzz",
		"zzzzzzzzzzzzz",
		"Purrrrrrrr",
		"Meoooowwwww!",
	}

	cat = Conversation{
		speaker: "cat",
		phrase:  catSays,
	}

	conversation = []*Conversation{&cat, &mouse, &mouse, &cat, &mouse}
)

func InjectConversations() {
	locus := "CONVERSATIONS"
	checkpoint(locus, "BEGIN")
	for {
		seconds := rand.Intn(config.maxSecondsBetweenConversations)
		time.Sleep(time.Duration(seconds) * time.Second)
		speaker, segments := findSomethingToSay()
		if config.debug {
			checkpoint(locus, "TIME-TO-SPEAK")
		}
		for _, line := range segments {
			if len(line) > 0 {
				MakeSmalltalk(line, speaker)
			}
			time.Sleep(time.Duration(pauseBetweenSegments) * time.Second)
			speaker = ""
		}
	}
}

func findSomethingToSay() (speaker string, lines []string) {
	c := conversation[rand.Intn(len(conversation))]
	if len(c.unsaid) < 1 {
		resetConversation(c)
	}
	speaker = c.speaker
	lines = strings.Split(c.phrase[c.unsaid[0]], "\n")
	c.unsaid = c.unsaid[1:]
	return
}

func resetConversation(c *Conversation) {
	locus := "CONVERSATION.RESET"
	checkpoint(locus, fmt.Sprintf("speaker=%s", c.speaker))
	n := len(c.phrase)
	taken := make([]bool, len(c.phrase))
	for {
		if n < 2 {
			break
		}
		i := rand.Intn(len(c.phrase))
		if taken[i] {
			continue
		}
		taken[i] = true
		c.unsaid = append(c.unsaid, i)
		n -= 1
	}
	// Find the final untaken phrase
	for i := 0; i < n; i += 1 {
		if !taken[i] {
			c.unsaid = append(c.unsaid, i)
			return
		}
	}
}
