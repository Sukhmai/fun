package questions

import "golang.org/x/exp/rand"

var questionSet1 = []string{
	"What would constitute a 'perfect' day for you?",
	"For what in your life do you feel most grateful?",
}

var questionSet2 = []string{
	"Is there something that you've dreamed of doing for a long time? Why haven't you done it?",
	"How do you feel about your relationship with your mother?",
}

var questionSet3 = []string{
	"If you were to die today with no opportunity to communicate with anyone, what would you most regret not having told someone? Why haven't you told them yet?",
	"What roles do love and affection play in your life?",
}

var questionSets = [][]string{
	questionSet1,
	questionSet2,
	questionSet3,
}

func GetRandomQuestions() []string {
	randomQuestions := make([]string, 3)
	for i, questionSet := range questionSets {
		randomQuestion := questionSet[rand.Intn(len(questionSet))]
		randomQuestions[i] = randomQuestion
	}
	return randomQuestions
}

func GetQuestions() []string {
	questions := make([]string, 6)
	for i, questionSet := range questionSets {
		for j, question := range questionSet {
			questions[2*i+j] = question
		}
	}
	return questions
}
