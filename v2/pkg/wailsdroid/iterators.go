package wailsdroid

// workarounds for arrays not supported in gomobile

type Iterator interface {
	Index() int
	HasValue() bool
	Values() int
	Next()
}

type StringIterator interface {
	Iterator
	Value() string
}

type stringIterator struct {
	i      int
	values []string
}

func (si *stringIterator) Index() int {
	return si.i
}

func (si *stringIterator) HasValue() bool {
	return si.i < len(si.values)
}

func (si *stringIterator) Values() int {
	return len(si.values)
}

func (si *stringIterator) Next() {
	si.i++
}

func (si *stringIterator) Value() string {
	return si.values[si.i]
}

func NewStringIterator(values []string) StringIterator {
	return &stringIterator{values: values}
}

type MapStream interface {
	StringIterator
	Key() string
}

type mapStream struct {
	StringIterator
	keys []string
}

func (ms *mapStream) Key() string {
	return ms.keys[ms.Index()]
}

func HeaderMapToMapStream(headers map[string][]string) MapStream {
	keys := make([]string, len(headers))
	values := make([]string, len(headers))

	i := 0
	for key, value := range headers {
		keys[i] = key
		values[i] = value[0]
		i++
	}

	return &mapStream{
		NewStringIterator(values),
		keys,
	}
}

/*
type ScreenIterator interface {
	Iterator
	Value() frontend.Screen
}

type screenIterator struct {
	i      int
	values []frontend.Screen
}

func (si *screenIterator) Index() int {
	return si.i
}

func (si *screenIterator) HasValue() bool {
	return si.i < len(si.values)
}

func (si *screenIterator) Next() {
	si.i++
}

func (si *screenIterator) Values() int {
	return len(si.values)
}

func (si *screenIterator) Value() frontend.Screen {
	return si.values[si.i]
}

func NewScreenIterator(values []frontend.Screen) ScreenIterator {
	return &screenIterator{values: values}
}
*/

/*
type FileFilterIterator interface {
	Iterator
	Value() frontend.FileFilter
}

type fileFilterIterator struct {
	i      int
	values []frontend.FileFilter
}

func (si *fileFilterIterator) Index() int {
	return si.i
}

func (si *fileFilterIterator) HasValue() bool {
	return si.i < len(si.values)
}

func (si *fileFilterIterator) Next() {
	si.i++
}

func (si *fileFilterIterator) Values() int {
	return len(si.values)
}

func (si *fileFilterIterator) Value() frontend.FileFilter {
	return si.values[si.i]
}

func NewFileFilterIterator(values []frontend.FileFilter) FileFilterIterator {
	return &fileFilterIterator{values: values}
}
*/
