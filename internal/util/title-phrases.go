// all phrases you see in the terminal are stored here :)

package util

import "math/rand"

var titlePhrasesSlice = []string{
	"I don't know how, but I will.",
	"One day at a time.",
	"Progress, not perfection.",
	"Start before you're ready.",
	"Done is better than perfect.",
	"Show up. Every day.",
	"Small steps. Big results.",
	"You've survived 100% of your worst days.",
	"Fall seven times, stand up eight.",
	"The only way out is through.",
	"Not yet is not never.",
	"Keep going.",
	"Discipline beats motivation.",
	"Do it scared.",
	"Hard now, easy later.",
	"Trust the process.",
	"Consistency compounds.",
	"Your future self is watching.",
	"It always seems impossible until it's done.",
	"Don't stop when you're tired. Stop when you're done.",
	"You are closer than you think.",
	"Every expert was once a beginner.",
	"The comeback is always stronger than the setback.",
	"Do something today your future self will thank you for.",
	"Earn it.",
	"Be the person you needed.",
	"Prove them wrong.",
	"Make it happen.",
	"Caffeine is my steroids.",
}

func GetRandTitlePhrase() string {
	return titlePhrasesSlice[rand.Intn(len(titlePhrasesSlice))]
}