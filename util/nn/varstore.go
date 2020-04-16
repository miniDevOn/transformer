package nn

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"

	G "gorgonia.org/gorgonia"
)

// The separator is used to separate path elements in the tensor names.
const SEP string = "."

// Variables holds variable pointer.
// When variable store is frozen, trainable is still set to tree.
// However, the tensor is not set to require gradients.
type Variables struct {
	NamedVariables     map[string]*G.Node
	TrainableVariables []*G.Node
	Mut                *sync.Mutex
}

/// VarStore is used to store variables used by one
// or multiple layers. It specifies a single device
// where all variables are stored.
type VarStore struct {
	// Variables Variables
	// device    Device
	G.ExprGraph
}

// Path is a variable store with an associated path
// for variable naming.
type Path struct {
	Path     []string
	VarStore VarStore
}

// Entry holds an entry corresponding to a give name in `Path`
type Entry struct {
	Name      string
	Variables Variables
	Path      Path
}

// Implement methods for `VarStore`
// ================================

// NewVarStore create a new variable store on the specified device
func NewVarStore(device Device) VarStore {
	variables := Variables{
		NamedVariables:     make(map[string]*G.Node),
		TrainableVariables: []*G.Node{},
		Mut:                &sync.Mutex{},
	}

	return VarStore{
		Variables: variables,
		device:    device,
	}
}

// Device get the device for this varstore
func (vs *VarStore) Device() Device {
	return vs.device
}

// Len returns number of tensors currently stored on this varstore
func (vs *VarStore) Len() uint {
	vs.Variables.Mut.Lock()
	variables := len(vs.Variables.NamedVariables)
	vs.Variables.Mut.Unlock()

	return uint(variables)
}

// IsEmpty checks whether there is no tensors are currently stored
// on this varstore
func (vs *VarStore) IsEmpty() bool {
	vs.Variables.Mut.Lock()
	defer vs.Variables.Mut.Unlock()

	if len(vs.Variables.NamedVariables) == 0 {
		return true
	}

	return false
}

// TrainableVariables returns all trainable variables for this varstore
func (vs *VarStore) TrainableVariables() []*G.Node {
	vs.Variables.Mut.Lock()
	defer vs.Variables.Mut.Unlock()

	variables := vs.Variables.TrainableVariables

	return variables
}

// GetVariables returns all variables currently stores in varstore
// along with their names.
func (vs *VarStore) GetVariables() map[string]*G.Node {
	vs.Variables.Mut.Lock()
	defer vs.Variables.Mut.Unlock()

	variables := vs.Variables.NamedVariables

	return variables
}

// Root returns the root path for this varstore
// Variables are named and organized using paths. This function
// returns the top level path for the varstore and can be combined
// with "/" to create sub-paths.
func (vs *VarStore) Root() Path {
	return Path{
		Path:     nil,
		VarStore: *vs,
	}
}

// Save saves the `varstore` variable values to a file
// Weight values for all the tensors currently stored in the `varstore`
// will be saved to a file.
func (vs *VarStore) Save(path string) error {
	vs.Variables.Mut.Lock()
	defer vs.Variables.Mut.Unlock()

	variables := vs.Variables.NamedVariables

	var namedTensors []*G.Node

	for _, t := range variables {
		namedTensors = append(namedTensors, t)
	}

	// TODO: save nameTensors to file
	return nil
}

// Load loads the `varstore` variable values from a file.
// weight values for all the tensors currently stored in
// the `varstore` gets loaded from the given file. The set
// of variables stored in the `varstore` is not changed, only
// the values for these tensors are modified.
func (vs *VarStore) Load(path string) error {
	// TODO: load multi with device
	/* let named_tensors = Tensor::load_multi_with_device(&path, self.device)?;
	 * let named_tensors: HashMap<_, _> = named_tensors.into_iter().collect();
	 * let mut variables = self.variables_.lock().unwrap();
	 * for (name, var) in variables.named_variables.iter_mut() {
	 *     match named_tensors.get(name) {
	 *         Some(src) => {
	 *             crate::no_grad(|| var.f_copy_(src).map_err(|e| format_err!("{}: {}", name, e)))?
	 *         }
	 *         None => return Err(format_err!("cannot find {} in {:?}", name, path.as_ref())),
	 *     }
	 * }
	 * Ok(()) */

	return nil
}

