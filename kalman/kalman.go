package kalman

import (
	"fmt"
	"log"
	"math"

	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
)

var (
	errNotSquare = errors.New("kalman: non-square matrix")
	errMismatch  = errors.New("kalman: dimension mismatch")
)

// Filter provides a Kalman filter.
type Filter struct {
	c  *mat.Dense // Output matrix
	q  *mat.Dense // Process noise covariance
	r  *mat.Dense // Measurement noise covariance
	p  *mat.Dense // Estimate error covariance
	k  *mat.Dense
	p0 *mat.Dense

	m, n int // dimensions of the system

	t0, t float64 // initial and current time
	dt    float64 // discrete time step

	x, xn *mat.VecDense

	id *mat.Dense
}

func New(dt float64, c, q, r, p *mat.Dense) *Filter {
	m, n := c.Dims()

	return &Filter{
		c:  c,
		q:  q,
		r:  r,
		p0: p,
		m:  m,
		n:  n,
		dt: dt,
		id: newIdentity(n),
	}
}

func (kf *Filter) State() *mat.VecDense {
	return kf.x
}

func (kf *Filter) Init(t0 float64, x0 *mat.VecDense) {
	kf.x = x0
	kf.xn = mat.NewVecDense(kf.n, nil)
	kf.p = kf.p0
	kf.t0 = t0
	kf.t = t0
}

// Update updates the Kalman filter.
// a: System dynamics matrix
func (kf *Filter) Update(y *mat.VecDense, dt float64, a *mat.Dense) error {
	if kf.x == nil {
		panic("kalman: filter not initialized")
	}

	kf.dt = dt

	kf.xn.MulVec(a, kf.x)
	kf.p.Mul(a, kf.p)
	kf.p.Mul(kf.p, a.T())
	kf.p.Add(kf.p, kf.q)

	ct := kf.c.T()
	var k mat.Dense
	log.Printf(">>> c: %v", dimsOf(kf.c))
	log.Printf(">>> p: %v", dimsOf(kf.p))
	log.Printf(">>> ct:%v", dimsOf(ct))
	k.Mul(kf.c, kf.p)
	log.Printf(">>> k: %v", dimsOf(&k))
	k.Mul(&k, ct)
	k.Add(&k, kf.r)
	var kinv mat.Dense
	err := kinv.Inverse(&k)
	if err != nil {
		return err
	}

	k.Product(kf.p, ct, &k)
	var tmp mat.VecDense
	tmp.MulVec(kf.c, kf.xn)
	tmp.SubVec(y, &tmp)
	tmp.MulVec(&k, &tmp)

	kf.xn.AddVec(kf.xn, &tmp)

	var m mat.Dense
	m.Mul(&k, kf.c)
	m.Sub(kf.id, &m)

	kf.p.Mul(&m, kf.p)

	kf.t += dt

	return nil
}

func dimsOf(m mat.Matrix) string {
	r, c := m.Dims()
	return fmt.Sprintf("(%d, %d)", r, c)
}

func newIdentity(n int) *mat.Dense {
	if n <= 0 {
		panic("kalman: invalid matrix identity dimension")
	}
	data := make([]float64, n*n)
	for i := 0; i < n; i++ {
		data[i*(i+1)] = 1
	}
	return mat.NewDense(n, n, data)
}

type Settings struct {
	F *mat.Dense
	G *mat.Dense
	Q *mat.Dense

	H *mat.Dense
	R *mat.Dense
}

type KF struct {
	Settings

	id *mat.Dense // identity matrix

	x *mat.VecDense
	v *mat.Dense
}

func NewKF(settings *Settings) *KF {
	rf, cf := settings.F.Dims()
	rg, cg := settings.G.Dims()
	rq, cq := settings.Q.Dims()
	rh, ch := settings.H.Dims()
	rr, cr := settings.R.Dims()

	switch {
	case rf != cf, rq != cq, rr != cr:
		panic(errNotSquare)
	case rf != rg:
		panic(errMismatch)
	case cg != rq:
		panic(errMismatch)
	case ch != cf:
		panic(errMismatch)
	case rh != rr:
		panic(errMismatch)
	}

	return &KF{
		Settings: *settings,
		id:       newIdentity(rf),
		x:        mat.NewVecDense(cf, nil),
		v:        mat.NewDense(rf, cf, nil),
	}
}

func (kf *KF) Init(x *mat.VecDense, v mat.Matrix) {
	rf, cf := kf.Settings.F.Dims()
	switch x {
	case nil:
		kf.x = mat.NewVecDense(cf, nil)
	default:
		rx, _ := x.Dims()
		if rx != cf {
			panic(errMismatch)
		}
		kf.x = x
	}

	switch v {
	case nil:
		var m mat.Dense
		m.Mul(kf.G, kf.Q)
		kf.v.Mul(&m, kf.G.T())
	default:
		rv, cv := v.Dims()
		switch {
		case rv != cv:
			panic(errNotSquare)
		case rv != rf:
			panic(errMismatch)
		}
		kf.v = mat.DenseCopyOf(v)
	}
}

func (kf *KF) Filter(out, sys *mat.Dense) (*mat.Dense, error) {
	var (
		rh, _  = kf.H.Dims()
		rs, cs = sys.Dims()
	)

	if out == nil {
		out = mat.NewDense(rh, cs, nil)
	}

	var (
		m0, m1, m2     mat.Dense
		di0, di, k, kh mat.Dense
		ke, e          mat.VecDense
		retv           mat.VecDense
		arr            []float64

		ft = kf.F.T()
		gt = kf.G.T()
		ht = kf.H.T()
	)

	for j := 0; j < cs; j++ {
		// x = F . x
		kf.x.MulVec(kf.F, kf.x)

		// v = F . v . F^T + G . Q . G^T
		m0.Mul(kf.F, kf.v)
		kf.v.Mul(&m0, ft)

		m1.Mul(kf.G, kf.Q)
		m2.Mul(&m1, gt)
		kf.v.Add(kf.v, &m2)

		// d = (H . v . H^T + R)^{-1}
		di0.Mul(kf.H, kf.v)
		di.Mul(&di0, ht)
		di.Add(&di, kf.R)
		err := di.Inverse(&di) // FIXME(sbinet): not stable.
		if err != nil {
			return nil, err
		}

		// v . H^T . d^-1
		k.Mul(kf.v, ht)
		k.Mul(&k, &di)

		// e = y - H . x
		y := sys.ColView(j)
		e.MulVec(kf.H, kf.x)
		e.SubVec(y, &e)

		hasNaN := false
		for i := 0; i < rs; i++ {
			v := y.At(i, 0)
			if math.IsNaN(v) {
				sys.Set(i, j, 0)
				hasNaN = true
			}
		}

		if !hasNaN {
			// x = x + K . e
			ke.MulVec(&k, &e)
			kf.x.AddVec(kf.x, &ke)

			// v = (I - K . H) . v
			kh.Mul(&k, kf.H)
			kh.Sub(kf.id, &kh)
			kf.v.Mul(&kh, kf.v)
		}

		retv.MulVec(kf.H, kf.x)
		arr = mat.Col(arr, 0, &retv)
		out.SetCol(j, arr)
	}

	return out, nil
}
