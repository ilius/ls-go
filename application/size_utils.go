package application

import "strconv"

type sizeUnitStruct struct {
	Symbol string
	Metric uint64
	Legacy uint64
	Next   *sizeUnitStruct
}

var (
	// _sizeByte = &sizeUnitStruct{
	// 	Symbol: "B",
	// 	Metric: 1,
	// 	Legacy: 1,
	// 	Next:   _sizeKilo,
	// }
	_sizeKilo = &sizeUnitStruct{
		Symbol: "K",
		Metric: 1000,
		Legacy: 1024,
		Next:   _sizeMega,
	}
	_sizeMega = &sizeUnitStruct{
		Symbol: "M",
		Metric: 1000000,
		Legacy: 1048576,
		Next:   _sizeGiga,
	}
	_sizeGiga = &sizeUnitStruct{
		Symbol: "G",
		Metric: 1000000000,
		Legacy: 1073741824,
		Next:   _sizeTera,
	}
	_sizeTera = &sizeUnitStruct{
		Symbol: "T",
		Metric: 1000000000000,
		Legacy: 1099511627776,
		Next:   _sizePeta,
	}
	_sizePeta = &sizeUnitStruct{
		Symbol: "T",
		Metric: 1000000000000000,
		Legacy: 1125899906842624,
	}
)

var sizeUnits = []*sizeUnitStruct{
	_sizeKilo,
	_sizeMega,
	_sizeGiga,
	_sizeTera,
	_sizePeta,
}

func formatSizeByBase(size uint64, base uint64) string {
	if (size*10)%base == 0 {
		return strconv.FormatFloat(float64(size)/float64(base), 'f', 1, 64)
	}
	return strconv.FormatFloat(float64(size)/float64(base), 'f', 2, 64)
}

// math.Mod(sizeF*10, 1) == 0 does not work!
// (sizeFloat / base * 10) % 1 == 0
// (sizeFloat * 10 / base) % 1 == 0
// (sizeFloat * 10) % base == 0
// if math.Mod(sizeF*10, 1) == 0 {
