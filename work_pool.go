package matrix

var (
	f32Pool = newf32Pool()
	f64Pool = newf64Pool()
)

type f32Bucket struct {
	vals []float32
}

type matf32Pool struct {
	pool chan *f32Bucket
}

func newf32Pool() *matf32Pool {
	return &matf32Pool{
		pool: make(chan *f32Bucket, 10),
	}
}

func (p *matf32Pool) get() *f32Bucket {
	var c *f32Bucket
	select {
	case c = <-p.pool:
	default:
		c = &f32Bucket{
			vals: make([]float32, 0),
		}
	}
	return c
}

func (p *matf32Pool) put(m *f32Bucket) {
	select {
	case p.pool <- m:
	default:
		return
	}
}

type f64Bucket struct {
	vals []float64
}

type matf64Pool struct {
	pool chan *f64Bucket
}

func newf64Pool() *matf64Pool {
	return &matf64Pool{
		pool: make(chan *f64Bucket, 10),
	}
}

func (p *matf64Pool) get() *f64Bucket {
	var c *f64Bucket
	select {
	case c = <-p.pool:
	default:
		c = &f64Bucket{
			vals: make([]float64, 0),
		}
	}
	return c
}

func (p *matf64Pool) put(m *f64Bucket) {
	select {
	case p.pool <- m:
	default:
		return
	}
}
