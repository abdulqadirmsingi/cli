package qa

import (
	"strings"
	"testing"

	"github.com/abdulqadirmsingi/pulse-cli/internal/db"
)

func TestSimilarTodayQuestionsResolveToTodayAnswer(t *testing.T) {
	store := fakeStore{
		today: &db.Stats{TotalCommands: 4, NoiseCommands: 1, TotalTimeMS: 120_000, SuccessRate: 100},
		todayProjects: []db.TopEntry{
			{Name: "cli", Count: 3, MS: 120_000},
		},
		hourly: []db.HourlyBucket{{Hour: 8, Count: 4}},
	}

	for _, question := range []string{
		"what did I work on today?",
		"what am I doing today",
		"what did i do today",
	} {
		answer, err := AnswerQuestion(store, question)
		if err != nil {
			t.Fatalf("AnswerQuestion(%q): %v", question, err)
		}
		if answer.Title != "today's pulse" {
			t.Fatalf("AnswerQuestion(%q) title = %q, want today's pulse", question, answer.Title)
		}
	}
}

func TestUnknownQuestionReturnsSuggestions(t *testing.T) {
	answer, err := AnswerQuestion(fakeStore{}, "can you predict the weather")
	if err != nil {
		t.Fatalf("AnswerQuestion: %v", err)
	}
	if answer.Title != "ask pulse" {
		t.Fatalf("title = %q, want ask pulse", answer.Title)
	}
	if !strings.Contains(strings.Join(answer.Lines, " "), "could not match") {
		t.Fatalf("fallback lines = %#v, want helpful mismatch message", answer.Lines)
	}
	if len(answer.Tips) == 0 {
		t.Fatalf("fallback should include suggestions")
	}
}

type fakeStore struct {
	stats         *db.Stats
	today         *db.Stats
	todayProjects []db.TopEntry
	topCommands   []db.TopEntry
	topProjects   []db.TopEntry
	hourly        []db.HourlyBucket
}

func (f fakeStore) GetStats(int) (*db.Stats, error) {
	if f.stats != nil {
		return f.stats, nil
	}
	return &db.Stats{}, nil
}

func (f fakeStore) GetTodayStats() (*db.Stats, error) {
	if f.today != nil {
		return f.today, nil
	}
	return &db.Stats{}, nil
}

func (f fakeStore) GetTodayTopProjects(int) ([]db.TopEntry, error) {
	return f.todayProjects, nil
}

func (f fakeStore) GetTopCommands(int, int) ([]db.TopEntry, error) {
	return f.topCommands, nil
}

func (f fakeStore) GetTopProjects(int, int) ([]db.TopEntry, error) {
	return f.topProjects, nil
}

func (f fakeStore) GetHourlyStats(string) ([]db.HourlyBucket, error) {
	return f.hourly, nil
}
