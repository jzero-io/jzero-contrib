package handlerx

import (
	"encoding/json"
	"io"
	"math"
	"strings"
	"unsafe"

	"github.com/json-iterator/go"
	"github.com/samber/lo"
)

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1

func RegisterPointerFuzzyDecoders() {
	jsoniter.RegisterTypeDecoder("*string", &fuzzyPointerStringDecoder{})
	jsoniter.RegisterTypeDecoder("*float32", &fuzzyFloat32Decoder{})
	jsoniter.RegisterTypeDecoder("*float64", &fuzzyFloat64Decoder{})
	jsoniter.RegisterTypeDecoder("*int", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(maxInt) || val < float64(minInt) {
				iter.ReportError("fuzzy decode *int", "exceed range")
				return
			}
			*((**int)(ptr)) = lo.ToPtr(int(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.StringValue {
				str := iter.ReadString()
				if str == "" {
					*((*any)(ptr)) = nil
				} else {
					iter.ReportError("fuzzy decode *int", "expected null or integer value")
				}
			} else if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**int)(ptr)) = lo.ToPtr(iter.ReadInt())
			}
		}
	}})

	jsoniter.RegisterTypeDecoder("*uint", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(maxUint) || val < 0 {
				iter.ReportError("fuzzy decode *uint", "exceed range")
				return
			}
			*((**uint)(ptr)) = lo.ToPtr(uint(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.StringValue {
				str := iter.ReadString()
				if str == "" {
					*((*any)(ptr)) = nil
				} else {
					iter.ReportError("fuzzy decode *uint", "expected null or integer value")
				}
			} else if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**uint)(ptr)) = lo.ToPtr(iter.ReadUint())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*int8", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt8) || val < float64(math.MinInt8) {
				iter.ReportError("fuzzy decode *int8", "exceed range")
				return
			}
			*((**int8)(ptr)) = lo.ToPtr(int8(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.StringValue {
				str := iter.ReadString()
				if str == "" {
					*((*any)(ptr)) = nil
				} else {
					iter.ReportError("fuzzy decode *int8", "expected null or integer value")
				}
			} else if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**int8)(ptr)) = lo.ToPtr(iter.ReadInt8())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*uint8", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint8) || val < 0 {
				iter.ReportError("fuzzy decode *uint8", "exceed range")
				return
			}
			*((*uint8)(ptr)) = uint8(val)
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.StringValue {
				str := iter.ReadString()
				if str == "" {
					*((*any)(ptr)) = nil
				} else {
					iter.ReportError("fuzzy decode *uint8", "expected null or integer value")
				}
			} else if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**uint8)(ptr)) = lo.ToPtr(iter.ReadUint8())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*int16", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt16) || val < float64(math.MinInt16) {
				iter.ReportError("fuzzy decode *int16", "exceed range")
				return
			}
			*((**uint16)(ptr)) = lo.ToPtr(uint16(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.StringValue {
				str := iter.ReadString()
				if str == "" {
					*((*any)(ptr)) = nil
				} else {
					iter.ReportError("fuzzy decode *int16", "expected null or integer value")
				}
			} else if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**int16)(ptr)) = lo.ToPtr(iter.ReadInt16())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*uint16", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint16) || val < 0 {
				iter.ReportError("fuzzy decode *uint16", "exceed range")
				return
			}
			*((**uint16)(ptr)) = lo.ToPtr(uint16(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.StringValue {
				str := iter.ReadString()
				if str == "" {
					*((*any)(ptr)) = nil
				} else {
					iter.ReportError("fuzzy decode *uint16", "expected null or integer value")
				}
			} else if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**uint16)(ptr)) = lo.ToPtr(iter.ReadUint16())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*int32", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt32) || val < float64(math.MinInt32) {
				iter.ReportError("fuzzy decode *int32", "exceed range")
				return
			}
			*((**int32)(ptr)) = lo.ToPtr(int32(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.StringValue {
				str := iter.ReadString()
				if str == "" {
					*((*any)(ptr)) = nil
				} else {
					iter.ReportError("fuzzy decode *int32", "expected null or integer value")
				}
			} else if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**int32)(ptr)) = lo.ToPtr(iter.ReadInt32())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*uint32", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint32) || val < 0 {
				iter.ReportError("fuzzy decode *uint32", "exceed range")
				return
			}
			*((**uint32)(ptr)) = lo.ToPtr(uint32(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.StringValue {
				str := iter.ReadString()
				if str == "" {
					// 当值为空字符串时，设置目标指针为 nil
					*((*any)(ptr)) = nil
				} else {
					// 如果需要支持非空字符串转整数，请在此添加逻辑
					iter.ReportError("fuzzy decode *uint32", "expected null or integer value")
				}
			} else if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**uint32)(ptr)) = lo.ToPtr(iter.ReadUint32())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*int64", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxInt64) || val < float64(math.MinInt64) {
				iter.ReportError("fuzzy decode *int64", "exceed range")
				return
			}
			*((**int64)(ptr)) = lo.ToPtr(int64(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.StringValue {
				str := iter.ReadString()
				if str == "" {
					*((*any)(ptr)) = nil
				} else {
					iter.ReportError("fuzzy decode *int64", "expected null or integer value")
				}
			} else if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**int64)(ptr)) = lo.ToPtr(iter.ReadInt64())
			}
		}
	}})
	jsoniter.RegisterTypeDecoder("*uint64", &fuzzyPointerIntegerDecoder{fun: func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(math.MaxUint64) || val < 0 {
				iter.ReportError("fuzzy decode *uint64", "exceed range")
				return
			}
			*((**uint64)(ptr)) = lo.ToPtr(uint64(val))
		} else {
			token := iter.WhatIsNext()
			if token == jsoniter.StringValue {
				str := iter.ReadString()
				if str == "" {
					// 当值为空字符串时，设置目标指针为 nil
					*((*any)(ptr)) = nil
				} else {
					// 如果需要支持非空字符串转整数，请在此添加逻辑
					iter.ReportError("fuzzy decode *uint64", "expected null or integer value")
				}
			} else if token == jsoniter.NilValue {
				*((*any)(ptr)) = nil
			} else {
				*((**uint64)(ptr)) = lo.ToPtr(iter.ReadUint64())
			}
		}
	}})
}

