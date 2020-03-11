package bpe_test

import (
	"reflect"
	"testing"

	bpe "github.com/sugarme/sermo/tokenizer/model/bpe"
)

func TestBpeTrainer_Train(t *testing.T) {
	var wordCounts map[string]uint32 = make(map[string]uint32)
	wordCounts["roses"] = 1
	wordCounts["are"] = 2
	wordCounts["red"] = 1
	wordCounts["voilets"] = 1
	wordCounts["blue"] = 1
	wordCounts["BERT"] = 1
	wordCounts["is"] = 2
	wordCounts["big"] = 1
	wordCounts["and"] = 1
	wordCounts["so"] = 1
	wordCounts["GPT-2"] = 1

	trainer := bpe.NewBpeTrainer(2, 100)

	model, _ := trainer.Train(wordCounts)

	got := model.Vocab

	var want map[string]uint32 = make(map[string]uint32)
	want["-"] = 0
	want["2"] = 1
	want["B"] = 2
	want["E"] = 3
	want["G"] = 4
	want["P"] = 5
	want["R"] = 6
	want["T"] = 7
	want["a"] = 8
	want["b"] = 9
	want["d"] = 10
	want["e"] = 11
	want["g"] = 12
	want["i"] = 13
	want["l"] = 14
	want["n"] = 15
	want["o"] = 16
	want["r"] = 17
	want["s"] = 18
	want["t"] = 19
	want["u"] = 20
	want["v"] = 21
	want["re"] = 22
	want["are"] = 23
	want["is"] = 24

	if !reflect.DeepEqual(want, got) {

		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}

}
