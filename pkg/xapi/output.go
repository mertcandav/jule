package xapi

// CxxMain is the entry point output of X-CXX code.
var CxxMain = `// region X_ENTRY_POINT
int main(void) {
    std::set_terminate(&x_terminate_handler);
    std::cout << std::boolalpha;
#ifdef _WINDOWS
    SetConsoleOutputCP(CP_UTF8);
    _setmode(_fileno(stdin), 0x00020000);
#endif

    ` + InitializerCaller + `();
    XID(main());

    return EXIT_SUCCESS;
}
// endregion X_ENTRY_POINT`

// CxxDefault is the default pre-cxx code output of X-CXX code.
var CxxDefault = `#if defined(WIN32) || defined(_WIN32) || defined(__WIN32__) || defined(__NT__)
#define _WINDOWS
#endif

// region X_STANDARD_IMPORTS
#include <iostream>
#include <cstring>
#include <string>
#include <sstream>
#include <functional>
#include <vector>
#include <map>
#include <thread>
#include <typeinfo>
#ifdef _WINDOWS
#include <codecvt>
#include <windows.h>
#include <fcntl.h>
#endif
// endregion X_STANDARD_IMPORTS

#define _CONCAT(_A, _B) _A ## _B
#define CONCAT(_A, _B) _CONCAT(_A, _B)
#define XID(_Identifier) CONCAT(_, _Identifier)

static inline void XID(panic)(const char *_Message);

// region X_CXX_API
// region X_BUILTIN_VALUES
#define nil nullptr
// endregion X_BUILTIN_VALUES

// region X_BUILTIN_TYPES
typedef std::size_t                       uint_xt;
typedef std::make_signed<uint_xt>::type   int_xt;
typedef int8_t                            i8_xt;
typedef int16_t                           i16_xt;
typedef int32_t                           i32_xt;
typedef int64_t                           i64_xt;
typedef uint8_t                           u8_xt;
typedef uint16_t                          u16_xt;
typedef uint32_t                          u32_xt;
typedef uint64_t                          u64_xt;
typedef float                             f32_xt;
typedef double                            f64_xt;
typedef unsigned char                     char_xt;
typedef bool                              bool_xt;
typedef void                              *voidptr_xt;
typedef intptr_t                          intptr_xt;
typedef uintptr_t                         uintptr_xt;

#define func std::function

// region X_STRUCTURES
template<typename _Item_t>
class array {
public:
    std::vector<_Item_t> _buffer{};

    array<_Item_t>(const uint_xt &_N) noexcept
    { this->_buffer = std::vector<_Item_t>(_N); }

    array<_Item_t>(void) noexcept {}

    array<_Item_t>(const std::nullptr_t) noexcept {}

    array<_Item_t>(const array<_Item_t>& _Src) noexcept
    { this->_buffer = _Src._buffer; }

    array<_Item_t>(const std::initializer_list<_Item_t> &_Src) noexcept {
        this->_buffer = std::vector<_Item_t>(_Src.size());
        const auto _Src_begin{_Src.begin()};
        for (uint_xt _index{0}; _index < _Src.size(); ++_index)
        { this->_buffer[_index] = *(_Item_t*)(_Src_begin+_index); }
    }

    ~array<_Item_t>(void) noexcept
    { this->_buffer.clear(); }

    typedef _Item_t       *iterator;
    typedef const _Item_t *const_iterator;

    inline constexpr
    iterator begin(void) noexcept
    { return &this->_buffer[0]; }

    inline constexpr
    const_iterator begin(void) const noexcept
    { return &this->_buffer[0]; }

    inline constexpr
    iterator end(void) noexcept
    { return &this->_buffer[this->_buffer.size()]; }

    inline constexpr
    const_iterator end(void) const noexcept
    { return &this->_buffer[this->_buffer.size()]; }

    inline constexpr
    uint_xt len(void) const noexcept
    { return this->_buffer.size(); }

    inline bool empty(void) const noexcept
    { return this->_buffer.empty(); }

    _Item_t *find(const _Item_t &_Item) noexcept {
        iterator _it{this->begin()};
        const iterator _end{this->end()};
        for (; _it < _end; ++_it)
        { if (_Item == *_it) { return _it; } }
        return nil;
    }

    _Item_t *rfind(const _Item_t &_Item) noexcept {
        iterator _it{this->end()};
        const iterator _begin{this->begin()};
        for (; _it >= _begin; --_it)
        { if (_Item == *_it) { return _it; } }
        return nil;
    }

    bool operator==(const array<_Item_t> &_Src) const noexcept {
        const uint_xt _length{this->_buffer.size()};
        const uint_xt _Src_length{_Src._buffer.size()};
        if (_length != _Src_length) { return false; }
        for (uint_xt _index{0}; _index < _length; ++_index)
        { if (this->_buffer[_index] != _Src._buffer[_index]) { return false; } }
        return true;
    }

    inline constexpr
    bool operator!=(const array<_Item_t> &_Src) const noexcept
    { return !this->operator==(_Src); }

    inline constexpr
    bool operator==(const std::nullptr_t) const noexcept
    { return this->_buffer.empty(); }

    inline constexpr
    bool operator!=(const std::nullptr_t) const noexcept
    { return !this->operator==(nil); }

    _Item_t& operator[](const uint_xt _Index) {
        if (this->len() <= _Index) { XID(panic)("index out of range"); }
        return this->_buffer[_Index];
    }

    friend std::ostream& operator<<(std::ostream &_Stream,
                                    const array<_Item_t> &_Src) {
        _Stream << '[';
        const uint_xt _length{_Src._buffer.size()};
        for (uint_xt _index{0}; _index < _length;) {
            _Stream << _Src._buffer[_index++];
            if (_index < _length) { _Stream << u8", "; }
        }
        _Stream << ']';
        return _Stream;
    }
};

template<typename _Key_t, typename _Value_t>
class map: public std::unordered_map<_Key_t, _Value_t> {
public:
    map<_Key_t, _Value_t>(void) noexcept                 {}
    map<_Key_t, _Value_t>(const std::nullptr_t) noexcept {}
    map<_Key_t, _Value_t>(const std::initializer_list<std::pair<_Key_t, _Value_t>> _Src)
    { for (const auto _data: _Src) { this->insert(_data); } }

    array<_Key_t> keys(void) const noexcept {
        array<_Key_t> _keys(this->size());
        uint_xt _index{0};
        for (const auto &_pair: *this)
        { _keys._buffer[_index++] = _pair.first; }
        return _keys;
    }

    array<_Value_t> values(void) const noexcept {
        array<_Value_t> _keys(this->size());
        uint_xt _index{0};
        for (const auto &_pair: *this)
        { _keys._buffer[_index++] = _pair.second; }
        return _keys;
    }

    inline constexpr
    bool has(const _Key_t _Key) const noexcept
    { return this->find(_Key) != this->end(); }

    inline void del(const _Key_t _Key) noexcept
    { this->erase(_Key); }

    inline bool operator==(const std::nullptr_t) const noexcept
    { return this->empty(); }

    inline bool operator!=(const std::nullptr_t) const noexcept
    { return !this->operator==(nil); }

    friend std::ostream& operator<<(std::ostream &_Stream,
                                    const map<_Key_t, _Value_t> &_Src) {
        _Stream << '{';
        uint_xt _length{_Src.size()};
        for (const auto _pair: _Src) {
            _Stream << _pair.first;
            _Stream << ':';
            _Stream << _pair.second;
            if (--_length > 0) { _Stream << u8", "; }
        }
        _Stream << '}';
        return _Stream;
    }
};
// endregion X_STRUCTURES

class str_xt {
public:
    std::string _buffer{};

    str_xt(void) noexcept                   {}
    str_xt(const char *_Src) noexcept       { this->_buffer = _Src ? _Src : ""; }
    str_xt(const std::string _Src) noexcept { this->_buffer = _Src; }
    str_xt(const str_xt &_Src) noexcept     { this->_buffer = _Src._buffer; }
    
    str_xt(const array<char> &_Src) noexcept
    { this->_buffer = std::string{_Src.begin(), _Src.end()}; }

    str_xt(const array<u8_xt> &_Src) noexcept
    { this->_buffer = std::string{_Src.begin(), _Src.end()}; }

    typedef char_xt       *iterator;
    typedef const char_xt *const_iterator;

    inline iterator begin(void) noexcept
    { return (iterator)(&this->_buffer[0]); }

    inline const_iterator begin(void) const noexcept
    { return (const_iterator)(&this->_buffer[0]); }

    inline iterator end(void) noexcept
    { return (iterator)(&this->_buffer[this->len()]); }

    inline const_iterator end(void) const noexcept
    { return (const_iterator)(&this->_buffer[this->len()]); }

    inline uint_xt len(void) const noexcept
    { return this->_buffer.length(); }

    inline bool empty(void) const noexcept
    { return this->_buffer.empty(); }

    inline str_xt sub(const uint_xt start, const uint_xt end) const noexcept
    { return this->_buffer.substr(start, end); }

    inline str_xt sub(const uint_xt start) const noexcept
    { return this->_buffer.substr(start); }

    inline bool has_prefix(const str_xt &_Sub) const noexcept
    { return this->len() >= _Sub.len() && this->sub(0, _Sub.len()) == _Sub._buffer; }

    inline bool has_suffix(const str_xt &_Sub) const noexcept
    { return this->len() >= _Sub.len() && this->sub(this->len()-_Sub.len()) == _Sub; }

    inline uint_xt find(const str_xt &_Sub) const noexcept
    { return this->_buffer.find(_Sub._buffer); }

    inline uint_xt rfind(const str_xt &_Sub) const noexcept
    { return this->_buffer.rfind(_Sub._buffer); }

    inline const char* cstr(void) const noexcept
    { return this->_buffer.c_str(); }

    str_xt trim(const str_xt &_Bytes) const noexcept {
        const_iterator _it{this->begin()};
        const const_iterator _end{this->end()};
        const_iterator _begin{this->begin()};
        for (; _it < _end; ++_it) {
            bool exist{false};
            const_iterator _bytes_it{_Bytes.begin()};
            const const_iterator _bytes_end{_Bytes.end()};
            for (; _bytes_it < _bytes_end; ++_bytes_it)
            { if ((exist = *_it == *_bytes_it)) { break; } }
            if (!exist) { return this->sub(_it-_begin); }
        }
        return str_xt{u8""};
    }

    str_xt rtrim(const str_xt &_Bytes) const noexcept {
        const_iterator _it{this->end()-1};
        const const_iterator _begin{this->begin()};
        for (; _it >= _begin; --_it) {
            bool exist{false};
            const_iterator _bytes_it{_Bytes.begin()};
            const const_iterator _bytes_end{_Bytes.end()};
            for (; _bytes_it < _bytes_end; ++_bytes_it)
            { if ((exist = *_it == *_bytes_it)) { break; } }
            if (!exist) { return this->sub(0, _it-_begin+1); }
        }
        return str_xt{u8""};
    }

    array<str_xt> split(const str_xt &_Sub, const i64_xt &_N) const noexcept {
        array<str_xt> _parts{};
        if (_N == 0) { return _parts; }
        const const_iterator _begin{this->begin()};
        std::string _s{this->_buffer};
        uint_xt _pos{std::string::npos};
        if (_N < 0) {
            while ((_pos = _s.find(_Sub._buffer)) != std::string::npos) {
                _parts._buffer.push_back(_s.substr(0, _pos));
                _s = _s.substr(_pos+_Sub.len());
            }
            if (!_parts.empty()) { _parts._buffer.push_back(str_xt{_s}); }
        } else {
            uint_xt _n{0};
            while ((_pos = _s.find(_Sub._buffer)) != std::string::npos) {
                _parts._buffer.push_back(_s.substr(0, _pos));
                _s = _s.substr(_pos+_Sub.len());
                if (++_n >= _N) { break; }
            }
            if (!_parts.empty() && _n < _N) { _parts._buffer.push_back(str_xt{_s}); }
        }
        return _parts;
    }

    str_xt replace(const str_xt &_Sub,
                   const str_xt &_New,
                   const i64_xt &_N) const noexcept {
        if (_N == 0) { return *this; }
        std::string _s{this->_buffer};
        uint_xt start_pos{0};
        if (_N < 0) {
            while((start_pos = _s.find(_Sub._buffer, start_pos)) != std::string::npos) {
                _s.replace(start_pos, _Sub.len(), _New._buffer);
                start_pos += _New.len();
            }
        } else {
            uint_xt _n{0};
            while((start_pos = _s.find(_Sub._buffer, start_pos)) != std::string::npos) {
                _s.replace(start_pos, _Sub.len(), _New._buffer);
                start_pos += _New.len();
                if (++_n >= _N) { break; }
            }
        }
        return str_xt{_s};
    }

    operator array<char>(void) const noexcept {
        array<char> _array{};
        _array._buffer = std::vector<char>{this->begin(), this->end()};
        return _array;
    }

    operator array<u8_xt>(void) const noexcept {
        array<u8_xt> _array{};
        _array._buffer = std::vector<u8_xt>{this->begin(), this->end()};
        return _array;
    }

    inline operator const char*(void) const noexcept
    { return this->_buffer.c_str(); }
    
    inline operator char*(void) const noexcept
    { return (char*)(this->_buffer.c_str()); }

    char &operator[](uint_xt _Index) {
        if (this->len() <= _Index) { XID(panic)("index out of range"); }
        return this->_buffer[_Index];
    }

    inline void operator+=(const str_xt &_Str) noexcept
    { this->_buffer += _Str._buffer; }

    inline str_xt operator+(const str_xt &_Str) const noexcept
    { return str_xt{this->_buffer + _Str._buffer}; }

    inline bool operator==(const str_xt &_Str) const noexcept
    { return this->_buffer == _Str._buffer; }

    inline bool operator!=(const str_xt &_Str) const noexcept
    { return !this->operator==(_Str); }

    friend std::ostream& operator<<(std::ostream &_Stream, const str_xt &_Src)
    { return _Stream << _Src._buffer; }
};

struct any_xt {
public:
    void *_expr{nil};
    char *_inf{nil};

    template<typename T>
    any_xt(const T &_Expr) noexcept
    { this->operator=(_Expr); }

    ~any_xt(void) noexcept
    { this->_delete(); }

    inline void _delete(void) noexcept {
        this->_expr = nil;
        this->_inf = nil;
    }

    template<typename T>
    inline bool type_is(void) const noexcept {
        if (!this->_expr)
        { return std::is_same<std::nullptr_t, T>::value; }
        return std::strcmp(this->_inf, typeid(T).name()) == 0;
    }

    template<typename T>
    void operator=(const T &_Expr) noexcept {
        this->_delete();
        this->_expr = (void*)&_Expr;
        this->_inf  = (char*)(typeid(T).name());
    }

    inline void operator=(const std::nullptr_t) noexcept
    { this->_delete(); }

    template<typename T>
    operator T(void) const noexcept {
        if (!this->_expr)
        { XID(panic)("casting failed because data is nil"); }
        if (std::strcmp(this->_inf, typeid(T).name()) != 0)
        { XID(panic)("incompatible type"); }
        return *(T*)(this->_expr);
    }

    template<typename T>
    inline bool operator==(const T &_Expr) const noexcept
    { return this->type_is<T>() && *(T*)(this->_expr) == _Expr; }

    template<typename T>
    inline constexpr
    bool operator!=(const T &_Expr) const noexcept
    { return !this->operator==(_Expr); }

    inline constexpr
    bool operator==(const any_xt &_Any) const noexcept
    { return this->_expr == _Any._expr; }

    inline constexpr
    bool operator!=(const any_xt &_Any) const noexcept
    { return !this->operator==(_Any); }

    friend std::ostream& operator<<(std::ostream &_Stream, const any_xt &_Src)
    { return _Stream << _Src._expr; }
};
// endregion X_BUILTIN_TYPES

// region X_MISC
template <typename _Enum_t, typename _Index_t, typename _Item_t>
static inline void foreach(const _Enum_t _Enum,
                           const func<void(_Index_t, _Item_t)> _Body) {
    _Index_t _index{0};
    for (auto _item: _Enum) { _Body(_index++, _item); }
}

template <typename _Enum_t, typename _Index_t>
static inline void foreach(const _Enum_t _Enum,
                           const func<void(_Index_t)> _Body) {
    _Index_t _index{0};
    for (auto begin = _Enum.begin(), end = _Enum.end(); begin < end; ++begin)
    { _Body(_index++); }
}

template <typename _Key_t, typename _Value_t>
static inline void foreach(const map<_Key_t, _Value_t> _Map,
                           const func<void(_Key_t)> _Body) {
    for (const auto _pair: _Map) { _Body(_pair.first); }
}

template <typename _Key_t, typename _Value_t>
static inline void foreach(const map<_Key_t, _Value_t> _Map,
                           const func<void(_Key_t, _Value_t)> _Body) {
    for (const auto _pair: _Map) { _Body(_pair.first, _pair.second); }
}

template<typename ...T>
static inline std::string strpol(const T... _Expressions) noexcept {
    return (std::stringstream{} << ... << _Expressions).str();
}

template<typename Type, unsigned N, unsigned Last>
struct tuple_ostream {
    static void arrow(std::ostream &_Stream, const Type &_Type) {
        _Stream << std::get<N>(_Type) << u8", ";
        tuple_ostream<Type, N + 1, Last>::arrow(_Stream, _Type);
    }
};

template<typename Type, unsigned N>
struct tuple_ostream<Type, N, N> {
    static void arrow(std::ostream &_Stream, const Type &_Type)
    { _Stream << std::get<N>(_Type); }
};

template<typename... Types>
std::ostream& operator<<(std::ostream &_Stream,
                         const std::tuple<Types...> &_Tuple) {
    _Stream << u8"(";
    tuple_ostream<std::tuple<Types...>, 0, sizeof...(Types)-1>::arrow(_Stream, _Tuple);
    _Stream << u8")";
    return _Stream;
}

template<typename _Function_t, typename _Tuple_t, size_t ... _I_t>
inline auto tuple_as_args(const _Function_t _Function,
                          const _Tuple_t _Tuple,
                          const std::index_sequence<_I_t ...>)
{ return _Function(std::get<_I_t>(_Tuple) ...); }

template<typename _Function_t, typename _Tuple_t>
inline auto tuple_as_args(const _Function_t _Function, const _Tuple_t _Tuple) {
    static constexpr auto _size{std::tuple_size<_Tuple_t>::value};
    return tuple_as_args(_Function, _Tuple, std::make_index_sequence<_size>{});
}

struct defer {
    typedef func<void(void)> _Function_t;
    template<class Callable>
    defer(Callable &&_function): _function(std::forward<Callable>(_function)) {}
    defer(defer &&_Src): _function(std::move(_Src._function))                 { _Src._function = nullptr; }
    ~defer() noexcept                                                         { if (this->_function) { this->_function(); } }
    defer(const defer &)          = delete;
    void operator=(const defer &) = delete;
    _Function_t _function;
};

std::ostream &operator<<(std::ostream &_Stream, const i8_xt &_Src)
{ return _Stream << (i32_xt)(_Src); }

std::ostream &operator<<(std::ostream &_Stream, const u8_xt &_Src)
{ return _Stream << (i32_xt)(_Src); }

std::ostream &operator<<(std::ostream &_Stream, const std::nullptr_t)
{ return _Stream << "<nil>"; }

template<typename _Obj_t>
str_xt tostr(const _Obj_t &_Obj) noexcept {
    std::stringstream _stream;
    _stream << _Obj;
    return str_xt{_stream.str()};
}

#define DEFER(_Expr) defer CONCAT(XXDEFER_, __LINE__){[&](void) mutable -> void { _Expr; }}
#define CO(_Expr) std::thread{[&](void) mutable -> void { _Expr; }}.detach()
// endregion X_MISC

// region PANIC_DEFINES
struct XID(Error) {
public:
    str_xt XID(message);
};

std::ostream &operator<<(std::ostream &_Stream, const XID(Error) &_Error)
{ return _Stream << _Error.XID(message); }

static inline void XID(panic)(const struct XID(Error) &_Error) { throw _Error; }
static inline void XID(panic)(const char *_Message) { XID(panic)(XID(Error){_Message}); }
// endregion PANIC_DEFINES

// region X_BUILTIN_FUNCTIONS
template<typename _Obj_t>
static inline void XID(out)(const _Obj_t _Obj) noexcept { std::cout << _Obj; }

template<typename _Obj_t>
static inline void XID(outln)(const _Obj_t _Obj) noexcept {
    XID(out)<_Obj_t>(_Obj);
    std::cout << std::endl;
}
// endregion X_BUILTIN_FUNCTIONS

// region BOTTOM_MISC
void x_terminate_handler(void) noexcept {
    try { std::rethrow_exception(std::current_exception()); }
    catch (const XID(Error) _error)
    { std::cout << "panic: " << _error.XID(message) << std::endl; }
    catch (...)
    { std::cout << "panic: <undefined panics>" << std::endl; }
    std::exit(EXIT_FAILURE);
}
// endregion BOTTOM_MIST
// endregion X_CXX_API`
