// Copyright 2022 The Jule Programming Language.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

#ifndef __JULE_ANY_HPP
#define __JULE_ANY_HPP

#include <typeinfo>
#include <stddef.h>
#include <cstdlib>
#include <cstring>
#include <ostream>

#include "str.hpp"
#include "builtin.hpp"
#include "ref.hpp"

namespace jule {

    // Built-in any type.
    class Any;

    class Any {
    private:
        template<typename T>
        struct DynamicType {
        public:
            static const char *type_id(void)
            { return typeid(T).name(); }

            static void dealloc(void *alloc)
            { delete static_cast<T*>(alloc); }

            static jule::Bool eq(void *alloc, void *other) {
                T *l{ static_cast<T*>(alloc) };
                T *r{ static_cast<T*>(other) };
                return *l == *r;
            }

            static const jule::Str to_str(const void *alloc) {
                const T *v{ static_cast<const T*>(alloc) };
                return jule::to_str(*v);
            }

            static void *alloc_new_copy(void *data) {
                T *heap{ new (std::nothrow) T };
                if (!heap)
                    return nullptr;

                *heap = *static_cast<T*>(data);
                return heap;
            }
        };

        struct Type {
        public:
            const char*(*type_id)(void)  ;
            void(*dealloc)(void *alloc)  ;
            jule::Bool(*eq)(void *alloc, void *other)  ;
            const jule::Str(*to_str)(const void *alloc)  ;
            void *(*alloc_new_copy)(void *data)  ;
        };

        template<typename T>
        static jule::Any::Type *new_type(void) {
            using t = typename std::decay<DynamicType<T>>::type;
            static jule::Any::Type table = {
                t::type_id,
                t::dealloc,
                t::eq,
                t::to_str,
                t::alloc_new_copy,
            };
            return &table;
        }

    public:
        jule::Ref<void*> data{};
        jule::Any::Type *type{ nullptr };

        Any(void) {}

        template<typename T>
        Any(const T &expr)
        { this->operator=(expr); }

        Any(const jule::Any &src)
        { this->operator=(src); }

        Any(const std::nullptr_t)
        { this->operator=(nullptr); }

        ~Any(void)
        { this->dealloc(); }

        void dealloc(void) {
#ifdef __JULE_DISABLE__REFERENCE_COUNTING
            this->type = nullptr;
            this->data.drop();
#else
            if (!this->data.ref) {
                this->type = nullptr;
                this->data.drop();
                return;
            }

            // Use jule::REFERENCE_DELTA, DON'T USE drop_ref METHOD BECAUSE
            // jule_ref does automatically this.
            // If not in this case:
            //   if this is method called from destructor, reference count setted to
            //   negative integer but reference count is unsigned, for this reason
            //   allocation is not deallocated.
            if (this->data.get_ref_n() != jule::REFERENCE_DELTA) {
                this->type = nullptr;
                this->data.drop();
                return;
            }

            this->type->dealloc(*this->data.alloc);
            *this->data.alloc = nullptr;
            this->type = nullptr;

            delete this->data.ref;
            this->data.ref = nullptr;
            std::free(this->data.alloc);
            this->data.alloc = nullptr;
            this->data.drop();
#endif // __JULE_DISABLE__REFERENCE_COUNTING
        }

        template<typename T>
        inline jule::Bool type_is(void) const {
            if (std::is_same<typename std::decay<T>::type, std::nullptr_t>::value)
                return false;

            if (this->operator==(nullptr))
                return false;

            return std::strcmp(this->type->type_id(), typeid(T).name()) == 0;
        }

        template<typename T>
        void operator=(const T &expr) {
            this->dealloc();

            T *alloc{ new(std::nothrow) T };
            if (!alloc)
                jule::panic(jule::ERROR_MEMORY_ALLOCATION_FAILED);

            void **main_alloc{ new(std::nothrow) void* };
            if (!main_alloc)
                jule::panic(jule::ERROR_MEMORY_ALLOCATION_FAILED);

            *alloc = expr;
            *main_alloc = static_cast<void*>(alloc);
#ifdef __JULE_DISABLE__REFERENCE_COUNTING
            this->data = jule::Ref<void*>::make(main_alloc, nullptr);
#else
            this->data = jule::Ref<void*>::make(main_alloc);
#endif
            this->type = jule::Any::new_type<T>();
        }

        void operator=(const jule::Any &src) {
            // Assignment to itself.
            if (this->data.alloc != nullptr && this->data.alloc == src.data.alloc)
                return;

            if (src.operator==(nullptr)) {
                this->operator=(nullptr);
                return;
            }

            this->dealloc();

            void *new_heap{ src.type->alloc_new_copy(*src.data.alloc) };
            if (!new_heap)
                jule::panic(jule::ERROR_MEMORY_ALLOCATION_FAILED);

#ifdef __JULE_DISABLE__REFERENCE_COUNTING
            this->data = jule::Ref<void*>::make(new_heap, nullptr);
#else
            this->data = jule::Ref<void*>::make(new_heap);
#endif
            this->type = src.type;
        }

        inline void operator=(const std::nullptr_t)
        { this->dealloc(); }

        template<typename T>
        operator T(void) const {
            if (this->operator==(nullptr))
                jule::panic(jule::ERROR_INVALID_MEMORY);

            if (!this->type_is<T>())
                jule::panic(jule::ERROR_INCOMPATIBLE_TYPE);

            return *static_cast<T*>(*this->data.alloc);
        }

        template<typename T>
        inline jule::Bool operator==(const T &expr) const
        { return this->type_is<T>() && this->operator T() == expr; }

        template<typename T>
        inline constexpr
        jule::Bool operator!=(const T &expr) const
        { return !this->operator==(expr); }

        inline jule::Bool operator==(const jule::Any &other) const {
            // Break comparison cycle.
            if (this->data.alloc != nullptr && this->data.alloc == other.data.alloc)
                return true;

            if (this->operator==(nullptr))
                return other.operator==(nullptr);

            if (other.operator==(nullptr))
                return false;

            if (std::strcmp(this->type->type_id(), other.type->type_id()) != 0)
                return false;

            return this->type->eq(*this->data.alloc, *other.data.alloc);
        }

        inline jule::Bool operator!=(const jule::Any &other) const
        { return !this->operator==(other); }

        inline jule::Bool operator==(std::nullptr_t) const
        { return !this->data.alloc; }

        inline jule::Bool operator!=(std::nullptr_t) const
        { return !this->operator==(nullptr); }

        friend std::ostream &operator<<(std::ostream &stream,
                                        const jule::Any &src) {
            if (src.operator!=(nullptr))
                stream << src.type->to_str(*src.data.alloc);
            else
                stream << 0;
            return stream;
        }
    };

} // namespace jule

#endif // ifndef __JULE_ANY_HPP
