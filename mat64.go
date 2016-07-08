/*
Package mat64 implements a "mat" object, which behaves like a 2D array
or list in other programming languages. Under the hood, the mat object is a
flat slice, which provides for optimal performance in Go, while the methods
and constructors provide for a higher level of performance and abstraction
when compared to the "2D" slices of go (slices of slices).

All errors encountered in this package, such as attempting to access an
element out of bounds are treated as critical error, and thus, the code
immediately exits with signal 1. In such cases, the function/method in
which the error was encountered is printed to the screen, in addition
to the full stack trace, in order to help fix the issue rapidly.
*/
package mat64

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

/*
Mat is the main struct of this library. Mat is a essentially a 1D slice
(a []float64) that contains two integers, representing rows and columns,
which allow it to behave as if it was a 2D slice. This allows for higher
performance and flexibility for the users of this library, at the expense
of some bookkeeping that is done here.

The fields of this struct are not directly accessible, and they may only
change by the use of the various methods in this library.
*/
type Mat struct {
	r, c int
	vals []float64
}

/*
New is the primary constructor for the "Mat" object. New is a variadic function,
expecting 0 to 3 integers, with differing behavior as follows:

	m := New()

m is now an empty &Mat{}, where the number of rows,
columns and the length and capacity of the underlying
slice are all zero. This is mostly for internal use.

	m := New(x)

m is a x by x (square) matrix, with the underlying
slice of length x, and capacity 2x.

	m := New(x, y)

m is an x by y matrix, with the underlying slice of
length xy, and capacity of 2xy. This is a good case
for when your matrix is going to expand in the
future. There is a negligible hit to performance
and a larger memory usage of your code. But in case
expanding matrices, many re-allocations are avoided.

	m := New(x, y, z)

m is a x by u matrix, with the underlying slice of
length xy, and capacity z. This is a good choice for
when the size of the matrix is static, or when the
application is memory constrained.

For most cases, we recommend using the New(x) or New(x, y) options, and
almost never the New() option.
*/
func New(dims ...int) *Mat {
	m := &Mat{}
	switch len(dims) {
	case 0:
		m = &Mat{
			0,
			0,
			make([]float64, 0),
		}
	case 1:
		m = &Mat{
			dims[0],
			dims[0],
			make([]float64, dims[0]*dims[0], 2*dims[0]*dims[0]),
		}
	case 2:
		m = &Mat{
			dims[0],
			dims[1],
			make([]float64, dims[0]*dims[1], 2*dims[0]*dims[1]),
		}
	case 3:
		m = &Mat{
			dims[0],
			dims[1],
			make([]float64, dims[0]*dims[1], dims[2]),
		}
	default:
		s := "\nIn mat64.%s, expected 0 to 3 arguments, but received %d arguments."
		s = fmt.Sprintf(s, "New()", len(dims))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	return m
}

/*
FromData creates a mat object from a []float64 or a [][]float64 slice.
This function is designed to do the "right thing" based on the type of
the slice passed to it. The "right thing" based on each possible case
is as follows:

Assume that s is a [][]float64, and v is a []float64 for the examples
below.

	x := mat64.FromData(v)

In this case, x.Dims() is (1, len(v)), and the values in x are the same
as the values in v. x is essentially a row vector.

Alternatively, this function can be invoked as:

	x := mat64.FromData(v, a

In this case, x.Dims() is (a, 1), and the values in x are the same
as the values in v. x is essentially a column vector. Note that a
must be equal to len(v).

Finally for the case where the data is a []float64, the function can be
invoked as:

	x := mat64.FromData(v, a, b)

In this case, x.Dims() is (a, b), and the values in x are the same as
the values in v. Note that a*b must be equal to len(v).

This function can also be invoked with data that is stored in a 2D
slice ([][]float64). Just as the []float64 case, there are three
possibilities:

	x := mat64.FromData(s)

In this case, x.Dims() is (len(s), len(s[0])), and the values in x
are the same as the values in s. It is assumed that s is not jagged.

Another form to call this function with a 2D slice of data is:

	x := mat64.FromData(s, a)

In this case, x.Dims() is (a, a), and the values in x are the same
as the values in s. Note that the total number of elements in s
must be exactly a*a.

Finally, this function can be called as:

	x := mat64.FromData(s, a, b)

In this case, x.Dims() is (a, b), and the values in x are the same
as the values in s. Note that the total number of elements in s
must be exactly a*b. Also note that this is equivalent to:

	x := mat64.FromData(s).Reshape(a, b)

Choose the format that suits your needs, as there is no performace
difference between the two forms.
*/
func FromData(oneOrTwoDSlice interface{}, dims ...int) *Mat {
	m := New()
	switch v := oneOrTwoDSlice.(type) {
	case []float64:
		switch len(dims) {
		case 0:
			m.vals = make([]float64, len(v), len(v)*2)
			copy(m.vals, v)
			m.r, m.c = 1, len(v)
		case 1:
			if dims[0] != len(v) {
				s := "\nIn mat64.%s, a 1D slice of data and a single int were passed.\n"
				s += "However the int (%d) is not equal to the length of the data (%d)."
				s = fmt.Sprintf(s, "FromData()", dims[0], len(v))
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			m.vals = make([]float64, dims[0], dims[0]*2)
			copy(m.vals, v)
			m.r, m.c = dims[0], 1
		case 2:
			if dims[0]*dims[1] != len(v) {
				s := "\nIn mat64.%s, a 1D slice of data and two ints were passed.\n"
				s += "However, the product of the two ints (%d, %d) does not equal\n"
				s += "the number of elements in the data slice, %d. They must be equal."
				s = fmt.Sprintf(s, "FromData()", dims[0]*dims[1], len(v))
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			m.vals = make([]float64, dims[0]*dims[1], dims[0]*dims[1]*2)
			copy(m.vals, v)
			m.r, m.c = dims[0], dims[1]
		default:
			s := "\nIn mat64.%s, a 1D slice of data and %d ints were passed.\n"
			s += "This function expects 0 to 2 integers. Please review the docs for\n"
			s += "this function and adjust the number of integers based on the\n"
			s += "desired output."
			s = fmt.Sprintf(s, "FromData()", len(dims))
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		} // switch len(dims) for case []float64
	case [][]float64:
		switch len(dims) {
		case 0:
			m.vals = make([]float64, len(v)*len(v[0]), len(v)*len(v[0])*2)
			for i := range v {
				for j := range v[i] {
					m.vals[i*len(v[0])+j] = v[i][j]
				}
			}
			m.r, m.c = len(v), len(v[0])
		case 1:
			if dims[0]*2 != len(v)*len(v[0]) {
				s := "\nIn mat64.%s, a 2D slice of data and 1 int were passed.\n"
				s += "This would generate a %d by %d Mat. However, %d*%d does not\n"
				s += "equal the number of elements in the passed 2D slice, %d.\n"
				s += "Note that this function expects a non-jagged 2D slice, and\n"
				s += "is assumed that every row in the passed 2D slice contains\n"
				s += "%d elements."
				s = fmt.Sprintf(s, "FromData()", dims[0], dims[0], dims[0], dims[0],
					len(v)*len(v[0]), len(v[0]))
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			m.vals = make([]float64, dims[0]*dims[0], dims[0]*dims[0]*2)
			for i := range v {
				for j := range v[i] {
					m.vals[i*len(v[0])+j] = v[i][j]
				}
			}
			m.r, m.c = len(v), len(v[0])
		case 2:
			if dims[0] != len(v) || dims[1] != len(v[0]) {
				s := "\nIn mat64.%s, a 2D slice of data and 2 ints were passed.\n"
				s += "However, the requested number of rows and columns (%d and %d)\n"
				s += "of the resultant Mat does not match the length and width of\n"
				s += "the data slice (%d and %d)."
				s = fmt.Sprintf(s, "FromData()", dims[0], dims[1], len(v), len(v[0]))
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			m.vals = make([]float64, dims[0]*dims[1], dims[0]*dims[1]*2)
			for i := range v {
				for j := range v[i] {
					m.vals[i*len(v[0])+j] = v[i][j]
				}
			}
			m.r, m.c = len(v), len(v[0])
		default:
			s := "\nIn mat64.%s, a 2D slice of data and %d ints were passed.\n"
			s += "However, this function expects 0 to 2 ints. Review the docs for\n"
			s += "this function and adjust the number of integers passed accordingly."
			s = fmt.Sprintf(s, "FromData()", len(dims))
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		} // switch len(dims) for case [][]float64
	default:
		s := "\nIn mat64.%s, expected input data of type []float64 or\n"
		s += "[][]float64, However, data of type \"%v\" was received."
		s = fmt.Sprintf(s, "FromData()", reflect.TypeOf(v))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	} // switch data.(type)
	return m
}

/*
FromCSV creates a mat object from a CSV (comma separated values) file. Here, we
assume that the number of rows of the resultant mat object is equal to the
number of lines, and the number of columns is equal to the number of entries
in each line. As before, we make sure that each line contains the same number
of elements.

The file to be read is assumed to be very large, and hence it is read one line
at a time. This results in some major inefficiencies, and it is recommended
that this function be used sparingly, and not as a major component of your
library/executable.

Unlike other mat creation functions in this package, the capacity of the mat
object created here is the same as its length since we assume the mat to
be very large.
*/
func FromCSV(filename string) *Mat {
	f, err := os.Open(filename)
	if err != nil {
		s := "\nIn mat64.%s, cannot open %s due to error: %v.\n"
		s = fmt.Sprintf(s, "FromCSV()", filename, err)
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	defer f.Close()
	r := csv.NewReader(f)
	// I am going with the assumption that a mat loaded from a CSV is going to
	// be large. So, we are going to read one line, and determine the number
	// of columns based on the number of comma separated entries in that line.
	// Then we will read the rest of the lines one at a time, checking that the
	// number of entries in each line is the same as the first line.
	str, err := r.Read()
	if err != nil {
		s := "\nIn mat64.%s, cannot read from %s due to error: %v.\n"
		s = fmt.Sprintf(s, "FromCSV()", filename, err)
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	line := 1
	m := New()
	// Start with one row, and set the number of entries per row
	m.r = 1
	m.c = len(str)
	row := make([]float64, len(str))
	for {
		for i := range str {
			row[i], err = strconv.ParseFloat(str[i], 64)
			if err != nil {
				s := "\nIn mat64.%s, item %d in line %d is %s, which cannot\n"
				s += "be converted to a float64 due to: %v"
				s = fmt.Sprintf(s, "FromCSV()", i, line, str[i], err)
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
		}
		m.vals = append(m.vals, row...)
		// Read the next line. If there is one, increment the number of rows
		str, err = r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			s := "\nIn mat64.%s, cannot read from %s due to error: %v.\n"
			s = fmt.Sprintf(s, "FromCSV()", filename, err)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
		line++
		if len(str) != len(row) {
			s := "\nIn mat64.%s, line %d in %s has %d entries. The first line\n"
			s += "(line 1) has %d entries.\n"
			s += "All the lines in the CSV file must contains the same number\n"
			s += "of entries.\n"
			s = fmt.Sprintf(s, "Load()", line, filename, len(str), len(row))
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
		m.r++
	}
	return m
}

func Rand(r, c int, args ...float64) *Mat {
	m := New(r, c)
	switch len(args) {
	case 0:
		for i := 0; i < m.r*m.c; i++ {
			m.vals[i] = rand.Float64()
		}
	case 1:
		to := args[0]
		for i := 0; i < m.r*m.c; i++ {
			m.vals[i] = rand.Float64() * to
		}
	case 2:
		from := args[0]
		to := args[1]
		if !(from < to) {
			s := "\nIn mat64.%s the first argument, %f, is not less than the\n"
			s += "second argument, %f. The first argument must be strictly\n"
			s += "less than the second.\n"
			s = fmt.Sprintf(s, "Rand()", from, to)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
		for i := 0; i < m.r*m.c; i++ {
			m.vals[i] = rand.Float64()*(to-from) + from
		}
	default:
		s := "\nIn mat64.%s expected 0 to 2 arguments, but received %d."
		s = fmt.Sprintf(s, "Rand()", len(args))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	return m
}

/*
Reshape changes the row and the columns of the mat object as long as the total
number of values contained in the mat object remains constant. The order and
the values of the mat does not change with this function.
*/
func (m *Mat) Reshape(rows, cols int) *Mat {
	if rows*cols != m.r*m.c {
		s := "\nIn mat64.%s, The total number of entries of the old and new shape\n"
		s += "must match.\n"
		s = fmt.Sprintf(s, "Reshape()")
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	} else {
		m.r = rows
		m.c = cols
	}
	return m
}

/*
Dims returns the number of rows and columns of a mat object.
*/
func (m *Mat) Dims() (int, int) {
	return m.r, m.c
}

/*
Vals returns the values contained in a mat object as a 1D slice of float64s.
*/
func (m *Mat) Vals() []float64 {
	s := make([]float64, len(m.vals))
	copy(s, m.vals)
	return s
}

/*
ToSlice returns the values of a mat object as a 2D slice of float64s.
*/
func (m *Mat) ToSlice() [][]float64 {
	s := make([][]float64, m.r)
	for i := range s {
		s[i] = make([]float64, m.c)
		for j := range s[i] {
			s[i][j] = m.vals[i*m.c+j]
		}
	}
	return s
}

/*
ToCSV creates a file with the passed name, and writes the content of a mat
object to it, by putting each row in a single comma separated line. The
number of entries in each line is equal to the columns of the mat object.
*/
func (m *Mat) ToCSV(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		s := "\nIn mat64.%s, cannot open %s due to error: %v.\n"
		s = fmt.Sprintf(s, "ToCSV()", fileName, err)
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	defer f.Close()
	str := ""
	idx := 0
	for i := 0; i < m.r; i++ {
		for j := 0; j < m.c; j++ {
			str += strconv.FormatFloat(m.vals[idx], 'e', 14, 64)
			if j+1 != m.c {
				str += ","
			}
			idx++
		}
		if i+1 != m.r {
			str += "\n"
		}
	}
	_, err = f.Write([]byte(str))
	if err != nil {
		s := "\nIn mat64.%s, cannot write to %s due to error: %v.\n"
		s = fmt.Sprintf(s, "ToCSV()", fileName, err)
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
}

/*
At returns a pointer to the float64 stored in the given row and column.
*/
func (m *Mat) At(r, c int) float64 {
	return m.vals[r*m.c+c]
}

/*
Foreach applies a given function to each element of a mat object. The given
function must take a pointer to a float64, and return nothing.
*/
func (m *Mat) Foreach(f func(*float64)) *Mat {
	for i := 0; i < m.r*m.c; i++ {
		f(&m.vals[i])
	}
	return m
}

/*
SetAll sets all values of a mat to the passed float64 value.
*/
func (m *Mat) SetAll(val float64) *Mat {
	for i := range m.vals {
		m.vals[i] = val
	}
	return m
}

/*
Set sets the value of a mat at a given row and column to a given
value.
*/
func (m *Mat) Set(r, c int, val float64) *Mat {
	m.vals[r*m.r+c] = val
	return m
}

/*
Col returns a new mat object whose values are equal to a column of the original
mat object. The number of Rows of the returned mat object is equal to the
number of rows of the original mat, and the number of columns is equal to 1.
*/
func (m *Mat) Col(x int) *Mat {
	if (x >= m.c) || (x < -m.c) {
		s := "\nIn mat64.%s the requested column %d is outside of bounds [%d, %d)\n"
		s = fmt.Sprintf(s, "Col()", x, m.c, m.c)
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	v := New(m.r, 1)
	if x >= 0 {
		for r := 0; r < m.r; r++ {
			v.vals[r] = m.vals[r*m.c+x]
		}
	} else {
		for r := 0; r < m.r; r++ {
			v.vals[r] = m.vals[r*m.c+(m.c+x)]
		}
	}
	return v
}

/*
Row returns a new mat object whose values are equal to a row of the original
mat object. The number of Rows of the returned mat object is equal to 1, and
the number of columns is equal to the number of columns of the original mat.
*/
func (m *Mat) Row(x int) *Mat {
	if (x >= m.r) || (x < -m.r) {
		s := "\nIn mat64.%s, row %d is outside of the bounds [-%d, %d)\n"
		s = fmt.Sprintf(s, "Row()", x, m.r, m.r)
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	v := New(1, m.c)
	if x >= 0 {
		for r := 0; r < m.c; r++ {
			v.vals[r] = m.vals[x*m.c+r]
		}
	} else {
		for r := 0; r < m.c; r++ {
			v.vals[r] = m.vals[(m.r+x)*m.c+r]
		}
	}
	return v
}

/*
Equals checks to see if two mat objects are equal. That mean that the two mats
have the same number of rows, same number of columns, and have the same float64
in each entry at a given index.
*/
func (m *Mat) Equals(n *Mat) bool {
	if m.r != n.r {
		return false
	}
	if m.c != n.c {
		return false
	}
	for i := 0; i < m.r*m.c; i++ {
		if m.vals[i] != n.vals[i] {
			return false
		}
	}
	return true
}

/*
Copy returns a duplicate of a mat object. The returned copy is "deep", meaning
that the object can be manipulated without effecting the original mat object.
*/
func (m *Mat) Copy() *Mat {
	n := New(m.r, m.c)
	copy(n.vals, m.vals)
	return n
}

/*
T returns the transpose of the original matrix. The transpose of a mat object
is defined in the usual manner, where every value at row x, and column y is
placed at row y, and column x. The number of rows and column of the transposed
mat are equal to the number of columns and rows of the original matrix,
respectively. This method creates a new mat object, and the original is
left intact.
*/
func (m *Mat) T() *Mat {
	n := New(m.c, m.r)
	idx := 0
	for i := 0; i < m.c; i++ {
		for j := 0; j < m.r; j++ {
			n.vals[idx] = m.vals[j*m.c+i]
			idx++
		}
	}
	return n
}

/*
All checks if a supplied function is true for all elements of a mat object.
For instance, consider

	positive := func(i *float64) bool {
		if i > 0.0 {
			return true
		}
		return false
	}

Then calling

	m.All(positive)

will return true if and only if all elements in m are positive.
*/
func (m *Mat) All(f func(*float64) bool) bool {
	for i := range m.vals {
		if !f(&m.vals[i]) {
			return false
		}
	}
	return true
}

/*
Any checks if a supplied function is true for one elements of a mat object.
For instance,

	positive := func(i *float64) bool {
		if i > 0.0 {
			return true
		}
		return false
	}

Then calling

	m.Any(positive)

would be true if at least one element of the mat object is positive.
*/
func (m *Mat) Any(f func(*float64) bool) bool {
	for i := range m.vals {
		if f(&m.vals[i]) {
			return true
		}
	}
	return false
}

func (m *Mat) Mul(float64OrMat64 interface{}) *Mat {
	switch v := float64OrMat64.(type) {
	case float64:
		for i := range m.vals {
			m.vals[i] *= v
		}
	case *Mat:
		if v.r != m.r {
			s := "\nIn mat64.%s, the number of the rows of the receiver is %d\n"
			s += "but the number of rows of the passed mat is %d. They must\n"
			s += "match.\n"
			s = fmt.Sprintf(s, "Mul()", m.r, v.r)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
		if v.c != m.c {
			s := "\nIn mat64.%s, the number of the columns of the receiver is %d\n"
			s += "but the number of columns of the passed mat is %d. They must\n"
			s += "match.\n"
			s = fmt.Sprintf(s, "Mul()", m.c, v.c)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
		for i := range m.vals {
			m.vals[i] *= v.vals[i]
		}
	default:
		s := "\nIn mat64.%s, the passed value must be a float64 or *Mat.\n"
		s += "However, value of type  \"%v\" was received.\n"
		s = fmt.Sprintf(s, "Mul()", reflect.TypeOf(v))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	return m
}

func (m *Mat) Add(float64OrMat64 interface{}) *Mat {
	switch v := float64OrMat64.(type) {
	case float64:
		for i := range m.vals {
			m.vals[i] += v
		}
	case *Mat:
		if v.r != m.r {
			s := "\nIn mat64.%s, the number of the rows of the receiver is %d\n"
			s += "but the number of rows of the passed mat is %d. They must\n"
			s += "match.\n"
			s = fmt.Sprintf(s, "Add()", m.r, v.r)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
		if v.c != m.c {
			s := "\nIn mat64.%s, the number of the columns of the receiver is %d\n"
			s += "but the number of columns of the passed mat is %d. They must\n"
			s += "match.\n"
			s = fmt.Sprintf(s, "Add()", m.c, v.c)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
		for i := range m.vals {
			m.vals[i] += v.vals[i]
		}
	default:
		s := "\nIn mat64.%s, the passed value must be a float64 or *Mat.\n"
		s += "However, value of type  \"%v\" was received.\n"
		s = fmt.Sprintf(s, "Add()", reflect.TypeOf(v))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	return m
}

func (m *Mat) Sub(float64OrMat64 interface{}) *Mat {
	switch v := float64OrMat64.(type) {
	case float64:
		for i := range m.vals {
			m.vals[i] -= v
		}
	case *Mat:
		if v.r != m.r {
			s := "\nIn mat64.%s, the number of the rows of the receiver is %d\n"
			s += "but the number of rows of the passed mat is %d. They must\n"
			s += "match.\n"
			s = fmt.Sprintf(s, "Sub()", m.r, v.r)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
		if v.c != m.c {
			s := "\nIn mat64.%s, the number of the columns of the receiver is %d\n"
			s += "but the number of columns of the passed mat is %d. They must\n"
			s += "match.\n"
			s = fmt.Sprintf(s, "Sub()", m.c, v.c)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
		for i := range m.vals {
			m.vals[i] -= v.vals[i]
		}
	default:
		s := "\nIn mat64.%s, the passed value must be a float64 or *Mat.\n"
		s += "However, value of type  \"%v\" was received.\n"
		s = fmt.Sprintf(s, "Sub()", reflect.TypeOf(v))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	return m
}

func (m *Mat) Div(float64OrMat64 interface{}) *Mat {
	switch v := float64OrMat64.(type) {
	case float64:
		for i := range m.vals {
			m.vals[i] /= v
		}
	case *Mat:
		if v.r != m.r {
			s := "\nIn mat64.%s, the number of the rows of the receiver is %d\n"
			s += "but the number of rows of the passed mat is %d. They must\n"
			s += "match.\n"
			s = fmt.Sprintf(s, "Div()", m.r, v.r)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
		if v.c != m.c {
			s := "\nIn mat64.%s, the number of the columns of the receiver is %d\n"
			s += "but the number of columns of the passed mat is %d. They must\n"
			s += "match.\n"
			s = fmt.Sprintf(s, "Div()", m.c, v.c)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
		for i := range m.vals {
			m.vals[i] /= v.vals[i]
		}
	default:
		s := "\nIn mat64.%s, the passed value must be a float64 or *Mat.\n"
		s += "However, value of type  \"%v\" was received.\n"
		s = fmt.Sprintf(s, "Div()", reflect.TypeOf(v))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	return m
}

func (m *Mat) Sum(args ...int) float64 {
	sum := 0.0
	switch len(args) {
	case 0:
		for i := range m.vals {
			sum += m.vals[i]
		}
	case 2:
		axis, slice := args[0], args[1]
		switch axis {
		case 0:
			if (slice >= m.r) || (slice < 0) {
				s := "\nIn mat64.%s the row %d is outside of bounds [0, %d)\n"
				s = fmt.Sprintf(s, "Sum()", slice, m.r)
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			for i := 0; i < m.c; i++ {
				sum += m.vals[slice*m.c+i]
			}
		case 1:
			if (slice >= m.c) || (slice < 0) {
				s := "\nIn mat64.%s the column %d is outside of bounds [0, %d)\n"
				s = fmt.Sprintf(s, "Sum()", slice, m.c)
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			for i := 0; i < m.r; i++ {
				sum += m.vals[i*m.c+slice]
			}
		default:
			s := "\nIn mat64.%s, the first argument must be 0 or 1, however %d "
			s += "was received.\n"
			s = fmt.Sprintf(s, "Sum()", axis)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
	default:
		s := "\nIn mat64.%s, 0 or 2 arguments expected, but %d was received.\n"
		s = fmt.Sprintf(s, "Sum()", len(args))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	return sum
}

func (m *Mat) Avg(args ...int) float64 {
	sum := 0.0
	switch len(args) {
	case 0:
		for i := range m.vals {
			sum += m.vals[i]
		}
		sum /= float64(len(m.vals))
	case 2:
		axis, slice := args[0], args[1]
		if axis == 0 {
			if (slice >= m.r) || (slice < 0) {
				s := "\nIn mat64.%s the row %d is outside of bounds [0, %d)\n"
				s = fmt.Sprintf(s, "Avg()", slice, m.r)
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			for i := 0; i < m.c; i++ {
				sum += m.vals[slice*m.c+i]
			}
			sum /= float64(m.c)
		} else if axis == 1 {
			if (slice >= m.c) || (slice < 0) {
				s := "\nIn mat64.%s the column %d is outside of bounds [0, %d)\n"
				s = fmt.Sprintf(s, "Avg()", slice, m.c)
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			for i := 0; i < m.r; i++ {
				sum += m.vals[i*m.c+slice]
			}
			sum /= float64(m.r)
		} else {
			s := "\nIn mat64.%s, the first argument must be 0 or 1, however %d "
			s += "was received.\n"
			s = fmt.Sprintf(s, "Avg()", axis)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
	default:
		s := "\nIn mat64.%s, 0 or 2 arguments expected, but %d was received.\n"
		s = fmt.Sprintf(s, "Avg()", len(args))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	return sum
}

func (m *Mat) Prd(args ...int) float64 {
	prd := 1.0
	switch len(args) {
	case 0:
		for i := range m.vals {
			prd *= m.vals[i]
		}
	case 2:
		axis, slice := args[0], args[1]
		if axis == 0 {
			if (slice >= m.r) || (slice < 0) {
				s := "\nIn mat64.%s the row %d is outside of bounds [0, %d)\n"
				s = fmt.Sprintf(s, "Prd()", slice, m.r)
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			for i := 0; i < m.c; i++ {
				prd *= m.vals[slice*m.c+i]
			}
		} else if axis == 1 {
			if (slice >= m.c) || (slice < 0) {
				s := "\nIn mat64.%s the column %d is outside of bounds [0, %d)\n"
				s = fmt.Sprintf(s, "Prd()", slice, m.c)
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			for i := 0; i < m.r; i++ {
				prd *= m.vals[i*m.c+slice]
			}
		} else {
			s := "\nIn mat64.%s, the first argument must be 0 or 1, however %d "
			s += "was received.\n"
			s = fmt.Sprintf(s, "Prd()", axis)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
	default:
		s := "\nIn mat64.%s, 0 or 2 arguments expected, but %d was received.\n"
		s = fmt.Sprintf(s, "Prd()", len(args))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	return prd
}

func (m *Mat) Std(args ...int) float64 {
	std := 0.0
	switch len(args) {
	case 0:
		avg := m.Avg()
		sum := 0.0
		for i := range m.vals {
			sum += ((avg - m.vals[i]) * (avg - m.vals[i]))
		}
		std = math.Sqrt(sum / float64(len(m.vals)))
	case 2:
		axis, slice := args[0], args[1]
		if axis == 0 {
			if (slice >= m.r) || (slice < 0) {
				s := "\nIn mat64.%s the row %d is outside of bounds [0, %d)\n"
				s = fmt.Sprintf(s, "Std()", slice, m.r)
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			avg := m.Avg(axis, slice)
			sum := 0.0
			for i := 0; i < m.c; i++ {
				sum += ((avg - m.vals[slice*m.c+i]) * (avg - m.vals[slice*m.c+i]))
			}
			std = math.Sqrt(sum / float64(len(m.vals)))
		} else if axis == 1 {
			if (slice >= m.c) || (slice < 0) {
				s := "\nIn mat64.%s the column %d is outside of bounds [0, %d)\n"
				s = fmt.Sprintf(s, "Std()", slice, m.c)
				color.Red(s)
				color.Yellow("\nStack trace for this error:\n\n")
				q := string(debug.Stack())
				w := strings.Split(q, "\n")
				fmt.Println(strings.Join(w[5:], "\n"))
				os.Exit(1)
			}
			avg := m.Avg(axis, slice)
			sum := 0.0
			for i := 0; i < m.r; i++ {
				sum += ((avg - m.vals[i*m.c+slice]) * (avg - m.vals[i*m.c+slice]))
			}
			std = math.Sqrt(sum / float64(len(m.vals)))
		} else {
			s := "\nIn mat64.%s, the first argument must be 0 or 1, however %d "
			s += "was received.\n"
			s = fmt.Sprintf(s, "Std()", axis)
			color.Red(s)
			color.Yellow("\nStack trace for this error:\n\n")
			q := string(debug.Stack())
			w := strings.Split(q, "\n")
			fmt.Println(strings.Join(w[5:], "\n"))
			os.Exit(1)
		}
	default:
		s := "\nIn mat64.%s, 0 or 2 arguments must be passed, but %d was received.\n"
		s = fmt.Sprintf(s, "Std()", len(args))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	return std
}

/*
Dot is the matrix multiplication of two mat objects. Consider the following two
mats:

	m := New(5, 6)
	n := New(6, 10)

then

	o := m.Dot(n)

is a 5 by 10 mat whose element at row i and column j is given by:

	Sum(m.Row(i).Mul(n.col(j))
*/
func (m *Mat) Dot(n *Mat) *Mat {
	if m.c != n.r {
		s := "\nIn mat64.%s the number of columns of the first mat is %d\n"
		s += "which is not equal to the number of rows of the second mat,\n"
		s += "which is %d. They must be equal.\n"
		s = fmt.Sprintf(s, "Dot()", m.c, n.r)
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	o := New(m.r, n.c)
	for i := 0; i < m.r; i++ {
		for j := 0; j < n.c; j++ {
			for k := 0; k < m.c; k++ {
				o.vals[i*o.c+j] += (m.vals[i*m.c+k] * n.vals[k*n.c+j])
			}
		}
	}
	return o
}

/*
String returns the string representation of a mat. This is done by putting
every row into a line, and separating the entries of that row by a space. note
that the last line does not contain a newline.
*/
func (m *Mat) String() string {
	var str string
	str += "["
	for i := 0; i < m.r; i++ {
		for j := 0; j < m.c; j++ {
			if j == 0 {
				str += "["
			}
			str += strconv.FormatFloat(m.vals[i*m.c+j], 'f', 14, 64)
			if j+1 != m.c {
				str += ",\t"
			}
		}
		if i+1 <= m.r {
			str += "]\n "
		}
	}
	str = str[:len(str)-2] // take out the last newline.
	str += "]\n"
	return str
}

/*
AppendCol appends a column to the right side of a Mat.
*/
func (m *Mat) AppendCol(v []float64) *Mat {
	if m.r != len(v) {
		s := "\nIn mat64.%s the number of rows of the reciever is %d, while\n"
		s += "the number of rows of the vector is %d. They must be equal.\n"
		s = fmt.Sprintf(s, "AppendCol()", m.r, len(v))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	// TODO: redo this by hand, instead of taking this shortcut... or check if
	// this is a huge bottleneck
	q := m.ToSlice()
	for i := range q {
		q[i] = append(q[i], v[i])
	}
	m.c++
	m.vals = append(m.vals, v...)
	for i := 0; i < m.r; i++ {
		for j := 0; j < m.c; j++ {
			m.vals[i*m.c+j] = q[i][j]
		}
	}
	return m
}

/*
AppendRow appends a row to the bottom of a Mat.
*/
func (m *Mat) AppendRow(v []float64) *Mat {
	if m.c != len(v) {
		s := "\nIn mat64.%s the number of cols of the receiver is %d, while\n"
		s += "the number of rows of the vector is %d. They must be equal.\n"
		s = fmt.Sprintf(s, "AppendRow()", m.c, len(v))
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	if cap(m.vals) < (len(m.vals) + len(v)) {
		newVals := make([]float64, len(m.vals)+len(v), len(m.vals)+len(v)*2)
		lastElem := len(m.vals)
		for i := range m.vals {
			newVals[i] = m.vals[i]
		}
		for i := range v {
			newVals[lastElem+i] = v[i]
		}
		m.vals = newVals
	} else {
		m.vals = append(m.vals, v...)
	}
	m.r++
	return m
}

/*
Concat concatenates the inner slices of two `[][]float64` arguments..

For example, if we have:

	m := [[1.0, 2.0], [3.0, 4.0]]
	n := [[5.0, 6.0], [7.0, 8.0]]
	o := m.Concat(n).Row(0) // [1.0, 2.0, 5.0, 6.0]

*/
func (m *Mat) Concat(n *Mat) *Mat {
	if m.r != n.r {
		s := "\nIn mat64.%s the number of rows of the receiver is %d, while\n"
		s += "the number of rows of the second Mat is %d. They must be equal.\n"
		s = fmt.Sprintf(s, "Concat()", m.r, n.r)
		color.Red(s)
		color.Yellow("\nStack trace for this error:\n\n")
		q := string(debug.Stack())
		w := strings.Split(q, "\n")
		fmt.Println(strings.Join(w[5:], "\n"))
		os.Exit(1)
	}
	q := m.ToSlice()
	t := n.Vals()
	r := n.ToSlice()
	m.vals = append(m.vals, t...)
	for i := range q {
		q[i] = append(q[i], r[i]...)
	}
	m.c += n.c
	for i := 0; i < m.r; i++ {
		for j := 0; j < m.c; j++ {
			m.vals[i*m.c+j] = q[i][j]
		}
	}
	return m
}

func (m *Mat) Tanh() *Mat {
	for i := range m.vals {
		m.vals[i] = math.Tanh(m.vals[i])
	}
	return m
}