// LoadPartial loads the `varstore` variable values from a file if it exists.
// Weight values for the tensors currently stored in the `varstore` and the given
// file get loaded from the given file. If a variable in the var store is not present
// in the given file, it is skipped and its values are not updated. This method should
// be used if pre-trained weight for only parts of the model are available.
// The set of variables stored in the `varstore` is not changed, only the values
// for these tensors are modified.
func (vs *VarStore) LoadPartial(path string) error {
	// TODO: implement
	/* let named_tensors = Tensor::load_multi_with_device(&path, self.device)?;
	 * let named_tensors: HashMap<_, _> = named_tensors.into_iter().collect();
	 * let mut variables = self.variables_.lock().unwrap();
	 * let mut missing_variables = Vec::new();
	 * for (name, var) in variables.named_variables.iter_mut() {
	 *     match named_tensors.get(name) {
	 *         Some(src) => {
	 *             crate::no_grad(|| var.f_copy_(src).map_err(|e| format_err!("{}: {}", name, e)))?
	 *         }
	 *         None => {
	 *             missing_variables.push(name.to_owned());
	 *         }
	 *     }
	 * }
	 * Ok(missing_variables) */

	return nil
}

// Freeze freezes a `varstore`
// Gradients for the variables in this store are not tracked anymore.
func (vs *VarStore) Freeze() {
	vs.Variables.Mut.Lock()
	defer vs.Variables.Mut.Unlock()

	variables := vs.Variables.NamedVariables

	for _, v := range variables {
		// TODO:
		fmt.Println(v)
		// v.SetRequiresGrad(false)
	}
}

// Unfreeze unfreezes a `varstore`
// Gradients for the variables in this store are tracked again.
func (vs *VarStore) Unfreeze() {
	vs.Variables.Mut.Lock()
	defer vs.Variables.Mut.Unlock()

	variables := vs.Variables.NamedVariables

	for _, v := range variables {
		// TODO:
		fmt.Println(v)
		// v.SetRequiresGrad(true)
	}
}

// Copy copies variable values from a source `varstore` to this `varstore`
// All the variables in this `varstore` have to exist with the same name
// in the source `varstore`, otherwise return an error.
func (vs *VarStore) Copy(src VarStore) error {
	vs.Variables.Mut.Lock()
	defer vs.Variables.Mut.Unlock()
	// variables := vs.Variables

	src.Variables.Mut.Lock()
	defer src.Variables.Mut.Unlock()
	// srcVariables := src.Variables

	/*   device := vs.device
	 *
	 *   for name, _ := range variables.NamedVariables {
	 *     if ok, _ := srcVariables.NamedVariables[name]; !ok {
	 *       err := fmt.Errorf("cannot find %v in the source var store", name)
	 *       continue
	 *     }
	 *
	 *     srcVar = srcVariables.NamedVariables[name]
	 *     // TODO: copy to device
	 *     // crate::no_grad(|| var.f_copy_(&src_var.to_device(device)))?;
	 *
	 *   } */

	return nil

}

// Implement methods for `Path` struct
// ===================================

// Sub gets a sub-path of the given path
func (p *Path) Sub(s string) Path {
	if strings.Contains(s, SEP) {
		log.Fatalf("sub name cannot contain %v %v", SEP, s)
	}

	path := p.Path
	path = append(path, s)

	return Path{
		Path:     path,
		VarStore: p.VarStore,
	}
}

// Device gets the device where the `varstore` variables are stored
func (p *Path) Device() Device {
	return p.VarStore.device
}

func (p *Path) path(name string) string {
	if strings.Contains(name, SEP) {
		log.Fatalf("variable name cannot contain %v %v", SEP, name)
	}

	if p.VarStore.IsEmpty() {
		return name
	}

	switch {
	case len(p.Path) == 0:
		return name
	case len(p.Path) == 1:
		return p.Path[0]
	case len(p.Path) > 1:
		return fmt.Sprintf("%v%v%v", strings.Join(p.Path, SEP), SEP, name)
	}

	return fmt.Sprintf("%v%v%v", strings.Join(p.Path, SEP), SEP, name)

}

func (p *Path) add(name string, tensor *G.Node, trainable bool) *G.Node {
	path := p.path(name)

	p.VarStore.Variables.Mut.Lock()
	defer p.VarStore.Variables.Mut.Unlock()

	variables := p.VarStore.Variables
	if _, ok := variables.NamedVariables[path]; ok {
		path = fmt.Sprintf("%v__%v", path, len(variables.NamedVariables))
	}

	if trainable {
		// TODO: Turn on `gradient tracking`
		// tensor.set_requires_grad(true)
		p.VarStore.Variables.TrainableVariables = append(p.VarStore.Variables.TrainableVariables, tensor)
	}

	variables.NamedVariables[path] = tensor

	return tensor

}

