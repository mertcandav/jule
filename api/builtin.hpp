// Copyright 2022 The Jule Programming Language.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

#ifndef __JULEC_BUILTIN_HPP
#define __JULEC_BUILTIN_HPP

typedef u8_julet   JULEC_ID(byte); // builtin: type byte: u8
typedef i32_julet  JULEC_ID(rune); // builtin: type rune: i32

// Update here manually if changed output of the JULEC_ID
#define _out(_EXPR) \
    (std::cout << _EXPR)

// Declarations
struct JULEC_ID(Error);
template<typename _Item_t>
int_julet JULEC_ID(copy)(const slice<_Item_t> &_Dest,
                         const slice<_Item_t> &_Src) noexcept;
template<typename _Item_t>
slice<_Item_t> JULEC_ID(append)(const slice<_Item_t> &_Src,
                                const slice<_Item_t> &_Components) noexcept;
template<typename T>
ptr<T> JULEC_ID(new)(void) noexcept;

// Definitions

/* Panic function defined at main header */

template<typename _Obj_t>
inline void JULEC_ID(outln)(const _Obj_t _Obj) noexcept {
    JULEC_ID(out)(_Obj);
    std::cout << std::endl;
}

struct JULEC_ID(Error) {
    virtual str_julet error(void) = 0;
};

template<typename _Item_t>
inline slice<_Item_t> JULEC_ID(make)(const int_julet &_N) noexcept
{ return _N < 0 ? nil : slice<_Item_t>(_N); }

template<typename _Item_t>
int_julet JULEC_ID(copy)(const slice<_Item_t> &_Dest,
                         const slice<_Item_t> &_Src) noexcept {
    if (_Dest.empty() || _Src.empty()) { return 0; }
    int_julet _len = _Dest.len() > _Src.len() ? _len = _Src.len()
                     : _Src.len() > _Dest.len() ? _len = _Dest.len()
                     : _len = _Src.len();
    for (int_julet _index{0}; _index < _len; ++_index)
    { _Dest._slice[_index] = _Src._slice[_index]; }
    return _len;
}

template<typename _Item_t>
slice<_Item_t> JULEC_ID(append)(const slice<_Item_t> &_Src,
                                const slice<_Item_t> &_Components) noexcept {
    const int_julet _N{_Src.len() + _Components.len()};
    slice<_Item_t> _buffer{JULEC_ID(make)<_Item_t>(_N)};
    JULEC_ID(copy)<_Item_t>(_buffer, _Src);
    for (int_julet _index{0}; _index < _Components.len(); ++_index)
    { _buffer[_Src.len()+_index] = _Components._slice[_index]; }
    return _buffer;
}

template<typename T>
ptr<T> JULEC_ID(new)(void) noexcept {
    ptr<T> _ptr;
    _ptr._heap = new(std::nothrow) bool*{__JULEC_PTR_HEAP_TRUE};
    if (!_ptr._heap)
    { JULEC_ID(panic)(__JULEC_ERROR_MEMORY_ALLOCATION_FAILED); }
    _ptr._ptr = new(std::nothrow) T*;
    if (!_ptr._ptr)
    { JULEC_ID(panic)(__JULEC_ERROR_MEMORY_ALLOCATION_FAILED); }
    *_ptr._ptr = new(std::nothrow) T;
    if (!*_ptr._ptr)
    { JULEC_ID(panic)(__JULEC_ERROR_MEMORY_ALLOCATION_FAILED); }
    _ptr._ref = new(std::nothrow) uint_julet{1};
    if (!_ptr._ref)
    { JULEC_ID(panic)(__JULEC_ERROR_MEMORY_ALLOCATION_FAILED); }
    return _ptr;
}

#endif // #ifndef __JULEC_BUILTIN_HPP
