package bert

import (
	"errors"

	"github.com/sugarme/gotch/nn"
	ts "github.com/sugarme/gotch/tensor"

	"github.com/sugarme/transformer/common"
)

// BertConfig defines the BERT model architecture (i.e., number of layers,
// hidden layer size, label mapping...)
type BertConfig struct {
	HiddenAct                 common.ActivationFn `json:"hidden_act"`
	AttentionProbsDropoutProb float64             `json:"attention_probs_dropout_prob"`
	HiddenDropoutProb         float64             `json:"hidden_dropout_prob"`
	HiddenSize                int64               `json:"hidden_size"`
	InitializerRange          float32             `json:"initializer_range"`
	IntermediateSize          int64               `json:"intermediate_size"`
	MaxPositionEmbeddings     int64               `json:"max_position_embeddings"`
	NumAttentionHeads         int64               `json:"num_attention_heads"`
	NumHiddenLayers           int64               `json:"num_hidden_layers"`
	TypeVocabSize             int64               `json:"type_vocab_size"`
	VocabSize                 int64               `json:"vocab_size"`
	OutputAttentions          bool                `json:"output_attentions"`
	OutputHiddenStates        bool                `json:"output_hidden_states"`
	IsDecoder                 bool                `json:"is_decoder"`
	Id2Label                  map[int64]string    `json:"id_2_label"`
	Label2Id                  map[string]int64    `json:"label_2_id"`
	NumLabels                 int64               `json:"num_labels"`
}

// BertModel defines base architecture for BERT models.
// `Task-specific` models can be built from this base model.
// `Embeddings`: for `token`, `position` and `segment` embeddings
// `Encoder`: is a vector of layers. Each layer compose of a `self-attention`,
// an `intermedate` (linear) and an output ( linear + layer norm) sub-layers.
// `Pooler`: linear layer applied to the first element of the sequence (`[MASK]` token)
// `IsDecoder`: whether model is used as a decoder. If set to `true`
// a casual mask will be applied to hide future positions that should be attended to.
type BertModel struct {
	Embeddings *BertEmbeddings
	Encoder    *BertEncoder
	Pooler     *BertPooler
	IsDecoder  bool
}

// NewBertModel builds a new `BertModel`
// * `p` Variable store path for the root of the BERT Model
// * `config` `BertConfig` configuration for model architecture and decoder status
// Example: TODO - create example
// let config_path = Path::new("path/to/config.json");
// let device = Device::Cpu;
// let p = nn::VarStore::new(device);
// let config = BertConfig::from_file(config_path);
// let bert: BertModel<BertEmbeddings> = BertModel::new(&(&p.root() / "bert"), &config);
func NewBertModel(p nn.Path, config *BertConfig) *BertModel {
	isDecoder := false
	if config.IsDecoder {
		isDecoder = true
	}

	embeddings := NewBertEmbedding(p.Sub("embeddings"), config)

	encoder := NewBertEncoder(p.Sub("encoder"), config)

	pooler := NewBertPooler(p.Sub("pooler"), config)
	bertModel := BertModel{embeddings, encoder, pooler, isDecoder}

	return &bertModel
}

func (b *BertModel) ForwardT(inputIds, mask, tokenTypeIds, positionIds, inputEmbeds, encoderHiddenStates, encoderMask *G.Node, train bool) (*G.Node, *G.Node, *G.Node, []*G.Node, string, error) {

	var (
		inputShape ts.Shape
		device     G.Device
		err        error
	)

	if inputIds != nil {
		if inputEmbeds != nil {
			err = errors.New("Only one of input ids or input embeddings may be set")
			return nil, nil, nil, nil, "", err
		} else {
			inputShape = inputIds.Shape()
			device = inputIds.Device()
		}
	} else {
		if inputEmbeds == nil {
			err = errors.New("At least one of input ids or input embeddings must be set")
			return nil, nil, nil, nil, "", err
		} else {
			// Check this.
			// Some(embeds) => (vec!(embeds.size()[0], embeds.size()[1]), embeds.device()),
			inputShape, err = inputEmbeds.Shape().S([]int{0, 1})
			device = inputEmbeds.Device()
		}

		if mask == nil {
			// TODO: create new node
			// None => Tensor::ones(&input_shape, (Kind::Int64, device))
			g := G.NewGraph()
			mask = G.NewTensor(g, G.Float64, inputShape)
		}
	}

}