func (p *Path) getOrAddWithLock(name string, tensor *G.Node, trainable bool, variables Variables) *G.Node {
	variables.Mut.Lock()
	variables.Mut.Unlock()

	path := p.path(name)
	if t, ok := variables.NamedVariables[path]; ok {
		return t
	}

	if trainable {
		// TODO: Set qruires grad
		// tensor.set_requires_grad(true)
		variables.TrainableVariables = append(variables.TrainableVariables, tensor)
		variables.NamedVariables[path] = tensor
	}

	return tensor
}

/*     fn get_or_add_with_lock(
 *         &self,
 *         name: &str,
 *         tensor: Tensor,
 *         trainable: bool,
 *         mut variables: MutexGuard<Variables>,
 *     ) -> Tensor {
 *         let path = self.path(name);
 *         if let Some(var) = variables.named_variables.get(&path) {
 *             return var.shallow_clone();
 *         }
 *
 *         let tensor = if trainable {
 *             tensor.set_requires_grad(true)
 *         } else {
 *             tensor
 *         };
 *         if trainable {
 *             variables.trainable_variables.push(tensor.shallow_clone());
 *         }
 *         variables
 *             .named_variables
 *             .insert(path, tensor.shallow_clone());
 *         tensor
 *     }
 *  */

// ZerosNoTrain creates a new variable initialized with zeroes.
// The new variable is named according to the name parameter
// and has the specified shape. The variable will not be trainable
// so gradients will not be tracked.
func (p *Path) ZerosNoTrain(name string, dims []int) *G.Node {

	// TODO: implement tensor
	z := ts.New()

	return p.add(name, z, false)
}

// OnesNoTrain create a new variable initialized with ones.
// The new variable is named according to the name parameter
// and has the specified shape. The variable will not be trainable
// so gradients will not be tracked.
// The variable uses a float tensor initialized with ones.
func (p *Path) OnesNoTrain(name string, dims []int) *G.Node {
	dt := ts.Float64
	o := ts.Ones(dt, dims...)
	return p.add(name, o, false)
}

// Var creates a new variable
// The new variable is named according to the name parameter
// and has the specified shape. The variable is trainable, its
// gradient will be tracked. The variable uses a float tensor
// initialized as per the related argument.
func (p *Path) Var(name string, dims []int, init InitT) *G.Node {
	v := Init(init, dims, p.VarStore.Device())

	return p.add(name, v, true)
}

// Zeros creates a new variables initialized with zeros
// The new variables is named after the name parameter
// and has the specified shape. The variable is trainable
// its gradient will be tracked.
// The variable uses a float tensor initialized with zeros
func (p *Path) Zeros(name string, dims []int) *G.Node {
	// TODO: check to make sure tensor of zero values.
	return p.Var(name, dims, InitFloat64(0.0))
}

// Ones creates a new variable initialized with ones.
// The new variable is named after the name parameter
// and has the specified shape. The variable is trainable
// its gradient will be tracked.
// The variable uses a float tensor initialized with ones.
func (p *Path) Ones(name string, dims []int) *G.Node {

	// TODO: check to make sure tensor of 1 values.
	return p.Var(name, dims, 1.0)
}

// Randn creates a new variable initialized randomly with normal distribution
// The new variable is named after the name parameter and has
// the specified shape. The variable is trainable, its gradient
// will be tracked.
// The variable uses a float tensor initialized randomly using a
// STANDARD normal distribution.
func (p *Path) RandnStandard(name string, dims []int) *G.Node {
	init := rand.NormFloat64()
	return p.Var(name, dims, init)
}

// Randn create a new variable initialed randomly with normal distribution.
// The new variable is named after the name parameter and has
// the specified shape. The variable is trainable, its gradient
// will be tracked.
// The variable uses a float tensor initialized randomly using a
// normal distribution with the specified mean and standard deviation.
func (p *Path) Randn(name string, dims []int, mean, stdev float64) *G.Node {
	init := rand.NormFloat64()*stdev + mean
	return p.Var(name, dims, init)
}

// Uniform creates a new variable initialized randomly with uniform distribution.
//
// The new variable is named according to the name parameter and
// has the specified shape. The variable is trainable, its gradient
// will be tracked.
// The variable uses a float tensor initialized randomly using a
// uniform distribution between the specified bounds.
func (p *Path) Uniform(name string, dims []int, lo, up float64) *G.Node {
	init := NewUniform(lo, up)
	return p.Var(name, dims, init)
}

