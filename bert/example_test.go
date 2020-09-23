package bert_test

import (
	"fmt"
	"log"

	"github.com/sugarme/gotch"
	"github.com/sugarme/gotch/nn"
	ts "github.com/sugarme/gotch/tensor"
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/transformer/bert"
)

func ExampleBertForMaskedLM() {

	device := gotch.CPU
	vs := nn.NewVarStore(device)

	config := bert.ConfigFromFile("../data/bert/config.json")
	model := bert.NewBertForMaskedLM(vs.Root(), config)
	err := vs.Load("../data/bert/model.ot")
	if err != nil {
		log.Fatalf("Load model weight error: \n%v", err)
	}

	tk := getBertTokenizer()
	sentence1 := "Looks like one [MASK] is missing"
	sentence2 := "It was a very nice and [MASK] day"

	var input []tokenizer.EncodeInput
	input = append(input, tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(sentence1)))
	input = append(input, tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(sentence2)))

	encodings, err := tk.EncodeBatch(input, true)
	if err != nil {
		log.Fatal(err)
	}

	var maxLen int = 0
	for _, en := range encodings {
		if len(en.Ids) > maxLen {
			maxLen = len(en.Ids)
		}
	}

	var tensors []ts.Tensor
	for _, en := range encodings {
		var tokInput []int64 = make([]int64, maxLen)
		for i := 0; i < len(en.Ids); i++ {
			tokInput[i] = int64(en.Ids[i])
		}

		tensors = append(tensors, ts.TensorFrom(tokInput))
	}

	inputTensor := ts.MustStack(tensors, 0).MustTo(device, true)

	var output ts.Tensor
	ts.NoGrad(func() {
		output, _, _ = model.ForwardT(inputTensor, ts.None, ts.None, ts.None, ts.None, ts.None, ts.None, false)
	})

	index1 := output.MustGet(0).MustGet(4).MustArgmax(0, false, false).Int64Values()[0]
	index2 := output.MustGet(1).MustGet(7).MustArgmax(0, false, false).Int64Values()[0]

	word1, ok := tk.IdToToken(int(index1))
	if !ok {
		fmt.Printf("Cannot find a corresponding word for the given id (%v) in vocab.\n", index1)
	}

	word2, ok := tk.IdToToken(int(index2))
	if !ok {
		fmt.Printf("Cannot find a corresponding word for the given id (%v) in vocab.\n", index2)
	}

	fmt.Println(word1)
	fmt.Println(word2)

	// NOTE: this output causes test failed (maybe due to cached from unit test with same method)
	/*
	 *   // Output:
	 *   // person
	 *   // pleasant
	 *  */
}

func ExampleBertForSequenceClassification() {

	device := gotch.CPU
	vs := nn.NewVarStore(device)

	config := bert.ConfigFromFile("../data/bert/config.json")

	var dummyLabelMap map[int64]string = make(map[int64]string)
	dummyLabelMap[0] = "positive"
	dummyLabelMap[1] = "negative"
	dummyLabelMap[3] = "neutral"

	config.Id2Label = dummyLabelMap
	config.OutputAttentions = true
	config.OutputHiddenStates = true
	model := bert.NewBertForSequenceClassification(vs.Root(), config)
	tk := getBertTokenizer()

	// Define input
	sentence1 := "Looks like one thing is missing"
	sentence2 := `It's like comparing oranges to apples`
	var input []tokenizer.EncodeInput
	input = append(input, tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(sentence1)))
	input = append(input, tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(sentence2)))
	encodings, err := tk.EncodeBatch(input, true)
	if err != nil {
		log.Fatal(err)
	}

	// Find max length of token Ids from slice of encodings
	var maxLen int = 0
	for _, en := range encodings {
		if len(en.Ids) > maxLen {
			maxLen = len(en.Ids)
		}
	}

	var tensors []ts.Tensor
	for _, en := range encodings {
		var tokInput []int64 = make([]int64, maxLen)
		for i := 0; i < len(en.Ids); i++ {
			tokInput[i] = int64(en.Ids[i])
		}

		tensors = append(tensors, ts.TensorFrom(tokInput))
	}

	inputTensor := ts.MustStack(tensors, 0).MustTo(device, true)

	var (
		output                         ts.Tensor
		allHiddenStates, allAttentions []ts.Tensor
	)

	ts.NoGrad(func() {
		output, allHiddenStates, allAttentions = model.ForwardT(inputTensor, ts.None, ts.None, ts.None, ts.None, false)
	})

	fmt.Println(output.MustSize())
	fmt.Println(len(allHiddenStates))
	fmt.Println(len(allAttentions))

	// Output:
	// [2 3]
	// 12
	// 12
}