type fuzzyPointerStringDecoder struct {
}

func (decoder *fuzzyPointerStringDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoniter.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		*((**string)(ptr)) = lo.ToPtr(string(number))
	case jsoniter.StringValue:
		*((**string)(ptr)) = lo.ToPtr(iter.ReadString())
	case jsoniter.NilValue:
		iter.Skip()
		*((**string)(ptr)) = nil
	default:
		iter.ReportError("fuzzyStringDecoder", "not number or string")
	}
}

type fuzzyPointerIntegerDecoder struct {
	fun func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator)
}

func (decoder *fuzzyPointerIntegerDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoniter.NumberValue:
		var number json.Number
		iter.ReadVal(&number)
		str = string(number)
	case jsoniter.StringValue:
		str = iter.ReadString()
		if str == "" {
			str = "null"
		}
	case jsoniter.BoolValue:
		if iter.ReadBool() {
			str = "1"
		} else {
			str = "0"
		}
	case jsoniter.NilValue:
		iter.Skip()
		str = "null"
	default:
		iter.ReportError("fuzzyPointerIntegerDecoder", "not number or string")
	}
	if len(str) == 0 {
		str = "0"
	}
	newIter := iter.Pool().BorrowIterator([]byte(str))
	defer iter.Pool().ReturnIterator(newIter)
	isFloat := strings.IndexByte(str, '.') != -1
	decoder.fun(isFloat, ptr, newIter)
	if newIter.Error != nil && newIter.Error != io.EOF {
		iter.Error = newIter.Error
	}
}

type fuzzyFloat32Decoder struct {
}

func (decoder *fuzzyFloat32Decoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoniter.NumberValue:
		*((**float32)(ptr)) = lo.ToPtr(iter.ReadFloat32())
	case jsoniter.StringValue:
		str = iter.ReadString()
		newIter := iter.Pool().BorrowIterator([]byte(str))
		defer iter.Pool().ReturnIterator(newIter)
		*((**float32)(ptr)) = lo.ToPtr(newIter.ReadFloat32())
		if newIter.Error != nil && newIter.Error != io.EOF {
			iter.Error = newIter.Error
		}
	case jsoniter.BoolValue:
		// support bool to float32
		if iter.ReadBool() {
			*((**float32)(ptr)) = lo.ToPtr(float32(1))
		} else {
			*((**float32)(ptr)) = lo.ToPtr(float32(0))
		}
	case jsoniter.NilValue:
		iter.Skip()
		*((**float32)(ptr)) = nil
	default:
		iter.ReportError("fuzzyPointerFloat32Decoder", "not number or string")
	}
}

type fuzzyFloat64Decoder struct {
}

func (decoder *fuzzyFloat64Decoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoniter.NumberValue:
		*((**float64)(ptr)) = lo.ToPtr(iter.ReadFloat64())
	case jsoniter.StringValue:
		str = iter.ReadString()
		newIter := iter.Pool().BorrowIterator([]byte(str))
		defer iter.Pool().ReturnIterator(newIter)
		*((**float64)(ptr)) = lo.ToPtr(newIter.ReadFloat64())
		if newIter.Error != nil && newIter.Error != io.EOF {
			iter.Error = newIter.Error
		}
	case jsoniter.BoolValue:
		// support bool to float64
		if iter.ReadBool() {
			*((**float64)(ptr)) = lo.ToPtr(float64(1))
		} else {
			*((**float64)(ptr)) = lo.ToPtr(float64(0))
		}
	case jsoniter.NilValue:
		iter.Skip()
		*((**float64)(ptr)) = nil
	default:
		iter.ReportError("fuzzyFloat64Decoder", "not number or string")
	}
}