// KaimingUniform creates a new variable initialized randomly with kaiming uniform.
//
// The new variable is named according to the name parameter and
// has the specified shape. The variable is trainable, its gradient
// will be tracked.
// The variable uses a float tensor initialized randomly using a
// uniform distribution which bounds follow Kaiming initialization.
func (p *Path) KaimingUniform(name string, dims []int) *G.Node {
	init := NewKaimingUniform(dims)

	return p.Var(name, dims, init)
}

// VarCopy creates a new variable initialized by copying an existing tensor.
//
// The new variable is named according to the name parameter and
// has the specified shape. The variable is trainable, its gradient
// will be tracked.
// The variable uses a float tensor initialized by copying some
// given tensor.
func (p *Path) VarCopy(name string, t *G.Node) *G.Node {
	v := p.Zeros(name, []int{t.Size()})
	return v
}

// Gets the tensor corresponding to a given name if present.
func (p *Path) Get(name string) (tensor *G.Node, ok bool) {
	p.VarStore.Variables.Mut.Lock()
	defer p.VarStore.Variables.Mut.Unlock()

	path := p.path(name)
	variables := p.VarStore.Variables

	tensor, ok = variables.NamedVariables[path]
	return tensor, ok
}

// Entry gets the entry corresponding to a given name for in-place manipulation.
func (p *Path) Entry(name string) Entry {
	p.VarStore.Variables.Mut.Lock()
	defer p.VarStore.Variables.Mut.Unlock()

	variables := p.VarStore.Variables
	return Entry{
		name,
		variables,
		*p,
	}
}

// Implement methods for `Entry`
// ================================

// OrVar returns the existing entry if, otherwise create a new variable.
//
// If this entry name matches the name of a variables stored in the
// var store, the corresponding tensor is returned. Otherwise a new
// variable is added to the var-store with the entry name and is
// initialized according to the init parameter.
func (e *Entry) OrVar(dims []int, init InitT) *G.Node {
	v := Init(init, dims, e.Path.Device())
	path := e.Path

	return path.getOrAddWithLock(e.Name, v, true, e.Variables)
}

// OrVarCopy returns the existing entry if, otherwise create a new variable.
func (e *Entry) OrVarCopy(tensor *G.Node) *G.Node {
	v := e.OrZeros([]int{tensor.Size()})
	return v
}

// Returns the existing entry if, otherwise create a new variable.
func (e *Entry) OrZeros(dims []int) *G.Node {
	return e.OrVar(dims, 0.0)
}

// OrKaimingUniform returns the existing entry if, otherwise create a new variable.
func (e *Entry) OrKaimingUniform(dims []int) *G.Node {
	return e.OrVar(dims, NewKaimingUniform(dims))
}

// OrOnes returns the existing entry if, otherwise create a new variable.
func (e *Entry) OrOnes(dims []int) *G.Node {
	return e.OrVar(dims, 1.0)
}

// OrOnesNoTrain returns the existing entry if, otherwise create a new variable.
func (e *Entry) OrOnesNoTrain(dims []int) *G.Node {

	shape := ts.WithShape(dims...)

	o := ts.New(shape)

	path := e.Path

	return path.getOrAddWithLock(e.Name, o, true, e.Variables)
}

// OrRandn returns the existing entry if, otherwise create a new variable.
func (e *Entry) OrRandn(dims []int, mean, stdev float64) *G.Node {
	init := rand.NormFloat64()*stdev + mean
	return e.OrVar(dims, init)
}

// OrRandnStandard returns the existing entry if, otherwise create a new variable.
func (e *Entry) OrRandnStandard(dims []int) *G.Node {
	init := rand.NormFloat64()
	return e.OrVar(dims, init)
}

// OrUniform returns the existing entry if, otherwise create a new variable.
func (e *Entry) OrUniform(dims []int, lo, up float64) *G.Node {
	init := NewUniform(lo, up)
	return e.OrVar(dims, init)
}

// OrZerosNoTrain returns the existing entry if, otherwise create a new variable.
func (e *Entry) OrZerosNoTrain(dims []int) *G.Node {

	shape := G.WithShape(dims...)
	z := G.NewTensor(e.Graph, G.Float64, shape)

	return e.Path.getOrAddWithLock(e.Name, z, true, e.Variables)
}
